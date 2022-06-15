RATCHET_VERSION=0.2.2
WORKFLOWS_PATH=.github/workflows

.PHONY: ratchet
ratchet:
	go install github.com/sethvargo/ratchet@v${RATCHET_VERSION} ; \
	chmod -R +x ${WORKFLOWS_PATH}/*

.PHONY: pin
pin: ratchet ## create a pinned workflow IF NOT already pinned
	ratchet pin "${WORKFLOWS_PATH}/${workflow}"

.PHONY: update
update: ratchet ## update pinned workflows
	ratchet update "${WORKFLOWS_PATH}/${workflow}"

.PHONY: check
check: ratchet ##  check pinned workflows
	ratchet check "${WORKFLOWS_PATH}/${workflow}"

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
