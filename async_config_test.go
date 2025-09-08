package goresponse

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewAsyncConfigManager tests AsyncConfigManager creation
func TestNewAsyncConfigManager(t *testing.T) {
	source := ConfigSource{
		Method: "file",
		Path:   "test.json",
	}
	interval := 5 * time.Second

	manager := NewAsyncConfigManager(source, interval)

	if manager == nil {
		t.Fatal("Expected manager to be created")
	}

	if manager.source != source {
		t.Errorf("Expected source %+v, got %+v", source, manager.source)
	}

	if manager.interval != interval {
		t.Errorf("Expected interval %v, got %v", interval, manager.interval)
	}

	if manager.ctx == nil {
		t.Error("Expected context to be initialized")
	}

	if manager.cancel == nil {
		t.Error("Expected cancel function to be initialized")
	}

	if manager.callbacks == nil {
		t.Error("Expected callbacks slice to be initialized")
	}

	if manager.isRunning {
		t.Error("Expected manager to not be running initially")
	}

	if manager.lastError != nil {
		t.Errorf("Expected no initial error, got %v", manager.lastError)
	}
}

// TestAsyncConfigManagerStart tests the Start method
func TestAsyncConfigManagerStart(t *testing.T) {
	tests := []struct {
		name          string
		source        ConfigSource
		setupFiles    map[string]string
		expectedError bool
		errorContains string
	}{
		{
			name: "Start with valid config",
			source: ConfigSource{
				Method: "file",
				Path:   "test_start.json",
			},
			setupFiles: map[string]string{
				"test_start.json": `{
					"default_language": "en",
					"languages": ["en", "id"],
					"translations": {
						"en": {
							"test": "Test message"
						}
					}
				}`,
			},
			expectedError: false,
		},
		{
			name: "Start with invalid config file",
			source: ConfigSource{
				Method: "file",
				Path:   "nonexistent.json",
			},
			expectedError: true,
			errorContains: "failed to load initial config",
		},
		{
			name: "Start with invalid JSON",
			source: ConfigSource{
				Method: "file",
				Path:   "invalid.json",
			},
			setupFiles: map[string]string{
				"invalid.json": `{ invalid json }`,
			},
			expectedError: true,
			errorContains: "failed to load initial config",
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

			manager := NewAsyncConfigManager(tt.source, 1*time.Second)
			defer manager.Stop()

			err := manager.Start()

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

			// Check if manager is running
			if !manager.IsRunning() {
				t.Error("Expected manager to be running")
			}

			// Check if config is loaded
			config := manager.GetConfig()
			if config == nil {
				t.Error("Expected config to be loaded")
			}
		})
	}
}

// TestAsyncConfigManagerStartTwice tests starting manager twice
func TestAsyncConfigManagerStartTwice(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_double_start.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_double_start.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_double_start.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	// Start first time
	err = manager.Start()
	if err != nil {
		t.Fatalf("First start failed: %v", err)
	}

	// Try to start second time
	err = manager.Start()
	if err == nil {
		t.Error("Expected error when starting twice, but got none")
	}

	if !strings.Contains(err.Error(), "already running") {
		t.Errorf("Expected error about already running, got: %v", err)
	}
}

// TestAsyncConfigManagerStop tests the Stop method
func TestAsyncConfigManagerStop(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_stop.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_stop.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_stop.json",
	}

	manager := NewAsyncConfigManager(source, 100*time.Millisecond)

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Verify it's running
	if !manager.IsRunning() {
		t.Error("Expected manager to be running")
	}

	// Stop manager
	manager.Stop()

	// Verify it's stopped
	if manager.IsRunning() {
		t.Error("Expected manager to be stopped")
	}

	// Stop again should not cause issues
	manager.Stop()
}

