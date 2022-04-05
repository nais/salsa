# nais slsa action

> in line with the best from abroad

## About

[SLSA](https://github.com/slsa-framework/slsa) is a framework intended to codify and promote
secure [software supply-chain](https://slsa.dev/) practices. This GitHub Action is used to create a SBOM
/ [in-toto attestation](https://github.com/in-toto/attestation) then upload, sign and verify a
generated [provenance](https://slsa.dev/provenance/v0.2) using [cosign](https://github.com/sigstore/cosign).  
All predicate payloads are signed using the [DSSE](https://github.com/secure-systems-lab/dsse).

### Materials

This actions creates attestation with [materials](https://slsa.dev/provenance/v0.2#example) based on dependencies, the
action digest over listed dependencies from a [supported](#support) build tool.

### Support

#### Build tools

> jvm: [gradle](https://gradle.org/) and [maven](https://maven.apache.org/)

> nodejs: [yarn](https://yarnpkg.com/) and [npm](https://www.npmjs.com/)

> [Golang](https://go.dev/)

> [PHP](https://www.php.net/) (with known limitation: no digest over dependencies)

___

* [Usage](#usage)
    * [Git context](#git-context)
    * [Runner context](#runner-context)
    * [Example](#example)
        * [Workflows](#workflows)
        * [Attestation](#attestation)
* [Customizing](#customizing)
    * [Inputs](#inputs)

## Usage

In the examples below we are also using 3 other `required` actions:

* Action to [check out of repository](https://github.com/actions/checkout)
* Action for Google Cloud credentials to establishes [authentication](https://github.com/google-github-actions/auth) to
  Google Cloud
* Action to build and push [Docker images](https://github.com/docker/build-push-action)

### Git context

The github context contains information about the workflow run and the event that triggered the run. By default, this
action uses the [Git context](https://docs.github.com/en/actions/learn-github-actions/contexts#github-context).

### Runner Context

The runner context contains information about the runner that is executing the current job. By default, this action uses
the [Runner context](https://docs.github.com/en/actions/learn-github-actions/contexts#runner-context).

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
  KEY: gcpkms://projects/plattformsikkerhet-dev-496e/locations/europe-north1/keyRings/cosign/cryptoKeys/cosign-test/versions/1

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

* github action [nais-salsa-demo.yml](.github/workflows/nais-salsa-demo.yml)
* generated [salsa provenance](pkg/dsse/testdata/salsa.provenance) with transitive dependencies
* signed [cosign dsse attestation](pkg/dsse/testdata/cosign-dsse-attestation.json)
    * decoded [cosign attestation](pkg/dsse/testdata/cosign-attestation.json)

## Customizing

### inputs

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