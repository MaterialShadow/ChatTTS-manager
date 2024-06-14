package cts

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	commonApi "github.com/tiger1103/gfast/v3/api/v1/common"
	"github.com/tiger1103/gfast/v3/internal/app/cts/model/entity"
)

type VoiceSearchReq struct {
	g.Meta `path:"/voice/list" tags:"声音管理" method:"GET" summary:"声音列表"`
	Gender string `p:"gender"`
	Name   string `p:"name"`
	commonApi.PageReq
}

type VoiceSearchRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
	VoiceList []*entity.CtsVoice `json:"voiceList"`
}

type VoiceAddReq struct {
	g.Meta    `path:"/voice/add" tags:"声音管理" method:"POST" summary:"添加声音"`
	Gender    string `p:"gender" json:"gender" v:"required#性别必须"`
	Name      string `p:"name" json:"name" v:"required#名称不能为空"`
	Describe  string `p:"describe" json:"describe" v:"required#描述不能为空"`
	Tensor    string `p:"tensor" json:"tensor" v:"required#tensor不能为空"`
	AudioName string `p:"audioName" json:"audioName"`
}

type VoiceAddRes struct {
	g.Meta `mime:"application/json"`
}

type VoiceRateReq struct {
	g.Meta  `path:"/voice/rate" tags:"声音管理" method:"PUT" summary:"添加声音"`
	VoiceId uint    `p:"voicdeId" json:"voiceId" v:"required#voicdeId不能为空"`
	Rate    float64 `p:"rate" json:"rate" v:"required#评分不能为空"`
}

type VoiceRateRes struct {
	g.Meta `mime:"application/json"`
}

type VoiceEditReq struct {
	g.Meta   `path:"/voice/list" tags:"声音管理" method:"PUT" summary:"修改声音"`
	VoiceId  uint   `p:"voicdeId" json:"voiceId" v:"required#voicdeId不能为空"`
	Gender   string `p:"gender" json:"gender" v:"required#性别必须"`
	Name     string `p:"name" json:"name" v:"required#名称不能为空"`
	Describe string `p:"describe" json:"describe" v:"required#描述不能为空"`
	Tensor   string `p:"tensor" json:"tensor" v:"required#tensor不能为空"`
}

type VoiceEditRes struct {
	g.Meta `mime:"application/json"`
}

type VoiceDeleteReq struct {
	g.Meta  `path:"/voice/delete" tags:"声音管理" method:"DELETE" summary:"删除声音"`
	VoiceId uint `p:"voiceId" json:"voiceId" v:"required#voicdeId不能为空"`
}

type VoiceDeleteRes struct {
	g.Meta `mime:"application/json"`
}

type VoiceBatchDeleteReq struct {
	g.Meta   `path:"/voice/batchDelete" tags:"声音管理" method:"DELETE" summary:"批量删除声音"`
	VoiceIds []uint `p:"voiceIds" json:"voiceIds" v:"required#voicdeIds不能为空"`
}

type VoiceBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

type VoiceGetReq struct {
	g.Meta  `path:"/voice/get" tags:"声音管理" method:"GET" summary:"获取声音"`
	VoiceId uint `p:"voiceId" json:"voiceId" v:"required#voicdeId不能为空"`
}

type VoiceGetRes struct {
	g.Meta `mime:"application/json"`
	*entity.CtsVoice
}

type AudioPlayReq struct {
	g.Meta  `path:"/voice/play" tags:"声音管理" method:"GET" summary:"播放声音"`
	VoiceId uint `p:"voiceId" json:"voiceId" v:"required#voicdeId不能为空"`
}

type AudioPlayRes struct {
	g.Meta `mime:"audio/x-wav" type:"stream" name:"文件流"`
}

type AudioTestGenReq struct {
	g.Meta  `path:"/voice/gen" tags:"声音管理" method:"POST" summary:"生成测试音频"`
	VoiceId uint `p:"voiceId" json:"voiceId" v:"required#voicdeId不能为空"`
}

type AudioTestGenRes struct {
	g.Meta `mime:"application/json"`
}

type VoiceDownloadReq struct {
	g.Meta  `path:"/voice/download" tags:"声音管理" method:"GET" summary:"音色下载"`
	VoiceId uint `p:"voiceId" json:"voiceId"`
}

type VoiceDownloadRes struct {
	//文件下载
	g.Meta `mime:"application/octet-stream" type:"attachment" name:"文件流"`
}

type VoiceUploadReq struct {
	g.Meta `path:"/voice/upload" tags:"声音管理" method:"POST" summary:"上传声音"`
	File   *ghttp.UploadFile `json:"file" type:"file" dc:"选择上传文件"`
}

type VoiceUploadRes struct {
	g.Meta `mime:"application/json"`
}
