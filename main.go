package main

import (
  "errors"
  "flag"
  "net/http"
  "log"
  "os"
  "strings"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/go-chi/docgen"
  "github.com/go-chi/docgen/raml"
  "github.com/go-chi/render"
  yaml "gopkg.in/yaml.v3"
)

var docs = flag.String("docs", "", "Generates routing documentation for RESTful API - markdown, json, or raml.") // To see the description, run the `make gen_docs_help` command.

func main() {
  flag.Parse()

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

  docgen.PrintRoutes(r)

  if *docs == "markdown" {
    if err := os.Remove("routes.md"); err != nil && !errors.Is(err, os.ErrNotExist) {
      log.Fatal(err)
    }

    f, err := os.Create("routes.md")

    if err != nil {
      log.Fatal(err)
    }

    defer f.Close()

    text := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
      ProjectPath: "github.com/newline-sandbox/go-chi-docgen-example",
      URLMap: map[string]string{
        "github.com/newline-sandbox/go-chi-docgen-example/vendor/github.com/go-chi/chi/v5/": "https://github.com/go-chi/chi/blob/master/",
      },
      ForceRelativeLinks: true,
      Intro: "Welcome to the documentation for the RESTful API.",
    })

    if _, err = f.Write([]byte(text)); err != nil {
      log.Fatal(err)
    }

    return
  } else if *docs == "json" {
    if err := os.Remove("routes.json"); err != nil && !errors.Is(err, os.ErrNotExist) {
      log.Fatal(err)
    }

    f, err := os.Create("routes.json")

    if err != nil {
      log.Fatal(err)
    }

    defer f.Close()

    json := docgen.JSONRoutesDoc(r)

    if _, err = f.Write([]byte(json)); err != nil {
      log.Fatal(err)
    }

    return
  } else if *docs == "raml" {
    if err := os.Remove("routes.raml"); err != nil && !errors.Is(err, os.ErrNotExist) {
      log.Fatal(err)
    }

    f, err := os.Create("routes.raml")

    if err != nil {
      log.Fatal(err)
    }

    defer f.Close()

    ramlDocs := &raml.RAML{
      Title:     "RAML Representation of RESTful API",
      BaseUri:   "http://api.go-chi-docgen-example.com/v1",
      Version:   "v1.0",
      MediaType: "application/json",
    }

    if err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
      handlerInfo := docgen.GetFuncInfo(handler)
      resource := &raml.Resource{
        DisplayName: strings.ToUpper(method) + " " + route,
        Description: "Handler Function: " + handlerInfo.Func + "\nComment: " + handlerInfo.Comment,
      }

      return ramlDocs.Add(method, route, resource)
    }); err != nil {
       log.Fatalf("error: %v", err)
    }

    raml, err := yaml.Marshal(ramlDocs)

    if err != nil {
      log.Fatal(err)
    }

    if _, err = f.Write(append([]byte("#%RAML 1.0\n---\n"), raml...)); err != nil {
      log.Fatal(err)
    }

		return
  }

	log.Fatal(http.ListenAndServe(":"+port, r))
}