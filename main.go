package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var greeting string
var fileServerFolder string

func main() {

	greeting = os.Getenv("GREETING")
	if len(greeting) == 0 {
		greeting = "Howdy"
	}
	fileServerFolder = os.Getenv("FILE_SERVER_FOLDER")
	if len(fileServerFolder) == 0 {
		fileServerFolder = "/var/www"
	}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	fs := http.FileServer(http.Dir(fileServerFolder))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/hi", hiHandler)
	http.HandleFunc("/greet/howdy", howdyHandler)
	http.HandleFunc("/greet/ciao", ciaoHandler)

	address := "0.0.0.0:" + port
	log.Printf("Starting web server, listening on [%s] ...\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}

}

func hiHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "%s\n", greeting)
}

func howdyHandler(rw http.ResponseWriter, req *http.Request) {
	patchGreet(rw, req, "Howdy\n")
}

func ciaoHandler(rw http.ResponseWriter, req *http.Request) {
	patchGreet(rw, req, "Ciao\n")
}

func patchGreet(rw http.ResponseWriter, req *http.Request, greet string) {
	if req.Method == "PATCH" {
		err := writeGreet(greet)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}
		rw.WriteHeader(http.StatusAccepted)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func writeGreet(greet string) error {
	err := os.MkdirAll(fileServerFolder, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = os.WriteFile(fileServerFolder+"/greet.txt", []byte(greet), 0640)
	if err != nil {
		return err
	}
	return nil
}
