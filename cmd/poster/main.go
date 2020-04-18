package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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
	router.HandleFunc("/generate", generate).Methods(http.MethodPost)
	router.HandleFunc("/download", download).Methods(http.MethodGet)

	server := &http.Server{Handler: router, Addr: ":" + os.Getenv("PORT")}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	poster, err := os.Open("cartel.png")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
	}
	defer poster.Close()

	w.Header().Set("Content-Type", "image/png")
	_, err = io.Copy(w, poster)
	if err != nil {
		log.Println(err)
	}
}

func computeSignature(payload []byte, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write(payload)
	if err != nil {
		return "", err
	}

	return "sha256=" + base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func verifySignature(payload []byte, secret, receivedSignature string) (bool, error) {
	signature, err := computeSignature(payload, secret)
	if err != nil {
		return false, err
	}

	return signature == receivedSignature, nil
}

func generate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer r.Body.Close()

	ok, err := verifySignature(body, os.Getenv("SECRET_TOKEN"), r.Header.Get("Typeform-Signature"))
	if err != nil || !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err, ok)
		return
	}

	var answers webhooks.Webhook
	err = json.Unmarshal(body, &answers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
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
		log.Println(err)
		return
	}

	w.Write([]byte("done!"))
}
