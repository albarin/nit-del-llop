package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/albarin/poster/pkg/poster"
	"github.com/albarin/poster/pkg/webhooks"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/generate", paint).Methods(http.MethodPost)
	router.HandleFunc("/download", download).Methods(http.MethodGet)

	server := &http.Server{Handler: router, Addr: ":" + os.Getenv("PORT")}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	poster, err := os.Open("cartel.png")
	if err != nil {
		log.Fatal(err)
	}
	defer poster.Close()

	w.Header().Set("Content-Type", "image/png")
	_, err = io.Copy(w, poster)
	if err != nil {
		log.Fatal(err)
	}
}

func paint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer r.Body.Close()

	var answers webhooks.Webhook
	err = json.Unmarshal(body, &answers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	cartel := poster.Parse(answers)

	err = poster.Run(
		cartel,
		"assets/images/background.png",
		"assets/images/logos.png",
		"assets/images/foto.png",
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Write([]byte("done!"))
}
