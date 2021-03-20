package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func main() {
	var err error
	nats_url := "localhost:4222"
	nc, err = nats.Connect(nats_url)
	if err != nil {
		log.Print(err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/publish/{subject}", natsPublish)
	r.Post("/request/{subject}", natsRequest)
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Print(err)
	}
}

func natsPublish(w http.ResponseWriter, r *http.Request) {
	subject := chi.URLParam(r, "subject")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = nc.Publish(subject, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func natsRequest(w http.ResponseWriter, r *http.Request) {
	subject := chi.URLParam(r, "subject")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	msg, err := nc.Request(subject, data, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		_, err := w.Write(msg.Data)
		if err != nil {
			log.Print(err)
		}
	}
}
