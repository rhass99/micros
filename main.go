package main

import "net/http"
import "fmt"
import "encoding/json"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	response := User{Name: "Rami", Email: "r@r.com"}
	encoder.Encode(response)
}

type wrapper struct {
	handler http.Handler
}

func Wrap(h http.Handler) http.Handler {
	return &wrapper{handler: h}
}

func (h *wrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Before")
	h.handler.ServeHTTP(w, r)
	fmt.Println("after")
}

func main() {
	handler := http.HandlerFunc(homeHandler)
	http.Handle("/", Wrap(handler))
	http.ListenAndServe(":8080", nil)
}
