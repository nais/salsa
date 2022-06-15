#!/bin/sh -l

REPO_NAME="${INPUT_REPO_NAME##*/}"
# REPO_DIR="${INPUT_REPO_DIR%/*}"

if [ -n "$INPUT_DOCKER_USER" ]; then
  export GITHUB_ACTOR=$INPUT_DOCKER_USER
fi

if [ -n "$INPUT_IMAGE" ]; then
  export IMAGE=$INPUT_IMAGE
fi

GITHUB=$(echo "${INPUT_GITHUB_CONTEXT}" | base64 -w 0)
RUNNER=$(echo "${INPUT_RUNNER_CONTEXT}" | base64 -w 0)
ENVS=$(jq -n env | base64 -w 0)
DOCKER_REGISTRY="${IMAGE%%/*}"

export JAVA_HOME=/opt/java/openjdk

echo "---------- Preparing pico-de-galo slsa for repository: $REPO_NAME ----------"
echo "$INPUT_DOCKER_PWD" | docker login "$DOCKER_REGISTRY" -u "$GITHUB_ACTOR" --password-stdin

salsa scan \
  --repo "$REPO_NAME" \
  --build-context "$GITHUB" \
  --runner-context "$RUNNER" \
  --env-context "$ENVS" \
  --subDir "$INPUT_REPO_SUB_DIR" \
  --with-deps="$INPUT_DEPENDENCIES" \
  --remote-run &&
  salsa attest \
    --repo "$REPO_NAME" \
    --subDir "$INPUT_REPO_SUB_DIR" \
    --remote-run \
    --key "$INPUT_KEY" \
    "$IMAGE" &&
  salsa attest \
    --verify \
    --repo "$REPO_NAME" \
    --subDir "$INPUT_REPO_SUB_DIR" \
    --remote-run \
    --key "$INPUT_KEY" \
    "$IMAGE" && echo "clean up..." && docker logout "$DOCKER_REGISTRY" &&
  if [ -n "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    rm "$GOOGLE_APPLICATION_CREDENTIALS"
  fi
