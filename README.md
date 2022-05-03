<br />
<div align="center">
<a href="https://github.com/nais/salsa">
    <img src="https://slsa.dev/images/SLSA-Badge-full-level2.svg"/>
 </a>
    <h2 align="center">nais SLSA Provenance Action</h2>
</div>

<div id="top"></div>

<div align="center">

[![Salsa build & release](https://github.com/nais/salsa/actions/workflows/main.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/main.yml)
[![Salsa integration](https://github.com/nais/salsa/actions/workflows/nais-salsa-integration.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/nais-salsa-integration.yml)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/nais/salsa?color=pink&label=release%40latest&logo=github)
![GitHub last commit](https://img.shields.io/github/last-commit/nais/salsa?color=yellow&logo=github)  
[![GitHub stars](https://img.shields.io/github/stars/nais/salsa.svg)](https://github.com/nais/salsa/stargazers?color=purple)
[![GitHub license](https://badgen.net/github/license/nais/salsa)](https://github.com/nais/salsa/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nais/salsa)](https://goreportcard.com/report/github.com/nais/salsa)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/nais/salsa/master)


</div>

## About

The project is started as an initiative by [nais](https://nais.io/)

> NAV's Application Infrastructure Service

to establishing a [Level 2](#level-2-after-the-build) cryptographic chain of custody between trusted builds and our
release and code-signing workflows. [SLSA](https://github.com/slsa-framework/slsa) is a framework intended to codify and
promote secure [software supply-chain](https://slsa.dev/) practices. This GitHub Action can be used to create, upload,
sign and verify a SBOM / [in-toto attestation](https://github.com/in-toto/attestation)
also called a [provenance](https://slsa.dev/provenance/v0.2) using [cosign](https://github.com/sigstore/cosign). All
predicate payloads are signed using the [DSSE](https://github.com/secure-systems-lab/dsse).

This is not an official GitHub Action set up and maintained by the SLSA team. This GitHub Action is built to provide
teams and developers with the ability to trace software back to the source and define the moving parts in a complex
supply chain.

### SLSA Security Levels

SLSA is organized into a series of [levels](https://slsa.dev/spec/v0.1/levels) that provide increasing integrity
guarantees. This gives the action user confidence that software hasn’t been tampered with and can be securely traced
back to its source.

#### Level 2 After the build

This Action fulfills the requirements for [level 2](https://slsa.dev/spec/v0.1/index) and shows more trustworthiness in
the build, builders are source-aware, and signatures are used to prevent provenance being tampered with.

### Materials

This actions creates attestation with [materials](https://slsa.dev/provenance/v0.2#example) based on both runtime and
transitive dependencies, the action digest over listed dependencies from a [supported](#support) build tool.

### Support

#### Build tools

* jvm
    * [gradle](https://gradle.org/)
    * [maven](https://maven.apache.org/)
* nodejs
    * [yarn](https://yarnpkg.com/)
    * [npm](https://www.npmjs.com/)
* [golang](https://go.dev/)
* [php](https://www.php.net/) (with
  known [limitation](https://github.com/composer/composer/issues/2540#issuecomment-850206846): there is no digest over
  dependencies)

___

* [Usage](#usage)
    * [Key Management](#key-management)
        * [Setup](#setup)
        * [Other KMS providers](#other-kms-providers)
    * [Example](#example)
        * [Workflows](#workflows)
        * [Attestation](#attestation)
* [Customizing](#customizing)
    * [Inputs](#inputs)
        * [Github context](#git-context)
        * [Runner context](#runner-context)
* [Release](#release)
    * [Checksums](#checksums)
    * [Verify signature](#verify-signature)
        * [Cosign](#cosign)

## Usage

In the examples below we are also using 2 other `required` actions:

* Action to [check out of repository](https://github.com/actions/checkout)
* Action for Google Cloud credentials to establishes [authentication](https://github.com/google-github-actions/auth) to
  Google. Customizing for [other providers](#other-kms-providers)

Not `required`:

* Action to build and push [Docker images](https://github.com/docker/build-push-action)

### Key Management

This action use [cosign](https://github.com/sigstore/cosign) with
supported [Google KMS](https://github.com/sigstore/cosign/blob/main/KMS.md) keys for signing and verifying the
attestation. Cosign supports all the
standard [key management systems](https://github.com/sigstore/cosign/blob/main/USAGE.md). If your project requires other
providers please feel free to submit an [issue](https://github.com/nais/salsa/issues)
or [pull request](https://github.com/nais/salsa/pulls).

#### Setup

KMS with cosign requires som pre-setup at you provider. In short for Google KMS:

1. KMS is enabled in your Google project
    1. create a keyring
        1. create HMS keys: `Elliptic Curve P-256 key SHA256 Digest`
2. Serviceuser in project has roles:
    1. Cloud KMS CryptoKey signer/verifier
    2. Cloud KMS viewer Role
3. Configure [Github](https://docs.github.com/en/actions/security-guides/encrypted-secrets) actions secret containing
   the serviceuser
4. Configure `with.key` with the right [URI format](https://github.com/sigstore/cosign/blob/main/KMS.md#gcp) for
   Google: `gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION`

##### Other KMS providers

It should be possible to switch [Key Management](#key-management) provider
from [Google](https://github.com/google-github-actions/auth) with for
example [Azure](https://github.com/marketplace/actions/azure-login). See
the [cosign KMS](https://github.com/sigstore/cosign/blob/main/KMS.md)
for more information about provider setup and key URI formats.

### Example

#### Workflows

```yaml
name: ci

on:
  push:
    branches:
      - 'main'

env:
  IMAGE: ttl.sh/nais/salsa-test:1h
  KEY: gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION

jobs:
  provenance:
    runs-on: ubuntu-20.04
    steps:

      - name: Checkout Code
        uses: actions/checkout@v3

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v0
        with:
          credentials_json: ${{ secrets.GCP_CREDENTIALS }}

      - name: Build and push
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: ${{ env.IMAGE }}

      - name: Provenance, upload and sign attestation
        uses: nais/salsa@v0.1
        with:
          key: ${{ env.KEY }}
          docker_pwd: ${{ secrets.GITHUB_TOKEN }}
```

#### Attestation

* Example of a Github action [nais-salsa-integration.yml](.github/workflows/nais-salsa-integration.yml)
* An example of a generated [salsa provenance](pkg/dsse/testdata/salsa.provenance) with transitive dependencies
* An example of a signed [cosign dsse attestation](pkg/dsse/testdata/cosign-dsse-attestation.json)
    * result after an decoded [cosign attestation](pkg/dsse/testdata/cosign-attestation.json)

## Customizing

### inputs

#### Github context

The github context contains information about the workflow run and the event that triggered the run. By default, this
action uses the [Github context](https://docs.github.com/en/actions/learn-github-actions/contexts#github-context).

#### Runner Context

The runner context contains information about the runner that is executing the current job. By default, this action uses
the [Runner context](https://docs.github.com/en/actions/learn-github-actions/contexts#runner-context).

The Following inputs can be used as `step.with` keys

| Name             | Type   | Default               | Description                                                                          | Required |
|------------------|--------|:----------------------|--------------------------------------------------------------------------------------|----------|
| `key`            | String | ""                    | Private key (cosign.key) or kms provider, for signing the attestation                | True     |
| `docker_pwd`     | String | ""                    | Password for docker                                                                  | True     |
| `image`          | String | $IMAGE                | Docker image to sign                                                                 | False    |
| `docker_user`    | String | github.actor          | User to login to docker                                                              | False    |
| `repo_name`      | String | github.repository     | This will name the generated provenance                                              | False    |
| `repo_sub_dir`   | String | ""                    | Specify a sub directory if build file not found in working root directory            | False    |
| `dependencies`   | Bool   | true                  | Set to false if action should not digest over dependencies                           | False    |
| `repo_dir`       | String | $GITHUB_WORKSPACE     | **Internal value (do notset):** Root of directory to look for build files            | False    |
| `github_context` | String | ${{ toJSON(github) }} | **Internal value (do notset):** the [github context](#git-context) object in json    | False    |
| `runner_context` | String | ${{ toJSON(runner) }} | **Internal value (do notset):** the [runner context](#runner-context) object in json | False    |

## Release

### Checksums

`nais salsa` generates a `checksums.txt` file and uploads it with the release, so users can validate if the downloaded
files are correct. All files are by default encoded with algorithm `sha256`.

### Verify signature

`nais salsa` sign release [artifacts](https://github.com/nais/salsa/releases) to ensures that the artifacts have been
generated by `nais salsa` and users can verify that by comparing the generated signature with
`nais salsa` [public signing key](https://github.com/nais/salsa/blob/main/cosign.pub).

#### Cosign

`nais salsa` sign artifacts with [cosign](https://github.com/sigstore/cosign).

Users can then verify the signature with:

```shell
cosign verify-blob --key cosign.pub --signature salsa.tar.gz.sig salsa.tar.gz
```
> Verified OK
