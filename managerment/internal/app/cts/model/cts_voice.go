package model

type CtsVoiceInfoRes struct {
	VoiceId   uint    `orm:"voice_id,primary"  json:"voiceId"`
	Gender    int     `orm:"gender"   json:"gender"`
	Name      string  `orm:"name"   json:"name"`
	Describe  string  `orm:"describe"   json:"describe"`
	VoicePath string  `orm:"voice_path"  json:"voicePath"`
	Tensor    string  `orm:"tensor"   json:"tensor"`
	Count     uint    `orm:"count"  json:"count"`
	Rate      float64 `orm:"rate"  json:"rate"`
}

type CtsVoiceSimple struct {
	Name     string      `orm:"name"     tag:"name"     json:"name"`
	Gender   int         `orm:"gender"   tag:"gender"   json:"gender"`
	Describe string      `orm:"describe" tag:"describe" json:"describe"`
	Tensor   interface{} `orm:"tensor"   tag:"tensor"   json:"tensor"`
}
