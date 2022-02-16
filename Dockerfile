# FROM gcr.io/projectsigstore/cosign:v1.5.1 AS COSIGN

FROM golang:1.17 AS builder

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /src
COPY . /src
RUN make salsa

FROM alpine:3.14

COPY --from=builder /src/bin/salsa /usr/local/bin/

RUN apk add --no-cache ca-certificates git curl
RUN curl -L -f https://github.com/sigstore/cosign/releases/download/v1.5.1/cosign-linux-amd64 > /usr/local/bin/cosign && chmod +x /usr/local/bin/cosign

#RUN export PATH=$PATH:/app

RUN chmod +x /usr/local/bin/salsa

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/entrypoint.sh"]