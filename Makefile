ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

run_service:
	go run .

install_deps:
	rm -f go.mod go.sum
	go mod init github.com/newline-sandbox/go-chi-docgen-example
	go mod tidy

gen_docs_md:
	go run . -docs=markdown

gen_docs_json:
	go run . -docs=json

gen_docs_raml:
	go run . -docs=raml

gen_docs_help:
	go run . -h