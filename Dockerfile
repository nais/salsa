FROM golang:1.17 AS builder

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /src
COPY . /src
RUN make salsa

FROM alpine:3.14

ENV COSIGN_VERSION=v1.6.0

COPY --from=builder /src/bin/salsa /usr/local/bin/
RUN chmod +x /usr/local/bin/salsa

RUN apk add --no-cache ca-certificates git curl
RUN curl -L -f https://github.com/sigstore/cosign/releases/download/$COSIGN_VERSION/cosign-linux-amd64 > /usr/local/bin/cosign && chmod +x /usr/local/bin/cosign

RUN apk add --no-cache jq httpie

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
