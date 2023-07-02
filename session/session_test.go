// ---

package session_test

import (
	"testing"

	"github.com/simbiont-runtime/graphengine"
	"github.com/simbiont-runtime/graphengine/session"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)
	assert.NotNil(db)

	s := session.New(db.Store(), db.Catalog())
	assert.NotNil(s)
}

func TestSession_OnClosed(t *testing.T) {
	assert := assert.New(t)
	db, err := graphengine.Open(t.TempDir(), nil)
	assert.Nil(err)
	assert.NotNil(db)

	s := session.New(db.Store(), db.Catalog())
	assert.NotNil(s)

	var closed bool
	s.OnClosed(func(session *session.Session) {
		closed = true
	})

	s.Close()
	assert.True(closed)
}
