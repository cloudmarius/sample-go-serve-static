package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	greeting         string
	fileServerFolder string
)

func main() {

	greeting = os.Getenv("GREETING")
	if len(greeting) == 0 {
		greeting = "Hello"
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
	http.HandleFunc("/", okHandler)
	http.HandleFunc("/hi", hiHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/greet/howdy", howdyHandler)
	http.HandleFunc("/greet/ciao", ciaoHandler)

	address := "0.0.0.0:" + port
	log.Printf("Starting web server, listening on [%s], serving files from [%s], greeting with [%s]\n", address, fileServerFolder, greeting)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}

}

func okHandler(rw http.ResponseWriter, req *http.Request) {
}

func hiHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "%s", greeting)
}

func configHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/x-yaml")
	fmt.Fprintf(rw, "---\ngreeting: %s\nfileServerFolder: %s", greeting, fileServerFolder)
}

func howdyHandler(rw http.ResponseWriter, req *http.Request) {
	patchGreet(rw, req, "Howdy")
}

func ciaoHandler(rw http.ResponseWriter, req *http.Request) {
	patchGreet(rw, req, "Ciao")
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
	for i := 0; i < 100; i++ {
		err = os.WriteFile(fileServerFolder+fmt.Sprintf("/greet-%d.txt", i), []byte(greet), 0664)
		if err != nil {
			return err
		}
	}
	return nil
}
