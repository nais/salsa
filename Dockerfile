FROM alpine:3.14

ENV COSIGN_VERSION=v1.6.0
ENV SALSA_VERSION=0.1.20220324094800

RUN apk add --no-cache ca-certificates git curl
RUN curl -L -f https://github.com/sigstore/cosign/releases/download/$COSIGN_VERSION/cosign-linux-amd64 > /usr/local/bin/cosign && chmod +x /usr/local/bin/cosign

RUN curl -o salsa.tar.gz -L -f https://github.com/nais/salsa/releases/download/$SALSA_VERSION/nais-salsa_${SALSA_VERSION}_linux_amd64.tar.gz
RUN tar -xzf salsa.tar.gz  && cp salsa /usr/local/bin/salsa && chmod +x /usr/local/bin/salsa

RUN apk add --no-cache jq httpie

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

