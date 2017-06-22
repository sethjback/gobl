package leveldb

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	s, err := testDB()
	if !assert.Nil(err) {
		return
	}
	defer s.Close()

	jid := uuid.New().String()
	i1 := index{itype: indexTypeJobDate, key: "start-5000.00" + jid, value: jid}
	i2 := index{itype: indexTypeJobDate, key: "start-5010.00" + jid, value: jid}
	i3 := index{itype: indexTypeJobDate, key: "start-5020.01" + jid, value: jid}
	i4 := index{itype: indexTypeJobDate, key: "start-5030.05" + jid, value: jid}

	e := s.NewIndex(i1)
	assert.Nil(e)
	e = s.NewIndex(i2)
	assert.Nil(e)
	e = s.NewIndex(i3)
	assert.Nil(e)
	e = s.NewIndex(i4)
	assert.Nil(e)

	is, e := s.indexRange(indexTypeJobDate, "start-0", "start-5021")
	assert.Nil(e)
	if assert.Len(is, 3) {
		assert.Contains(is, i1)
		assert.Contains(is, i2)
		assert.Contains(is, i3)
	}

	is, e = s.indexQuery(indexTypeJobDate, "")
	assert.Nil(e)
	if assert.Len(is, 4) {
		assert.Contains(is, i1)
		assert.Contains(is, i2)
		assert.Contains(is, i3)
		assert.Contains(is, i4)
	}

	i, e := s.GetIndex(indexTypeJobDate, i1.key)
	if assert.Nil(e) {
		assert.Equal(i1, *i)
	}

	assert.Nil(s.DeleteIndex(i1))
	_, e = s.GetIndex(indexTypeJobDate, i1.key)
	if assert.NotNil(e) {
		assert.Contains(e.Error(), "not found")
	}
}
