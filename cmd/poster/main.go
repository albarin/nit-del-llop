package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/albarin/poster/pkg/poster"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	port        = "PORT"
	secretToken = "SECRET_TOKEN"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	router := mux.NewRouter()
	router.HandleFunc("/generate", generate).Methods(http.MethodPost)
	router.HandleFunc("/download", download).Methods(http.MethodGet)

	server := &http.Server{Handler: router, Addr: ":" + os.Getenv(port)}
	if err := server.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("server failed")
	}
}

func generate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("could not read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	ok, err := verifySignature(body, os.Getenv(secretToken), r.Header.Get("Typeform-Signature"))
	if err != nil || !ok {
		log.WithFields(log.Fields{"error": err}).Error("could not verify signature")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var wh poster.Webhook
	err = json.Unmarshal(body, &wh)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("could not unmarshal webhook")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = poster.Draw(wh.Parse())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("could not generate poster")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	image, err := os.Open("cartel.png")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("could not open image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer image.Close()

	w.Header().Set("Content-Type", "image/png")
	_, err = io.Copy(w, image)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("could not write image")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func verifySignature(payload []byte, secret, receivedSignature string) (bool, error) {
	if secret == "" {
		return true, nil
	}

	signature, err := computeSignature(payload, secret)
	if err != nil {
		return false, err
	}

	return signature == receivedSignature, nil
}

func computeSignature(payload []byte, secret string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write(payload)
	if err != nil {
		return "", err
	}

	return "sha256=" + base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
