package pexcutor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getMockedProcess(p string, args ...string) *Process {
	return &Process{
		ctx:  context.TODO(),
		rc:   2,
		crc:  0,
		path: p,
		args: args,
	}
}

func TestProcessNew(t *testing.T) {
	p := New(context.TODO(), 3, "ls", "-alh")
	assert.NotNil(t, p)
	assert.Equal(t, p.rc, 3)
}

func TestProcess_Start(t *testing.T) {
	p := getMockedProcess("ls", "-alh")
	err := p.Start()
	assert.NoError(t, err, "error")
}

func TestIntegration(t *testing.T) {
	p := getMockedProcess("date")
	p.SetEnv("TEST_ENV=VALUE")
	err := p.Start()
	assert.NoError(t, err, "error")
	_, _, err = p.GetResult()
	assert.NoError(t, err, "error")
	err = p.Stop()
	assert.Error(t, err, "process should already be stopped")
	p = getMockedProcess("date", "invalid args")
	err = p.Start()
	assert.NoError(t, err, "error")
	_, _, err = p.GetResult()
	assert.Error(t, err, "error")
}
