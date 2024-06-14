// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// CtsVoice is the golang structure for table cts_voice.
type CtsVoice struct {
	VoiceId   uint        `json:"voiceId"   description:"参数主键"`
	Gender    int         `json:"gender"    description:"性别（男1 女0）"`
	Describe  string      `json:"describe"  description:"描述"`
	Name      string      `json:"name"      description:"名称"`
	VoicePath string      `json:"voicePath" description:"语音路径"`
	Tensor    string      `json:"tensor"    description:"tensor"`
	Count     uint         `json:"count"     description:"使用次数"`
	Rate      float64     `json:"rate"      description:"评分"`
	LastAccessTime    *gtime.Time `json:"lastAccessTime"    description:"最后访问时间"`
	CreateBy  uint        `json:"createBy"  description:"创建者"`
	UpdateBy  uint        `json:"updateBy"  description:"更新者"`
	Remark    string      `json:"remark"    description:"备注"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"修改时间"`
}
