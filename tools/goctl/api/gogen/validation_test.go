package gogen

import (
	"os"
	"testing"
)

func TestInjectValidateTag_NewTag(t *testing.T) {
	result := injectValidateTag("`json:\"username\"`", "required,min=2,max=32")
	expected := "`json:\"username\" validate:\"required,min=2,max=32\"`"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestInjectValidateTag_ExistingValidate(t *testing.T) {
	result := injectValidateTag("`json:\"phone\" validate:\"optional\"`", "regexp=^1[3-9]\\d{9}$")
	expected := "`json:\"phone\" validate:\"optional,regexp=^1[3-9]\\d{9}$\"`"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestInjectValidateTag_NoBackticks(t *testing.T) {
	result := injectValidateTag("", "required")
	expected := "`validate:\"required\"`"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestLoadValidationConfig_FileNotFound(t *testing.T) {
	err := LoadValidationConfig("/path/to/nonexistent.yaml")
	if err != nil {
		t.Errorf("expected nil for missing file, got %v", err)
	}
}

func TestLoadValidationConfig_ValidFile(t *testing.T) {
	content := []byte(`
patterns:
  phone: "^1[3-9]\\d{9}$"
types:
  CreateUserReq:
    Phone: "regexp=$phone"
`)
	tmpFile := "/tmp/test-validate-config.yaml"
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	if err := LoadValidationConfig(tmpFile); err != nil {
		t.Fatalf("LoadValidationConfig failed: %v", err)
	}

	rule := resolveValidateTag("CreateUserReq", "Phone")
	expected := "regexp=^1[3-9]\\d{9}$"
	if rule != expected {
		t.Errorf("got %q, want %q", rule, expected)
	}

	// pattern 替换
	if validationRules.Patterns["phone"] != "^1[3-9]\\d{9}$" {
		t.Errorf("pattern not preserved correctly")
	}
}

func TestResolveValidateTag_NoMatch(t *testing.T) {
	ResetValidationRules()
	rule := resolveValidateTag("NonexistentType", "Field")
	if rule != "" {
		t.Errorf("expected empty, got %q", rule)
	}
}
