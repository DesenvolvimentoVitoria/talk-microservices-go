package main

import (
	"encoding/json"
	"feedbacks/feedback"
	"github.com/codegangsta/negroni"
	"github.com/eminetto/talk-microservices-go/pkg/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fService := feedback.NewService()
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.NewLogger(),
	)
	r.Handle("/v1/feedback", n.With(
		negroni.HandlerFunc(middleware.IsAuthenticated()),
		negroni.Wrap(storeFeedback(fService)),
	)).Methods("POST", "OPTIONS")


	http.Handle("/", r)
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8082",//@TODO usar variável de ambiente
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func storeFeedback(fService feedback.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var param feedback.Feedback
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		param.Email = r.Header.Get("email")
		err = fService.Store(param)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	})
}