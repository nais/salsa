<br />
<div align="center">
<a href="https://github.com/nais/salsa">
    <img src="https://slsa.dev/images/SLSA-Badge-full-level3.svg"/>
 </a>
    <h2 align="center">nais SLSA Provenance Action</h2>
</div>

<div id="top"></div>

<div align="center">

[![Salsa build & release](https://github.com/nais/salsa/actions/workflows/main.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/main.yml)
[![Salsa integration](https://github.com/nais/salsa/actions/workflows/nais-salsa-integration.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/nais-salsa-integration.yml)
[![Check pinned workflows](https://github.com/nais/salsa/actions/workflows/ratchet.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/ratchet.yml)  
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/nais/salsa?color=pink&label=release%40latest&logo=github)
![GitHub last commit](https://img.shields.io/github/last-commit/nais/salsa?color=yellow&logo=github)
[![GitHub license](https://badgen.net/github/license/nais/salsa)](https://github.com/nais/salsa/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nais/salsa)](https://goreportcard.com/report/github.com/nais/salsa)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/nais/salsa/master)


</div>

## About

This is a Github Action for generating signed [provenance](https://slsa.dev/provenance/v0.2) about a build and its
related artifacts. Provenance is an attestation (a "software bill of materials") about a software artifact or collection
of artifacts, documenting how an artifact was produced - all in a common format.

Supply chain Levels for Software Artifacts, or [SLSA](https://slsa.dev) (pronounced: *salsa*), is a security framework (
standards, guidelines etc.) to prevent tampering, improve integrity, and secure packages and infrastructure in your
projects, businesses or enterprises.

The action implements the [level 3](https://slsa.dev/spec/v0.1/levels) requirements of
the [SLSA Framework](https://slsa.dev) (as long as it is run in an ephemeral environment) 
producing a signed software [attestation](https://github.com/slsa-framework/slsa/blob/main/controls/attestations.md) of
your build and dependencies. The attestation is signed and uploaded to your container registry
using [cosign](https://github.com/sigstore/cosign)
and can be verified by the salsa cli (also provided in this repo) or using the `cosign verify-attestation` command. 

Verification requires access to the corresponding public keys. The keys for the [navikt](https://nais.io/slsakeys/navikt.pub) and 
[nais](https://nais.io/slsakeys/nais.pub) organizations can be found at our website (for now), if you use this action at other 
organizations you need to host your keys somewhere appropriate.

> Disclaimer:
> This is not an official GitHub Action maintained by the SLSA team. It is created by the [nais.io](https://nais.io) team for the purpose of securing supply chains in [NAV](https://github.com/navikt). However we encourage other organizations/users to use it and even contribute as it is built with open source in mind.

### Built with

* [golang](https://golang.org)
* [cosign](https://github.com/sigstore/cosign)
* [Github Actions](https://github.com/features/actions)

### Formats/Standards implemented

* Statement type: [in-toto v0.1](https://github.com/in-toto/attestation/)
* Signing envelope: [DSSE](https://github.com/secure-systems-lab/dsse/)
* Predicate type: [Provenance v0.2](https://slsa.dev/provenance/v0.2)

### Materials

This actions creates attestations with [materials](https://slsa.dev/provenance/v0.2#example) based on both runtime and
transitive dependencies, using a supported build tool.

#### Supported build tools

* jvm
    * [gradle](https://gradle.org/)
    * [maven](https://maven.apache.org/)
* js
    * [yarn](https://yarnpkg.com/)
    * [npm](https://www.npmjs.com/)
* [golang](https://go.dev/)
* [php](https://www.php.net/) (with
  known [limitation](https://github.com/composer/composer/issues/2540#issuecomment-850206846): there is no digest over
  dependencies)

## Getting started

* [How to use](#how-to-use)
    * [Requirements](#requirements)
    * [Key Management](#key-management)
        * [Setup](#setup)
        * [Other KMS providers](#other-kms-providers)
        * [Workflow](#workflow-with-service-account-secrets)
    * [Keyless](#keyless)
        * [Workflow](#workflow-with-workload-identity)
* [Customizing](#customizing)
    * [Inputs](#inputs)
        * [GitHub context](#github-context)
        * [Runner context](#runner-context)
* [Release](#release)
    * [Checksums](#checksums)
    * [Verify signature](#verify-signature)

## How to use

### Requirements

We support `KMS providers` and `Cosign Keyless` for signing the attestation and for uploading the attestation to the registry.

* Read the [cosign documentation](https://docs.sigstore.dev/cosign/kms_support)
  for more details. Your workflow must set up an explicit authentication step for your KMS provider before the nais
  salsa action.

* In the workflow examples we use [google-github-actions/auth](https://github.com/google-github-actions/auth)
  to authenticate with Google KMS or workload identity for signing the attestation.

* [actions/checkout](https://github.com/actions/checkout) is required prior to using this action as `nais salsa`
  must have access to your [build manifest](#supported-build-tools) to digest over dependencies.

### Key Management

The salsa action use [cosign](https://github.com/sigstore/cosign)
and [KMS](https://github.com/sigstore/cosign/blob/main/KMS.md) for signing and verifying the attestation. Cosign
supports all the standard [key management systems](https://github.com/sigstore/cosign/blob/main/USAGE.md). If your
project requires other providers please feel free to submit an [issue](https://github.com/nais/salsa/issues)
or [pull request](https://github.com/nais/salsa/pulls).

#### Setup

KMS with cosign requires some setup at you provider. In short for Google KMS:

1. KMS is enabled in your Google project
    1. create a keyring
        1. create keys: `Elliptic Curve P-256 key SHA256 Digest`
2. Service accounnt in project has roles:
    1. Cloud KMS CryptoKey signer/verifier
    2. Cloud KMS viewer Role
3. Configure [GitHub](https://docs.github.com/en/actions/security-guides/encrypted-secrets) actions secret containing
   the serviceuser
4. Configure `with.key` with the right [URI format](https://github.com/sigstore/cosign/blob/main/KMS.md#gcp) for
   Google: `gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION`

##### Other KMS providers

It is possible to use other KMS providers.
Read the [cosign KMS](https://github.com/sigstore/cosign/blob/main/KMS.md) documentation for more information about provider setup and key URI formats.

##### workflow with service account secrets

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

      - name: Build and push
        uses: docker/build-push-action@v3
        with:
          push: true
          tags: ${{ env.IMAGE }}

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v0.8.1
        with:
          credentials_json: ${{ secrets.GCP_CREDENTIALS }}

      - name: Provenance, upload and sign attestation
        uses: nais/salsa@v0.1
        with:
          key: ${{ env.KEY }}
          docker_pwd: ${{ secrets.GITHUB_TOKEN }}
```

### Keyless

It is possible to use [keyless](https://github.com/sigstore/cosign/blob/main/KEYLESS.md) signing of attestations with cosign. 

#### workflow with workload identity

The workflow need som setup for Google check
out [google-github-actions/auth](https://github.com/google-github-actions/auth) for more information.

You need to create a Workload Identity Federation.   
Follow this [steps](https://github.com/google-github-actions/auth#setting-up-workload-identity-federation)
either with the `gcloud` cli or in the browser and the [Google Console](https://cloud.google.com/iam/docs/workload-identity-federation).

When the federation and service account is created and added to the federation,
you need to assign the service account with the correct roles:

- roles/iam.serviceAccountTokenCreator

```yaml
name: slsa keyless
on:
  push:
    branches:
      - 'main'
env:
  IMAGE: ttl.sh/nais/salsa-keyless-test:1h
jobs:
  keyless:
    permissions:
      contents: 'read'
      id-token: 'write'
    runs-on: ubuntu-20.04
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3

      - name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: integration-test
          push: true
          tags: ${{ env.IMAGE }}

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v0.8.1
        id: google
        with:
          workload_identity_provider: ${{ secrets.SLSA_WORKLOAD_IDENTITY_PROVIDER }}
          service_account: cosign-kms@plattformsikkerhet-dev-496e.iam.gserviceaccount.com
          token_format: "id_token"
          id_token_audience: ${{ secrets.SLSA_WORKLOAD_IDENTITY_PROVIDER }}
          id_token_include_email: true
          audience: sigstore

      - name: Generate provenance, upload and sign image
        uses: ./
        with:
          identity_token: ${{ steps.google.outputs.id_token }}
          docker_pwd: ${{ secrets.GITHUB_TOKEN }}
        env:
          COSIGN_EXPERIMENTAL: "true"
```


Required `with` fields to work with workload identity and cosign

###### workload_identity_provider

The value is retrieved from the federation created in the Google Console.

Format: `projects/$PROJECT/locations/$LOCATION/workloadIdentityPools/$POOL/providers/$PROVIDER`

###### service_account

The value is retrieved from the service account created in the Google Console.

Format: `name@project-id.iam.gserviceaccount.com`

###### token_format

This is always: "id_token"

###### id_token_audience

Google Auth Action requires this to be set, and it can be the same value as the `workload_identity_provider` field.

###### id_token_include_email

Parameter of whether to include the service account email in the generated token.

###### audience

Additional `aud` claims for the generated `id_token`. This field must contain sigstore.

##### Salsa Action

Required `with` fields for salsa action.

###### identity_token

The output `identity_token` from the Google Auth Action.
Format: `steps.steps-id.outputs.id_token`

###### docker_pwd

This is used by the salsa action to authenticate with the docker registry to download the image for cosign to sign.

###### env

Cosign expects the environment variable `COSIGN_EXPERIMENTAL=1` to be set.

> Note: this is an experimental feature

To publish signed artifacts to a Rekor transparency log and verify their existence in the log.

##### Example output from the workflow

* An example of a generated [slsa provenance](pkg/dsse/testdata/salsa.provenance) with transitive dependencies
* An example of a signed [cosign dsse attestation](pkg/dsse/testdata/cosign-dsse-attestation.json)
    * result after an decoded [cosign attestation](pkg/dsse/testdata/cosign-attestation.json)

## Customizing

### inputs

#### GitHub context

The GitHub context contains information about the workflow run and the event that triggered the run. By default, this
action uses the [GitHub context](https://docs.github.com/en/actions/learn-github-actions/contexts#github-context).

#### Runner Context

The runner context contains information about the runner that is executing the current job. By default, this action uses
the [Runner context](https://docs.github.com/en/actions/learn-github-actions/contexts#runner-context).

The Following inputs can be used as `step.with` keys

| Name             | Type   | Default               | Description                                                                                                                                        | Required |
|------------------|--------|:----------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `docker_pwd`     | String | ""                    | Password for docker                                                                                                                                | True     |
| `key`            | String | ""                    | Private key (cosign.key) or kms provider, used for signing the attestation                                                                         | False    |
| `identity_token` | String | ""                    | Identity token used for cosign keyless authentication                                                                                              | False    |
| `image`          | String | $IMAGE                | The container image to create a attestation for                                                                                                    | False    |
| `docker_user`    | String | github.actor          | User to login to container registry                                                                                                                | False    |
| `repo_name`      | String | github.repository     | The name of the repo/project                                                                                                                       | False    |
| `repo_sub_dir`   | String | ""                    | Specify a subdirectory if build file not found in working root directory                                                                           | False    |
| `dependencies`   | Bool   | true                  | Set to false if action should not create materials for dependencies (e.g. if build tool is unsupported or repo uses internal/private dependencies) | False    |
| `repo_dir`       | String | $GITHUB_WORKSPACE     | **Internal value (do not set):** Root of directory to look for build files                                                                         | False    |
| `github_context` | String | ${{ toJSON(github) }} | **Internal value (do not set):** the [github context](#github-context) object in json                                                              | False    |
| `runner_context` | String | ${{ toJSON(runner) }} | **Internal value (do not set):** the [runner context](#runner-context) object in json                                                              | False    |


## Release

### Checksums

We generate a `checksums.txt` file and upload it with the release, so users can validate if the downloaded files are
correct. All files are by default digested with algorithm `sha256`.

### Verify signature

The release [artifacts](https://github.com/nais/salsa/releases) are signed with cosign
and can be verified by using the [public signing key](https://github.com/nais/salsa/blob/main/cosign.pub).

```shell
cosign verify-blob --key cosign.pub --signature salsa.tar.gz.sig salsa.tar.gz
```

> Verified OK