// TestAsyncConfigManagerRefreshLoop tests the refresh loop functionality
func TestAsyncConfigManagerRefreshLoop(t *testing.T) {
	// Create initial config file
	initialConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Initial message"
			}
		}
	}`

	err := ioutil.WriteFile("test_refresh.json", []byte(initialConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_refresh.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_refresh.json",
	}

	manager := NewAsyncConfigManager(source, 100*time.Millisecond)
	defer manager.Stop()

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Verify initial config
	config := manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	translation, exists := config.GetTranslation("en", "test")
	if !exists {
		t.Fatal("Expected translation to exist")
	}
	if translation != "Initial message" {
		t.Errorf("Expected 'Initial message', got '%s'", translation)
	}

	// Update config file
	updatedConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Updated message"
			}
		}
	}`

	err = ioutil.WriteFile("test_refresh.json", []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update test config file: %v", err)
	}

	// Wait for refresh
	time.Sleep(200 * time.Millisecond)

	// Verify updated config
	config = manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	translation, exists = config.GetTranslation("en", "test")
	if !exists {
		t.Fatal("Expected translation to exist")
	}
	if translation != "Updated message" {
		t.Errorf("Expected 'Updated message', got '%s'", translation)
	}
}

// TestAsyncConfigManagerCallbacks tests callback functionality
func TestAsyncConfigManagerCallbacks(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_callbacks.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_callbacks.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_callbacks.json",
	}

	manager := NewAsyncConfigManager(source, 100*time.Millisecond)
	defer manager.Stop()

	// Track callback calls
	var callbackCalls int
	var lastOldConfig, lastNewConfig *ResponseConfig
	var callbackMutex sync.Mutex

	callback := func(oldConfig, newConfig *ResponseConfig) {
		callbackMutex.Lock()
		callbackCalls++
		lastOldConfig = oldConfig
		lastNewConfig = newConfig
		callbackMutex.Unlock()
	}

	// Add callback
	manager.AddCallback(callback)

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Update config file to trigger callback
	updatedConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Updated message"
			}
		}
	}`

	err = ioutil.WriteFile("test_callbacks.json", []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update test config file: %v", err)
	}

	// Wait for refresh and callback
	time.Sleep(200 * time.Millisecond)

	// Check callback was called
	callbackMutex.Lock()
	if callbackCalls == 0 {
		t.Error("Expected callback to be called")
	}
	if lastOldConfig == nil {
		t.Error("Expected old config to be passed to callback")
	}
	if lastNewConfig == nil {
		t.Error("Expected new config to be passed to callback")
	}
	callbackMutex.Unlock()

	// Test multiple callbacks
	var secondCallbackCalls int
	secondCallback := func(oldConfig, newConfig *ResponseConfig) {
		secondCallbackCalls++
	}

	manager.AddCallback(secondCallback)

	// Update config again
	err = ioutil.WriteFile("test_callbacks.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update test config file: %v", err)
	}

	// Wait for refresh
	time.Sleep(200 * time.Millisecond)

	// Check both callbacks were called
	callbackMutex.Lock()
	if callbackCalls < 2 {
		t.Errorf("Expected at least 2 callback calls, got %d", callbackCalls)
	}
	callbackMutex.Unlock()

	if secondCallbackCalls == 0 {
		t.Error("Expected second callback to be called")
	}

	// Test RemoveAllCallbacks
	manager.RemoveAllCallbacks()

	// Update config again
	err = ioutil.WriteFile("test_callbacks.json", []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update test config file: %v", err)
	}

	// Wait for refresh
	time.Sleep(200 * time.Millisecond)

	// Check callbacks were not called
	callbackMutex.Lock()
	initialCalls := callbackCalls
	callbackMutex.Unlock()

	time.Sleep(100 * time.Millisecond)

	callbackMutex.Lock()
	if callbackCalls != initialCalls {
		t.Error("Expected callbacks to not be called after removal")
	}
	callbackMutex.Unlock()
}

// TestAsyncConfigManagerThreadSafety tests thread safety of various methods
func TestAsyncConfigManagerThreadSafety(t *testing.T) {
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

	err := ioutil.WriteFile("test_thread_safety.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_thread_safety.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_thread_safety.json",
	}

	manager := NewAsyncConfigManager(source, 50*time.Millisecond)
	defer manager.Stop()

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(100 * time.Millisecond)

	// Test concurrent access
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Test various read operations
				manager.GetConfig()
				manager.GetTranslation("en", "test")
				manager.GetTranslationWithFallback("en", "test")
				manager.GetMessageTemplate("welcome")
				manager.GetSupportedLanguages()
				manager.GetDefaultLanguage()
				manager.IsRunning()
				manager.GetLastError()
				manager.GetMessageTemplateTranslation("welcome", "en")
				manager.GetMessageTemplateTranslationWithFallback("welcome", "en")
			}
		}()
	}

	// Concurrent writers (update source and interval)
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				manager.UpdateSource(source)
				manager.UpdateInterval(time.Duration(j) * time.Millisecond)
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify manager is still running
	if !manager.IsRunning() {
		t.Error("Expected manager to still be running after concurrent access")
	}
}

// TestAsyncConfigManagerErrorHandling tests error handling during refresh
func TestAsyncConfigManagerErrorHandling(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_error.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_error.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_error.json",
	}

	manager := NewAsyncConfigManager(source, 100*time.Millisecond)
	defer manager.Stop()

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Verify initial config is loaded
	config := manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	// Remove config file to cause error
	err = os.Remove("test_error.json")
	if err != nil {
		t.Fatalf("Failed to remove config file: %v", err)
	}

	// Wait for refresh attempt
	time.Sleep(200 * time.Millisecond)

	// Check that error is recorded
	lastError := manager.GetLastError()
	if lastError == nil {
		t.Error("Expected error to be recorded")
	}

	// Check that old config is still available
	config = manager.GetConfig()
	if config == nil {
		t.Error("Expected old config to still be available after error")
	}

	// Recreate config file
	err = ioutil.WriteFile("test_error.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to recreate config file: %v", err)
	}

	// Wait for successful refresh
	time.Sleep(200 * time.Millisecond)

	// Check that error is cleared
	lastError = manager.GetLastError()
	if lastError != nil {
		t.Errorf("Expected error to be cleared, got: %v", lastError)
	}
}

// TestAsyncConfigManagerForceRefresh tests ForceRefresh method
func TestAsyncConfigManagerForceRefresh(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Initial message"
			}
		}
	}`

	err := ioutil.WriteFile("test_force_refresh.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_force_refresh.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_force_refresh.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Verify initial config
	config := manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	translation, exists := config.GetTranslation("en", "test")
	if !exists {
		t.Fatal("Expected translation to exist")
	}
	if translation != "Initial message" {
		t.Errorf("Expected 'Initial message', got '%s'", translation)
	}

	// Update config file
	updatedConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Updated message"
			}
		}
	}`

	err = ioutil.WriteFile("test_force_refresh.json", []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update test config file: %v", err)
	}

	// Force refresh
	err = manager.ForceRefresh()
	if err != nil {
		t.Errorf("ForceRefresh failed: %v", err)
	}

	// Verify updated config
	config = manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	translation, exists = config.GetTranslation("en", "test")
	if !exists {
		t.Fatal("Expected translation to exist")
	}
	if translation != "Updated message" {
		t.Errorf("Expected 'Updated message', got '%s'", translation)
	}
}

// TestAsyncConfigManagerMessageTemplateManagement tests message template management
func TestAsyncConfigManagerMessageTemplateManagement(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_template_management.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_template_management.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_template_management.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Test AddMessageTemplate
	template := &MessageTemplate{
		Key:      "manual_test",
		Template: "Manual template",
		CodeMappings: map[string]int{
			"http": 201,
		},
		Translations: map[string]string{
			"en": "Manual template",
			"id": "Template manual",
		},
	}

	manager.AddMessageTemplate(template)

	// Verify template was added
	config := manager.GetConfig()
	if config == nil {
		t.Fatal("Expected config to be loaded")
	}

	retrievedTemplate, exists := config.GetMessageTemplate("manual_test")
	if !exists {
		t.Error("Expected manual template to exist")
	}
	if retrievedTemplate.Template != "Manual template" {
		t.Errorf("Expected 'Manual template', got '%s'", retrievedTemplate.Template)
	}

	// Test AddMessageTemplates
	template1 := &MessageTemplate{
		Key:      "template1",
		Template: "Template 1",
	}
	template2 := &MessageTemplate{
		Key:      "template2",
		Template: "Template 2",
	}

	manager.AddMessageTemplates(template1, template2)

	// Verify templates were added
	_, exists = config.GetMessageTemplate("template1")
	if !exists {
		t.Error("Expected template1 to exist")
	}
	_, exists = config.GetMessageTemplate("template2")
	if !exists {
		t.Error("Expected template2 to exist")
	}

	// Test UpdateMessageTemplate
	updatedTemplate := &MessageTemplate{
		Key:      "manual_test",
		Template: "Updated manual template",
		CodeMappings: map[string]int{
			"http": 202,
		},
	}

	manager.UpdateMessageTemplate(updatedTemplate)

	// Verify template was updated
	retrievedTemplate, exists = config.GetMessageTemplate("manual_test")
	if !exists {
		t.Error("Expected manual template to exist")
	}
	if retrievedTemplate.Template != "Updated manual template" {
		t.Errorf("Expected 'Updated manual template', got '%s'", retrievedTemplate.Template)
	}

	// Test RemoveMessageTemplate
	manager.RemoveMessageTemplate("template1")

	// Verify template was removed
	_, exists = config.GetMessageTemplate("template1")
	if exists {
		t.Error("Expected template1 to be removed")
	}
}

// TestAsyncConfigManagerPrinterMethods tests printer methods
func TestAsyncConfigManagerPrinterMethods(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_printer.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_printer.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_printer.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	// Test printer methods before starting
	printer := manager.Printer()
	if printer != nil {
		t.Error("Expected printer to be nil before starting")
	}

	err = manager.PrintConfig()
	if err == nil {
		t.Error("Expected error when printing before starting")
	}

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Test printer methods after starting
	printer = manager.Printer()
	if printer == nil {
		t.Error("Expected printer to be created after starting")
	}

	// Test PrintConfig
	err = manager.PrintConfig()
	if err != nil {
		t.Errorf("PrintConfig failed: %v", err)
	}

	// Test PrintConfigWithIndent
	err = manager.PrintConfigWithIndent(true)
	if err != nil {
		t.Errorf("PrintConfigWithIndent(true) failed: %v", err)
	}

	err = manager.PrintConfigWithIndent(false)
	if err != nil {
		t.Errorf("PrintConfigWithIndent(false) failed: %v", err)
	}

	// Test ExportConfig
	jsonStr, err := manager.ExportConfig()
	if err != nil {
		t.Errorf("ExportConfig failed: %v", err)
	}

	// Verify it's valid JSON
	var exportedConfig ResponseConfig
	err = json.Unmarshal([]byte(jsonStr), &exportedConfig)
	if err != nil {
		t.Errorf("Exported config is not valid JSON: %v", err)
	}

	// Test ExportConfigToFile
	testFile := "test_export_async.json"
	defer os.Remove(testFile)

	err = manager.ExportConfigToFile(testFile)
	if err != nil {
		t.Errorf("ExportConfigToFile failed: %v", err)
	}

	// Verify file was created and contains valid JSON
	data, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read exported file: %v", err)
	}

	err = json.Unmarshal(data, &exportedConfig)
	if err != nil {
		t.Errorf("Exported file is not valid JSON: %v", err)
	}
}

// TestAsyncConfigManagerWithNilConfig tests methods with nil config
func TestAsyncConfigManagerWithNilConfig(t *testing.T) {
	manager := &AsyncConfigManager{}

	// Test all methods with nil config
	config := manager.GetConfig()
	if config != nil {
		t.Error("Expected config to be nil")
	}

	translation, exists := manager.GetTranslation("en", "test")
	if exists {
		t.Error("Expected translation to not exist")
	}
	if translation != "" {
		t.Errorf("Expected empty string, got '%s'", translation)
	}

	fallbackTranslation := manager.GetTranslationWithFallback("en", "test")
	if fallbackTranslation != "test" {
		t.Errorf("Expected 'test' (key), got '%s'", fallbackTranslation)
	}

	template, exists := manager.GetMessageTemplate("test")
	if exists {
		t.Error("Expected template to not exist")
	}
	if template != nil {
		t.Error("Expected nil template")
	}

	languages := manager.GetSupportedLanguages()
	if len(languages) != 0 {
		t.Errorf("Expected empty languages slice, got %v", languages)
	}

	defaultLang := manager.GetDefaultLanguage()
	if defaultLang != "" {
		t.Errorf("Expected empty default language, got '%s'", defaultLang)
	}

	// Test message template management with nil config
	templateToAdd := &MessageTemplate{
		Key:      "test",
		Template: "Test",
	}
	manager.AddMessageTemplate(templateToAdd)
	manager.AddMessageTemplates(templateToAdd)
	manager.RemoveMessageTemplate("test")
	manager.UpdateMessageTemplate(templateToAdd)

	// These should not panic
	manager.GetMessageTemplateTranslation("test", "en")
	manager.GetMessageTemplateTranslationWithFallback("test", "en")
}

// TestAsyncConfigManagerUpdateSourceAndInterval tests updating source and interval
func TestAsyncConfigManagerUpdateSourceAndInterval(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_update.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_update.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_update.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Update source
	newSource := ConfigSource{
		Method: "file",
		Path:   "test_update.json",
	}
	manager.UpdateSource(newSource)

	// Update interval
	newInterval := 500 * time.Millisecond
	manager.UpdateInterval(newInterval)

	// Verify updates (these are internal, so we can't directly test them)
	// But we can verify the manager still works
	config := manager.GetConfig()
	if config == nil {
		t.Error("Expected config to be loaded")
	}
}

// TestAsyncConfigManagerContextCancellation tests context cancellation
func TestAsyncConfigManagerContextCancellation(t *testing.T) {
	// Create test config file
	testConfig := `{
		"default_language": "en",
		"languages": ["en"],
		"translations": {
			"en": {
				"test": "Test message"
			}
		}
	}`

	err := ioutil.WriteFile("test_context.json", []byte(testConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	defer os.Remove("test_context.json")

	source := ConfigSource{
		Method: "file",
		Path:   "test_context.json",
	}

	manager := NewAsyncConfigManager(source, 100*time.Millisecond)

	// Start manager
	err = manager.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	// Verify it's running
	if !manager.IsRunning() {
		t.Error("Expected manager to be running")
	}

	// Stop manager
	manager.Stop()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Verify it's stopped
	if manager.IsRunning() {
		t.Error("Expected manager to be stopped")
	}

	// Try to start again
	err = manager.Start()
	if err != nil {
		t.Errorf("Expected to be able to start again, got error: %v", err)
	}

	// Verify it's running again
	if !manager.IsRunning() {
		t.Error("Expected manager to be running again")
	}
}

// BenchmarkAsyncConfigManagerGetConfig benchmarks GetConfig performance
func BenchmarkAsyncConfigManagerGetConfig(b *testing.B) {
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
		}
	}`

	err := ioutil.WriteFile("benchmark_async_config.json", []byte(testConfig), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark config file: %v", err)
	}
	defer os.Remove("benchmark_async_config.json")

	source := ConfigSource{
		Method: "file",
		Path:   "benchmark_async_config.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	err = manager.Start()
	if err != nil {
		b.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.GetConfig()
	}
}

// BenchmarkAsyncConfigManagerGetTranslation benchmarks GetTranslation performance
func BenchmarkAsyncConfigManagerGetTranslation(b *testing.B) {
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
		}
	}`

	err := ioutil.WriteFile("benchmark_async_translation.json", []byte(testConfig), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark config file: %v", err)
	}
	defer os.Remove("benchmark_async_translation.json")

	source := ConfigSource{
		Method: "file",
		Path:   "benchmark_async_translation.json",
	}

	manager := NewAsyncConfigManager(source, 1*time.Second)
	defer manager.Stop()

	err = manager.Start()
	if err != nil {
		b.Fatalf("Start failed: %v", err)
	}

	// Wait for initial load
	time.Sleep(50 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.GetTranslation("en", "test")
	}
}
