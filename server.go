package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/flosch/pongo2"  // use django like template
	"github.com/gorilla/csrf" // add csrf protection
	"github.com/julienschmidt/httprouter" // use http router
)

var UUID string = "8d8ee6fc9b8a4861a277ce1184cd51b"

func CsrfDecorator(f func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		f(w, r, ps)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL, r.Proto)
	io.WriteString(w, "/")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(r.URL, r.Proto, ps.ByName("name"))
	io.WriteString(w, "Hello World! "+ps.ByName("name")+"!")
}

var home = pongo2.Must(pongo2.FromFile("home.html"))

func HomePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := home.ExecuteWriter(pongo2.Context{}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/app", HomePage)
	router.ServeFiles("/usercreatedata/*filepath", http.Dir("./public"))

	port := ":3010"

	fmt.Println("create https server")

	// log.Fatal(server.ListenAndServeTLS(port, cert, key, router))
	// fmt.Println(csrf.Protect([]byte(UUID)))
	log.Fatal(http.ListenAndServe(port, csrf.Protect([]byte(UUID))(router)))

	fmt.Println("Server listening on http://localhost" + port)
}
