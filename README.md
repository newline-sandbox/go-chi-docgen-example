# Go and `chi` `docgen` Documentation

This project demonstrates how to generate routing documentation for a Go `chi` RESTful API.

## Get Started

Install the dependencies...

```shell
$ make install_deps
```

...then generate the routing documentation:

```shell
$ make gen_docs_md # For pre-formatted Markdown file that documents the routes of the RESTful API (`routes.md`).
$ make gen_docs_json # For JSON representation of the RESTful API (`routes.json`).
$ make gen_docs_raml # For RAML representation of the RESTful API (`routes.raml`).
```

Feel free to clone this project!
