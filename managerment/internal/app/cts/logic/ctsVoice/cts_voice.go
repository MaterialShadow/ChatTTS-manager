package ctsVoice

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/tiger1103/gfast/v3/api/v1/cts"
	commonDao "github.com/tiger1103/gfast/v3/internal/app/common/dao"
	commonEntity "github.com/tiger1103/gfast/v3/internal/app/common/model/entity"
	"github.com/tiger1103/gfast/v3/internal/app/cts/common"
	"github.com/tiger1103/gfast/v3/internal/app/cts/dao"
	"github.com/tiger1103/gfast/v3/internal/app/cts/model"
	"github.com/tiger1103/gfast/v3/internal/app/cts/model/do"
	"github.com/tiger1103/gfast/v3/internal/app/cts/model/entity"
	"github.com/tiger1103/gfast/v3/internal/app/cts/service"
	"github.com/tiger1103/gfast/v3/internal/app/cts/utils"
	"github.com/tiger1103/gfast/v3/library/libUtils"
	"github.com/tiger1103/gfast/v3/library/liberr"
	"io"
)

func init() {
	service.RegisterCtsVoice(New())
}

func New() *sCtsVoice {
	return &sCtsVoice{}
}

type sCtsVoice struct {
}

func (s sCtsVoice) List(ctx context.Context, req *cts.VoiceSearchReq) (total interface{}, voiceList []*entity.CtsVoice, err error) {
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.CtsVoice.Ctx(ctx)
		if req.Gender != "" && req.Gender != "0" {
			m = m.Where("gender=?", req.Gender)
		}
		if req.Name != "" {
			m = m.Where(fmt.Sprintf("%s like ?", dao.CtsVoice.Columns().Name), "%"+req.Name+"%")
		}
		total, err = m.Count()
		liberr.ErrIsNil(ctx, err, "获取声音列表失败")
		orderBy := req.OrderBy
		if orderBy == "" {
			orderBy = "created_at desc"
		}
		err = m.Page(req.PageNum, req.PageSize).Order(orderBy).Scan(&voiceList)
		liberr.ErrIsNil(ctx, err, "获取声音列表失败")
	})
	return
}

func (s sCtsVoice) Add(ctx context.Context, req *cts.VoiceAddReq) (err error) {
	voice := (*entity.CtsVoice)(nil)
	err = g.Try(ctx, func(ctx context.Context) {
		m := dao.CtsVoice.Ctx(ctx)
		m = m.Where(fmt.Sprintf("%s=?", dao.CtsVoice.Columns().Name), req.Name)
		err := m.Limit(1).Scan(&voice)
		liberr.ErrIsNil(ctx, err, "获取声音列表失败")
		if voice != nil {
			err = fmt.Errorf("声音名称已存在")
			return
		}
		audioName := req.AudioName
		fileName := ""
		filePath := ""
		if audioName != "" {
			year, month, day := libUtils.GetDateInfo()
			fileName = guid.S() + ".wav"
			filePath = "/" + year + "/" + month + "/" + day + "/" + fileName
			wavaTempPath := common.AUDIO_TEMP_PATH + "/" + audioName + ".wav"
			saveDir := common.AUDIO_SAVE_PATH + "/" + year + "/" + month + "/" + day + "/"
			if !gfile.Exists(saveDir) {
				libUtils.MkDir(saveDir)
			}
			savePath := saveDir + "/" + fileName
			//复制wave_temp_path到save_path
			err := gfile.CopyFile(wavaTempPath, savePath)
			liberr.ErrIsNil(ctx, err, "音频复制失败")
		}

		// add
		_, err = dao.CtsVoice.Ctx(ctx).Insert(do.CtsVoice{
			Gender:    req.Gender,
			Name:      req.Name,
			Describe:  req.Describe,
			Tensor:    req.Tensor,
			VoicePath: filePath,
			Count:     0,
			Rate:      0,
		})
		fmt.Println(err)
		liberr.ErrIsNil(ctx, err, "新增声音失败")
	})
	return
}

func (s sCtsVoice) Edit(ctx context.Context, req *cts.VoiceEditReq) (err error) {
	g.Try(ctx, func(ctx context.Context) {
		_, err = s.GetById(ctx, req.VoiceId)
		liberr.ErrIsNil(ctx, err, "获取声音信息失败")
		// 查询Name是否存在
		voice := (*entity.CtsVoice)(nil)
		err = dao.CtsVoice.Ctx(ctx).Where(fmt.Sprintf("%s=? and %s!=?", dao.CtsVoice.Columns().Name, dao.CtsVoice.Columns().VoiceId), req.Name, req.VoiceId).Limit(1).Scan(&voice)
		if voice != nil {
			panic("声音名称已存在")
		}
		//edit
		_, err = dao.CtsVoice.Ctx(ctx).WherePri(req.VoiceId).Update(do.CtsVoice{
			Gender:   req.Gender,
			Name:     req.Name,
			Describe: req.Describe,
			Tensor:   req.Tensor,
		})
		liberr.ErrIsNil(ctx, err, "修改声音失败")
	})
	return
}

