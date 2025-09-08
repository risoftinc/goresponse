package goresponse

import (
	"encoding/json"
	"os"
	"testing"
)

// TestConfigSource tests ConfigSource struct
func TestConfigSource(t *testing.T) {
	tests := []struct {
		name     string
		source   ConfigSource
		expected ConfigSource
	}{
		{
			name: "File source",
			source: ConfigSource{
				Method: "file",
				Path:   "/path/to/config.json",
			},
			expected: ConfigSource{
				Method: "file",
				Path:   "/path/to/config.json",
			},
		},
		{
			name: "URL source",
			source: ConfigSource{
				Method: "url",
				Path:   "https://example.com/config.json",
			},
			expected: ConfigSource{
				Method: "url",
				Path:   "https://example.com/config.json",
			},
		},
		{
			name: "Empty source",
			source: ConfigSource{
				Method: "",
				Path:   "",
			},
			expected: ConfigSource{
				Method: "",
				Path:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.source.Method != tt.expected.Method {
				t.Errorf("Expected method %s, got %s", tt.expected.Method, tt.source.Method)
			}
			if tt.source.Path != tt.expected.Path {
				t.Errorf("Expected path %s, got %s", tt.expected.Path, tt.source.Path)
			}
		})
	}
}

// TestConfigSourceJSON tests ConfigSource JSON marshaling/unmarshaling
func TestConfigSourceJSON(t *testing.T) {
	tests := []struct {
		name    string
		source  ConfigSource
		jsonStr string
	}{
		{
			name: "File source JSON",
			source: ConfigSource{
				Method: "file",
				Path:   "/path/to/config.json",
			},
			jsonStr: `{"method":"file","path":"/path/to/config.json"}`,
		},
		{
			name: "URL source JSON",
			source: ConfigSource{
				Method: "url",
				Path:   "https://example.com/config.json",
			},
			jsonStr: `{"method":"url","path":"https://example.com/config.json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.source)
			if err != nil {
				t.Errorf("Failed to marshal ConfigSource: %v", err)
				return
			}

			if string(jsonData) != tt.jsonStr {
				t.Errorf("Expected JSON %s, got %s", tt.jsonStr, string(jsonData))
			}

			// Test unmarshaling
			var unmarshaled ConfigSource
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal ConfigSource: %v", err)
				return
			}

			if unmarshaled.Method != tt.source.Method {
				t.Errorf("Expected method %s, got %s", tt.source.Method, unmarshaled.Method)
			}
			if unmarshaled.Path != tt.source.Path {
				t.Errorf("Expected path %s, got %s", tt.source.Path, unmarshaled.Path)
			}
		})
	}
}

// TestTranslationSource tests TranslationSource struct
func TestTranslationSource(t *testing.T) {
	tests := []struct {
		name     string
		source   TranslationSource
		expected TranslationSource
	}{
		{
			name: "File translation source",
			source: TranslationSource{
				Method: "file",
				Path:   "/path/to/translations.json",
			},
			expected: TranslationSource{
				Method: "file",
				Path:   "/path/to/translations.json",
			},
		},
		{
			name: "URL translation source",
			source: TranslationSource{
				Method: "url",
				Path:   "https://example.com/translations.json",
			},
			expected: TranslationSource{
				Method: "url",
				Path:   "https://example.com/translations.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.source.Method != tt.expected.Method {
				t.Errorf("Expected method %s, got %s", tt.expected.Method, tt.source.Method)
			}
			if tt.source.Path != tt.expected.Path {
				t.Errorf("Expected path %s, got %s", tt.expected.Path, tt.source.Path)
			}
		})
	}
}

// TestTranslationSourceJSON tests TranslationSource JSON marshaling/unmarshaling
func TestTranslationSourceJSON(t *testing.T) {
	tests := []struct {
		name    string
		source  TranslationSource
		jsonStr string
	}{
		{
			name: "File translation source JSON",
			source: TranslationSource{
				Method: "file",
				Path:   "/path/to/translations.json",
			},
			jsonStr: `{"method":"file","path":"/path/to/translations.json"}`,
		},
		{
			name: "URL translation source JSON",
			source: TranslationSource{
				Method: "url",
				Path:   "https://example.com/translations.json",
			},
			jsonStr: `{"method":"url","path":"https://example.com/translations.json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.source)
			if err != nil {
				t.Errorf("Failed to marshal TranslationSource: %v", err)
				return
			}

			if string(jsonData) != tt.jsonStr {
				t.Errorf("Expected JSON %s, got %s", tt.jsonStr, string(jsonData))
			}

			// Test unmarshaling
			var unmarshaled TranslationSource
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal TranslationSource: %v", err)
				return
			}

			if unmarshaled.Method != tt.source.Method {
				t.Errorf("Expected method %s, got %s", tt.source.Method, unmarshaled.Method)
			}
			if unmarshaled.Path != tt.source.Path {
				t.Errorf("Expected path %s, got %s", tt.source.Path, unmarshaled.Path)
			}
		})
	}
}

