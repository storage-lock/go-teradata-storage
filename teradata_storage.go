package teradata_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
)

// TeradataStorage 基于 Teradata 的存储引擎实现
// Teradata 是企业级数据仓库，支持SQL标准
type TeradataStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewTeradataStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*TeradataStorage, error) {
	baseStorage, err := sql_based_storage.NewSqlBasedStorage(options)
	if err != nil {
		return nil, err
	}
	s := &TeradataStorage{SqlBasedStorage: baseStorage}
	if err := s.Init(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

const StorageName = "teradata-storage"

func (x *TeradataStorage) GetName() string {
	return StorageName
}

func (x *TeradataStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{
		storage.CapabilityCAS,
		storage.CapabilityReliableTime,
	}
}
