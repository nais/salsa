#!/bin/sh -l

echo "------- Preparing spicy Slsa in container with repository: $SALSA_SCAN_REPO -------"

ls -l
# Run salsa commands
# --repoDir "$SALSA_SCAN_REPO_DIR" \
salsa scan \
  --repo "$SALSA_SCAN_REPO" \
  --github_context "$SALSA_SCAN_GITHUB_CONTEXT" \
  --runner_context "$SALSA_SCAN_RUNNER_CONTEXT" \
  --env_context "" &&
  echo "------- Slsa provenance for repository $SALSA_SCAN_REPO generated -------" &&
  salsa attest \
    --repo "$SALSA_SCAN_REPO" \
    --config salsa-sample.yaml "$SALSA_REPO_IMAGE" &&
  echo "------- Slsa provenance $SALSA_SCAN_REPO signed and uploaded -------" &&
  salsa attest \
    --verify \
    --repo "$SALSA_SCAN_REPO" \
    --config salsa-sample.yaml "$SALSA_REPO_IMAGE" &&
  echo "------- Attest verified for $SALSA_SCAN_REPO -------"
