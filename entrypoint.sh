#!/bin/sh -l

echo "------- Running container with repo $SALSA_SCAN_REPO -------"

# Run salsa commands
salsa scan \
  --repo "$SALSA_SCAN_REPO" \
  --github_context "$SALSA_SCAN_GITHUB_CONTEXT" \
  --runner_context "$SALSA_SCAN_RUNNER_CONTEXT" \
  --env_context "" &&

echo "------- Salsa provenance $SALSA_SCAN_REPO generated -------" &&

salsa attest \
  --repo "$SALSA_SCAN_REPO" \
  --config salsa-sample.yaml "$SALSA_REPO_IMAGE" &&

echo "------- Salsa provenance $SALSA_SCAN_REPO uploaded -------" &&

salsa attest \
  --verify \
  --repo "$SALSA_SCAN_REPO" \
  --config salsa-sample.yaml "$SALSA_REPO_IMAGE" &&

echo "------- Attest $SALSA_SCAN_REPO fetch and generated -------"