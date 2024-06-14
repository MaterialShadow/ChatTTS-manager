package cts

import "github.com/gogf/gf/v2/frame/g"

type AudioGenReq struct {
	g.Meta              `path:"/tts/gen" tags:"TTS管理" method:"POST" summary:"音频生成"`
	Text                string  `p:"text" json:"text"v:"required#文本不能为空"`
	Speed               int     `p:"speed" json:"speed"v:"required#语速不能为空"`
	Temperature         float64 `p:"temperature" json:"temperature" v:"required#温度不能为空"`
	TopP                float64 `p:"top_P" json:"top_P" v:"required#top_P不能为空"`
	TopK                int     `p:"top_K" json:"top_K" v:"required#top_K不能为空"`
	RefineOral          int     `p:"refine_oral" json:"refine_oral" v:"required#refine_oral不能为空"`
	RefineLaugh         int     `p:"refine_laugh" json:"refine_laugh" v:"required#refine_laugh不能为空"`
	RefineBreak         int     `p:"refine_break" json:"refine_break" v:"required#refine_break不能为空"`
	AudioSeedInput      int     `p:"audio_seed_input" json:"audio_seed_input" v:"required#audio_seed_input不能为空"`
	TextSeedInput       int     `p:"text_seed_input" json:"text_seed_input" v:"required#text_seed_input不能为空"`
	RefineTextFlag      bool    `p:"refine_text_flag" json:"refine_text_flag" v:"required#refine_text_flag不能为空"`
	CustomerSpeakerFlag bool    `p:"customer_speaker_flag" json:"customer_speaker_flag" v:"required#customer_speaker_flag不能为空"`
	CustomerSpeakerId   uint    `p:"customer_speaker_id" json:"customer_speaker_id"`
	CustomerPromotFlag  bool    `p:"customer_promot_flag" json:"customer_promot_flag"`
	CustomerPromot      string  `p:"customer_prompt" json:"customer_promot"`
}

type AudioGenRes struct {
	g.Meta     `mime:"application/octet-stream" type:"stream" name:"文件流"`
	OutPutText string        `json:"out_put_text"`
	AduioName  string        `json:"audio_name"`
	Tensor     []interface{} `json:"tensor"`
}

type AduioDownloadReq struct {
	g.Meta    `path:"/tts/download" tags:"TTS管理" method:"GET" summary:"音频下载"`
	AudioName string `p:"audio_name" json:"audio_name" v:"required#audio_name不能为空"`
}

type AduioDownloadRes struct {
	g.Meta `mime:"audio/x-wav" type:"stream" name:"文件流"`
}

type AudioTempCleanReq struct {
	g.Meta `path:"/tts/clean" tags:"TTS管理" method:"DELETE" summary:"清除temp缓存文件"`
}

type AudioTempCleanRes struct {
	g.Meta `mime:"application/json"`
}
