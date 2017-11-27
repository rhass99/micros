package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/matryer/respond"
)

type OK interface {
	OK() error
}

type User struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (u *User) OK() error {
	if u.Password != u.PasswordConfirm {
		err := errors.New("Passords don't match")
		return err
	}
	return nil
}

func decoder(r *http.Request, v OK) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return v.OK()
}

func handleIncoming(w http.ResponseWriter, r *http.Request) {
	var incomingUser User
	if err := decoder(r, &incomingUser); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
	}
	respond.With(w, r, http.StatusOK, &incomingUser)
}

type Adapter func(h http.Handler) http.Handler

func Logging(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Println(r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
}
func WithHeader(key string, value string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		})
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func main() {
	handler := http.HandlerFunc(handleIncoming)
	logger := log.New(os.Stdout, "server: ", log.Lshortfile)
	http.Handle("/", Adapt(handler, WithHeader("Rami", "Ahmed"), Logging(logger)))
	http.ListenAndServe(":8000", nil)
}
