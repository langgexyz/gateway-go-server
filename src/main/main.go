package main

import (
	"context"
	"fmt"
	"gateway-go-server-main/src/api"
	"gateway-go-server-main/src/proxy"
	"os"
	"path/filepath"

	"github.com/xpwu/go-cmd/arg"
	"github.com/xpwu/go-cmd/cmd"
	"github.com/xpwu/go-cmd/exe"
	"github.com/xpwu/go-config/configs"
	"github.com/xpwu/go-log/log"
	"github.com/xpwu/go-stream/push"
	"github.com/xpwu/go-stream/websocket"
	"github.com/xpwu/go-tinyserver/http"
)

func main() {
	cmd.RegisterCmd(cmd.DefaultCmdName, "start server", func(args *arg.Arg) {
		arg.ReadConfig(args)
		args.Parse()

		_, logger := log.WithCtx(context.Background())
		currentDirectory, _ := os.Getwd()
		logger.Debug(fmt.Sprintf("Current Directory: %s", currentDirectory))

		api.AddAPI()
		proxy.AddAPI()
		http.Start()
		websocket.Start()
		push.Start()

		// block
		block := make(chan struct{})
		<-block
	})

	argR := "config.json.default"
	cmd.RegisterCmd("print", "print config with json", func(args *arg.Arg) {
		args.String(&argR, "c", "the file name of config file")
		args.Parse()
		if !filepath.IsAbs(argR) {
			argR = filepath.Join(exe.Exe.AbsDir, argR)
		}
		configs.SetConfigurator(&configs.JsonConfig{PrintFile: argR})
		err := configs.Print()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	})

	cmd.Run()
}
