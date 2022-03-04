#!/bin/sh -l

# DIRECTORY="$1"
REPOSITORY="$2"
REPO_NAME="${REPOSITORY##*/}"
# REPO_DIR="${DIRECTORY%/*}"
GITHUB=$(echo "${3}" | base64 -w 0)
RUNNER=$(echo "${4}" | base64 -w 0)
IMAGE="$5"
ENVS=$(echo "${6}" | base64 -w 0)

echo "---------- Preparing pico-de-galo Slsa for repository: $REPO_NAME ----------"
salsa scan \
  --repo "$REPO_NAME" \
  --github-context "$GITHUB" \
  --runner-context "$RUNNER" \
  --env-context "$ENVS" \
  --remote-run &&
  salsa attest \
    --repo "$REPO_NAME" \
    --config salsa-sample.yaml "$IMAGE" \
    --remote-run &&
  salsa attest \
    --verify \
    --repo "$REPO_NAME" \
    --config salsa-sample.yaml "$IMAGE" \
    --remote-run
