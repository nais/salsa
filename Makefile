
salsa:
	go build -o bin/salsa cmd/main.go
test: fmt vet
	go test ./... -coverprofile cover.out -short
fmt:
	go fmt ./...
vet:
	go vet ./...

cover-out:
	go test -race -v -count=1 -covermode=atomic -coverprofile=coverage.out ./... || true

cover-html: cover-out
	go tool cover -html=$<
