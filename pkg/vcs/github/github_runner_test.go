package github

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseRunnerContext(t *testing.T) {
	encodedContext := base64.StdEncoding.EncodeToString([]byte(RunnerTestContext))
	context, err := ParseRunner(&encodedContext)
	assert.NoError(t, err)
	assert.Equal(t, "Hosted Agent", context.Name)
	assert.Equal(t, "Linux", context.OS)
	assert.Equal(t, "X64", context.Arch)
	assert.Equal(t, "/opt/hostedtoolcache", context.ToolCache)
	assert.Equal(t, "/home/runner/work/_temp", context.Temp)
}

func TestParseRunnerNoContext(t *testing.T) {
	encodedContext := base64.StdEncoding.EncodeToString([]byte(""))
	context, err := ParseRunner(&encodedContext)
	assert.NoError(t, err)
	assert.Nil(t, context)
}

func TestParseRunnerFailContext(t *testing.T) {
	data := "yolo"
	context, err := ParseRunner(&data)
	assert.Nil(t, context)
	assert.EqualError(t, err, "unmarshal runner context json: invalid character 'Ê' looking for beginning of value")
}

var RunnerTestContext = `{
		"os": "Linux",
		"arch": "X64",
		"name": "Hosted Agent",
		"tool_cache": "/opt/hostedtoolcache",
		"temp": "/home/runner/work/_temp",
		"workspace": "/home/runner/work/nais-salsa-action"
	  }`