func (s sCtsVoice) Delete(ctx context.Context, id uint) (err error) {
	g.Try(ctx, func(ctx context.Context) {
		//删除音频文件
		_, realPath, _ := s.getVoicePath(ctx, id)
		_, err = dao.CtsVoice.Ctx(ctx).WherePri(id).Delete()
		liberr.ErrIsNil(ctx, err, "删除声音失败")
		if realPath != "" && gfile.Exists(realPath) {
			err = gfile.Remove(realPath)
			liberr.ErrIsNil(ctx, err, "删除声音音频文件失败")
		}
	})
	return
}

func (s sCtsVoice) BatchDelete(ctx context.Context, ids []uint) (err error) {
	g.Try(ctx, func(ctx context.Context) {
		_, err = dao.CtsVoice.Ctx(ctx).Where(dao.CtsVoice.Columns().VoiceId+" in(?)", ids).Delete()
		liberr.ErrIsNil(ctx, err, "批量删除声音失败")
	})
	return
}

func (s sCtsVoice) GetById(ctx context.Context, id uint) (res *model.CtsVoiceInfoRes, err error) {
	g.Try(ctx, func(ctx context.Context) {
		err = dao.CtsVoice.Ctx(ctx).Where(fmt.Sprintf("%s=?", dao.CtsVoice.Columns().VoiceId), id).Scan(&res)
		liberr.ErrIsNil(ctx, err, "获取声音信息失败")
	})
	return
}

func (s sCtsVoice) VoiceRate(ctx context.Context, req *cts.VoiceRateReq) (err error) {
	g.Try(ctx, func(ctx context.Context) {
		_, err = dao.CtsVoice.Ctx(ctx).Update(do.CtsVoice{
			Rate: req.Rate,
		}, fmt.Sprintf("%s=?", dao.CtsVoice.Columns().VoiceId), req.VoiceId)
		liberr.ErrIsNil(ctx, err, "修改声音评分失败")
	})
	return
}

func (s sCtsVoice) UpdateUsage(ctx context.Context, id uint) (err error) {
	g.Try(ctx, func(ctx context.Context) {
		voice, err := s.GetById(ctx, id)
		if err != nil {
			panic(err)
		}
		_, err = dao.CtsVoice.Ctx(ctx).Update(do.CtsVoice{
			Count:          voice.Count + 1,
			LastAccessTime: gtime.Now(),
		}, fmt.Sprintf("%s=?", dao.CtsVoice.Columns().VoiceId), id)
	})
	return
}

func (s sCtsVoice) getVoicePath(ctx context.Context, voiceId uint) (voicePath string, realVoicePath string, err error) {
	voicePath = ""
	realVoicePath = ""
	g.Try(ctx, func(ctx context.Context) {
		voice, err := s.GetById(ctx, voiceId)
		liberr.ErrIsNil(ctx, err, "获取声音信息失败")
		voicePath = voice.VoicePath
		realVoicePath = gfile.Join(common.AUDIO_SAVE_PATH, voicePath)
	})
	return voicePath, realVoicePath, err
}

func (s sCtsVoice) GenTestAudio(ctx context.Context, req *cts.AudioTestGenReq) (err error) {
	g.Try(ctx, func(ctx context.Context) {
		voiceId := req.VoiceId
		voice, getErr := s.GetById(ctx, voiceId)
		liberr.ErrIsNil(ctx, getErr, "获取声音信息失败")
		tensor := voice.Tensor
		if tensor == "" {
			err = gerror.New("tensor为空")
			return
		}
		// 查询默认生成文本测试文案
		defaultText := ""
		sysConfig := new(commonEntity.SysConfig)
		err = commonDao.SysConfig.Ctx(ctx).Where(fmt.Sprintf("%s=?", commonDao.SysConfig.Columns().ConfigKey), "chatts.default.text").Scan(&sysConfig)
		if err != nil {
			defaultText = sysConfig.ConfigValue
		}
		if defaultText == "" {
			defaultText = "曾经有一段真诚的爱情摆在我的面前,我没有珍惜。"
		}
		audioReq, reqErr := utils.BuildTestAudiotReq(defaultText, tensor)
		liberr.ErrIsNil(ctx, reqErr, "测试音频请求构建失败")
		r, resErr := utils.SendGenerateAudioReq(ctx, audioReq)
		defer r.Close()
		if r == nil || r.Response == nil || resErr != nil {
			err = gerror.New("测试音频生成失败")
			return
		}
		if r.Response.StatusCode != 200 || err != nil {
			err = gerror.New("测试音频生成失败")
			return
		}
		if j, jsonErr := gjson.DecodeToJson(r.ReadAllString()); jsonErr != nil {
			err = gerror.New("测试音频生成失败")
			return
		} else {
			audioName := j.Get("audio_name").String()
			tempPath := common.AUDIO_TEMP_PATH + "/" + audioName
			if !gfile.Exists(tempPath) {
				err = gerror.New("测试音频生成失败")
				return
			}
			//复制到保存目录
			fileName, saveErr := utils.CopyAduio(audioName)
			liberr.ErrIsNil(ctx, saveErr, "测试音频生成失败")
			_, err = dao.CtsVoice.Ctx(ctx).Update(do.CtsVoice{
				VoicePath: fileName,
			}, fmt.Sprintf("%s=?", dao.CtsVoice.Columns().VoiceId), voiceId)
			liberr.ErrIsNil(ctx, err, "测试音频生成失败")
		}
	})
	return
}

