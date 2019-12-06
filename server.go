package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type DeletableFile struct {
	Name  string
	Added time.Time
}

func NewDeletableFile(name string) *DeletableFile {
	return &DeletableFile{Name: name, Added: time.Now()}
}

type Server struct {
	Port   string
	Router *mux.Router
	// key is the id of a stream
	Streams       map[string]chan string
	MaxUploadSize int64
	DeletionQueue chan *DeletableFile
	TmpDir        string
}

func NewServer() *Server {
	s := Server{
		Port:          "8081",
		Router:        mux.NewRouter(),
		Streams:       make(map[string]chan string),
		MaxUploadSize: 5 * 1024 * 1024, // 5mb
		DeletionQueue: make(chan *DeletableFile, 500),
		TmpDir:        "tmp/"}

	s.Router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		s.versionHandler(w, r)
	}).Methods("GET")

	s.Router.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		s.runHandler(w, r)
	}).Methods("POST")

	s.Router.HandleFunc("/out/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.streamHandler(w, r)
	}).Methods("GET")

	s.Router.HandleFunc("/log/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.fileHandler(w, r)
	}).Methods("GET")

	s.Router.
		PathPrefix("/").
		Handler(http.FileServer(http.Dir("./public/")))

	return &s
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func (s *Server) Start() {
	go cleaner(s)
	Info("Starting server on port", s.Port)
	http.ListenAndServe(":"+s.Port, s.Router)
}

func (s *Server) AddStream(id string, stream chan string) {
	if _, ok := s.Streams[id]; ok {
		// TODO return error
	}

	s.Streams[id] = stream
}

func (s *Server) RemoveStream(id string) {
	delete(s.Streams, id)
}

func (s *Server) GetStream(id string) chan string {
	return s.Streams[id]
}

func cleaner(s *Server) {
	for {
		file, isOpen := <-s.DeletionQueue

		if !isOpen {
			return
		}

		duration := time.Since(file.Added)
		if duration.Minutes() > 15 {
			Info("Removing file", file.Name)
			os.Remove(file.Name)
		} else {
			s.DeletionQueue <- file
			// only wait if there was nothing to do
			// if the channel is empty the call will block there
			time.Sleep(5 * time.Minute)
		}
	}
}

func (s *Server) versionHandler(w http.ResponseWriter, r *http.Request) {
    Info("getting version")
	cmd := exec.Command("go", "version")

	out, err := cmd.StdoutPipe()
	if err != nil {
		Error("Unable to get stdout:", err.Error())
		return
	}

	scanner := bufio.NewScanner(out)
	go func() {
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			write(w, "unable to get version: "+err.Error(), false)
		} else {
			msg := scanner.Text()
			words := strings.Split(msg, " ")
			version := words[len(words)-2]
            Debugf("go version: %s", version)
			write(w, version, true)
		}
	}()

	if cmd.Start(); err != nil {
		write(w, "unable to get version: "+err.Error(), false)
	}

	if err = cmd.Wait(); err != nil {
		write(w, "unable to get version: "+err.Error(), false)
	}
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	file := "./tmp/" + id + ".log"

	if ok, _ := exists(file); !ok {
		http.Error(w, "Requested file does not exists", 404)
		return
	}

	log, err := os.Open(file)
	if err != nil {
		Error("File", file, "does not exist")
	}
	defer log.Close()

	fi, err := log.Stat()
	if err != nil {
		Error("Unable to get file stats:", err.Error())
	}

	w.Header().Set("Content-Disposition", "attachment; filename=output.log")
	w.Header().Set("Content-Type", "data:text/plain")
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

	_, err = io.Copy(w, log)
	if err != nil {
		Error("Unable to copy file:", err.Error())
	}
}

func (s *Server) streamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, ok := s.Streams[id]; !ok {
		write(w, "Stream "+id+" does not exists", false)
		return
	}

	stream := s.GetStream(id)
	line, hasNext := <-stream

	if hasNext {
		write(w, line, true)
	} else {
		Info("Stream", id, "is closed")
		// TODO there must be a difference between error and finished
		write(w, "all messages received", false)
		Info("Removing stream", id)
		s.RemoveStream(id)
	}
}

