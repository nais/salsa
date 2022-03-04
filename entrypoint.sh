#!/bin/sh -l

DIRECTORY="$1"
echo "$DIRECTORY"
REPO_DIR="${DIRECTORY%/*}"
echo "$REPO_DIR"
REPOSITORY="$2"
REPO_NAME="${REPOSITORY##*/}"
GITHUB=$(echo "${3}" | base64 -w 0)
RUNNER=$(echo "${4}" | base64 -w 0)
IMAGE="$5"
ENVS=$(echo "${6}" | base64 -w 0)

echo "---------- Preparing pico-de-galo Slsa for repository: $REPO_NAME ----------"
salsa scan \
  --repoDir "$REPO_DIR" \
  --repo "$REPO_NAME" \
  --github_context "$GITHUB" \
  --runner_context "$RUNNER" \
  --env_context "$ENVS" &&
  salsa attest \
    --repoDir "$REPO_DIR" \
    --repo "$REPO_NAME" \
    --config salsa-sample.yaml "$IMAGE" &&
  salsa attest \
    --verify \
    --repoDir "$REPO_DIR" \
    --repo "$REPO_NAME" \
    --config salsa-sample.yaml "$IMAGE"
