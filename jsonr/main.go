package main

import "net/http"
import "encoding/json"
import "log"
import "fmt"

//Port number
const PORT = 8080

type helloRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type helloResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

func helloRequestHandler(w http.ResponseWriter, r *http.Request) {
	var request helloRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	response := helloResponse{Message: "Hello " + request.Name, Email: request.Email}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func helloResponseHandler(w http.ResponseWriter, r *http.Request) {
	response := helloResponse{Message: "Hello world", Email: "rr@r.com"}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func main() {
	staticHandler := http.FileServer(http.Dir("../images"))
	http.HandleFunc("/", helloResponseHandler)
	http.HandleFunc("/adduser", helloRequestHandler)
	http.HandleFunc("/notfound", http.NotFound)
	http.Handle("/cat/", http.StripPrefix("/cat/", staticHandler))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil))
}
