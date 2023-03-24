#!/bin/sh -l

setup() {
  echo "---------- Preparing pico de gallo SLSA ----------"

  REPO_NAME="${INPUT_REPO_NAME##*/}"
  if [ -z "$REPO_NAME" ]; then
    echo "REPO_NAME is empty"
    exit 1
  fi

  if [ -z "$INPUT_REGISTRY" ]; then
    echo "INPUT_REGISTRY is empty"
    exit 1
  fi

  if [ -n "$INPUT_DOCKER_USER" ]; then
    export GITHUB_ACTOR="$INPUT_DOCKER_USER"
  fi

  if [ -z "$GITHUB_ACTOR" ]; then
    echo "GITHUB_ACTOR is not set. Please set it to your GitHub username."
    exit 1
  fi

  if [ -n "$INPUT_IMAGE" ]; then
    export IMAGE="$INPUT_IMAGE"
  fi

  if [ -z "$INPUT_IMAGE_DIGEST" ] || [ -z "$IMAGE" ]; then
    echo "IMAGE and IMAGE_DIGEST is not set. Please set it."
    exit 1
  fi

  export IMAGE="$IMAGE@$INPUT_IMAGE_DIGEST"

  if [ -z "$INPUT_GITHUB_CONTEXT" ] || [ -z "$INPUT_RUNNER_CONTEXT" ]; then
    echo "GITHUB_CONTEXT and RUNNER_CONTEXT are required"
    exit 1
  fi

  if [ "$INPUT_VERIFY_ATTESTATION" = "false" ] && [ -z "$INPUT_KEY" ]; then
    echo "When running keyless salsa you must verify the attestation. Please set the verify_attestation flag to 'true'.
    
    (This is also the default value, and may instead be omitted)."
    exit 1
  fi

  GITHUB=$(echo "${INPUT_GITHUB_CONTEXT}" | base64 -w 0) &&
    RUNNER=$(echo "${INPUT_RUNNER_CONTEXT}" | base64 -w 0) &&
    ENVS=$(jq -n env | base64 -w 0)

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
  else
    export GITHUB_TOKEN
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
  echo "---------- Logging in to Docker registry: $INPUT_REGISTRY ----------"
  if [ -n "$INPUT_REGISTRY_ACCESS_TOKEN" ]; then
    echo "$INPUT_REGISTRY_ACCESS_TOKEN" | docker login "$INPUT_REGISTRY" -u "$GITHUB_ACTOR" --password-stdin
  else
    echo "$GITHUB_TOKEN" | docker login "$INPUT_REGISTRY" -u "$GITHUB_ACTOR" --password-stdin
  fi
}

logoutDocker() {
  echo "---------- Logging out from Docker registry: $INPUT_REGISTRY ----------"
  docker logout "$INPUT_REGISTRY"
}

scan() {
  echo "---------- Running Salsa scan for deps ----------" &&
    salsa scan \
      --repo "$REPO_NAME" \
      --build-context "$GITHUB" \
      --runner-context "$RUNNER" \
      --env-context "$ENVS" \
      --subDir "$INPUT_REPO_SUB_DIR" \
      --mvn-opts "$INPUT_MVN_OPTS" \
      --build-started-on "$INPUT_BUILD_STARTED_ON" \
      --remote-run
}

attest() {
  echo "---------- Creating and Uploading Salsa attestation ----------" &&
    salsa attest \
      --repo "$REPO_NAME" \
      --subDir "$INPUT_REPO_SUB_DIR" \
      --remote-run \
      --identity-token "$INPUT_IDENTITY_TOKEN" \
      --key "$INPUT_KEY" \
      "$IMAGE"
}

attestVerify() {
  echo "---------- Verifying Salsa attestation ----------" &&
    salsa attest \
      --verify \
      --repo "$REPO_NAME" \
      --subDir "$INPUT_REPO_SUB_DIR" \
      --remote-run \
      --key "$INPUT_KEY" \
      "$IMAGE"
}

runSalsa() {
  echo "---------- Running Salsa for repository: $REPO_NAME ----------"
  if [ "$INPUT_VERIFY_ATTESTATION" = "true" ]; then
    scan && attest
  elif [ "$INPUT_VERIFY_ATTESTATION" = "false" ]; then
    scan && attest && attestVerify
  fi

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