// TestResponseConfig tests ResponseConfig struct
func TestResponseConfig(t *testing.T) {
	config := ResponseConfig{
		MessageTemplates: map[string]MessageTemplate{
			"welcome": {
				Key:      "welcome",
				Template: "Welcome $name",
				CodeMappings: map[string]int{
					"http": 200,
					"grpc": 0,
				},
				Translations: map[string]string{
					"en": "Welcome $name",
					"id": "Selamat datang $name",
				},
			},
		},
		ManualMessageTemplates: map[string]MessageTemplate{
			"manual": {
				Key:      "manual",
				Template: "Manual template",
				CodeMappings: map[string]int{
					"http": 201,
				},
			},
		},
		DefaultLanguage: "en",
		Languages:       []string{"en", "id", "es"},
		Translations: map[string]map[string]string{
			"en": {"hello": "Hello", "goodbye": "Goodbye"},
			"id": {"hello": "Halo", "goodbye": "Selamat tinggal"},
			"es": {"hello": "Hola"},
		},
		TranslationSources: map[string]TranslationSource{
			"en": {
				Method: "file",
				Path:   "en_translations.json",
			},
			"id": {
				Method: "url",
				Path:   "https://example.com/id_translations.json",
			},
		},
	}

	// Test MessageTemplates
	if len(config.MessageTemplates) != 1 {
		t.Errorf("Expected 1 message template, got %d", len(config.MessageTemplates))
	}
	if config.MessageTemplates["welcome"].Key != "welcome" {
		t.Errorf("Expected welcome template key, got %s", config.MessageTemplates["welcome"].Key)
	}

	// Test ManualMessageTemplates
	if len(config.ManualMessageTemplates) != 1 {
		t.Errorf("Expected 1 manual message template, got %d", len(config.ManualMessageTemplates))
	}
	if config.ManualMessageTemplates["manual"].Key != "manual" {
		t.Errorf("Expected manual template key, got %s", config.ManualMessageTemplates["manual"].Key)
	}

	// Test DefaultLanguage
	if config.DefaultLanguage != "en" {
		t.Errorf("Expected default language 'en', got %s", config.DefaultLanguage)
	}

	// Test Languages
	if len(config.Languages) != 3 {
		t.Errorf("Expected 3 languages, got %d", len(config.Languages))
	}

	// Test Translations
	if len(config.Translations) != 3 {
		t.Errorf("Expected 3 translation languages, got %d", len(config.Translations))
	}

	// Test TranslationSources
	if len(config.TranslationSources) != 2 {
		t.Errorf("Expected 2 translation sources, got %d", len(config.TranslationSources))
	}
}

