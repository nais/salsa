# Status

Proof of Concept for a SLSA github action / cli.

## Relevant concepts to test

* upload attestation somewhere
* how to make attestations searchable
* * upload attestation somewhere
* explore tools like cosign, Fulcio and Reko from sigstore to see where they can fit in
    * https://github.com/sigstore/fulcio:
        * Fulcio is a work in progress. There's working code and a running instance and a plan, but you should not
          attempt to try to actually use it for anything
* how to make attestations searchable
* Handle the ability to resolve packages that's private

## Concepts tested so far

Created simple CLI to test concepts:

* sign attestation using DSSE (leverage some of sigstore functionality)
* create a SBOM / in-toto attestation
* clone github project
* Should contain a Predicate for SLSA Provenance
* list all dependencies in a gradle project
* get all dependencies (including transitive) for a given repo and language
    * One language at a time
* create attestation with materials based on dependencies
* sign attestation with DSSE
* sign docker image and put into attestation, using cosign
* digest over dependencies etc in attestation
* include build steps from workflow
* create a pipeline where a "provenance" action can be used
* how to get/add the digest for dependency artifacts for all build tools

# Relevant links

* https://github.com/in-toto/attestation/blob/main/spec/README.md
* https://github.com/slsa-framework/slsa/blob/main/controls/attestations.md
* https://github.com/secure-systems-lab/dsse
* https://slsa.dev/provenance/v0.2
* Mostly cosign, reko and fulcio: https://docs.sigstore.dev/