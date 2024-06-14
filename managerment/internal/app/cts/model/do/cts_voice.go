// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CtsVoice is the golang structure of table cts_voice for DAO operations like Where/Data.
type CtsVoice struct {
	g.Meta    `orm:"table:cts_voice, do:true"`
	VoiceId   interface{} // 参数主键
	Gender    interface{} // 性别（男1 女0）
	Name      interface{}      // 名称
	Describe  interface{} // 描述
	Tensor    interface{} // tensor
	VoicePath interface{} // 音频路径
	Count     interface{} // 使用次数
	Rate      interface{} // 评分
	CreateBy  interface{} // 创建者
	UpdateBy  interface{} // 更新者
	Remark    interface{} // 备注
	LastAccessTime       *gtime.Time // 最后访问时间
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 修改时间
	CreatedBy uint        // 创建者
	UpdatedBy uint        // 更新者
}
