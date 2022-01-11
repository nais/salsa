# salsa

>in line with the best from abroad

## Status

Proof of Concept for a SLSA github action / cli.

### Relevant concepts to test

* get all dependencies (including transitive) for a given repo and language
    * One language at a time
* create a SBOM / in-toto attestation
    * Should contain a Predicate for SLSA Provenance
* sign attestation using DSSE (leverage some of sigstore functionality)
* upload attestation somewhere
* explore tools like cosign, Fulcio and Reko from sigstore to see where they can fit in
* how to make attestations searchable


### Concepts tested so far

Created simple CLI to test concepts:

* clone github project
* list all dependencies in a gradle project
* create attestation with materials based on dependencies
* sign attestation with DSSE
* sign docker image and put into attestation, using cosign
* digest over dependencies etc in attestation

### Stuff we should explore

* include build steps from workflow
* create a pipeline where a "provenance" action can be used
* upload attestation somewhere
* explore tools like cosign, Fulcio and Reko from sigstore to see where they can fit in
* how to make attestations searchable
* Handle the ability to resolve packages that's private

## Relevant links
* https://github.com/in-toto/attestation/blob/main/spec/README.md
* https://github.com/slsa-framework/slsa/blob/main/controls/attestations.md
* https://github.com/secure-systems-lab/dsse
* https://slsa.dev/provenance/v0.2
* Mostly cosign, reko and fulcio: https://docs.sigstore.dev/
