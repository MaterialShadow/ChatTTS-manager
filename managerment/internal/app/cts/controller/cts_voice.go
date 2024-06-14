package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/tiger1103/gfast/v3/api/v1/cts"
	"github.com/tiger1103/gfast/v3/internal/app/cts/common"
	"github.com/tiger1103/gfast/v3/internal/app/cts/service"
	"github.com/tiger1103/gfast/v3/internal/app/system/consts"
	"github.com/tiger1103/gfast/v3/library/liberr"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var Voice = voiceController{}

type voiceController struct {
	BaseController
}

func (c *voiceController) List(ctx context.Context, req *cts.VoiceSearchReq) (res *cts.VoiceSearchRes, err error) {
	res = new(cts.VoiceSearchRes)
	if req.PageSize == 0 {
		req.PageSize = consts.PageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	total, voices, err := service.CtsVoice().List(ctx, req)
	res.Total = total
	res.CurrentPage = req.PageNum
	res.VoiceList = voices
	return
}

func (c *voiceController) GetVoice(ctx context.Context, req *cts.VoiceGetReq) (res *cts.VoiceGetRes, err error) {
	res = new(cts.VoiceGetRes)
	service.CtsVoice().GetById(ctx, req.VoiceId)
	return
}

func (c *voiceController) Add(ctx context.Context, req *cts.VoiceAddReq) (res *cts.VoiceAddRes, err error) {
	res = new(cts.VoiceAddRes)
	err = service.CtsVoice().Add(ctx, req)
	return
}

func (c *voiceController) Edit(ctx context.Context, req *cts.VoiceEditReq) (res *cts.VoiceEditRes, err error) {
	err = service.CtsVoice().Edit(ctx, req)
	return
}

func (c *voiceController) Delete(ctx context.Context, req *cts.VoiceDeleteReq) (res *cts.VoiceDeleteRes, err error) {
	err = service.CtsVoice().Delete(ctx, req.VoiceId)
	return
}

func (c *voiceController) BatchDelete(ctx context.Context, req *cts.VoiceBatchDeleteReq) (res *cts.VoiceBatchDeleteRes, err error) {
	service.CtsVoice().BatchDelete(ctx, req.VoiceIds)
	return
}

func (c *voiceController) VoiceRate(ctx context.Context, req *cts.VoiceRateReq) (res *cts.VoiceRateRes, err error) {
	err = service.CtsVoice().VoiceRate(ctx, req)
	return
}

func (c *voiceController) AudioPlay(ctx context.Context, req *cts.AudioPlayReq) (res *cts.AudioPlayRes, err error) {
	res = new(cts.AudioPlayRes)
	// 查询语音
	voice, _ := service.CtsVoice().GetById(ctx, req.VoiceId)
	if voice == nil {
		err = gerror.New("语音不存在")
		return
	}
	voicePath := voice.VoicePath
	filePath := filepath.Join(common.AUDIO_SAVE_PATH, voicePath) // 使用filepath.Join来确保路径正确

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err.Error())
		g.RequestFromCtx(ctx).Response.WriteStatusExit(http.StatusInternalServerError)
		return nil, err // 直接返回错误
	}
	defer file.Close() // 确保文件在函数结束时被关闭

	fileInfo, err := file.Stat()
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteStatusExit(http.StatusInternalServerError)
		return nil, err // 直接返回错误
	}

	header := g.RequestFromCtx(ctx).Response.Header()
	header.Set("Content-Type", "audio/x-wav")
	header.Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	_, err = io.Copy(g.RequestFromCtx(ctx).Response.Writer, file)
	if err != nil {
		g.RequestFromCtx(ctx).Response.WriteStatusExit(http.StatusInternalServerError)
		return nil, err // 直接返回错误
	}
	return
}

func (c *voiceController) GenTestAudio(ctx context.Context, req *cts.AudioTestGenReq) (res *cts.AudioTestGenRes, err error) {
	res = new(cts.AudioTestGenRes)
	err = service.CtsVoice().GenTestAudio(ctx, req)
	return
}

func (c *voiceController) Download(ctx context.Context, req *cts.VoiceDownloadReq) (res *cts.AduioDownloadRes, err error) {
	res = new(cts.AduioDownloadRes)
	voiceList, queryErr := service.CtsVoice().DownloadList(ctx, req)
	liberr.ErrIsNil(ctx, queryErr, "查询失败")
	log.Print("voiceList", voiceList)
	// 将结构体列表转换为JSON
	jsonBytes, marshalErr := json.MarshalIndent(voiceList, "", "    ")
	if marshalErr != nil {
		// 错误处理
		err = gerror.New("文件处理失败")
		return nil, err
	}
	var jsonBuffer bytes.Buffer
	jsonBuffer.Write(jsonBytes)
	writer := g.RequestFromCtx(ctx).Response.ResponseWriter
	writer.Header().Set("Content-Type", "application/json")
	// 获得当前时间戳 并转换为字符串 YYYY_MM_DD_HH_MM_SS.json
	fileName := "音色_" + gtime.Now().Format("Y_m_d_H_i_s") + ".json"
	writer.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)
	// 将内存中的JSON数据发送给客户端
	io.Copy(writer, &jsonBuffer)
	return res, nil
}

func (c *voiceController) Upload(ctx context.Context, req *cts.VoiceUploadReq) (res *cts.VoiceUploadRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "请选择要上柴的json")
	}
	err = service.CtsVoice().UploadJson(ctx, req)
	return
}
