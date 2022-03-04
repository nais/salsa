#!/bin/sh -l

DIRECTORY="$1"
REPOSITORY="$2"
REPO_NAME="${REPOSITORY##*/}"
REPO_DIR="${DIRECTORY%/*}"
# GITHUB_WORKSPACE: "/home/runner/work/salsa/salsa"
if test "${REPO_DIR#*$REPO_NAME}" != "$REPO_DIR"; then
  REPO_DIR="${REPO_DIR%/*}"
fi
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
