// Code scaffolded by goctl. Safe to edit.
// goctl {{.version}}

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	gobasecfg "github.com/qkja/gobase/config"

	{{.importPackages}}
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "{{.serviceName}}",
}

var appCmd = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置并启动热加载（rest.RestConf 的 default tag 由 go-zero NewServer 补齐）
		if err := gobasecfg.Init(&config.AppCfg, cfgFile, nil); err != nil {
			panic(fmt.Sprintf("failed to load config: %v", err))
		}

		// 构建 RestConf（服务名从 application.name 带入）并启动
		rc := config.AppCfg.Get().Rest
		rc.ServiceConf.Name = config.AppCfg.Get().Application.Name
		server := rest.MustNewServer(rc)
		defer server.Stop()

		ctx := svc.NewServiceContext(config.AppCfg.Get())
		handler.RegisterHandlers(server, ctx)

		fmt.Printf("Starting server at %s:%d...\n",
			config.AppCfg.Get().Rest.Host, config.AppCfg.Get().Rest.Port)
		server.Start()
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"配置文件目录（默认 ./config/application.toml）")
	rootCmd.AddCommand(appCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
