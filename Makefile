salsa:
	go build -o bin/salsa cmd/main.go
test: fmt vet
	go test ./... -coverprofile cover.out -short
fmt:
	go fmt ./...
vet:
	go vet ./...

coverage.out:
	go test -race -v -count=1 -covermode=atomic -coverprofile=coverage.out ./... || true

cover-html: coverage.out
	go tool cover -html=$<