package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/remisb/go-nethttp/middleware"
)

const port = 9000

func productHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Details for Product #" + id))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page content"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page content"))
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /product/{id}", productHandler)
	mux.HandleFunc("GET /home", homeHandler)
	mux.HandleFunc("GET /", indexHandler)
}

func main() {
	mux := http.NewServeMux()
	setupRoutes(mux)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	stack := middleware.NewStack(
		middleware.Logging,
	)
	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: stack(mux),
	}

	log.Printf("HTTP Server started at port: %d\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// HTTP middleware setting a value on the request context
func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), "user", "123")

		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
