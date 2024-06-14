package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/tiger1103/gfast/v3/api/v1/cts"
	"github.com/tiger1103/gfast/v3/internal/app/cts/common"
	"github.com/tiger1103/gfast/v3/internal/app/cts/service"
	"github.com/tiger1103/gfast/v3/internal/app/cts/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var TTS = ttsController{}

type ttsController struct {
	BaseController
}

func (c *ttsController) Generate(ctx context.Context, req *cts.AudioGenReq) (res *cts.AudioGenRes, err error) {
	res = &cts.AudioGenRes{}
	var tensorList []interface{}
	if req.CustomerSpeakerFlag {
		//查询tensor
		customerSpeakerId := req.CustomerSpeakerId
		if customerSpeakerId == 0 {
			err = gerror.New("请选择音色")
			return res, err
		}
		voiceRes, err := service.CtsVoice().GetById(ctx, customerSpeakerId)
		if err != nil {
			log.Print(err.Error())
			err = gerror.New("该音色不存在")
			return res, err
		}
		tensor := voiceRes.Tensor
		err = json.Unmarshal([]byte(tensor), &tensorList)
		if err != nil {
			log.Print(err.Error())
			err = gerror.New("TTS生成失败")
			return res, err
		}
		//更新音色使用信息
		err = service.CtsVoice().UpdateUsage(ctx, customerSpeakerId)
		if err != nil {
			log.Print(err.Error())
			err = gerror.New("TTS生成失败")
			return res, err
		}
	}
	apiRequest, _ := utils.BuildAudioGenReq(req, tensorList)
	r, err := utils.SendGenerateAudioReq(ctx, apiRequest)
	defer r.Close()
	if r == nil || r.Response == nil {
		err = gerror.New("TTS生成失败")
		return res, err
	}
	if r.Response.StatusCode != 200 || err != nil {
		err = gerror.New("TTS生成失败")
		return res, err
	}
	if j, err := gjson.DecodeToJson(r.ReadAllString()); err != nil {
		err = gerror.New("TTS生成失败")
		return res, err
	} else {
		res.OutPutText = j.Get("out_put_text").String()
		res.Tensor = j.Get("tensor").Slice()
		res.AduioName = j.Get("audio_name").String()
	}
	return
}

func (c *ttsController) Download(ctx context.Context, req *cts.AduioDownloadReq) (res *cts.AduioDownloadRes, err error) {
	res = &cts.AduioDownloadRes{}
	tempPath := common.AUDIO_TEMP_PATH
	filePath := tempPath + "/" + req.AudioName + ".wav"

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err.Error())
		g.RequestFromCtx(ctx).Response.WriteStatusExit(http.StatusInternalServerError)
		return nil, err
	}
	defer file.Close() // 确保文件在函数结束时被关闭

	fileInfo, err := file.Stat()
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteStatusExit(http.StatusInternalServerError)
		return nil, err
	}

	header := g.RequestFromCtx(ctx).Response.Header()
	header.Set("Content-Type", "audio/x-wav")
	header.Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	_, err = io.Copy(g.RequestFromCtx(ctx).Response.Writer, file)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteStatusExit(http.StatusInternalServerError)
		return nil, err
	}
	return res, nil
}

func (c *ttsController) TempClean(ctx context.Context, req *cts.AudioTempCleanReq) (res *cts.AudioTempCleanRes, err error) {
	r, _ := utils.SendClenTempAudioReq(ctx)
	if r == nil || r.Response == nil {
		err = gerror.New("清理缓存音频失败")
		return res, err
	} else if r.Response.StatusCode != 200 {
		err = gerror.New("清理缓存音频失败")
		return res, err
	}
	return
}
