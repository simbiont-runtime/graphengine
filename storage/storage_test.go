// ---

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMVCCStorage_Open(t *testing.T) {
	s, err := Open(t.TempDir())
	assert.Nil(t, err)
	assert.NotNil(t, s)
	err = s.Close()
	assert.Nil(t, err)
}

func TestMVCCStorage_CurrentVersion(t *testing.T) {
	s, err := Open(t.TempDir())
	assert.Nil(t, err)
	assert.NotNil(t, s)
	ver := s.CurrentVersion()
	assert.NotZero(t, ver)
}
