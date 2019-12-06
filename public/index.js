function sendData() {
    var form = document.getElementById("form");

    var file = document.getElementById("file");
    if (file.value === "") {
        alert("No file selected");
        return;
    }

    var request = new XMLHttpRequest();
    // Bind the FormData object and the form element
    var FD = new FormData(form);

    request.onload = function () {
        if (request.status === 200) {
            response = JSON.parse(request.responseText);
            if (response.success) {
                listener(response.message);
            } else {
                alert(response.message);
            }
        } else {
            alert(request.responseText);
        }
    };

    request.upload.onprogress = function(e) {
        var v = (e.loaded/e.total)*100 ;
        var p = document.createElement('progress');
        p.max = 100;
        p.value = v;
        document.getElementById("status").innerHTML = "uploading ";
        document.getElementById("status").appendChild(p);
    };

    request.onerror = function(e) {
        alert('Error while sending map data');
    };

    request.open("POST", "run");
    request.send(FD);
}

function onLoad() {
    form.addEventListener("submit", function (event) {
        event.preventDefault();
        sendData();
    });
    getVersion();
}

function getBaseUrl() {
    var getUrl = window.location;
    return getUrl.protocol + "//" + getUrl.host + "/" + getUrl.pathname.split('/')[1];
}

function getVersion(id) {
    var request = new XMLHttpRequest();
    request.open("GET", "version");
    request.onload = function () {
        if (request.status === 200) {
            response = JSON.parse(request.responseText);

            if (response.success) {
                var out = document.getElementById("version");
                out.appendChild(document.createTextNode("version " + response.message));
            } else {
                alert(response.message);
            }
        } else {
            alert(request.responseText);
        }
    };
    request.send();
}

function listener(id) {
    var request = new XMLHttpRequest();
    request.open("GET", "out/" + id);
    request.onload = function () {
        if (request.status === 200) {
            response = JSON.parse(request.responseText);

            if (response.success) {
                document.getElementById("output").appendChild(document.createTextNode(response.message))
                document.getElementById("output").appendChild(document.createElement("br"));
                listener(id);
            } else {
                if (response.message == "all messages received") {
                    var a = document.createElement('a');
                    a.href = getBaseUrl() + "/log/" + id;
                    a.target = "_blank";
                    a.download = "output.log";
                    a.appendChild(document.createTextNode("Download log file"));
                    document.getElementById("status").innerHTML = "";
                    document.getElementById("status").appendChild(a);
                } else {
                    alert(response.message);
                    document.getElementById("status").innerHTML = "";
                }
            }
        } else {
            alert(request.responseText);
            document.getElementById("status").innerHTML = "";
        }
    };

    request.send();
}
