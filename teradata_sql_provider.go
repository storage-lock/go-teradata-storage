package teradata_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
)

// TeradataSqlProvider Teradata 数据库的 SQL 方言提供者
// Teradata 不支持 MySQL 的 UNIX_TIMESTAMP(NOW())。
type TeradataSqlProvider struct {
	*sql_based_storage.Sql92Provider
}

func NewTeradataSqlProvider() *TeradataSqlProvider {
	return &TeradataSqlProvider{
		Sql92Provider: sql_based_storage.NewSql92Provider(),
	}
}

// NowTimestampSql Teradata 获取当前时间的 Unix 时间戳
// 修复：原来直接继承 Sql92Provider 的 "SELECT UNIX_TIMESTAMP(NOW())"，Teradata 不支持这两个
// MySQL 专有函数，执行即报语法错误，导致 GetTime 完全不可用。Teradata 支持 DATE 相减直接得到
// 天数差（ANSI 日期算术），再加上 CURRENT_TIME 的时分秒换算成秒数，合计为 Unix 时间戳；外层
// CAST AS BIGINT 截断 EXTRACT(SECOND) 可能带出的小数秒。
// 注意：未接入真实 Teradata 驱动验证，是按 Teradata SQL 官方文档做的最佳努力实现。
func (x *TeradataSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return `SELECT CAST(
    (CURRENT_DATE - DATE '1970-01-01') * 86400
    + (EXTRACT(HOUR FROM CURRENT_TIME) * 3600
       + EXTRACT(MINUTE FROM CURRENT_TIME) * 60
       + EXTRACT(SECOND FROM CURRENT_TIME))
    AS BIGINT)`, nil
}
