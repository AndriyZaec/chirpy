package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app", http.StripPrefix("/app", fileServer))

	assetServer := http.FileServer(http.Dir("./assets"))
	mux.Handle("/app/assets", http.StripPrefix("/app/assets", assetServer))
	mux.HandleFunc("/healthz", healtzHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