// TestResponseConfigJSON tests ResponseConfig JSON marshaling/unmarshaling
func TestResponseConfigJSON(t *testing.T) {
	config := ResponseConfig{
		MessageTemplates: map[string]MessageTemplate{
			"welcome": {
				Key:      "welcome",
				Template: "Welcome $name",
				CodeMappings: map[string]int{
					"http": 200,
				},
				Translations: map[string]string{
					"en": "Welcome $name",
					"id": "Selamat datang $name",
				},
			},
		},
		DefaultLanguage: "en",
		Languages:       []string{"en", "id"},
		Translations: map[string]map[string]string{
			"en": {"hello": "Hello"},
			"id": {"hello": "Halo"},
		},
		TranslationSources: map[string]TranslationSource{
			"en": {
				Method: "file",
				Path:   "en_translations.json",
			},
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Errorf("Failed to marshal ResponseConfig: %v", err)
		return
	}

	// Test unmarshaling
	var unmarshaled ResponseConfig
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ResponseConfig: %v", err)
		return
	}

	// Verify unmarshaled data
	if unmarshaled.DefaultLanguage != config.DefaultLanguage {
		t.Errorf("Expected default language %s, got %s", config.DefaultLanguage, unmarshaled.DefaultLanguage)
	}

	if len(unmarshaled.Languages) != len(config.Languages) {
		t.Errorf("Expected %d languages, got %d", len(config.Languages), len(unmarshaled.Languages))
	}

	if len(unmarshaled.MessageTemplates) != len(config.MessageTemplates) {
		t.Errorf("Expected %d message templates, got %d", len(config.MessageTemplates), len(unmarshaled.MessageTemplates))
	}

	if len(unmarshaled.Translations) != len(config.Translations) {
		t.Errorf("Expected %d translation languages, got %d", len(config.Translations), len(unmarshaled.Translations))
	}

	if len(unmarshaled.TranslationSources) != len(config.TranslationSources) {
		t.Errorf("Expected %d translation sources, got %d", len(config.TranslationSources), len(unmarshaled.TranslationSources))
	}
}

// TestMessageTemplate tests MessageTemplate struct
func TestMessageTemplate(t *testing.T) {
	template := MessageTemplate{
		Key:      "welcome",
		Template: "Welcome $name to $place",
		CodeMappings: map[string]int{
			"http": 200,
			"grpc": 0,
			"tcp":  1,
		},
		Translations: map[string]string{
			"en": "Welcome $name to $place",
			"id": "Selamat datang $name di $place",
			"es": "Bienvenido $name a $place",
		},
	}

	// Test Key
	if template.Key != "welcome" {
		t.Errorf("Expected key 'welcome', got %s", template.Key)
	}

	// Test Template
	if template.Template != "Welcome $name to $place" {
		t.Errorf("Expected template 'Welcome $name to $place', got %s", template.Template)
	}

	// Test CodeMappings
	if len(template.CodeMappings) != 3 {
		t.Errorf("Expected 3 code mappings, got %d", len(template.CodeMappings))
	}
	if template.CodeMappings["http"] != 200 {
		t.Errorf("Expected HTTP code 200, got %d", template.CodeMappings["http"])
	}

	// Test Translations
	if len(template.Translations) != 3 {
		t.Errorf("Expected 3 translations, got %d", len(template.Translations))
	}
	if template.Translations["en"] != "Welcome $name to $place" {
		t.Errorf("Expected English translation, got %s", template.Translations["en"])
	}
}

// TestMessageTemplateJSON tests MessageTemplate JSON marshaling/unmarshaling
func TestMessageTemplateJSON(t *testing.T) {
	template := MessageTemplate{
		Key:      "welcome",
		Template: "Welcome $name",
		CodeMappings: map[string]int{
			"http": 200,
		},
		Translations: map[string]string{
			"en": "Welcome $name",
			"id": "Selamat datang $name",
		},
	}

	// Test marshaling
	jsonData, err := json.Marshal(template)
	if err != nil {
		t.Errorf("Failed to marshal MessageTemplate: %v", err)
		return
	}

	// Test unmarshaling
	var unmarshaled MessageTemplate
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal MessageTemplate: %v", err)
		return
	}

	// Verify unmarshaled data
	if unmarshaled.Key != template.Key {
		t.Errorf("Expected key %s, got %s", template.Key, unmarshaled.Key)
	}
	if unmarshaled.Template != template.Template {
		t.Errorf("Expected template %s, got %s", template.Template, unmarshaled.Template)
	}
	if len(unmarshaled.CodeMappings) != len(template.CodeMappings) {
		t.Errorf("Expected %d code mappings, got %d", len(template.CodeMappings), len(unmarshaled.CodeMappings))
	}
	if len(unmarshaled.Translations) != len(template.Translations) {
		t.Errorf("Expected %d translations, got %d", len(template.Translations), len(unmarshaled.Translations))
	}
}

