package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/qkja/go-zero/tools/goctl/api/spec"
	"github.com/qkja/go-zero/tools/goctl/config"
	"github.com/qkja/go-zero/tools/goctl/internal/version"
	"github.com/qkja/go-zero/tools/goctl/util/pathx"
	"github.com/qkja/go-zero/tools/goctl/vars"
)

//go:embed main.tpl
var mainTemplate string

func genMain(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	name := strings.ToLower(api.Service.Name)
	// 服务名用于配置/日志标识；main 文件名固定为 main.go，避免随服务名变化
	configName := strings.TrimSuffix(name, "-api")

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        "main.go",
		templateName:    "mainTemplate",
		category:        category,
		templateFile:    mainTemplateFile,
		builtinTemplate: mainTemplate,
		data: map[string]string{
			"importPackages": genMainImports(rootPkg),
			"serviceName":    configName,
			"projectPkg":     projectPkg,
			"version":        version.BuildVersion,
		},
	})
}

func genMainImports(parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, handlerDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, contextDir)))
	imports = append(imports, fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}
