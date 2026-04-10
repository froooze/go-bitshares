.PHONY: fmt test tidy

fmt:
	gofmt -w $$(find . -name '*.go')

test:
	go test ./...

tidy:
	go mod tidy