// TestConfigChangeCallback tests ConfigChangeCallback function type
func TestConfigChangeCallback(t *testing.T) {
	var callback ConfigChangeCallback
	var oldConfig, newConfig *ResponseConfig

	// Test callback assignment
	callback = func(old, new *ResponseConfig) {
		oldConfig = old
		newConfig = new
	}

	// Test callback execution
	old := &ResponseConfig{
		DefaultLanguage: "en",
		Languages:       []string{"en"},
	}
	new := &ResponseConfig{
		DefaultLanguage: "id",
		Languages:       []string{"en", "id"},
	}

	callback(old, new)

	if oldConfig != old {
		t.Error("Expected old config to be set")
	}
	if newConfig != new {
		t.Error("Expected new config to be set")
	}
}

// TestMessageTemplateBuilder tests MessageTemplateBuilder
func TestMessageTemplateBuilder(t *testing.T) {
	t.Run("NewMessageTemplateBuilder", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test_key")
		if builder == nil {
			t.Fatal("Expected builder to be created")
		}
		if builder.template == nil {
			t.Fatal("Expected template to be initialized")
		}
		if builder.template.Key != "test_key" {
			t.Errorf("Expected key 'test_key', got %s", builder.template.Key)
		}
		if builder.template.Template != "" {
			t.Errorf("Expected empty template, got %s", builder.template.Template)
		}
		if builder.template.CodeMappings == nil {
			t.Error("Expected CodeMappings to be initialized")
		}
		if len(builder.template.CodeMappings) != 0 {
			t.Errorf("Expected empty CodeMappings, got %d", len(builder.template.CodeMappings))
		}
	})

	t.Run("WithTemplate", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		result := builder.WithTemplate("Hello $name")

		if result != builder {
			t.Error("Expected method chaining to return same builder")
		}
		if builder.template.Template != "Hello $name" {
			t.Errorf("Expected template 'Hello $name', got %s", builder.template.Template)
		}
	})

	t.Run("WithTranslation", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		result := builder.WithTranslation("en", "Hello $name")

		if result != builder {
			t.Error("Expected method chaining to return same builder")
		}
		if builder.template.Translations == nil {
			t.Error("Expected Translations to be initialized")
		}
		if builder.template.Translations["en"] != "Hello $name" {
			t.Errorf("Expected English translation, got %s", builder.template.Translations["en"])
		}

		// Test multiple translations
		builder.WithTranslation("id", "Halo $name")
		if builder.template.Translations["id"] != "Halo $name" {
			t.Errorf("Expected Indonesian translation, got %s", builder.template.Translations["id"])
		}
	})

	t.Run("WithTranslations", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		translations := map[string]string{
			"en": "Hello $name",
			"id": "Halo $name",
			"es": "Hola $name",
		}

		result := builder.WithTranslations(translations)

		if result != builder {
			t.Error("Expected method chaining to return same builder")
		}
		if len(builder.template.Translations) != 3 {
			t.Errorf("Expected 3 translations, got %d", len(builder.template.Translations))
		}
		for lang, translation := range translations {
			if builder.template.Translations[lang] != translation {
				t.Errorf("Expected %s translation %s, got %s", lang, translation, builder.template.Translations[lang])
			}
		}
	})

	t.Run("WithCodeMapping", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		result := builder.WithCodeMapping("http", 200)

		if result != builder {
			t.Error("Expected method chaining to return same builder")
		}
		if builder.template.CodeMappings["http"] != 200 {
			t.Errorf("Expected HTTP code 200, got %d", builder.template.CodeMappings["http"])
		}

		// Test multiple code mappings
		builder.WithCodeMapping("grpc", 0)
		if builder.template.CodeMappings["grpc"] != 0 {
			t.Errorf("Expected gRPC code 0, got %d", builder.template.CodeMappings["grpc"])
		}
	})

	t.Run("WithCodeMappings", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		codeMappings := map[string]int{
			"http": 200,
			"grpc": 0,
			"tcp":  1,
		}

		result := builder.WithCodeMappings(codeMappings)

		if result != builder {
			t.Error("Expected method chaining to return same builder")
		}
		if len(builder.template.CodeMappings) != 3 {
			t.Errorf("Expected 3 code mappings, got %d", len(builder.template.CodeMappings))
		}
		for mappingType, code := range codeMappings {
			if builder.template.CodeMappings[mappingType] != code {
				t.Errorf("Expected %s code %d, got %d", mappingType, code, builder.template.CodeMappings[mappingType])
			}
		}
	})

	t.Run("Build", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		builder.WithTemplate("Hello $name")
		builder.WithTranslation("en", "Hello $name")
		builder.WithCodeMapping("http", 200)

		template := builder.Build()

		if template == nil {
			t.Fatal("Expected template to be returned")
		}
		if template.Key != "test" {
			t.Errorf("Expected key 'test', got %s", template.Key)
		}
		if template.Template != "Hello $name" {
			t.Errorf("Expected template 'Hello $name', got %s", template.Template)
		}
		if template.CodeMappings["http"] != 200 {
			t.Errorf("Expected HTTP code 200, got %d", template.CodeMappings["http"])
		}
		if template.Translations["en"] != "Hello $name" {
			t.Errorf("Expected English translation, got %s", template.Translations["en"])
		}
	})

	t.Run("Method chaining", func(t *testing.T) {
		template := NewMessageTemplateBuilder("welcome").
			WithTemplate("Welcome $name to $place").
			WithTranslation("en", "Welcome $name to $place").
			WithTranslation("id", "Selamat datang $name di $place").
			WithCodeMapping("http", 200).
			WithCodeMapping("grpc", 0).
			Build()

		if template.Key != "welcome" {
			t.Errorf("Expected key 'welcome', got %s", template.Key)
		}
		if template.Template != "Welcome $name to $place" {
			t.Errorf("Expected template, got %s", template.Template)
		}
		if len(template.Translations) != 2 {
			t.Errorf("Expected 2 translations, got %d", len(template.Translations))
		}
		if len(template.CodeMappings) != 2 {
			t.Errorf("Expected 2 code mappings, got %d", len(template.CodeMappings))
		}
	})
}

