// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"

	"server/internal/config"
	"server/internal/handler"
	"server/internal/pkg/xerr"
	"server/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/todo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册自定义错误处理
	httpx.SetErrorHandlerCtx(xerr.ErrorResponse)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
