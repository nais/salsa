#!/bin/sh -l

#IMAGE=ttl.sh/nais-salsa-"$INPUT_REPO_NAME":1h

sh -c "echo $*"
# Run salsa commands
salsa scan \
  --repoDir "$1" \
  --repo "$2" \
  --github_context "$3" \
  --runner_context "$4" \
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
