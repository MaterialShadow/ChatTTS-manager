package utils

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/tiger1103/gfast/v3/api/v1/cts"
	"github.com/tiger1103/gfast/v3/internal/app/cts/common"
	"github.com/tiger1103/gfast/v3/library/libUtils"
	"log"
)

type ApiAudioGenReq struct {
	Text                  string      `p:"text" json:"text"`
	Speed                 int         `p:"speed" json:"speed"`
	Temperature           float64     `p:"temperature" json:"temperature"`
	TopP                  float64     `p:"top_P" json:"top_P"`
	TopK                  int         `p:"top_K" json:"top_K"`
	RefineOral            int         `p:"refine_oral" json:"refine_oral"`
	RefineLaugh           int         `p:"refine_laugh" json:"refine_laugh"`
	RefineBreak           int         `p:"refine_break" json:"refine_break"`
	AudioSeedInput        int         `p:"audio_seed_input" json:"audio_seed_input"`
	TextSeedInput         int         `p:"text_seed_input" json:"text_seed_input"`
	RefineTextFlag        bool        `p:"refine_text_flag" json:"refine_text_flag"`
	CustomerSpeakerFlag   bool        `p:"customer_speaker_flag" json:"customer_speaker_flag"`
	CustomerSpeakerTensor interface{} `p:"customer_speaker_tensor" json:"customer_speaker_tensor"`
	CustomerPromotFlag    bool        `p:"customer_promot_flag" json:"customer_promot_flag"`
	CustomerPromot        string      `p:"customer_promot" json:"customer_promot"`
}

var TTSHelper = ttsHelper{}

type ttsHelper struct {
}

func BuildTestAudiotReq(text string, tensor string) (req *ApiAudioGenReq, err error) {
	var tensorList []interface{}
	jsonErr := json.Unmarshal([]byte(tensor), &tensorList)
	if jsonErr != nil {
		err = gerror.New("TTS请求构造失败")
		return nil, err
	}
	return &ApiAudioGenReq{
		Text:                  text,
		Speed:                 2,
		Temperature:           0.3,
		TopP:                  0.7,
		TopK:                  20,
		RefineOral:            2,
		RefineLaugh:           0,
		RefineBreak:           0,
		AudioSeedInput:        2,
		TextSeedInput:         42,
		RefineTextFlag:        false,
		CustomerSpeakerFlag:   true,
		CustomerSpeakerTensor: tensorList,
		CustomerPromotFlag:    false,
		CustomerPromot:        "",
	}, err
}

func BuildAudioGenReq(req *cts.AudioGenReq, tensorList []interface{}) (resp *ApiAudioGenReq, err error) {
	return &ApiAudioGenReq{
		Text:                  req.Text,
		Speed:                 req.Speed,
		Temperature:           req.Temperature,
		TopP:                  req.TopP,
		TopK:                  req.TopK,
		RefineOral:            req.RefineOral,
		RefineLaugh:           req.RefineLaugh,
		RefineBreak:           req.RefineBreak,
		AudioSeedInput:        req.AudioSeedInput,
		TextSeedInput:         req.TextSeedInput,
		RefineTextFlag:        req.RefineTextFlag,
		CustomerSpeakerFlag:   req.CustomerSpeakerFlag,
		CustomerSpeakerTensor: tensorList,
		CustomerPromot:        req.CustomerPromot,
		CustomerPromotFlag:    req.CustomerPromotFlag,
	}, nil
}

func SendGenerateAudioReq(ctx context.Context, apiReq *ApiAudioGenReq) (resp *gclient.Response, err error) {
	/**
	发送生成音频请求
	*/
	reqStr, _ := gjson.EncodeString(apiReq)
	baseUrl := common.API_URL
	url := baseUrl + "/chat/generate_audio"
	return g.Client().Post(ctx, url, reqStr)
}

func SendClenTempAudioReq(ctx context.Context) (resp *gclient.Response, err error) {
	/**
	清理临时音频文件
	*/
	baseUrl := common.API_URL
	url := baseUrl + "/chat/clean/temp/wav"
	return g.Client().Post(ctx, url)
}

func CopyAduio(audioName string) (filePath string, err error) {
	/**
	复制音频文件
	*/
	year, month, day := libUtils.GetDateInfo()
	fileName := guid.S() + ".wav"
	filePath = "/" + year + "/" + month + "/" + day + "/" + fileName
	wavaTempPath := common.AUDIO_TEMP_PATH + "/" + audioName
	saveDir := common.AUDIO_SAVE_PATH + "/" + year + "/" + month + "/" + day + "/"
	if !gfile.Exists(saveDir) {
		libUtils.MkDir(saveDir)
	}
	savePath := saveDir + "/" + fileName
	//复制wave_temp_path到save_path
	err = gfile.CopyFile(wavaTempPath, savePath)
	return filePath, err
}

func ConverTensorStrToList(tensorStr string) (tensorList []interface{}, err error) {
	//将字符串转换为list
	err = json.Unmarshal([]byte(tensorStr), &tensorList)
	if err != nil {
		log.Print(err.Error())
		err = gerror.New("TTS生成失败")
		return tensorList, err
	}
	return tensorList, err
}
