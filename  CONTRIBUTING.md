# Salsa

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

```text
export GOOGLE_APPLICATION_CREDENTIALS=~/path/to/file/cosign-private-key.json
```

* Install
    * Cosign: https://github.com/sigstore/cosign

## Local install

```
make salsa
```

## Commands

clone: `clones the given project into user defined path`

```
bin/salsa clone --repo salsa --url https://github.com/nais/salsa
```

scan: `Scan files and dependencies for a given project and generate provenance`

```
bin/salsa scan --repo salsa
```

attest: `sign and upload in-toto attestation`

```
bin/salsa attest --verify --repo salsa --predicate salsa.provenance --key gcpkms://projects/$PROJECT/locations/$LOCATION/keyRings/$KEYRING/cryptoKeys/$KEY/versions/$KEY_VERSION  ttl.sh/salsax:1h
```

find: `find artifact from attestations`

```
bin/salsa find go-crypto
```

## Info

Image can be pushed to ttl.sh, who offers free, short-lived (ie: hours), anonymous container image hosting if you just
want to try these commands out.  
**Quick Start** [Cosign](https://github.com/sigstore/cosign#quick-start)  
**ttl.sh** [tt.sh info](https://ttl.sh/)

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