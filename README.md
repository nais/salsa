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
[![Salsa integration](https://github.com/nais/salsa/actions/workflows/service-account-salsa-integration.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/service-account-salsa-integration.yml)
[![Salsa keyless integration](https://github.com/nais/salsa/actions/workflows/keyless-salsa-integration.yaml/badge.svg)](https://github.com/nais/salsa/actions/workflows/keyless-salsa-integration.yaml)
[![Check pinned workflows](https://github.com/nais/salsa/actions/workflows/ratchet.yml/badge.svg)](https://github.com/nais/salsa/actions/workflows/ratchet.yml)  
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/nais/salsa?color=pink&label=release%40latest&logo=github)
![GitHub last commit](https://img.shields.io/github/last-commit/nais/salsa?color=yellow&logo=github)
[![GitHub license](https://badgen.net/github/license/nais/salsa)](https://github.com/nais/salsa/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nais/salsa)](https://goreportcard.com/report/github.com/nais/salsa)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/nais/salsa/master)


</div>

## About

This is a GitHub Action for generating signed [provenance](https://slsa.dev/provenance/v0.2) about a build and its
related artifacts. Provenance is an attestation (a "software bill of materials") about a software artifact or collection
of artifacts, documenting how an artifact was produced - all in a common format.

Supply chain Levels for Software Artifacts, or [SLSA](https://slsa.dev) (pronounced: *salsa*), is a security framework (
standards, guidelines etc.) to prevent tampering, improve integrity, and secure packages and infrastructure in your
projects, businesses or enterprises.

The action implements the [level 3](https://slsa.dev/spec/v0.1/levels) requirements of
the [SLSA Framework](https://slsa.dev) (as long as it is run in an ephemeral environment)
producing a signed software [attestation](https://github.com/slsa-framework/slsa/blob/main/controls/attestations.md) of
your build and dependencies. The attestation is signed and uploaded to your container registry
using [Cosign](https://github.com/sigstore/cosign)
and can be verified by the salsa cli (also provided in this repo) or using the `cosign verify-attestation` command.

> Signing attestation with KMS (Key Management Service), verification requires access to the corresponding public keys.

- The keys for the [navikt](https://nais.io/slsakeys/navikt.pub) and [nais](https://nais.io/slsakeys/nais.pub)
  organizations can be found at our website (for now), if you use this action at other
  organizations you need to host your keys somewhere appropriate.

> Signing attestation with Cosign Keyless, verification do not require access to the corresponding public keys.

```bash
cosign verify-attestation --type=slsaprovenance image:tag
```

> Disclaimer:
> This is not an official GitHub Action maintained by the SLSA team. It is created by the [nais.io](https://nais.io) team for the purpose of securing supply chains in [NAV](https://github.com/navikt). However, we encourage other organizations/users to use it and even contribute as it is built with open source in mind.

### Built with

[golang](https://golang.org)  
[Cosign](https://github.com/sigstore/cosign)  
[GitHub Actions](https://github.com/features/actions)

### Formats/Standards implemented

Statement type: [in-toto v0.1](https://github.com/in-toto/attestation/)  
Signing envelope: [DSSE](https://github.com/secure-systems-lab/dsse/)  
Predicate type: [Provenance v0.2](https://slsa.dev/provenance/v0.2)

### Materials

This actions creates attestations with [materials](https://slsa.dev/provenance/v0.2#example) based on both runtime and
transitive dependencies, using a supported build tool.

#### Supported build tools

##### JVM

[gradle](https://gradle.org/)  
[maven](https://maven.apache.org/)

##### JS

[yarn](https://yarnpkg.com/)  
[npm](https://www.npmjs.com/)

##### Other

[golang](https://go.dev/)  
[php](https://www.php.net/) (with
known [limitation](https://github.com/composer/composer/issues/2540#issuecomment-850206846): there is no digest over
dependencies)

## Getting started

* [How to use](#how-to-use)
    * [Requirements](#requirements)
    * [Key Management](#kms---key-management-service)
        * [Setup](#google-kms-setup)
        * [Other KMS providers](#other-kms-providers)
        * [Workflow](#workflow-with-service-account-secrets)
    * [Keyless](#keyless-signatures)
        * [Workload identity](#workload-identity)
        * [Workflow](#workflow-with-workload-identity-and-keyless)
* [Customizing](#customizing)
    * [Inputs](#inputs)
        * [Access Private Repositories](#access-private-repositories)
        * [GitHub context](#github-context)
        * [Runner context](#runner-context)
    * [Outputs](#outputs)
* [Release](#release)
    * [Checksums](#checksums)
    * [Verify signature](#verify-signature)

## How to use

### Requirements

The `nais salsa` action supports `KMS providers` or `Cosign Keyless` for signing and/or upload of the attestation to the
registry.  

An authentication step in the Workflow must be set up explicit before the `nais salsa` action. Configure
a [KMS provider](https://docs.sigstore.dev/cosign/kms_support) or a Workload Identity Federation before the
`nais salsa` is run.  

In the workflow examples we use [google-github-actions/auth](https://github.com/google-github-actions/auth)
to authenticate with Google KMS or with a Workload identity.  

[actions/checkout](https://github.com/actions/checkout) is required prior to using this action as `nais salsa`
must have access to your [build manifest](#supported-build-tools) to digest over dependencies.  

### KMS - Key Management Service

The `nais salsa` action use [Cosign](https://github.com/sigstore/cosign) with support
of [KMS](https://github.com/sigstore/cosign/blob/main/KMS.md) to sign and verify the attestation. Cosign
supports all the standard [key management systems](https://github.com/sigstore/cosign/blob/main/USAGE.md).

#### Google KMS Setup

> KMS with Cosign requires some setup at the provider.

KMS is enabled in your Google project:  

* create a keyring
* create key: `Elliptic Curve P-256 key SHA256 Digest`

Service account in project has roles:  

* `Cloud KMS CryptoKey signer/verifier`
* `Cloud KMS viewer Role`

##### Other KMS providers

It is possible to use other KMS providers (this will probably require another GitHub action to be configured).
Read the [Cosign KMS](https://github.com/sigstore/cosign/blob/main/KMS.md) documentation for more information about
providers, their specific setup and key URI formats.

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
        uses: nais/salsa@v0.x
        with:
          key: ${{ env.KEY }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
```

##### Google Authentication

`with.credentials_json` is the [GitHub](https://docs.github.com/en/actions/security-guides/encrypted-secrets) service
account json key.

##### Nais Salsa

`with.key` is the key [URI format](https://github.com/sigstore/cosign/blob/main/KMS.md#gcp) for Google KMS.
Format: `gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION`

`with.github_token` is the GitHub token to authenticate with the registry.

### Keyless Signatures

`nais salsa` supports [Cosign Keyless Signatures](https://github.com/sigstore/cosign/blob/main/KEYLESS.md) signing and
verification of attestations.

> Note: Cosign Keyless this is an experimental feature and is not recommended for production use.

#### Workload identity

Pre-requisites before using Keyless Signatures:

Create a Workload Identity Federation and
follow [steps](https://github.com/google-github-actions/auth#setting-up-workload-identity-federation) to configure
workload identity federation. This can be done with commands using the Google `gcloud` cli or in the
browser [Google Console](https://cloud.google.com/iam/docs/workload-identity-federation).

#### Workflow with workload identity and keyless

```yaml
name: slsa keyless signatures
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
          service_account: name@project-id.iam.gserviceaccount.com
          token_format: "id_token"
          id_token_audience: sigstore
          id_token_include_email: true

      - name: Generate provenance, sign and upload image
        uses: nais/salsa@v0.x
        with:
          identity_token: ${{ steps.google.outputs.id_token }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
        env:
          COSIGN_EXPERIMENTAL: "true"
```

##### Google Authentication

The described `with` fields is required to enable Federation with workload identity and `Cosign` keyless signatures.

`with.workload_identity_provider` is the workload identity provider. The value is retrieved from the Federation
instance created. Format: `projects/$PROJECT/locations/$LOCATION/workloadIdentityPools/$POOL/providers/$PROVIDER`

`with.service_account` is the service account to use for the workload identity. The value is retrieved from the
service account created. Format: `name@project-id.iam.gserviceaccount.com`

`with.token_format` is the token format to use. Cosign expects "id_token".

`with.id_token_audience` is the audience to use for the `id_token`. Cosign expects `sigstore`. `sigstore` audience
must be added to the workload identity provider as an allowed audience.

`with.id_token_include_email` Cosign expects the email to be included in the token.

##### Nais Salsa

The described `with` fields is required for `nais salsa`.

`with.identity_token` is the output `identity_token` from the Google Auth Action.
Format: `steps.steps-id.outputs.id_token`

`with.github_token` is the GitHub token to authenticate with the registry. The password is used by `nais salsa` to
authenticate with the registry to download the image for Cosign to sign and push attestation to the registry.

`with.env.COSIGN_EXPERIMENTAL` is required to be set to `true` for Cosign to enable keyless signatures.

### Signature repository

Cosign defaults to store signatures in the same repo as the image it is signing.
It is possible to specify a different repo for signatures, you can set the `COSIGN_REPOSITORY` environment variable to
store the cosign signatures and attestations, see more specification in
the [cosign docs](https://github.com/sigstore/cosign#specifying-registry)

```yaml
- name: Generate provenance, sign and upload image
  uses: nais/salsa@v0.x
  with:
    key: ${{ secrets.SALSA_KMS_KEY }}
    github_token: ${{ secrets.GITHUB_TOKEN }}
  env:
    COSIGN_REPOSITORY: "registry.io/signatures"
```

Actor must be sure that `with.github_token` has access to the signature repository.

## Customizing

### Inputs

#### Access private repositories

Salsa builds your application to retrieve running dependencies, when the build configuration contains private packages,
the build needs a token with the proper
access. [Maven](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-apache-maven-registry#authenticating-with-a-personal-access-token)
and [gradle](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-gradle-registry)
build tool can authenticate with a `PAT`. Use the `with.github_token` field to authenticate with the registry.

`with.token_key_pattern` can be used to specify a key pattern, other than default `GITHUB_TOKEN`.

#### Maven Options

`with.mvn_opts` - (optional) additional maven options in a comma-delimited string.

**NB!**  
Currently only supports the maven command cli option `-s`, specifying a settings.xml file.

Useful when your project depends on a custom maven settings file or use dependencies from a private repository.
If project depends on dependencies from a private repository, actor need to set GitHub [private token](#access-private-repositories) with proper access right.

```yaml
 - name: Generate provenance, sign and upload image
   uses: nais/salsa@v0.x
   with:
     mvn_opts: "-s .mvn/settings.xml"
     github_token: ${{ secrets.PAT }}
```

#### GitHub context

`with.github_context` - (required) default to `true` to include the github context in the provenance.

The github context contains information about the workflow run and the event that triggered the run. By default, this
action uses the [GitHub context](https://docs.github.com/en/actions/learn-github-actions/contexts#github-context).

#### Runner Context

`with.runner_context` - (required) default to `true` to include the runner context in the provenance.

The runner context contains information about the runner that is executing the current job. By default, this action uses
the [Runner context](https://docs.github.com/en/actions/learn-github-actions/contexts#runner-context).

The Following inputs can be used as `step.with` keys

| Name                | Type   | Default                           | Description                                                                                                                                               | Required |
|---------------------|--------|:----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `key`               | String | ""                                | Private key (cosign.key) or kms provider, used for signing the attestation (Not required for keyless)                                                     | true     |
| `github_token`      | String | ""                                | Token to authenticate and read private packages, the token must have read:packages scope                                                                  | true     |
| `identity_token`    | String | ""                                | Identity token used for Cosign keyless authentication                                                                                                     | False    |
| `image`             | String | $IMAGE                            | The container image to create a attestation for                                                                                                           | False    |
| `docker_user`       | String | github.actor                      | User to login to container registry                                                                                                                       | False    |
| `repo_name`         | String | github.repository                 | The name of the repo/project                                                                                                                              | False    |
| `repo_sub_dir`      | String | ""                                | Specify a subdirectory if build file not found in working root directory                                                                                  | False    |
| `dependencies`      | Bool   | true                              | Set to false if action should not create materials for dependencies (e.g. if build tool is unsupported or repo uses internal/private dependencies)        | False    |
| `token_key_pattern` | String | $GITHUB_TOKEN                     | If a token is provided but the the key pattern is different from the default key pattern "GITHUB_TOKEN"                                                   | False    |
| `build_started_on`  | String | "event.(type if any).head.commit" | Specify a workflow build start time. Default is set to github_context `event.head_commit` or `event.workflow_run.head_commit` depending on workflow usage | False    |
| `mvn_opts`          | String | ""                                | A comma-delimited string with additional maven cli options for the dependence build                                                                       | False    |
| `repo_dir`          | String | $GITHUB_WORKSPACE                 | **Internal value (do not set):** Root of directory to look for build files                                                                                | False    |
| `github_context`    | String | ${{ toJSON(github) }}             | **Internal value (do not set):** the [github context](#github-context) object in json                                                                     | False    |
| `runner_context`    | String | ${{ toJSON(runner) }}             | **Internal value (do not set):** the [runner context](#runner-context) object in json                                                                     | False    |

### Outputs

* `provenance_file_path` [SLSA provenance](pkg/dsse/testdata/salsa.provenance)
* `raw_file_path` [Signed Cosign dsse attestation](pkg/dsse/testdata/cosign-dsse-attestation.json)

## Release

### Checksums

We generate a `checksums.txt` file and upload it with the release, so users can validate if the downloaded files are
correct. All files are by default digested with algorithm `sha256`.

### Verify signature

The release [artifacts](https://github.com/nais/salsa/releases) are signed with Cosign
and can be verified by using the [public signing key](https://github.com/nais/salsa/blob/main/cosign.pub).

```shell
cosign verify-blob --key cosign.pub --signature salsa.tar.gz.sig salsa.tar.gz
```

> Verified OK