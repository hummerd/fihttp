package main

import (
	"code.google.com/p/hummerdk-log4go"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
)

type Config struct {
	Root       string
	Port       string
	FileRef    string
	StaticPath string
}

type PathData struct {
	Path string
}

type FileData struct {
	Name  string
	Size  int64
	IsDir bool
	Type  int
	Time  int64
}

const (
	ftDir   = 1
	ftVideo = 2
	ftOther = 10000
)

var conf Config
var videoSuffix = [...]string{"mov", "mkv", "avi", "mp4"}

func main() {
	if len(os.Args) < 2 {
		log4go.Exit("Specify config path as first argument")
	}

	confFile := os.Args[1]
	_, err := os.Stat(confFile)
	if err != nil {
		log4go.Exitf("Cannot find conf file: %s args: %s", err, os.Args)
	}

	file, err := os.Open(confFile)
	if err != nil {
		log4go.Exitf("Unable to open conf file: %s args: %s", err, os.Args)
	}

	decoder := json.NewDecoder(file)
	conf = Config{}
	err = decoder.Decode(&conf)

	file.Close()
	if _, ok := err.(*json.InvalidUnmarshalError); err != nil && !ok {
		log4go.Exitf("Unable to parse conf: %s args: %s", err, os.Args)
	}

	http.HandleFunc("/", app)
	http.HandleFunc("/cont", cont)
	http.HandleFunc("/stop", stop)
	err = http.ListenAndServe(":"+conf.Port, nil)
	if err != nil {
		log4go.Exitf("ListenAndServe fails %s", err)
	}
}

func app(w http.ResponseWriter, r *http.Request) {
	ap := path.Join(conf.StaticPath, "index.html")
	http.ServeFile(w, r, ap)
}

func stop(w http.ResponseWriter, r *http.Request) {
	os.Exit(12)
}

func cont(w http.ResponseWriter, r *http.Request) {
	argDecoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	reqParams := PathData{}
	err := argDecoder.Decode(&reqParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log4go.Error("Can not decode request params %s", err)
		return
	}

	ap := conf.Root
	reqParams.Path = strings.Trim(reqParams.Path, "/")
	if reqParams.Path != "" && reqParams.Path != "root" {
		rp := strings.Replace(reqParams.Path, "root/", "", 1)
		ap = path.Join(ap, rp)
	} else {
		reqParams.Path = "root"
	}

	f, err := os.Open(ap)
	defer f.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log4go.Error("Can not open dir %s", err)
		return
	}

	firstTurn := true
	decoder := json.NewEncoder(w)

	fmt.Fprint(w, "{\n\"Current\": ")
	decoder.Encode(reqParams.Path)

	fmt.Fprint(w, ",\n\"FileRef\": ")
	decoder.Encode(conf.FileRef)

	fmt.Fprint(w, ",\n\"Files\": [")

	for infos, err := f.Readdir(10); err == nil; infos, err = f.Readdir(10) {
		if !firstTurn {
			fmt.Fprint(w, ",")
		}
		rc := len(infos)
		for i, info := range infos {
			fn := info.Name()
			decoder.Encode(FileData{fn, info.Size(), info.IsDir(), fileType(fn), info.ModTime().Unix()})
			if i+1 < rc {
				fmt.Fprint(w, ",")
			}
		}
		firstTurn = false
	}

	fmt.Fprint(w, "]}")
	if err != nil && err == io.EOF {
		log4go.Error("Can not read dir %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fileType(fileName string) int {
	for _, s := range videoSuffix {
		if strings.HasSuffix(fileName, s) {
			return ftVideo
		}
	}
	return ftOther
}
