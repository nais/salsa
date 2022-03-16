#!/bin/sh -l

# DIRECTORY="$1"
REPO_NAME="${INPUT_REPO_NAME##*/}"
# REPO_DIR="${INPUT_REPO_DIR%/*}"
GITHUB=$(echo "${INPUT_GITHUB_CONTEXT}" | base64 -w 0)
RUNNER=$(echo "${INPUT_RUNNER_CONTEXT}" | base64 -w 0)
ENVS=$(jq -n env | base64 -w 0)

echo "---------- Preparing pico-de-galo Slsa for repository: $REPO_NAME ----------"
salsa scan \
  --repo "$REPO_NAME" \
  --build-context "$GITHUB" \
  --runner-context "$RUNNER" \
  --env-context "$ENVS" \
  --subDir "$INPUT_REPO_SUB_DIR" \
  --remote-run &&
  salsa attest \
    --repo "$REPO_NAME" \
    --config salsa-sample.yaml "$INPUT_IMAGE" \
    --subDir "$INPUT_REPO_SUB_DIR" \
    --remote-run &&
  salsa attest \
    --verify \
    --repo "$REPO_NAME" \
    --config salsa-sample.yaml "$INPUT_IMAGE" \
    --subDir "$INPUT_REPO_SUB_DIR" \
    --remote-run
