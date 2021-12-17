package intoto

import (
	"bytes"
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sigstore/cosign/pkg/oci"
	"github.com/sigstore/cosign/pkg/oci/static"
	"github.com/sigstore/cosign/pkg/types"
	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/dsse"
	signatureoptions "github.com/sigstore/sigstore/pkg/signature/options"
)

type SignOptions struct {
	CertPath string
	PrivKey  *crypto.PrivateKey
	HashFunc crypto.Hash
}

type StatementOptions struct {
	StatementPath string
	PredicateURI  string
	Statement     *Statement
}

func sign(ctx context.Context, timeout time.Duration, so StatementOptions, opts SignOptions) (oci.Signature, error) {

	if timeout != 0 {
		var cancelFn context.CancelFunc
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
		defer cancelFn()
	}

	sv, err := signature.LoadSignerVerifier(opts.PrivKey, opts.HashFunc)
	if err != nil {
		return nil, errors.Wrap(err, "getting signer")
	}

	wrapped := dsse.WrapSigner(sv, so.PredicateURI)

	fmt.Fprintln(os.Stderr, "Using payload from:", so.StatementPath)
	statement, err := os.Open(so.StatementPath)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	payload, err := json.Marshal(statement)
	if err != nil {
		return nil, err
	}
	signedPayload, err := wrapped.SignMessage(bytes.NewReader(payload), signatureoptions.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "signing")
	}

	opt := []static.Option{static.WithLayerMediaType(types.DssePayloadType)}
	/*if sv.Cert != nil {
	    opts = append(opt, static.WithCertChain(sv.Cert, sv.Chain))
	}*/

	sig, err := static.NewAttestation(signedPayload, opt...)
	if err != nil {
		return nil, err
	}
	return sig, nil
}
