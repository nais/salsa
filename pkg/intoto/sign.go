package intoto

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/dsse"
	signatureoptions "github.com/sigstore/sigstore/pkg/signature/options"
)

type SignOptions struct {
	CertPath string
	PrivKey  crypto.PrivateKey
	HashFunc crypto.Hash
}

type StatementOptions struct {
	StatementPath string
	PredicateURI  string
	Statement     *Statement
}

func sign(ctx context.Context, timeout time.Duration, so StatementOptions, opts SignOptions) ([]byte, error) {

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
	statement, err := ioutil.ReadFile(so.StatementPath)
	if err != nil {
		return nil, err
	}

	signedPayload, err := wrapped.SignMessage(bytes.NewReader(statement), signatureoptions.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "signing")
	}
	return signedPayload, nil
}
