package gogen

import (
	_ "embed"
	"strconv"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
)

const defaultPort = 8888

//go:embed application.toml.tpl
var etcTemplate string

// genEtc 生成 gobase 配置目录下的 application.toml（config/application.toml）。
func genEtc(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	service := api.Service
	host := "0.0.0.0"
	port := strconv.Itoa(defaultPort)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "config",
		filename:        "application.toml",
		templateName:    "etcTemplate",
		category:        category,
		templateFile:    etcTemplateFile,
		builtinTemplate: etcTemplate,
		data: map[string]string{
			"serviceName": service.Name,
			"host":        host,
			"port":        port,
		},
	})
}
