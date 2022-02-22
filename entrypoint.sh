#!/bin/sh -l

echo "Running container with repo $SALSA_SCAN_REPO"

# Run salsa commands
salsa scan \
  --github_context "$SALSA_SCAN_GITHUB_CONTEXT" \
  --runner_context "$SALSA_SCAN_RUNNER_CONTEXT" \
  --env_context ""

#salsa attest \
#  --repoDir "$INPUT_REPO_DIR" \
#  --repo "$INPUT_REPO_NAME" \
#  --config salsa-sample.yaml "$IMAGE"
