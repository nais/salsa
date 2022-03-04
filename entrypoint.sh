#!/bin/sh -l

WORKING_DIRECTORY="$1"
REPOSITORY="$2"
REPO_NAME="${REPOSITORY##*/}"
REPO_DIR=$(echo "$WORKING_DIRECTORY" | sed 's|\(.*\)/.*|\1|')
#REPO_DIR="${WORKING_DIRECTORY%/*}"
GITHUB=$(echo "${3}" | base64 -w 0)
RUNNER=$(echo "${4}" | base64 -w 0)
IMAGE="$5"
ENVS=$(echo "${6}" | base64 -w 0)

# Run salsa commands
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