func (s sCtsVoice) DownloadList(ctx context.Context, req *cts.VoiceDownloadReq) (res []*model.CtsVoiceSimple, err error) {
	// 创建返回值
	g.Try(ctx, func(ctx context.Context) {
		res = make([]*model.CtsVoiceSimple, 0)
		voiceId := req.VoiceId
		if voiceId != 0 {
			//查询
			voice, voiceErr := s.GetById(ctx, voiceId)
			liberr.ErrIsNil(ctx, voiceErr, "获取声音信息失败")
			tensorList, convertErr := utils.ConverTensorStrToList(voice.Tensor)
			liberr.ErrIsNil(ctx, convertErr, "tensor转换失败")
			downloadRes := &model.CtsVoiceSimple{
				Name:     voice.Name,
				Gender:   voice.Gender,
				Describe: voice.Describe,
				Tensor:   tensorList,
			}
			res = append(res, downloadRes)
		} else {
			voiceList := make([]*entity.CtsVoice, 0)
			err = dao.CtsVoice.Ctx(ctx).Scan(&voiceList)
			liberr.ErrIsNil(ctx, err, "获取声音信息失败")
			for _, v := range voiceList {
				tensorList, convertErr := utils.ConverTensorStrToList(v.Tensor)
				liberr.ErrIsNil(ctx, convertErr, "tensor转换失败")
				downloadRes := &model.CtsVoiceSimple{
					Name:     v.Name,
					Gender:   v.Gender,
					Describe: v.Describe,
					Tensor:   tensorList,
				}
				res = append(res, downloadRes)
			}
		}
	})
	return
}

func (s sCtsVoice) UploadJson(ctx context.Context, req *cts.VoiceUploadReq) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return g.Try(ctx, func(ctx context.Context) {
			file, fileErr := req.File.Open()
			liberr.ErrIsNil(ctx, fileErr, "文件打开失败")
			contentByte, readErr := io.ReadAll(file)
			liberr.ErrIsNil(ctx, readErr, "文件读取失败")
			voiceList := make([]*model.CtsVoiceSimple, 0)
			j, decodeErr := gjson.DecodeToJson(string(contentByte))
			liberr.ErrIsNil(ctx, decodeErr, "json解析失败")
			err = j.Scan(&voiceList)
			liberr.ErrIsNil(ctx, err, "json解析失败")
			//err = json.Unmarshal(contentByte, &voiceList)
			//liberr.ErrIsNil(ctx, err, "json解析失败")
			nameList := new([]string)
			for _, v := range voiceList {
				voiceName := v.Name
				if gstr.InArray(*nameList, voiceName) {
					err = gerror.New("声音名称重复")
					return
				}
				//添加到nameList
				*nameList = append(*nameList, voiceName)
			}
			// 查询nameList在数据库中是否存在
			existVoiceList := make([]*entity.CtsVoice, 0)
			err = dao.CtsVoice.Ctx(ctx).Where(fmt.Sprintf("%s in (?)", dao.CtsVoice.Columns().Name), *nameList).Scan(&existVoiceList)
			liberr.ErrIsNil(ctx, err, "查询声音信息失败")
			if len(existVoiceList) > 0 {
				existNameList := make([]string, 0)
				for _, v := range existVoiceList {
					existNameList = append(existNameList, v.Name)
				}
				err = gerror.New(fmt.Sprintf("声音名称重复,已存在,请修改后再试: %s", gstr.Join(existNameList, ",")))
				return
			}
			//添加到数据库
			saveList := make([]*do.CtsVoice, 0)
			for _, v := range voiceList {
				saveList = append(saveList, &do.CtsVoice{
					Gender:    v.Gender,
					Name:      v.Name,
					Describe:  v.Describe,
					Tensor:    v.Tensor,
					VoicePath: "",
					Count:     0,
					Rate:      0,
				})
			}
			//批量插入
			_, err = dao.CtsVoice.Ctx(ctx).Insert(saveList)
		})
	})
	return
}
