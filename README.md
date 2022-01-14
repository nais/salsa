# salsa

> in line with the best from abroad

## Usage

`Prerequisites to run locally`
* Google Setup
  * KMS is enabled in project
    * create keyring
      * create keys: `Elliptic Curve P-256 key SHA256 Digest`
  * Serviceuser in project has roles:
    * Cloud KMS CryptoKey signer/verifier
    * Cloud KMS viewer
* Logged in to Google
* Set: `GOOGLE_APPLICATION_CREDENTIALS` with path to .json file containing serviceuser credentials.

* Install
  * Cosign: https://github.com/sigstore/cosign
  
## Local install

```
make
```

### Commands

clone: `clones the given project into user defined path`
```
bin/salsa clone --repo salsa --url https://github.com/nais/salsa
```

scan: `Scan files and dependencies for a given project`
```
bin/salsa scan --repo salsa
```

attest: `sign and upload in-toto attestation`
```
bin/salsa attest --repo salsa --predicate salsa.provenance --no-upload --key gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION  ttl.sh/salsax:1h
```

Info:
Image can to be pushed to ttl.sh, who offers free, short-lived (ie: hours), anonymous container image hosting if you just want to try these commands out.
**Quick Start** [Cosign](https://github.com/sigstore/cosign#quick-start)  
**ttl.sh** [tt.sh info](https://ttl.sh/)

find: `find artifact from attestations`
```
bin/salsa find go-crypto
```

Instead of setting a bunch of flags, in home directory create a config file with name ".salsa" (without extension).
```yml
attest:
  key: gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION
...
```

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
    * https://github.com/sigstore/fulcio:
        * Fulcio is a work in progress. There's working code and a running instance and a plan, but you should not
          attempt to try to actually use it for anything
* how to make attestations searchable
* how to get/add the digest for dependency artifacts for all build tools
  * currently, implemented only in golang
* Handle the ability to resolve packages that's private

## Relevant links

* https://github.com/in-toto/attestation/blob/main/spec/README.md
* https://github.com/slsa-framework/slsa/blob/main/controls/attestations.md
* https://github.com/secure-systems-lab/dsse
* https://slsa.dev/provenance/v0.2
* Mostly cosign, reko and fulcio: https://docs.sigstore.dev/
