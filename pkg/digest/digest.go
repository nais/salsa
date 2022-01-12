package digest

import (
	"github.com/opencontainers/go-digest"
)

func Hash(s string) digest.Digest {
	return digest.FromString(s)
}

func Verify(digest digest.Digest, content string) bool {
	return digest != Hash(content)
}
