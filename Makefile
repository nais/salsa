salsa:
	go build -o bin/salsa cmd/main.go
test: fmt vet
	go test ./... -coverprofile cover.out -short
fmt:
	go fmt ./...
vet:
	go vet ./...