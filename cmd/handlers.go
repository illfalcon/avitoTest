package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/illfalcon/avitoTest/internal/validation"

	"github.com/illfalcon/avitoTest/pkg/weberrors"
	"github.com/pkg/errors"

	"github.com/go-chi/chi"

	"github.com/illfalcon/avitoTest/internal/interactions"
)

const jsonType string = "application/json"

func switchErrors(w http.ResponseWriter, err error) {
	//log.Println(w.Header())
	var e weberrors.Err
	switch errors.Cause(err) {
	case weberrors.ErrBadRequest:
		e.Error = errors.Cause(err).Error()
		jsonErr, _ := json.Marshal(e)
		//http.Error(w, string(jsonErr), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(jsonErr)
		return
	case weberrors.ErrForbidden:
		e.Error = errors.Cause(err).Error()
		jsonErr, _ := json.Marshal(e)
		//http.Error(w, string(jsonErr), http.StatusForbidden)
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write(jsonErr)
	case weberrors.ErrNotFound:
		e.Error = errors.Cause(err).Error()
		jsonErr, _ := json.Marshal(e)
		//http.Error(w, string(jsonErr), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(jsonErr)
	case weberrors.ErrConflict:
		e.Error = errors.Cause(err).Error()
		jsonErr, _ := json.Marshal(e)
		//http.Error(w, string(jsonErr), http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write(jsonErr)
	default:
		e.Error = errors.Cause(err).Error()
		jsonErr, _ := json.Marshal(e)
		//http.Error(w, string(jsonErr), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(jsonErr)
		return
	}
}

func isJSONRequest(w http.ResponseWriter, r *http.Request) bool {
	var e weberrors.Err
	ct := r.Header.Get("Content-Type")
	log.Println(ct)
	if ct != jsonType {
		e.Error = "wrong content type"
		jsonErr, _ := json.Marshal(e)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(jsonErr)
		return false
	}
	return true
}

type Handlers struct {
	interactor interactions.Interactor
}

func (h Handlers) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonType)
	if !isJSONRequest(w, r) {
		return
	}
	var nu validation.NewUser
	err := validation.Decode(r, &nu)
	if err != nil {
		log.Println(err)
		switchErrors(w, err)
		return
	}
	resp, err := h.interactor.AddUser(nu.Username)
	if err != nil {
		log.Println(err)
		switchErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h Handlers) CreateChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonType)
	if !isJSONRequest(w, r) {
		return
	}
	var nc validation.NewChat
	err := validation.Decode(r, &nc)
	if err != nil {
		switchErrors(w, err)
		return
	}
	resp, err := h.interactor.CreateChat(nc.Name, nc.DecodedUsers)
	if err != nil {
		switchErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h Handlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonType)
	if !isJSONRequest(w, r) {
		return
	}
	var sm validation.NewMessage
	err := validation.Decode(r, &sm)
	if err != nil {
		switchErrors(w, err)
		return
	}
	resp, err := h.interactor.SendMessage(sm.DecodedChat, sm.DecodedAuthor, sm.Text)
	if err != nil {
		switchErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h Handlers) GetChatsOfUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonType)
	if !isJSONRequest(w, r) {
		return
	}
	var gc validation.GetChats
	err := validation.Decode(r, &gc)
	if err != nil {
		switchErrors(w, err)
		return
	}
	resp, err := h.interactor.GetChatsOfUser(gc.DecodedUser)
	if err != nil {
		switchErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h Handlers) GetMessagesInChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonType)
	if !isJSONRequest(w, r) {
		return
	}
	var gm validation.GetMessages
	err := validation.Decode(r, &gm)
	if err != nil {
		switchErrors(w, err)
		return
	}
	resp, err := h.interactor.GetMessagesInChat(gm.DecodedChat)
	if err != nil {
		switchErrors(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h Handlers) Start(cl func() error) {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/users/add", h.AddUser)
		r.Post("/chats/add", h.CreateChat)
		r.Post("/messages/add", h.SendMessage)
		r.Post("/chats/get", h.GetChatsOfUser)
		r.Post("/messages/get", h.GetMessagesInChat)
	})
	srv := &http.Server{Addr: ":9000", Handler: r}
	stopAppCh := make(chan struct{})
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigquit
		fmt.Printf("captured signal: %v\n", s)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("could not shutdown server: %s", err)
		}
		fmt.Println("shut down gracefully")
		_ = cl()
		stopAppCh <- struct{}{}
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-stopAppCh
}
