<br />
<div align="center">
<a href="https://github.com/nais/salsa">
    <img src="https://slsa.dev/images/SLSA-Badge-full-level2.svg"/>
 </a>
    <h2 align="center">NAIS SLSA Provenance Action</h2>
</div>

<div id="top"></div>

<div align="center">

[![Salsa CI](https://github.com/nais/salsa/actions/workflows/main.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/main.yml)
[![GitHub license](https://badgen.net/github/license/nais/salsa)](https://github.com/nais/salsa/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/nais/salsa.svg)](https://github.com/nais/salsa/stargazers/)
[![Go Report Card](https://goreportcard.com/badge/github.com/nais/salsa)](https://goreportcard.com/report/github.com/nais/salsa)
[![Github all releases](https://img.shields.io/github/downloads/nais/salsa/total.svg)](https://github.com/nais/salsa/releases/)
[![Github tag](https://badgen.net/github/tag/nais/salsa)](https://github.com/nais/salsa/tags/)

</div>

## About

The project is started as an initiative by the [NAIS](https://nais.io/) team

- `NAV's Application Infrastructure Service`
  to establishing a [Level 2](#level-2-after-the-build) cryptographic chain of custody between trusted builds and our
  release and code-signing workflows. [SLSA](https://github.com/slsa-framework/slsa) is a framework intended to codify
  and promote secure [software supply-chain](https://slsa.dev/) practices. This GitHub Action can be used to create,
  upload, sign and verify a SBOM / [in-toto attestation](https://github.com/in-toto/attestation)
  also called a  [provenance](https://slsa.dev/provenance/v0.2) using [cosign](https://github.com/sigstore/cosign).  
  All predicate payloads are signed using the [DSSE](https://github.com/secure-systems-lab/dsse).

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

This actions creates attestation with [materials](https://slsa.dev/provenance/v0.2#example) based on dependencies, the
action digest over listed dependencies from a [supported](#support) build tool.

### Support

#### Build tools

* JVM: [gradle](https://gradle.org/) and [maven](https://maven.apache.org/)
* NODEJS: [yarn](https://yarnpkg.com/) and [npm](https://www.npmjs.com/)
* [golang](https://go.dev/)
* [PHP](https://www.php.net/) (with
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
        * [Git context](#git-context)
        * [Runner context](#runner-context)

## Usage

In the examples below we are also using 2 other `required` actions:

* Action to [check out of repository](https://github.com/actions/checkout)
* Action for Google Cloud credentials to establishes [authentication](https://github.com/google-github-actions/auth) to
  Google Cloud

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

KMS with cosign requires som pre-setup at you provider, in short for Google KMS:

1. KMS is enabled in Google project
    1. create keyring
        1. create keys: `Elliptic Curve P-256 key SHA256 Digest`
2. Serviceuser in project has roles:
    1. Cloud KMS CryptoKey signer/verifier
    2. Cloud KMS viewer Role
3. Set actions secret in github.com containing the serviceuser credentials.
4. Set `with.key` to the right [URI format](https://github.com/sigstore/cosign/blob/main/KMS.md#gcp) for
   Google: `gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION`

##### Other KMS providers

NB! **This is not tested**, but theoretically it should be possible to switch [Key Management](#key-management) provider
from Google Action with for example [Azure Action](https://github.com/marketplace/actions/azure-login). Please
see [cosign KMS](https://github.com/sigstore/cosign/blob/main/KMS.md)
for more information about setup and URI formats.

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

      - name: 'Authenticate to Google Cloud'
        id: 'google'
        uses: 'google-github-actions/auth@v0'
        with:
          credentials_json: ${{ secrets.GCP_CREDENTIALS }}

      - name: Build and push
        uses: docker/build-push-action@v2
        with:
          push: true
          tags: ${{ env.IMAGE }}

      - name: Provenance, upload and sign attestation
        uses: nais/salsa@v0.0.1-alpha-10
        with:
          image: ${{ env.IMAGE }}
          key: ${{ env.KEY }}
          docker_user: ${{ github.actor }}
          docker_pwd: ${{ secrets.GITHUB_TOKEN }}
```

#### Attestation

* Example of a Github action [nais-salsa-demo.yml](.github/workflows/nais-salsa-demo.yml)
* An example of a generated [salsa provenance](pkg/dsse/testdata/salsa.provenance) with transitive dependencies
* An example of a signed [cosign dsse attestation](pkg/dsse/testdata/cosign-dsse-attestation.json)
    * result after an decoded [cosign attestation](pkg/dsse/testdata/cosign-attestation.json)

## Customizing

### inputs

#### Git context

The github context contains information about the workflow run and the event that triggered the run. By default, this
action uses the [Git context](https://docs.github.com/en/actions/learn-github-actions/contexts#github-context).

#### Runner Context

The runner context contains information about the runner that is executing the current job. By default, this action uses
the [Runner context](https://docs.github.com/en/actions/learn-github-actions/contexts#runner-context).

The Following inputs can be used as `step.with` keys

| Name             | Type   | Description                                                                                                                 | Required |
|------------------|--------|-----------------------------------------------------------------------------------------------------------------------------|----------|
| `key`            | String | The key used to sign the attestation                                                                                        | True     |
| `docker_user`    | String | User to login to docker                                                                                                     | True     |
| `docker_pwd`     | String | Pwd to login to docker                                                                                                      | True     |
| `image`          | String | Docker image to sign. Defaults to $ENV_IMAGE.                                                                               | True     |
| `repo_name`      | String | Name of the file and path to provenance. Used as an relative path under $GITHUB_WORKSPACE. Defaults to "github.repository". | False    |
| `repo_sub_dir`   | String | Specify a sub directory if build file not found in working root directory                                                   | False    |
| `dependencies`   | Bool   | If the provenance should contain dependencies                                                                               | False    |
| `repo_dir`       | String | Internal value (do not set): Root of directory to look for build files. Defaults to $GITHUB_WORKSPACE                       | False    |
| `github_context` | String | Internal value (do not set): the "github" context object in json. The context is used when generating provenance            | False    |
| `runner_context` | String | Internal value (do not set): the "runner" context object in json. The context is used when generating provenance.           | False    |