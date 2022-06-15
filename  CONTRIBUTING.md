# Salsa

## About

Salsa CLI is a command line tool to generate, sign and upload a [provenance](https://slsa.dev/provenance/v0.2)

## Developer installation

If you have Go 1.17+, you can setup a development environment

```text
$ git clone https://github.com/nais/salsa
$ cd salsa
$ make salsa
$ $(go env GOPATH)/bin/salsa
```

## Prerequisites

* Google Setup
    * KMS is enabled in project
        * create keyring
            * create keys: `Elliptic Curve P-256 key SHA256 Digest`
    * Serviceuser in project has roles:
        * Cloud KMS CryptoKey signer/verifier
        * Cloud KMS viewer
* Logged in to Google
* Set: `GOOGLE_APPLICATION_CREDENTIALS` with path to .json file containing serviceuser credentials.

```text
export GOOGLE_APPLICATION_CREDENTIALS=~/path/to/file/cosign-private-key.json
```

* Install
    * Cosign: https://github.com/sigstore/cosign

## Commands

clone: `clones the given project into user defined path`

```
salsa clone --repo salsa --owner nais
```

scan: `Scan files and dependencies for a given project and generate provenance`

```
salsa scan --repo salsa
```

attest: `sign and upload in-toto attestation`

```
salsa attest --repo salsa --key gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION  ttl.sh/salsax:1h
```

attest: `verify and download in-toto attestation`

```
salsa attest --verify --repo salsa --key gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION  ttl.sh/salsax:1h
```

find: `find artifact from attestations`

```
salsa find go-crypto
```

## Info

When testing locally the image can be pushed to registry: ttl.sh, who offers free, short-lived (ie: hours), anonymous
container image hosting if you just want to try these commands out. Check
out [Cosign](https://github.com/sigstore/cosign#quick-start) and ttl.sh [tt.sh info](https://ttl.sh/)

Instead of setting a bunch of flags, in home directory create a config file with name ".salsa" (without extension)

```yml
attest:
  key: gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION
...
```

Another possibility is to set Environment variables with prefix `SALSA`

```
SALSA_ATTEST_KEY
```

## Ratchet

Ratchet is a tool for improving the security of CI/CD workflows by automating the process of pinning and unpinning
upstream versions.

More info in [ratchet](https://github.com/sethvargo/ratchet) repo.

### Commands

The workflow run [ratchet.yml](.github/workflows/ratchet.yml) checks if all actions and containers are pinned,
except `ratchet:exclude` and only runs if there are changes in the `.github/workflows` folder.

#### pin

To get a pin a workflow, add the `new workflow` file or `new action` to an existing workflow;

```
make pin workflow=my_workflow.yml
```

#### update

Only applies when dependabot comes with major update for an action e.g. v2 too v3, you manually have to edit the
corresponding ratchet constraint;

from;

```
ratchet:actions/some-action@v2
```

too;

```
ratchet:actions/some-action@v3
```

and then run update and commit your changes;

```
make update workflow=my_workflow.yml
```

Otherwise, just approve the dependabot PR.

#### check

Check if workflows is pinned;

```
make check workflow=my_workflow.yml
```