#!/bin/sh -l

echo "---------- Preparing pico-de-galo Slsa container with repository: $SALSA_SCAN_REPO ----------"

# Run salsa commands
# --repoDir "$SALSA_SCAN_REPO_DIR" \
salsa scan \
  --repo "$SALSA_SCAN_REPO" \
  --github_context "$SALSA_SCAN_GITHUB_CONTEXT" \
  --runner_context "$SALSA_SCAN_RUNNER_CONTEXT" \
  --env_context "" &&
  echo "---------- Slsa roja provenance for repository: $SALSA_SCAN_REPO generated ----------" &&
  salsa attest \
    --repo "$SALSA_SCAN_REPO" \
    --config salsa-sample.yaml "$SALSA_REPO_IMAGE" &&
  echo "---------- Slsa verde provenance for: $SALSA_SCAN_REPO signed and uploaded to remote host ----------" &&
  salsa attest \
    --verify \
    --repo "$SALSA_SCAN_REPO" \
    --config salsa-sample.yaml "$SALSA_REPO_IMAGE" &&
  echo "---------- Slsa aguacate attest verified for: $SALSA_SCAN_REPO ----------"
