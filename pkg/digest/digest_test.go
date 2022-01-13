package digest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDigest(t *testing.T) {
	hello := "hello"
	result, err := Hash(hello)
	assert.NoError(t, err)
	assert.Falsef(t, result.Verify(hello), "Should return false for verifying OK")

	result, err = Hash("yello")
	assert.NoError(t, err)
	assert.True(t, result.Verify(hello), "Should return true for verifying NOT OK")

	checksum, err := Digest("aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=").DecodeToString()
	assert.NoError(t, err)
	assert.Equal(t, "690518917cd5b2e7ccf83c05d5a13ed317dc53ee7a27009a2e2724d02966313c", checksum)
}
