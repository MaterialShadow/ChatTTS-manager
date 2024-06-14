package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	_ "github.com/tiger1103/gfast/v3/internal/app/boot"

	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/tiger1103/gfast/v3/internal/app/system/packed"
	"github.com/tiger1103/gfast/v3/internal/cmd"
)

// https://github.com/gogf/gf/blob/master/contrib/drivers/README.zh_CN.MD
func main() {
	cmd.Main.Run(gctx.New())
}
