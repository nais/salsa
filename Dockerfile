# FROM gcr.io/projectsigstore/cosign:v1.5.1 AS COSIGN

FROM golang:1.17 AS builder

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /app
COPY . /app
RUN make salsa

FROM alpine:3.14

ARG repo
ENV SALSA_SCAN_REPO=$repo
ARG dir
ENV SALSA_SCAN_REPO_DIR=$dir
ARG image
ENV SALSA_REPO_IMAGE=$image
ARG github
ENV SALSA_SCAN_GITHUB_CONTEXT=$github
ARG runner
ENV SALSA_SCAN_RUNNER_CONTEXT=$runner
ARG gcp_cred_path
ENV GOOGLE_APPLICATION_CREDENTIALS=$gcp_cred_path

COPY --from=builder /app/bin/salsa /usr/local/bin/

RUN apk add --no-cache ca-certificates git curl
RUN curl -L -f https://github.com/sigstore/cosign/releases/download/v1.5.1/cosign-linux-amd64 > /usr/local/bin/cosign && chmod +x /usr/local/bin/cosign

RUN chmod +x /usr/local/bin/salsa

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /salsa
COPY . /salsa

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/entrypoint.sh"]