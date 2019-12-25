package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// returns true if the given file or path exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// returns a random string of given length
func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// returns the given relative path as absolute path
func absPath(path string) string {
	if strings.Contains(path, "~") {
		usr, _ := user.Current()
		path = strings.Replace(path, "~", usr.HomeDir, 1)
	}

	path, _ = filepath.Abs(path)
	return path
}

// returns the progress as float if the line contains the progress
func getProgress(line string) (float32, error) {
	if strings.Contains(line, "validating tile") {
		words := strings.Split(line, " ")

		for i := 0; i < len(words); i++ {
			if strings.Contains(words[i], "/") {
				numbers := strings.Split(words[i], "/")
				a, err := strconv.Atoi(numbers[0])

				if err != nil {
					return 0, err
				}

				b, err := strconv.Atoi(numbers[1])

				if err != nil {
					return 0, err
				}

				return float32(a) / float32(b), nil
			}
		}
	}

	return 0, errors.New("no progress line")
}
