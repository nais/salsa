#!/bin/sh -l

REPOSITORY="$1"
REPO_NAME="${REPOSITORY##*/}"
GITHUB=$(echo "${2}" | base64 -w 0)
RUNNER=$(echo "${3}" | base64 -w 0)
IMAGE="$4"
ENVS=$(echo "${5}" | base64 -w 0)

# Run salsa commands
salsa scan \
  --repo "$REPO_NAME" \
  --github_context "$GITHUB" \
  --runner_context "$RUNNER" \
  --env_context "$ENVS"

salsa attest \
  --repo "$REPO_NAME" \
  --config salsa-sample.yaml "$IMAGE"

salsa attest \
  --verify \
  --repo "$REPO_NAME" \
  --config salsa-sample.yaml "$IMAGE"
