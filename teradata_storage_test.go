package teradata_storage

import (
	"context"
	"testing"

	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestTeradataStorage_GetName(t *testing.T) {
	s := &TeradataStorage{}
	assert.Equal(t, StorageName, s.GetName())
}

func TestTeradataStorage_Capabilities(t *testing.T) {
	s := &TeradataStorage{}
	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityReliableTime)
}

func TestTeradataStorage_Close(t *testing.T) {
	s := &TeradataStorage{}
	err := s.Close(context.Background())
	assert.NoError(t, err)
}
