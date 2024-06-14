// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CtsVoiceDao is the data access object for table cts_voice.
type CtsVoiceDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns CtsVoiceColumns // columns contains all the column names of Table for convenient usage.
}

// CtsVoiceColumns defines and stores column names for table cts_voice.
type CtsVoiceColumns struct {
	VoiceId   string // 参数主键
	Gender    string // 性别（男1 女0）
	Name      string // 名称
	Describe  string // 描述
	Tensor    string // tensor
	VoicePath string // VoicePath
	Count     string // Count
	Rate      string // Rate
	LastAccessTime string //最后访问时间
	CreateBy  string // 创建者
	UpdateBy  string // 更新者
	Remark    string // 备注
	CreatedAt string // 创建时间
	UpdatedAt string // 修改时间
}

// ctsVoiceColumns holds the columns for table cts_voice.
var ctsVoiceColumns = CtsVoiceColumns{
	VoiceId:   "voice_id",
	Gender:    "gender",
	Name:      "name",
	Describe:  "describe",
	Tensor:    "tensor",
	Count:     "count",
	Rate:      "rate",
	CreateBy:  "create_by",
	UpdateBy:  "update_by",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewCtsVoiceDao creates and returns a new DAO object for table data access.
func NewCtsVoiceDao() *CtsVoiceDao {
	return &CtsVoiceDao{
		group:   "default",
		table:   "cts_voice",
		columns: ctsVoiceColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CtsVoiceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CtsVoiceDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CtsVoiceDao) Columns() CtsVoiceColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CtsVoiceDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CtsVoiceDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CtsVoiceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
