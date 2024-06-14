package router

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/tiger1103/gfast/v3/internal/app/cts/common"
	"github.com/tiger1103/gfast/v3/internal/app/cts/controller"
	_ "github.com/tiger1103/gfast/v3/internal/app/cts/logic"
	"github.com/tiger1103/gfast/v3/internal/app/system/service"
	"github.com/tiger1103/gfast/v3/library/libRouter"
	"log"
)

var R = new(Router)

type Router struct{}

func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/cts", func(group *ghttp.RouterGroup) {

		//登录验证拦截
		service.GfToken().Middleware(group)
		//context拦截器
		group.Middleware(service.Middleware().Ctx, service.Middleware().Auth)
		//后台操作日志记录
		group.Hook("/*", ghttp.HookAfterOutput, service.OperateLog().OperationLog)
		group.Bind(
			controller.TTS,
			controller.Voice,
		)
		//自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			panic(err)
		}
		// 读取配置文件
		apiUrl, _ := g.Cfg().Get(ctx, "chatTTS.apiUrl")
		audioTempPath, _ := g.Cfg().Get(ctx, "chatTTS.audioTempPath")
		audioSavePath, _ := g.Cfg().Get(ctx, "chatTTS.audioSavePath")
		exists := gfile.Exists(audioTempPath.String())
		if !exists {
			panic("audioTempPath路径不存在")
		}
		exists = gfile.Exists(audioSavePath.String())
		if !exists {
			err := gfile.Mkdir(audioSavePath.String())
			if err != nil {
				log.Println(err)
			}
		}
		common.API_URL = apiUrl.String()
		common.AUDIO_TEMP_PATH = audioTempPath.String()
		common.AUDIO_SAVE_PATH = audioSavePath.String()

	})
}
