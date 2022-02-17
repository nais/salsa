#!/bin/sh -l

#IMAGE=ttl.sh/nais-salsa-"$INPUT_REPO_NAME":1h
REPO_DIR="$1"
REPO_NAME="$2"
GITHUB=$(echo "${3}" | base64 -w 0)
RUNNER=$(echo "${4}" | base64 -w 0)
#ENVS=$(echo "${5}" | base64 -w 0)

# Run salsa commands
salsa scan \
  --repoDir "$REPO_DIR" \
  --repo "$REPO_NAME" \
  --github_context "$GITHUB" \
  --runner_context "$RUNNER" \
  --env_context ""

# For private repo
# echo "$INPUT_PASSWORD" | docker login --username foo --password-stdin

#docker pull "$INPUT_IMAGE"
#docker tag "$INPUT_IMAGE" "$IMAGE"
#docker push ttl.sh/nais-salsa-"$INPUT_REPO_NAME":1h

#salsa attest \
#  --repoDir "$INPUT_REPO_DIR" \
#  --repo "$INPUT_REPO_NAME" \
#  --config salsa-sample.yaml "$IMAGE"
