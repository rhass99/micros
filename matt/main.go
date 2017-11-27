package main

import (
	"encoding/json"
	"errors"
	"net/http"

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

func decoder(r *http.Request, v interface{}) error {
	if obj, ok := v.(OK); ok {
		err := obj.OK()
		if err != nil {
			if err = json.NewDecoder(r.Body).Decode(obj); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleIncoming(w http.ResponseWriter, r *http.Request) {
	var incomingUser User
	if err := decoder(r, &incomingUser); err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
	}
	respond.With(w, r, http.StatusOK, &incomingUser)
}

func main() {
	http.HandleFunc("/", handleIncoming)
	http.ListenAndServe(":8080", nil)
}
