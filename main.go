package main

import (
	"bytes"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	mux := http.NewServeMux()
	webpageRoot := "./webpages/"
	assetsRoot := "./assets/"
	webpageDirHandler := http.FileServer(http.Dir(webpageRoot))
	assetsDirHandler := http.FileServer(http.Dir(assetsRoot))
	mux.Handle("/app/", http.StripPrefix("/app/", webpageDirHandler))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", assetsDirHandler))
	mux.HandleFunc("/healthz", healthCheckHandler)

	corsMux := middlewareCors(mux)

	address := "localhost"
	port := "8080"
	srv := &http.Server{
		Addr:    address + ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// Middleware functions

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareRequestIncrementor(next http.Handler) http.Handler {

}

// Route Handlers

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ok := []byte("OK")
	okBuffer := bytes.NewBuffer(ok)
	_, err := w.Write(okBuffer.Bytes())
	if err != nil {
		log.Println(err)
		_, err := w.Write([]byte("Internal Server Error"))
		if err != nil {
			log.Println(err)
		}
	}
}
