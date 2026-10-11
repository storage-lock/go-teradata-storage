package teradata_storage

import (
	"context"
	"strings"
	"time"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
)

// TeradataStorage 基于 Teradata 的存储引擎实现
// Teradata 是企业级数据仓库，支持SQL标准
type TeradataStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewTeradataStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*TeradataStorage, error) {
	// 如果调用方没有显式指定方言（还是默认的 Sql92Provider），换成 Teradata 自己的方言
	// 提供者：Sql92Provider.NowTimestampSql 是 MySQL 专有语法，在 Teradata 上会直接报语法错误
	if _, ok := options.SqlProvider.(*sql_based_storage.Sql92Provider); ok || options.SqlProvider == nil {
		options.SetSqlProvider(NewTeradataSqlProvider())
	}
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

// CreateWithVersion Teradata 唯一主索引冲突时返回 ErrLockAlreadyExists
// 注意：未接入真实 Teradata 驱动验证，错误消息按 Teradata 官方文档做的最佳努力匹配。
func (x *TeradataStorage) CreateWithVersion(ctx context.Context, lockId string, version storage.Version, lockInformation *storage.LockInformation) error {
	err := x.SqlBasedStorage.CreateWithVersion(ctx, lockId, version, lockInformation)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "Duplicate unique prime index") || strings.Contains(msg, "Duplicate unique") {
			return storage_lock.ErrLockAlreadyExists
		}
	}
	return err
}

const StorageName = "teradata-storage"

func (x *TeradataStorage) GetName() string {
	return StorageName
}

func (x *TeradataStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{
		storage.CapabilityCAS,
		storage.CapabilityReliableTime,
		storage.CapabilityWatch,
	}
}

func (x *TeradataStorage) Watch(ctx context.Context, lockId string) (<-chan *storage.LockInformation, error) {
	return storage.NewPollingLockWatcher(x, 500*time.Millisecond).Watch(ctx, lockId)
}
