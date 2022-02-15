FROM golang:1.17 AS builder
ENV APP_ROOT=/opt/app-root
ENV GOPATH=$APP_ROOT

WORKDIR $APP_ROOT/src/
ADD ../go.mod go.sum $APP_ROOT/src/
RUN go mod download

# Install cosign
RUN git clone https://github.com/sigstore/cosign && cd cosign && go install ./cmd/cosign && $(go env GOPATH)/bin/cosign

# Add source code
ADD ../.. $APP_ROOT/src/

RUN make salsa

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/entrypoint.sh"]