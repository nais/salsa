FROM golang:1.17 AS builder

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY . /src
RUN make salsa

FROM maven:3.8.4-eclipse-temurin-17-alpine

# ADD gradle binary
# Downloading and installing Gradle
# 1- Define a constant with the version of gradle you want to install
ARG GRADLE_VERSION=7.3.3

# 2- Define the URL where gradle can be downloaded from
ARG GRADLE_BASE_URL=https://services.gradle.org/distributions

# 3- Define the SHA key to validate the gradle download
#    obtained from here https://gradle.org/release-checksums/
ARG GRADLE_SHA=b586e04868a22fd817c8971330fec37e298f3242eb85c374181b12d637f80302

# 4- Create the directories, download gradle, validate the download, install it, remove downloaded file and set links
RUN mkdir -p /usr/share/gradle /usr/share/gradle/ref \
  && echo "Downloading gradle hash" \
  && curl -fsSL -o /tmp/gradle.zip ${GRADLE_BASE_URL}/gradle-${GRADLE_VERSION}-bin.zip \
  \
  && echo "Checking download hash" \
  && echo "${GRADLE_SHA}  /tmp/gradle.zip" | sha256sum -c - \
  \
  && echo "Unziping gradle" \
  && unzip -d /usr/share/gradle /tmp/gradle.zip \
   \
  && echo "Cleaning and setting links" \
  && rm -f /tmp/gradle.zip \
  && ln -s /usr/share/gradle/gradle-${GRADLE_VERSION} /usr/bin/gradle

# 5- Define environmental variables required by gradle
ENV GRADLE_VERSION 7.3.3
ENV GRADLE_HOME /usr/bin/gradle
ENV GRADLE_USER_HOME /cache

ENV PATH $PATH:$GRADLE_HOME/bin

RUN apk add --no-cache ca-certificates git curl docker jq httpie

# SALSA
COPY --from=builder /src/bin/salsa /usr/local/bin/
COPY --from=builder /src/salsa-sample.yaml .salsa.yaml
RUN chmod +x /usr/local/bin/salsa

# COSIGN
ENV COSIGN_VERSION=v1.8.0
ENV COSIGN_BINARY=cosign-linux-amd64
ENV COSIGN_CHECKSUM=cosign_checksums.txt
ARG COSIGN_BASE_URL=https://github.com/sigstore/cosign/releases/download/$COSIGN_VERSION
ARG COSIGN_CHECKSUM_URL=${COSIGN_BASE_URL}/${COSIGN_CHECKSUM}
ARG COSIGN_BINARY_URL=${COSIGN_BASE_URL}/${COSIGN_BINARY}

RUN echo "Download cosign checksum" \
  && curl -fsSL -o /tmp/${COSIGN_CHECKSUM} ${COSIGN_CHECKSUM_URL} \
  \
  && echo "Extract cosign checksum" \
  && export COSIGN_SHA256=$(grep -w ${COSIGN_BINARY} tmp/${COSIGN_CHECKSUM} | cut -d ' ' -f1) \
  \
  && echo "Download cosign binary version: ${COSIGN_VERSION}" \
  && curl -fsSL -o /tmp/${COSIGN_BINARY} ${COSIGN_BINARY_URL} \
  \
  && echo "Checking downloaded checksum ${COSIGN_SHA256} with ${COSIGN_BINARY}" \
  && sha256sum /tmp/${COSIGN_BINARY} \
  && echo "${COSIGN_SHA256}  /tmp/${COSIGN_BINARY}" | sha256sum -c \
  \
  && echo "Copy cosign" \
  && chmod +x /tmp/${COSIGN_BINARY} \
  && mv /tmp/${COSIGN_BINARY} /usr/local/bin/cosign \
   \
  && echo "Cleaning and setting rights" \
  && rm -f /tmp/${COSIGN_BINARY} \
  && rm -f /tmp/${COSIGN_CHECKSUM} \
  && chmod +x /usr/local/bin/cosign

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
