package goresponse

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestLoadConfig tests the main LoadConfig function
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		source         ConfigSource
		expectedError  bool
		errorContains  string
		setupFiles     map[string]string // filename -> content
		expectedConfig *ResponseConfig
	}{
		{
			name: "Load from valid file",
			source: ConfigSource{
				Method: "file",
				Path:   "test_config.json",
			},
			expectedError: false,
			setupFiles: map[string]string{
				"test_config.json": `{
					"message_templates": {
						"success": {
							"key": "success",
							"template": "Operation successful",
							"code_mappings": {
								"http": 200
							}
						}
					},
					"default_language": "en",
					"languages": ["en", "id"],
					"translations": {
						"en": {
							"success": "Operation successful"
						},
						"id": {
							"success": "Operasi berhasil"
						}
					}
				}`,
			},
			expectedConfig: &ResponseConfig{
				DefaultLanguage: "en",
				Languages:       []string{"en", "id"},
				Translations: map[string]map[string]string{
					"en": {"success": "Operation successful"},
					"id": {"success": "Operasi berhasil"},
				},
				MessageTemplates: map[string]MessageTemplate{
					"success": {
						Key:          "success",
						Template:     "Operation successful",
						CodeMappings: map[string]int{"http": 200},
					},
				},
			},
		},
		{
			name: "Load from URL",
			source: ConfigSource{
				Method: "url",
				Path:   "http://example.com/config.json",
			},
			expectedError: false,
			setupFiles: map[string]string{
				"config.json": `{
					"default_language": "en",
					"languages": ["en"],
					"translations": {
						"en": {
							"test": "Test message"
						}
					}
				}`,
			},
		},
		{
			name: "Unsupported method",
			source: ConfigSource{
				Method: "database",
				Path:   "test",
			},
			expectedError: true,
			errorContains: "unsupported method",
		},
		{
			name: "File not found",
			source: ConfigSource{
				Method: "file",
				Path:   "nonexistent.json",
			},
			expectedError: true,
			errorContains: "failed to open config file",
		},
		{
			name: "Invalid JSON",
			source: ConfigSource{
				Method: "file",
				Path:   "invalid.json",
			},
			expectedError: true,
			errorContains: "failed to unmarshal config",
			setupFiles: map[string]string{
				"invalid.json": `{ invalid json }`,
			},
		},
		{
			name: "With translation sources",
			source: ConfigSource{
				Method: "file",
				Path:   "config_with_sources.json",
			},
			expectedError: false,
			setupFiles: map[string]string{
				"config_with_sources.json": `{
					"default_language": "en",
					"languages": ["en", "id"],
					"translations": {
						"en": {
							"existing": "Existing translation"
						}
					},
					"translation_source": {
						"en": {
							"method": "file",
							"path": "en_translations.json"
						},
						"id": {
							"method": "file",
							"path": "id_translations.json"
						}
					}
				}`,
				"en_translations.json": `{
					"new_key": "New English translation",
					"existing": "Override existing"
				}`,
				"id_translations.json": `{
					"new_key": "Terjemahan Indonesia baru"
				}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test files
			for filename, content := range tt.setupFiles {
				err := ioutil.WriteFile(filename, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
				defer os.Remove(filename)
			}

			// Setup HTTP server for URL tests
			if tt.source.Method == "url" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if content, exists := tt.setupFiles["config.json"]; exists {
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte(content))
					} else {
						w.WriteHeader(http.StatusNotFound)
					}
				}))
				defer server.Close()
				tt.source.Path = server.URL
			}

			// Test LoadConfig
			config, err := LoadConfig(tt.source)

			// Check error expectations
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check config if expected
			if tt.expectedConfig != nil {
				if config.DefaultLanguage != tt.expectedConfig.DefaultLanguage {
					t.Errorf("Expected default language %s, got %s", tt.expectedConfig.DefaultLanguage, config.DefaultLanguage)
				}
				if len(config.Languages) != len(tt.expectedConfig.Languages) {
					t.Errorf("Expected %d languages, got %d", len(tt.expectedConfig.Languages), len(config.Languages))
				}
				if len(config.Translations) != len(tt.expectedConfig.Translations) {
					t.Errorf("Expected %d translation languages, got %d", len(tt.expectedConfig.Translations), len(config.Translations))
				}
			}
		})
	}
}

// TestLoadFromFile tests the loadFromFile function
func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		fileContent   string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Valid file",
			filePath:      "test_file.json",
			fileContent:   `{"test": "data"}`,
			expectedError: false,
		},
		{
			name:          "File not found",
			filePath:      "nonexistent.json",
			expectedError: true,
			errorContains: "failed to open config file",
		},
		{
			name:          "Empty file",
			filePath:      "empty.json",
			fileContent:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file
			if tt.fileContent != "" || tt.name == "Empty file" {
				err := ioutil.WriteFile(tt.filePath, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer os.Remove(tt.filePath)
			}

			// Test loadFromFile
			data, err := loadFromFile(tt.filePath)

			// Check error expectations
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check data
			if tt.fileContent != "" && string(data) != tt.fileContent {
				t.Errorf("Expected data %s, got %s", tt.fileContent, string(data))
			}
		})
	}
}

// TestLoadFromURL tests the loadFromURL function
func TestLoadFromURL(t *testing.T) {
	tests := []struct {
		name          string
		responseCode  int
		responseBody  string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Valid response",
			responseCode:  http.StatusOK,
			responseBody:  `{"test": "data"}`,
			expectedError: false,
		},
		{
			name:          "Not found",
			responseCode:  http.StatusNotFound,
			responseBody:  "",
			expectedError: true,
			errorContains: "status code: 404",
		},
		{
			name:          "Server error",
			responseCode:  http.StatusInternalServerError,
			responseBody:  "",
			expectedError: true,
			errorContains: "status code: 500",
		},
		{
			name:          "Empty response",
			responseCode:  http.StatusOK,
			responseBody:  "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Test loadFromURL
			data, err := loadFromURL(server.URL)

			// Check error expectations
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check data
			if string(data) != tt.responseBody {
				t.Errorf("Expected data %s, got %s", tt.responseBody, string(data))
			}
		})
	}
}

// TestLoadTranslationSources tests the loadTranslationSources function
func TestLoadTranslationSources(t *testing.T) {
	tests := []struct {
		name          string
		config        *ResponseConfig
		setupFiles    map[string]string
		expectedError bool
		errorContains string
		expectedLangs []string
	}{
		{
			name: "No translation sources",
			config: &ResponseConfig{
				Translations: map[string]map[string]string{
					"en": {"test": "test"},
				},
			},
			expectedError: false,
		},
		{
			name: "Valid translation sources",
			config: &ResponseConfig{
				Translations: map[string]map[string]string{
					"en": {"existing": "existing"},
				},
				TranslationSources: map[string]TranslationSource{
					"en": {
						Method: "file",
						Path:   "en.json",
					},
					"id": {
						Method: "file",
						Path:   "id.json",
					},
				},
			},
			setupFiles: map[string]string{
				"en.json": `{"new_key": "New English", "existing": "Override"}`,
				"id.json": `{"new_key": "Baru Indonesia"}`,
			},
			expectedError: false,
			expectedLangs: []string{"en", "id"},
		},
		{
			name: "Translation source file not found",
			config: &ResponseConfig{
				TranslationSources: map[string]TranslationSource{
					"en": {
						Method: "file",
						Path:   "nonexistent.json",
					},
				},
			},
			expectedError: true,
			errorContains: "failed to load translations for language en",
		},
		{
			name: "Invalid translation source method",
			config: &ResponseConfig{
				TranslationSources: map[string]TranslationSource{
					"en": {
						Method: "database",
						Path:   "test",
					},
				},
			},
			expectedError: true,
			errorContains: "unsupported translation source method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test files
			for filename, content := range tt.setupFiles {
				err := ioutil.WriteFile(filename, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
				defer os.Remove(filename)
			}

			// Test loadTranslationSources
			err := loadTranslationSources(tt.config)

			// Check error expectations
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check expected languages
			if tt.expectedLangs != nil {
				for _, lang := range tt.expectedLangs {
					if _, exists := tt.config.Translations[lang]; !exists {
						t.Errorf("Expected language %s to be loaded", lang)
					}
				}
			}
		})
	}
}

// TestResponseConfigMethods tests ResponseConfig methods
func TestResponseConfigMethods(t *testing.T) {
	config := &ResponseConfig{
		DefaultLanguage: "en",
		Languages:       []string{"en", "id", "es"},
		Translations: map[string]map[string]string{
			"en": {"hello": "Hello", "goodbye": "Goodbye"},
			"id": {"hello": "Halo", "goodbye": "Selamat tinggal"},
			"es": {"hello": "Hola"},
		},
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
	}

	t.Run("GetTranslation", func(t *testing.T) {
		// Test existing translation
		translation, exists := config.GetTranslation("en", "hello")
		if !exists {
			t.Error("Expected translation to exist")
		}
		if translation != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", translation)
		}

		// Test non-existing translation
		_, exists = config.GetTranslation("en", "nonexistent")
		if exists {
			t.Error("Expected translation to not exist")
		}

		// Test non-existing language
		_, exists = config.GetTranslation("fr", "hello")
		if exists {
			t.Error("Expected translation to not exist")
		}
	})

	t.Run("GetMessageTemplate", func(t *testing.T) {
		// Test manual template (priority)
		template, exists := config.GetMessageTemplate("manual")
		if !exists {
			t.Error("Expected manual template to exist")
		}
		if template.Template != "Manual template" {
			t.Errorf("Expected 'Manual template', got '%s'", template.Template)
		}

		// Test async template
		template, exists = config.GetMessageTemplate("welcome")
		if !exists {
			t.Error("Expected async template to exist")
		}
		if template.Template != "Welcome $name" {
			t.Errorf("Expected 'Welcome $name', got '%s'", template.Template)
		}

		// Test non-existing template
		_, exists = config.GetMessageTemplate("nonexistent")
		if exists {
			t.Error("Expected template to not exist")
		}
	})

	t.Run("GetSupportedLanguages", func(t *testing.T) {
		languages := config.GetSupportedLanguages()
		expected := []string{"en", "id", "es"}
		if len(languages) != len(expected) {
			t.Errorf("Expected %d languages, got %d", len(expected), len(languages))
		}
	})

	t.Run("GetDefaultLanguage", func(t *testing.T) {
		lang := config.GetDefaultLanguage()
		if lang != "en" {
			t.Errorf("Expected 'en', got '%s'", lang)
		}
	})

	t.Run("GetTranslationWithFallback", func(t *testing.T) {
		// Test existing translation
		translation := config.GetTranslationWithFallback("en", "hello")
		if translation != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", translation)
		}

		// Test fallback to default language
		translation = config.GetTranslationWithFallback("fr", "hello")
		if translation != "Hello" {
			t.Errorf("Expected 'Hello' (fallback), got '%s'", translation)
		}

		// Test fallback to key
		translation = config.GetTranslationWithFallback("en", "nonexistent")
		if translation != "nonexistent" {
			t.Errorf("Expected 'nonexistent' (key), got '%s'", translation)
		}
	})

	t.Run("GetMessageTemplateTranslation", func(t *testing.T) {
		// Test manual template translation (priority)
		_, exists := config.GetMessageTemplateTranslation("manual", "en")
		if !exists {
			t.Error("Expected manual template to not have translations")
		}

		// Test async template translation
		translation, exists := config.GetMessageTemplateTranslation("welcome", "en")
		if !exists {
			t.Error("Expected async template translation to exist")
		}
		if translation != "Welcome $name" {
			t.Errorf("Expected 'Welcome $name', got '%s'", translation)
		}

		// Test fallback to default language
		translation, exists = config.GetMessageTemplateTranslation("welcome", "fr")
		if !exists {
			t.Error("Expected fallback translation to exist")
		}
		if translation != "Welcome $name" {
			t.Errorf("Expected 'Welcome $name' (fallback), got '%s'", translation)
		}
	})

	t.Run("GetMessageTemplateTranslationWithFallback", func(t *testing.T) {
		// Test existing translation
		translation := config.GetMessageTemplateTranslationWithFallback("welcome", "en")
		if translation != "Welcome $name" {
			t.Errorf("Expected 'Welcome $name', got '%s'", translation)
		}

		// Test fallback to default language
		translation = config.GetMessageTemplateTranslationWithFallback("welcome", "fr")
		if translation != "Welcome $name" {
			t.Errorf("Expected 'Welcome $name' (fallback), got '%s'", translation)
		}

		// Test fallback to template
		translation = config.GetMessageTemplateTranslationWithFallback("manual", "en")
		if translation != "Manual template" {
			t.Errorf("Expected 'Manual template' (template), got '%s'", translation)
		}

		// Test fallback to key
		translation = config.GetMessageTemplateTranslationWithFallback("nonexistent", "en")
		if translation != "nonexistent" {
			t.Errorf("Expected 'nonexistent' (key), got '%s'", translation)
		}
	})
}

// TestMessageTemplateManagement tests message template management methods
func TestMessageTemplateManagement(t *testing.T) {
	config := &ResponseConfig{}

	t.Run("AddMessageTemplate", func(t *testing.T) {
		template := &MessageTemplate{
			Key:      "test",
			Template: "Test template",
			CodeMappings: map[string]int{
				"http": 200,
			},
			Translations: map[string]string{
				"en": "Test template",
				"id": "Template uji",
			},
		}

		config.AddMessageTemplate(template)

		// Check if template was added
		if config.ManualMessageTemplates == nil {
			t.Error("Expected ManualMessageTemplates to be initialized")
		}
		if _, exists := config.ManualMessageTemplates["test"]; !exists {
			t.Error("Expected template to be added")
		}
	})

	t.Run("AddMessageTemplates", func(t *testing.T) {
		template1 := &MessageTemplate{
			Key:      "test1",
			Template: "Test template 1",
		}
		template2 := &MessageTemplate{
			Key:      "test2",
			Template: "Test template 2",
		}

		config.AddMessageTemplates(template1, template2)

		// Check if templates were added
		if _, exists := config.ManualMessageTemplates["test1"]; !exists {
			t.Error("Expected template1 to be added")
		}
		if _, exists := config.ManualMessageTemplates["test2"]; !exists {
			t.Error("Expected template2 to be added")
		}
	})

	t.Run("RemoveMessageTemplate", func(t *testing.T) {
		config.RemoveMessageTemplate("test")

		// Check if template was removed
		if _, exists := config.ManualMessageTemplates["test"]; exists {
			t.Error("Expected template to be removed")
		}
	})

	t.Run("UpdateMessageTemplate", func(t *testing.T) {
		// Add initial template
		originalTemplate := &MessageTemplate{
			Key:      "update_test",
			Template: "Original template",
		}
		config.AddMessageTemplate(originalTemplate)

		// Update template
		updatedTemplate := &MessageTemplate{
			Key:      "update_test",
			Template: "Updated template",
		}
		config.UpdateMessageTemplate(updatedTemplate)

		// Check if template was updated
		if config.ManualMessageTemplates["update_test"].Template != "Updated template" {
			t.Error("Expected template to be updated")
		}
	})
}

// TestConfigManager tests ConfigManager functionality
func TestConfigManager(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en", "id"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_manager_config.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_manager_config.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_manager_config.json",
	}

	manager := NewConfigManager(source)

	t.Run("Load", func(t *testing.T) {
		err := manager.Load()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		config := manager.GetConfig()
		if config == nil {
			t.Error("Expected config to be loaded")
		}
		if config.DefaultLanguage != "en" {
			t.Errorf("Expected default language 'en', got '%s'", config.DefaultLanguage)
		}
	})

	t.Run("Reload", func(t *testing.T) {
		err := manager.Reload()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("GetTranslationWithFallback", func(t *testing.T) {
		translation := manager.GetTranslationWithFallback("en", "test")
		if translation != "Test message" {
			t.Errorf("Expected 'Test message', got '%s'", translation)
		}

		// Test with nil config
		emptyManager := &ConfigManager{}
		translation = emptyManager.GetTranslationWithFallback("en", "test")
		if translation != "" {
			t.Errorf("Expected empty string, got '%s'", translation)
		}
	})
}

// TestConfigPrinterIntegration tests ConfigPrinter integration
func TestConfigPrinterIntegration(t *testing.T) {
	config := &ResponseConfig{
		DefaultLanguage: "en",
		Languages:       []string{"en", "id"},
		Translations: map[string]map[string]string{
			"en": {"test": "Test"},
		},
	}

	t.Run("Printer", func(t *testing.T) {
		printer := config.Printer()
		if printer == nil {
			t.Error("Expected printer to be created")
		}
	})

	t.Run("PrintConfig", func(t *testing.T) {
		// This test just ensures the method doesn't panic
		// In a real test, you might want to capture stdout
		err := config.PrintConfig()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("PrintConfigWithIndent", func(t *testing.T) {
		err := config.PrintConfigWithIndent(true)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		err = config.PrintConfigWithIndent(false)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("ExportConfig", func(t *testing.T) {
		jsonStr, err := config.ExportConfig()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify it's valid JSON
		var exportedConfig ResponseConfig
		err = json.Unmarshal([]byte(jsonStr), &exportedConfig)
		if err != nil {
			t.Errorf("Exported config is not valid JSON: %v", err)
		}

		if exportedConfig.DefaultLanguage != config.DefaultLanguage {
			t.Errorf("Expected default language %s, got %s", config.DefaultLanguage, exportedConfig.DefaultLanguage)
		}
	})

	t.Run("ExportConfigToFile", func(t *testing.T) {
		testFile := "test_export.json"
		defer os.Remove(testFile)

		err := config.ExportConfigToFile(testFile)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify file was created and contains valid JSON
		data, err := ioutil.ReadFile(testFile)
		if err != nil {
			t.Errorf("Failed to read exported file: %v", err)
		}

		var exportedConfig ResponseConfig
		err = json.Unmarshal(data, &exportedConfig)
		if err != nil {
			t.Errorf("Exported file is not valid JSON: %v", err)
		}
	})
}

// TestEdgeCases tests edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("Empty config", func(t *testing.T) {
		config := &ResponseConfig{}

		// Test methods with empty config
		translation, exists := config.GetTranslation("en", "test")
		if exists {
			t.Error("Expected translation to not exist")
		}
		if translation != "" {
			t.Errorf("Expected empty string, got '%s'", translation)
		}

		_, exists = config.GetMessageTemplate("test")
		if exists {
			t.Error("Expected template to not exist")
		}

		languages := config.GetSupportedLanguages()
		if languages != nil {
			t.Error("Expected nil languages")
		}

		defaultLang := config.GetDefaultLanguage()
		if defaultLang != "" {
			t.Errorf("Expected empty default language, got '%s'", defaultLang)
		}
	})

	t.Run("Nil maps", func(t *testing.T) {
		config := &ResponseConfig{
			Translations:           nil,
			MessageTemplates:       nil,
			ManualMessageTemplates: nil,
		}

		// Test methods with nil maps
		_, exists := config.GetTranslation("en", "test")
		if exists {
			t.Error("Expected translation to not exist")
		}

		_, exists = config.GetMessageTemplate("test")
		if exists {
			t.Error("Expected template to not exist")
		}

		// Test adding template to nil map
		templateToAdd := &MessageTemplate{
			Key:      "test",
			Template: "Test",
		}
		config.AddMessageTemplate(templateToAdd)

		if config.ManualMessageTemplates == nil {
			t.Error("Expected ManualMessageTemplates to be initialized")
		}
	})

	t.Run("Invalid JSON in translation source", func(t *testing.T) {
		// Create invalid JSON file
		err := ioutil.WriteFile("invalid_translation.json", []byte(`{ invalid json }`), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid JSON file: %v", err)
		}
		defer os.Remove("invalid_translation.json")

		config := &ResponseConfig{
			TranslationSources: map[string]TranslationSource{
				"en": {
					Method: "file",
					Path:   "invalid_translation.json",
				},
			},
		}

		err = loadTranslationSources(config)
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
		if !strings.Contains(err.Error(), "failed to unmarshal translations") {
			t.Errorf("Expected error about unmarshaling, got: %v", err)
		}
	})
}

// BenchmarkLoadConfig benchmarks LoadConfig performance
func BenchmarkLoadConfig(b *testing.B) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en", "id"],
		"translations": {
			"en": {
				"test": "Test message"
			},
			"id": {
				"test": "Pesan uji"
			}
		},
		"message_templates": {
			"welcome": {
				"key": "welcome",
				"template": "Welcome $name",
				"code_mappings": {
					"http": 200
				}
			}
		}
	}`

	err := ioutil.WriteFile("benchmark_config.json", []byte(testConfig), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark config file: %v", err)
	}
	defer os.Remove("benchmark_config.json")

	source := ConfigSource{
		Method: "file",
		Path:   "benchmark_config.json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadConfig(source)
		if err != nil {
			b.Fatalf("LoadConfig failed: %v", err)
		}
	}
}

// BenchmarkGetTranslation benchmarks GetTranslation performance
func BenchmarkGetTranslation(b *testing.B) {
	config := &ResponseConfig{
		Translations: map[string]map[string]string{
			"en": {
				"hello": "Hello",
				"world": "World",
				"test":  "Test",
			},
			"id": {
				"hello": "Halo",
				"world": "Dunia",
				"test":  "Uji",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = config.GetTranslation("en", "hello")
	}
}

// BenchmarkGetMessageTemplate benchmarks GetMessageTemplate performance
func BenchmarkGetMessageTemplate(b *testing.B) {
	config := &ResponseConfig{
		MessageTemplates: map[string]MessageTemplate{
			"welcome": {
				Key:      "welcome",
				Template: "Welcome $name",
				CodeMappings: map[string]int{
					"http": 200,
					"grpc": 0,
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
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = config.GetMessageTemplate("welcome")
	}
}
