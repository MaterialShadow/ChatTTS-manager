package service

import (
	"context"
	"github.com/tiger1103/gfast/v3/api/v1/cts"
	"github.com/tiger1103/gfast/v3/internal/app/cts/model"
	"github.com/tiger1103/gfast/v3/internal/app/cts/model/entity"
)

type ICtsVoice interface {
	List(ctx context.Context, req *cts.VoiceSearchReq) (total interface{}, res []*entity.CtsVoice, err error)
	Add(ctx context.Context, req *cts.VoiceAddReq) (err error)
	Edit(ctx context.Context, req *cts.VoiceEditReq) (err error)
	Delete(ctx context.Context, id uint) (err error)
	BatchDelete(ctx context.Context, ids []uint) (err error)
	GetById(ctx context.Context, id uint) (res *model.CtsVoiceInfoRes, err error)
	VoiceRate(ctx context.Context, req *cts.VoiceRateReq) (err error)
	UpdateUsage(ctx context.Context, id uint) (err error)
	GenTestAudio(ctx context.Context, req *cts.AudioTestGenReq) (err error)
	DownloadList(ctx context.Context, req *cts.VoiceDownloadReq) (res []*model.CtsVoiceSimple, err error)
	UploadJson(ctx context.Context, req *cts.VoiceUploadReq) (err error)
}

var localCtsVoice ICtsVoice

func CtsVoice() ICtsVoice {
	if localCtsVoice == nil {
		panic("implement not found for interface ICtsVoice, forgot register?")
	}
	return localCtsVoice
}

func RegisterCtsVoice(i ICtsVoice) {
	localCtsVoice = i
}
