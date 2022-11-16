#!/bin/sh -l

setup() {
  echo "---------- Preparing pico-de-galo SLSA ----------"

  REPO_NAME="${INPUT_REPO_NAME##*/}"
  if [ -z "$REPO_NAME" ]; then
    echo "REPO_NAME is empty"
    exit 1
  fi

  if [ -n "$INPUT_DOCKER_USER" ]; then
    export GITHUB_ACTOR=$INPUT_DOCKER_USER
  fi

  if [ -z "$GITHUB_ACTOR" ]; then
    echo "GITHUB_ACTOR is not set. Please set it to your GitHub username."
    exit 1
  fi

  if [ -n "$INPUT_IMAGE" ]; then
    export IMAGE=$INPUT_IMAGE
  fi

  if [ -z "$IMAGE" ]; then
    echo "IMAGE is not set"
    exit 1
  fi

  if [ -z "$INPUT_GITHUB_CONTEXT" ] || [ -z "$INPUT_RUNNER_CONTEXT" ]; then
    echo "GITHUB_CONTEXT and RUNNER_CONTEXT are required"
    exit 1
  fi

  GITHUB=$(echo "${INPUT_GITHUB_CONTEXT}" | base64 -w 0) &&
    RUNNER=$(echo "${INPUT_RUNNER_CONTEXT}" | base64 -w 0) &&
    ENVS=$(jq -n env | base64 -w 0)

  DOCKER_REGISTRY="${IMAGE%%/*}"
  if [ -z "$DOCKER_REGISTRY" ]; then
    echo "DOCKER_REGISTRY is not set"
    exit 1
  fi

  exportCosignEnvironment
  exportGithubToken

  export JAVA_HOME=/opt/java/openjdk
}

exportGithubToken() {
  if [ -n "$INPUT_GITHUB_TOKEN" ]; then
    if [ -n "$INPUT_TOKEN_KEY_PATTERN" ]; then
      export "$INPUT_TOKEN_KEY_PATTERN"="$INPUT_GITHUB_TOKEN"
    else
      export GITHUB_TOKEN="$INPUT_GITHUB_TOKEN"
    fi
  fi
}

exportCosignEnvironment() {
  if [ -n "$COSIGN_EXPERIMENTAL" ]; then
    export COSIGN_EXPERIMENTAL
  fi

  if [ -n "$COSIGN_REPOSITORY" ]; then
    export COSIGN_REPOSITORY
  fi
}

loginDocker() {
  echo "---------- Logging in to Docker registry: $DOCKER_REGISTRY ----------"
  echo "$INPUT_DOCKER_PWD" | docker login "$DOCKER_REGISTRY" -u "$GITHUB_ACTOR" --password-stdin
}

logoutDocker() {
  echo "---------- Logging out from Docker registry: $DOCKER_REGISTRY ----------"
  docker logout "$DOCKER_REGISTRY"
}

scan() {
  salsa scan \
    --repo "$REPO_NAME" \
    --build-context "$GITHUB" \
    --runner-context "$RUNNER" \
    --env-context "$ENVS" \
    --subDir "$INPUT_REPO_SUB_DIR" \
    --with-deps="$INPUT_DEPENDENCIES" \
    --mvn-opts "$INPUT_MVN_OPTS" \
    --build-started-on "$INPUT_BUILD_STARTED_ON" \
    --remote-run
}

attest() {
  echo "create and upload attestation" &&
    salsa attest \
      --repo "$REPO_NAME" \
      --subDir "$INPUT_REPO_SUB_DIR" \
      --remote-run \
      --identity-token "$INPUT_IDENTITY_TOKEN" \
      --key "$INPUT_KEY" \
      "$IMAGE"
}

attestVerify() {
  echo "verify attestation" &&
    salsa attest \
      --verify \
      --repo "$REPO_NAME" \
      --subDir "$INPUT_REPO_SUB_DIR" \
      --remote-run \
      --key "$INPUT_KEY" \
      "$IMAGE"
}

runSalsa() {
  echo "---------- Running Salsa for repository: $REPO_NAME ----------" &&
    scan && attest && attestVerify
}

cleanUpGoogle() {
  echo "---------- Clean up Google Cloud stuff ----------"
  if
    [ -n "$GOOGLE_APPLICATION_CREDENTIALS" ] ||
      [ -n "$CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE" ] ||
      [ -n "$GOOGLE_GHA_CREDS_PATH" ]
  then
    rm -rvf "$GOOGLE_APPLICATION_CREDENTIALS" "$CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE" "$GOOGLE_GHA_CREDS_PATH"
  fi
}

setOutput() {
  echo "---------- Setting output ----------"
  {
    echo "provenance_file_path=$REPO_NAME.provenance"
    echo "raw_file_path=$REPO_NAME.raw.txt"
  } >>"$GITHUB_OUTPUT"
}

setup && loginDocker && runSalsa && logoutDocker && setOutput
cleanUpGoogle
