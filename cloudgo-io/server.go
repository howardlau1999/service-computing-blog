package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"github.com/unrolled/render"
	"github.com/gorilla/schema"
)

type Person struct {
	Id   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

type CreatePersonRequest struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

type User struct {
	Username string
	Password string
}

func main() {
	people := []Person{{1, 18, "Alice"}, {2, 22, "Bob"}}
	router := mux.NewRouter()
	r := render.New()

	router.HandleFunc("/unknown", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = fmt.Fprintf(w, "Method not implemented!\n")
	})

	router.HandleFunc("/api/people", func(w http.ResponseWriter, req *http.Request) {
		data, _ := json.Marshal(people)
		_, _ = w.Write(data)
	}).Methods("GET")

	router.HandleFunc("/api/people", func(w http.ResponseWriter, req *http.Request) {
		var request CreatePersonRequest

		_ = json.NewDecoder(req.Body).Decode(&request)

		var person = Person{
			Id: len(people) + 1,
			Age:  request.Age,
			Name: request.Name,
		}

		people = append(people, person)
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		var user User
		_ = schema.NewDecoder().Decode(&user, req.PostForm)
		_ = r.HTML(w, http.StatusOK, "login", user)
	})

	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrieve working directory")
		} else {
			webRoot = root
			//fmt.Println(root)
		}
	}

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintf(w, "Welcome to the home page!\n")
	})

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/assets/"))))

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(router)

	_ = http.ListenAndServe(":3000", n)
}
