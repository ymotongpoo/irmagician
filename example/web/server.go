package main

import (
	"log"
	"net/http"

	"bitbucket.org/ymotongpoo/irmagician"
)

type ApiHandler struct {
	ir *irmagician.NewIrMagician
}

func (a *ApiHandler) NewApiHandler() ApiHandler {
	a.ir = irmagician.NewIrMagician("/dev/ttyACM0", 9600, irmagician.DefaultTimeout)
}

func (a *ApiHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func init() {

	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", http.FileServer(http.Dir("/static")))
}

func main() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponserWriter, r *http.Request) {

}
