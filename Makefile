run_service:
	go run .

install_deps:
	rm -f go.mod go.sum
	go mod init github.com/newline-sandbox/go-chi-docgen-example
	go mod tidy