func (s *Server) runHandler(w http.ResponseWriter, r *http.Request) {
	// maximum file size
	r.Body = http.MaxBytesReader(w, r.Body, s.MaxUploadSize)
	if err := r.ParseMultipartForm(s.MaxUploadSize); err != nil {
		Error("File too large:", err.Error())
		http.Error(w, "file too large: "+err.Error(), 400)
		return
	}

	// parse and validate file and post parameters
	file, header, err := r.FormFile("file")
	Info("Received file:", header.Filename, "from", r.RemoteAddr)
	if err != nil {
		Error("Invalid file:", err.Error())
		http.Error(w, "invalid file: "+err.Error(), 400)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		Error("Invalid file", err.Error())
		http.Error(w, "invalid file: "+err.Error(), 400)
		return
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	filetype := http.DetectContentType(fileBytes)
	if !strings.Contains(filetype, "text/plain") && !strings.Contains(filetype, "application/octet-stream") {
		http.Error(w, "invalid file type: "+filetype, 400)
		return
	}

	id := randToken(16)
	newPath := filepath.Join(s.TmpDir, id+".go")
	if ok, _ := exists(s.TmpDir); !ok {
		if err := os.MkdirAll(s.TmpDir, os.ModePerm); err != nil {
			Errorf("unable for create %s: %s", s.TmpDir, err.Error())
			return
		}
	}

	// write file
	Infof("Creating path: %s", newPath)
	newFile, err := os.Create(newPath)
	if err != nil {
		Error("Unable to store file:", err.Error())
		http.Error(w, "unable to store file: "+err.Error(), 500)
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil {
		Error("Unable to store file:", err.Error())
		http.Error(w, "unable to store file: "+err.Error(), 500)
		return
	}

	stream := make(chan string, 100)
	s.AddStream(id, stream)
	Info("running file:", header.Filename)
	go worker(id, s)
	write(w, id, true)
}

func streamReader(scanner *bufio.Scanner, logFile *os.File, stream chan string) {
	for scanner.Scan() {
		line := fmt.Sprintf("%s\n", scanner.Text())
		fmt.Fprintf(logFile, line)
		stream <- line
	}

	if err := scanner.Err(); err != nil {
		Error("Error reading:", err.Error())
	}
}

func worker(id string, s *Server) {
	stream := s.GetStream(id)
	if stream == nil {
		Errorf("Stream %s is not available", id)
		return
	}
    defer close(stream)

	file := s.TmpDir + "/" + id + ".go"
	log, err := os.Create(s.TmpDir + "/" + id + ".log")
	if err != nil {
		Error("Cannot create file:", err.Error())
	}
	defer log.Close()

	cmd := exec.Command("go", "run", file)

	out, err := cmd.StdoutPipe()
	if err != nil {
		Error("Unable to get stdout:", err.Error())
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		Error("Unable to get stderr:", err.Error())
		return
	}

	go streamReader(bufio.NewScanner(out), log, stream)
	go streamReader(bufio.NewScanner(stderr), log, stream)

	if cmd.Start(); err != nil {
		Error("Unable to run command:", err.Error())
	}
	cmd.Wait()

	glob := file + "*"
	files, err := filepath.Glob(glob)
	if err != nil {
		Error(err.Error())
		Info("Removing file", file)
		if err := os.Remove(file); err != nil {
			Errorf("Error removing file %s: %s", file, err.Error())
		}
        return
	}

    if len(files) == 0 {
        Warning("No files matching: ", glob)
    }

    for _, f := range files {
        Info("Removing file", f)
        if err := os.Remove(f); err != nil {
            Errorf("Error removing file %s: %s", file, err.Error())
        }
    }

	Infof("Adding %s to the deletion queue", log.Name())
	s.DeletionQueue <- NewDeletableFile(log.Name())
}

func write(w http.ResponseWriter, msg string, rc bool) {
	json, err := NewResponse(msg, rc).toJSON()

	if err != nil {
		Error("Error converting to JSON:", err.Error())
	}

    if !rc {
        Error(msg)
    }
	fmt.Fprintf(w, string(json))
}
