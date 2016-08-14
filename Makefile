.PHONY: default deps dev-deps doc webdoc format test coverage markdown

default: format test

deps:
	go get -u -v ./...

dev-deps:
	go get -u github.com/campoy/embedmd

doc:
	godoc `pwd`

webdoc:
	godoc -http=:44444

format:
	go fmt

test:
	go test -v -race

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html="coverage.out"

markdown:
	embedmd -w README.md
