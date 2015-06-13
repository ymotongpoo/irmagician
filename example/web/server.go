package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ymotongpoo/irmagician"
)

var (
	mux       *http.ServeMux
	fileCache = make(map[string][]byte)
	files     = []string{"on.json", "off.json"}
)

type APIRequest struct {
	Power bool `json:"power"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/push/v1", pushV1Handler)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		data, err := ioutil.ReadFile(filepath.Join(filepath.Join(cwd, "json"), f))
		if err != nil {
			log.Fatal(err)
		}
		fileCache[f] = data
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func pushV1Handler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println(string(data))
	var req APIRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ir, err := irmagician.NewIrMagician("/dev/ttyACM0", 9600, 1*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var resp []byte
	switch {
	case req.Power:
		resp, err = ir.PlayData(fileCache["on.json"])
	case !req.Power:
		resp, err = ir.PlayData(fileCache["off.json"])
	}
	fmt.Fprintf(w, string(resp))
}