// TestConfigPrinter tests ConfigPrinter
func TestConfigPrinter(t *testing.T) {
	config := &ResponseConfig{
		MessageTemplates: map[string]MessageTemplate{
			"welcome": {
				Key:      "welcome",
				Template: "Welcome $name",
				CodeMappings: map[string]int{
					"http": 200,
				},
			},
		},
		ManualMessageTemplates: map[string]MessageTemplate{
			"manual": {
				Key:      "manual",
				Template: "Manual template",
				CodeMappings: map[string]int{
					"http": 201,
				},
			},
		},
		DefaultLanguage: "en",
		Languages:       []string{"en", "id"},
		Translations: map[string]map[string]string{
			"en": {"hello": "Hello"},
			"id": {"hello": "Halo"},
		},
		TranslationSources: map[string]TranslationSource{
			"en": {
				Method: "file",
				Path:   "en_translations.json",
			},
		},
	}

	t.Run("NewConfigPrinter", func(t *testing.T) {
		printer := NewConfigPrinter(config)
		if printer == nil {
			t.Fatal("Expected printer to be created")
		}
		if printer.config != config {
			t.Error("Expected config to be set")
		}
		if !printer.indent {
			t.Error("Expected indent to be true by default")
		}
	})

	t.Run("WithIndent", func(t *testing.T) {
		printer := NewConfigPrinter(config)
		result := printer.WithIndent(false)

		if result != printer {
			t.Error("Expected method chaining to return same printer")
		}
		if printer.indent {
			t.Error("Expected indent to be false")
		}

		// Test setting back to true
		printer.WithIndent(true)
		if !printer.indent {
			t.Error("Expected indent to be true")
		}
	})

	t.Run("Export", func(t *testing.T) {
		printer := NewConfigPrinter(config)

		// Test with indent
		jsonStr, err := printer.Export()
		if err != nil {
			t.Errorf("Export failed: %v", err)
		}
		if jsonStr == "" {
			t.Error("Expected non-empty JSON string")
		}

		// Verify it's valid JSON
		var exportedConfig map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &exportedConfig)
		if err != nil {
			t.Errorf("Exported JSON is invalid: %v", err)
		}

		// Test without indent
		printer.WithIndent(false)
		jsonStrNoIndent, err := printer.Export()
		if err != nil {
			t.Errorf("Export without indent failed: %v", err)
		}
		if jsonStrNoIndent == "" {
			t.Error("Expected non-empty JSON string")
		}

		// Indented JSON should be longer than non-indented
		if len(jsonStr) <= len(jsonStrNoIndent) {
			t.Error("Expected indented JSON to be longer than non-indented")
		}
	})

	t.Run("ExportToFile", func(t *testing.T) {
		printer := NewConfigPrinter(config)
		testFile := "test_config_export.json"
		defer os.Remove(testFile)

		err := printer.ExportToFile(testFile)
		if err != nil {
			t.Errorf("ExportToFile failed: %v", err)
		}

		// Verify file was created and contains valid JSON
		data, err := os.ReadFile(testFile)
		if err != nil {
			t.Errorf("Failed to read exported file: %v", err)
		}

		var exportedConfig map[string]interface{}
		err = json.Unmarshal(data, &exportedConfig)
		if err != nil {
			t.Errorf("Exported file contains invalid JSON: %v", err)
		}
	})

	t.Run("Print", func(t *testing.T) {
		printer := NewConfigPrinter(config)

		// This test just ensures the method doesn't panic
		// In a real test, you might want to capture stdout
		err := printer.Print()
		if err != nil {
			t.Errorf("Print failed: %v", err)
		}
	})

	t.Run("toJSON with manual templates priority", func(t *testing.T) {
		printer := NewConfigPrinter(config)
		jsonData, err := printer.toJSON()
		if err != nil {
			t.Errorf("toJSON failed: %v", err)
		}

		var exportedConfig map[string]interface{}
		err = json.Unmarshal(jsonData, &exportedConfig)
		if err != nil {
			t.Errorf("toJSON result is invalid JSON: %v", err)
		}

		// Check that message_templates contains both async and manual templates
		messageTemplates, ok := exportedConfig["message_templates"].(map[string]interface{})
		if !ok {
			t.Error("Expected message_templates to be a map")
		}

		// Should contain both welcome (async) and manual (manual) templates
		if _, exists := messageTemplates["welcome"]; !exists {
			t.Error("Expected welcome template to be in exported config")
		}
		if _, exists := messageTemplates["manual"]; !exists {
			t.Error("Expected manual template to be in exported config")
		}
	})
}

