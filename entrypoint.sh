#!/bin/sh -l

#IMAGE=ttl.sh/nais-salsa-"$INPUT_REPO_NAME":1h

# Run salsa commands
salsa scan \
  --repoDir "$INPUT_REPO_DIR" \
  --repo "$INPUT_REPO_NAME" \
  --github_context "$INPUT_GITHUB_CONTEXT" \
  --runner_context "$INPUT_RUNNER_CONTEXT" \
  --env_context "$INPUT_ENV_CONTEXT"

# For private repo
# echo "$INPUT_PASSWORD" | docker login --username foo --password-stdin

#docker pull "$INPUT_IMAGE"
#docker tag "$INPUT_IMAGE" "$IMAGE"
#docker push ttl.sh/nais-salsa-"$INPUT_REPO_NAME":1h

#salsa attest \
#  --repoDir "$INPUT_REPO_DIR" \
#  --repo "$INPUT_REPO_NAME" \
#  --config salsa-sample.yaml "$IMAGE"
