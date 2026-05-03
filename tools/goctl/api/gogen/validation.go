package gogen

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// ValidationConfig 验证规则配置
type ValidationConfig struct {
	Patterns map[string]string            `yaml:"patterns"`
	Types    map[string]map[string]string `yaml:"types"`
}

var validationRules *ValidationConfig

// LoadValidationConfig 加载验证配置文件
// 格式: etc/validate.yaml
// types:
//   CreateUserReq:
//     Phone: "required,regexp=^1[3-9]\\d{9}$"
//     Email: "regexp=^[a-zA-Z0-9._%+-]+@..."
//   RegisterTenantReq:
//     TenantName: "required,min=1,max=128"
func LoadValidationConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	config := &ValidationConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("parse validation config failed: %w", err)
	}
	validationRules = config

	// 把 patterns 中定义的变量替换到 types 中
	// $phone → 替换为正则
	if config.Patterns != nil {
		for typeName, fields := range config.Types {
			for fieldName, rule := range fields {
				for key, pattern := range config.Patterns {
					placeholder := "$" + key
					if strings.Contains(rule, placeholder) {
						config.Types[typeName][fieldName] = strings.ReplaceAll(rule, placeholder, pattern)
					}
				}
			}
		}
	}
	return nil
}

// resolveValidateTag 根据类型名和字段名获取 validate 规则
// 如果 typeName 为空或没有匹配规则，返回空字符串
func resolveValidateTag(typeName, fieldName string) string {
	if validationRules == nil || typeName == "" {
		return ""
	}
	// 先用类型精确匹配
	if fields, ok := validationRules.Types[typeName]; ok {
		if rule, ok := fields[fieldName]; ok {
			return rule
		}
	}
	return ""
}

// injectValidateTag 把 validate 规则注入到已有的 tag 字符串中
// 如果 tag 已有 validate 则追加，否则新增 validate tag
func injectValidateTag(tag, validateRule string) string {
	// 去掉反引号
	raw := tag
	raw = strings.TrimPrefix(raw, "`")
	raw = strings.TrimSuffix(raw, "`")

	// 解析现有 tag
	pairs := strings.Fields(raw)
	hasValidate := false
	newPairs := make([]string, 0, len(pairs)+1)

	for _, p := range pairs {
		if strings.HasPrefix(p, "validate:") {
			hasValidate = true
			// 已有 validate, 追加规则
			existing := strings.TrimPrefix(p, "validate:\"")
			existing = strings.TrimSuffix(existing, "\"")
			newPairs = append(newPairs, fmt.Sprintf("validate:\"%s,%s\"", existing, validateRule))
		} else {
			newPairs = append(newPairs, p)
		}
	}

	if !hasValidate {
		newPairs = append(newPairs, fmt.Sprintf("validate:\"%s\"", validateRule))
	}

	return "`" + strings.Join(newPairs, " ") + "`"
}

// loadValidateConfig 从项目目录加载 validate.yaml
func loadValidateConfig(projectDir string) {
	path := projectDir + "/etc/validate.yaml"
	if err := LoadValidationConfig(path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: load validation config failed: %v\n", err)
	}
}

// ResetValidationRules 重置规则（用于测试）
func ResetValidationRules() {
	validationRules = nil
}
