package main

import (
  "log"
  "net/http"
  "os"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/go-chi/render"
)

func main() {
	port := "8080"

  if fromEnv := os.Getenv("PORT"); fromEnv != "" {
    port = fromEnv
  }

  log.Printf("Starting up on http://localhost:%s", port)

  r := chi.NewRouter()

  r.Use(middleware.RequestID)
  r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"))
  })

	r.Mount("/posts", postsResource{}.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}