package intoto

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
)

func TestSign(t *testing.T) {

	statementPath := "testdata/test-statement.json"

	ctx := context.Background()
	timeout := 10 * time.Second
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	privateKey := crypto.PrivateKey(key)

	signOptions := SignOptions{
		CertPath: "",
		PrivKey:  privateKey,
		HashFunc: crypto.SHA256,
	}

	statementOpt := StatementOptions{
		StatementPath: statementPath,
		PredicateURI:  SlsaPredicateType,
		Statement:     nil,
	}

	signedStatement, err := sign(ctx, timeout, statementOpt, signOptions)
	if err != nil {
		panic(err)
	}

	println(string(signedStatement))
}
