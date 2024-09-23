package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	repo := newRepository()
	err := repo.Save(mock{})
	assert.Nil(t, err)
}

func TestGetAll(t *testing.T) {
	agg1 := mock{ID: "1"}
	agg2 := mock{ID: "2"}
	repo := newRepository()
	repo.Save(agg1)
	repo.Save(agg2)
	resp := repo.GetAll()
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp))
}