// TestConfigPrinterWithNilConfig tests ConfigPrinter with nil config
func TestConfigPrinterWithNilConfig(t *testing.T) {
	printer := NewConfigPrinter(nil)

	// These should not panic
	_, err := printer.Export()
	if err == nil {
		t.Error("Expected error when exporting nil config")
	}

	err = printer.ExportToFile("test.json")
	if err == nil {
		t.Error("Expected error when exporting nil config to file")
	}

	err = printer.Print()
	if err == nil {
		t.Error("Expected error when printing nil config")
	}
}

// TestConfigPrinterWithEmptyConfig tests ConfigPrinter with empty config
func TestConfigPrinterWithEmptyConfig(t *testing.T) {
	config := &ResponseConfig{}
	printer := NewConfigPrinter(config)

	jsonStr, err := printer.Export()
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	// Should be valid JSON even with empty config
	var exportedConfig map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &exportedConfig)
	if err != nil {
		t.Errorf("Exported JSON is invalid: %v", err)
	}

	// Should contain default values
	if exportedConfig["default_language"] != "" {
		t.Errorf("Expected empty default language, got %v", exportedConfig["default_language"])
	}
}

// TestMessageTemplateBuilderEdgeCases tests edge cases for MessageTemplateBuilder
func TestMessageTemplateBuilderEdgeCases(t *testing.T) {
	t.Run("WithTranslation with nil Translations map", func(t *testing.T) {
		builder := &MessageTemplateBuilder{
			template: &MessageTemplate{
				Key:          "test",
				Translations: nil,
			},
		}

		builder.WithTranslation("en", "Hello")
		if builder.template.Translations == nil {
			t.Error("Expected Translations to be initialized")
		}
		if builder.template.Translations["en"] != "Hello" {
			t.Errorf("Expected English translation, got %s", builder.template.Translations["en"])
		}
	})

	t.Run("WithTranslations with nil Translations map", func(t *testing.T) {
		builder := &MessageTemplateBuilder{
			template: &MessageTemplate{
				Key:          "test",
				Translations: nil,
			},
		}

		translations := map[string]string{
			"en": "Hello",
			"id": "Halo",
		}
		builder.WithTranslations(translations)
		if builder.template.Translations == nil {
			t.Error("Expected Translations to be initialized")
		}
		if len(builder.template.Translations) != 2 {
			t.Errorf("Expected 2 translations, got %d", len(builder.template.Translations))
		}
	})

	t.Run("WithTranslations with empty map", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		builder.WithTranslations(map[string]string{})

		if builder.template.Translations == nil {
			t.Error("Expected Translations to be initialized")
		}
		if len(builder.template.Translations) != 0 {
			t.Errorf("Expected 0 translations, got %d", len(builder.template.Translations))
		}
	})

	t.Run("WithCodeMappings with empty map", func(t *testing.T) {
		builder := NewMessageTemplateBuilder("test")
		builder.WithCodeMappings(map[string]int{})

		if len(builder.template.CodeMappings) != 0 {
			t.Errorf("Expected 0 code mappings, got %d", len(builder.template.CodeMappings))
		}
	})
}

// BenchmarkMessageTemplateBuilder benchmarks MessageTemplateBuilder performance
func BenchmarkMessageTemplateBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewMessageTemplateBuilder("test").
			WithTemplate("Hello $name").
			WithTranslation("en", "Hello $name").
			WithTranslation("id", "Halo $name").
			WithCodeMapping("http", 200).
			WithCodeMapping("grpc", 0).
			Build()
	}
}

// BenchmarkConfigPrinterExport benchmarks ConfigPrinter Export performance
func BenchmarkConfigPrinterExport(b *testing.B) {
	config := &ResponseConfig{
		MessageTemplates: map[string]MessageTemplate{
			"welcome": {
				Key:      "welcome",
				Template: "Welcome $name",
				CodeMappings: map[string]int{
					"http": 200,
				},
				Translations: map[string]string{
					"en": "Welcome $name",
					"id": "Selamat datang $name",
				},
			},
		},
		DefaultLanguage: "en",
		Languages:       []string{"en", "id"},
		Translations: map[string]map[string]string{
			"en": {"hello": "Hello"},
			"id": {"hello": "Halo"},
		},
	}

	printer := NewConfigPrinter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = printer.Export()
	}
}
