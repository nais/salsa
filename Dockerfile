FROM golang:1.22.1 AS builder

ENV GOOS=linux
ENV CGO_ENABLED=0

# Make Salsa
WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY . /src
RUN make salsa

FROM maven:3.9.1-eclipse-temurin-17-alpine

RUN apk add --no-cache ca-certificates docker jq httpie

# Define a constant with the version of gradle you want to install
ARG GRADLE_VERSION=7.5.1
# Define the URL where gradle can be downloaded from
ARG GRADLE_BASE_URL=https://services.gradle.org/distributions
# Define the SHA key to validate the gradle download
# obtained from here https://gradle.org/release-checksums/
ARG GRADLE_SHA=f6b8596b10cce501591e92f229816aa4046424f3b24d771751b06779d58c8ec4

# Create the directories, download gradle, validate the download, install it, remove downloaded file and set links
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

# Define environmental variables required by gradle
ENV GRADLE_VERSION ${GRADLE_VERSION}
ENV GRADLE_HOME /usr/bin/gradle
ENV GRADLE_USER_HOME /cache
ENV PATH $PATH:$GRADLE_HOME/bin

# Import Salsa
COPY --from=builder /src/bin/salsa /usr/local/bin/
COPY --from=builder /src/salsa-sample.yaml .salsa.yaml
RUN chmod +x /usr/local/bin/salsa

# Verify and install Cosign
ARG COSIGN_VERSION=v2.0.2
ENV COSIGN_BINARY=cosign-linux-amd64
ENV COSIGN_CHECKSUM=cosign_checksums.txt
ENV COSIGN_PUBLIC_KEY=release-cosign.pub
ENV COSIGN_SIG=cosign-linux-amd64.sig

# Cosign urls
ARG COSIGN_BASE_URL=https://github.com/sigstore/cosign/releases/download/$COSIGN_VERSION
ARG COSIGN_CHECKSUM_URL=${COSIGN_BASE_URL}/${COSIGN_CHECKSUM}
ARG COSIGN_BINARY_URL=${COSIGN_BASE_URL}/${COSIGN_BINARY}
ARG COSIGN_PUBLIC_KEY_URL=${COSIGN_BASE_URL}/${COSIGN_PUBLIC_KEY}
ARG COSIGN_SIG_URL=${COSIGN_BASE_URL}/${COSIGN_SIG}

RUN echo "Download cosign checksum" \
  && curl -fsSL -o /tmp/${COSIGN_CHECKSUM} ${COSIGN_CHECKSUM_URL} \
  \
  && echo "Extract current checksum from: ${COSIGN_CHECKSUM}" \
  && export COSIGN_SHA256=$(grep -w ${COSIGN_BINARY} tmp/${COSIGN_CHECKSUM} | cut -d ' ' -f1) \
  \
  && echo "Download cosign ${COSIGN_BINARY} version: ${COSIGN_VERSION}" \
  && curl -fsSL -o /tmp/${COSIGN_BINARY} ${COSIGN_BINARY_URL} \
  \
  && echo "Verify checksum ${COSIGN_SHA256} with ${COSIGN_BINARY}" \
  && sha256sum /tmp/${COSIGN_BINARY} \
  && echo "${COSIGN_SHA256}  /tmp/${COSIGN_BINARY}" | sha256sum -c - \
  \
  && echo "Move cosign to folder and make cosign executable" \
  && chmod +x /tmp/${COSIGN_BINARY} \
  && mkdir "tmp2" \
  && cp /tmp/${COSIGN_BINARY} tmp2/${COSIGN_BINARY} \
  && mv /tmp/${COSIGN_BINARY} /usr/local/bin/cosign \
  && chmod +x /usr/local/bin/cosign \
  \
  && echo "Verify ${COSIGN_BINARY} with public key and signature" \
  && curl -fsSL -o /tmp/${COSIGN_PUBLIC_KEY} ${COSIGN_PUBLIC_KEY_URL} \
  && curl -fsSL -o /tmp/${COSIGN_SIG} ${COSIGN_SIG_URL} \
  && cosign \
  && cosign verify-blob --key /tmp/${COSIGN_PUBLIC_KEY} --signature /tmp/${COSIGN_SIG} /tmp2/${COSIGN_BINARY}

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
