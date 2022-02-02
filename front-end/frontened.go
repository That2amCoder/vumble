//Make a basic front-end server that will server the static files in the static directory.
//The server should listen on port 440
//The server should serve the index.html file
//The server should serve the static directory
// it should use the cert.pem and key.pem files

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Listening on port :440")
	log.Fatal(http.ListenAndServeTLS(":440", "cert.pem", "key.pem", nil))
}
