package example

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.risoftinc.com/goresponse"
)

// File Structure:
// - types.go: Shared types (ConfigSource, ResponseConfig, MessageTemplate, ConfigChangeCallback, TranslationSource)
// - config.go: Sync loading functions and ConfigManager
// - async_config.go: Async loading with AsyncConfigManager
// - example_usage.go: Usage examples for both sync and async
// - translations/: Separate translation files per language

// ExampleUsage shows how to use simplified configuration functions
func ExampleUsage() {
	// Example 1: Load configuration from file
	fmt.Println("=== Example 1: Load Config from File ===")
	fileSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	config, err := goresponse.LoadConfig(fileSource)
	if err != nil {
		log.Printf("Error loading config from file: %v", err)
		return
	}

	// Display configuration information
	fmt.Printf("Default Language: %s\n", config.GetDefaultLanguage())
	fmt.Printf("Supported Languages: %v\n", config.GetSupportedLanguages())

	// Get translations
	translation, exists := config.GetTranslation("en", "success")
	if exists {
		fmt.Printf("English 'success' translation: %s\n", translation)
	}

	translation, exists = config.GetTranslation("id", "success")
	if exists {
		fmt.Printf("Indonesian 'success' translation: %s\n", translation)
	}

	// Get message template
	template, exists := config.GetMessageTemplate("validation_failed")
	if exists {
		fmt.Printf("Validation failed template: %s\n", template.Template)
		fmt.Printf("HTTP code mapping: %d\n", template.CodeMappings["http"])
	}

	fmt.Println()

	// Example 2: Load configuration from URL
	fmt.Println("=== Example 2: Load Config from URL ===")
	urlSource := goresponse.ConfigSource{
		Method: "url",
		Path:   "https://api.example.com/config.json",
	}

	// Uncomment to test URL (make sure URL is valid)
	// config2, err := LoadConfig(urlSource)
	// if err != nil {
	//     log.Printf("Error loading config from URL: %v", err)
	//     return
	// }
	// fmt.Printf("Config loaded from URL - Default Language: %s\n", config2.GetDefaultLanguage())

	// To avoid unused variable warning
	_ = urlSource

	fmt.Println()

	// Example 3: Using translations with parameters and fallback
	fmt.Println("=== Example 3: Using Translations with Parameters ===")

	// Simulate using translations with parameters
	userName := "John Doe"
	date := "2024-01-15"

	// Get template and perform manual substitution
	template, exists = config.GetMessageTemplate("welcome_message")
	if exists {
		// In real implementation, you would create a function for parameter substitution
		fmt.Printf("Welcome message template: %s\n", template.Template)
		fmt.Printf("With parameters - Welcome %s! Your account was created on %s\n", userName, date)
	}

	// Using fallback translation
	translation = config.GetTranslationWithFallback("fr", "success") // French not available, fallback to English
	fmt.Printf("French translation with fallback: %s\n", translation)

	translation = config.GetTranslationWithFallback("id", "success") // Indonesian available
	fmt.Printf("Indonesian translation: %s\n", translation)

	// Get translation with parameters
	translation, exists = config.GetTranslation("id", "user_created")
	if exists {
		fmt.Printf("Indonesian user_created: %s\n", translation)
		fmt.Printf("With parameters - User %s berhasil dibuat\n", userName)
	}
}

// ExampleConfigSourceJSON shows how to create ConfigSource from JSON
func ExampleConfigSourceJSON() {
	fmt.Println("=== Example ConfigSource from JSON ===")

	// JSON string for ConfigSource
	jsonStr := `{"method": "file", "path": "config.json"}`

	// Parse JSON to ConfigSource
	var source goresponse.ConfigSource
	if err := json.Unmarshal([]byte(jsonStr), &source); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("ConfigSource: Method=%s, Path=%s\n", source.Method, source.Path)

	// Load config using source
	config, err := goresponse.LoadConfig(source)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	fmt.Printf("Config loaded successfully - Default Language: %s\n", config.GetDefaultLanguage())
}

// ExampleConfigManagerUsage shows how to use ConfigManager (sync)
func ExampleConfigManagerUsage() {
	fmt.Println("=== Example ConfigManager ===")

	// Create config manager with file source
	fileSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}
	configManager := goresponse.NewConfigManager(fileSource)

	// Load configuration
	if err := configManager.Load(); err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// Using config manager
	config := configManager.GetConfig()
	fmt.Printf("Config loaded - Default Language: %s\n", config.GetDefaultLanguage())

	// Using fallback translation
	translation := configManager.GetTranslationWithFallback("fr", "success") // French not available, fallback to English
	fmt.Printf("Translation with fallback: %s\n", translation)

	translation = configManager.GetTranslationWithFallback("id", "success") // Indonesian available
	fmt.Printf("Indonesian translation: %s\n", translation)
}

// ExampleMessageTemplateBuilderUsage shows how to use MessageTemplateBuilder
func ExampleMessageTemplateBuilderUsage() {
	fmt.Println("=== Example MessageTemplateBuilder ===")

	// Example 1: Creating message template with method chaining
	fmt.Println("1. Creating message template with method chaining:")

	template1 := goresponse.NewMessageTemplateBuilder("user_created").
		WithTemplate("User $name has been created successfully").
		WithTranslation("en", "User $name has been created successfully").
		WithTranslation("id", "User $name berhasil dibuat").
		WithCodeMapping("http", 201).
		WithCodeMapping("grpc", 0).
		Build()

	fmt.Printf("Template 1: %+v\n", template1)

	// Example 2: Creating template with code mappings only
	template2 := goresponse.NewMessageTemplateBuilder("validation_error").
		WithTemplate("Validation failed for field: $field").
		WithCodeMapping("http", 422).
		WithCodeMapping("grpc", 3).
		Build()

	fmt.Printf("Template 2: %+v\n", template2)

	// Example 3: Creating template with multiple translations (using variadic)
	translations := map[string]string{
		"en": "Welcome $name! Your account was created on $date",
		"id": "Selamat datang $name! Akun Anda dibuat pada $date",
		"es": "¡Bienvenido $name! Tu cuenta fue creada el $date",
		"fr": "Bienvenue $name! Votre compte a été créé le $date",
	}

	codeMappings := map[string]int{
		"http": 200,
		"grpc": 0,
		"rest": 200,
	}

	template3 := goresponse.NewMessageTemplateBuilder("welcome_message").
		WithTemplate("Welcome $name! Your account was created on $date").
		WithTranslations(translations).
		WithCodeMappings(codeMappings).
		Build()

	fmt.Printf("Template 3: %+v\n", template3)

	fmt.Println()

	// Example 4: Using template with ResponseConfig
	fmt.Println("2. Using template with ResponseConfig:")

	// Load config
	source := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	config, err := goresponse.LoadConfig(source)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// Add templates to config
	config.AddMessageTemplates(template1, template2, template3)

	// Test template access
	if template, exists := config.GetMessageTemplate("user_created"); exists {
		fmt.Printf("Found template: %s\n", template.Template)
		fmt.Printf("HTTP code: %d\n", template.CodeMappings["http"])
	}

	// Test template translation
	translation := config.GetMessageTemplateTranslationWithFallback("user_created", "id")
	fmt.Printf("User created translation (ID): %s\n", translation)

	translation = config.GetMessageTemplateTranslationWithFallback("welcome_message", "es")
	fmt.Printf("Welcome translation (ES): %s\n", translation)

	// Test fallback
	translation = config.GetMessageTemplateTranslationWithFallback("user_created", "fr")
	fmt.Printf("User created translation (FR - fallback): %s\n", translation)

	fmt.Println()

	// Example 5: Using with AsyncConfigManager
	fmt.Println("3. Using with AsyncConfigManager:")

	asyncManager := goresponse.NewAsyncConfigManager(source, 30*time.Second)

	if err := asyncManager.Start(); err != nil {
		log.Printf("Error starting async manager: %v", err)
		return
	}

	// Add template to async manager
	asyncManager.AddMessageTemplate(template1)

	// Test async template access
	translation = asyncManager.GetMessageTemplateTranslationWithFallback("user_created", "en")
	fmt.Printf("Async template translation: %s\n", translation)

	asyncManager.Stop()
}

// ExampleManualTemplatePriorityUsage shows manual vs async template priority
func ExampleManualTemplatePriorityUsage() {
	fmt.Println("=== Example Manual Template Priority ===")

	// Load config with async templates
	source := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	config, err := goresponse.LoadConfig(source)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// Check template from async (config file)
	if template, exists := config.GetMessageTemplate("validation_failed"); exists {
		fmt.Printf("Async template: %s\n", template.Template)
	}

	// Add manual template with same key
	manualTemplate := goresponse.NewMessageTemplateBuilder("validation_failed").
		WithTemplate("Custom validation failed for field: $field (MANUAL OVERRIDE)").
		WithTranslations(map[string]string{
			"en": "Custom validation failed for field: $field (MANUAL OVERRIDE)",
			"id": "Validasi kustom gagal untuk field: $field (OVERRIDE MANUAL)",
		}).
		WithCodeMappings(map[string]int{
			"http": 400, // Override from 422 to 400
			"grpc": 3,
		}).
		Build()

	config.AddMessageTemplate(manualTemplate)

	// Check template after manual override
	if template, exists := config.GetMessageTemplate("validation_failed"); exists {
		fmt.Printf("Manual template (priority): %s\n", template.Template)
		fmt.Printf("HTTP code: %d\n", template.CodeMappings["http"])
	}

	// Test translation priority
	translation := config.GetMessageTemplateTranslationWithFallback("validation_failed", "id")
	fmt.Printf("Translation (manual priority): %s\n", translation)

	// Simulate async reload
	fmt.Println("\nSimulating async reload...")

	// Load new config (simulate reload)
	newConfig, err := goresponse.LoadConfig(source)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Copy manual templates to new config
	if config.ManualMessageTemplates != nil {
		newConfig.ManualMessageTemplates = make(map[string]goresponse.MessageTemplate)
		for key, template := range config.ManualMessageTemplates {
			newConfig.ManualMessageTemplates[key] = template
		}
	}

	// Check template after reload
	if template, exists := newConfig.GetMessageTemplate("validation_failed"); exists {
		fmt.Printf("After reload (manual still priority): %s\n", template.Template)
	}

	// Test with template that only exists in async
	if template, exists := newConfig.GetMessageTemplate("welcome_message"); exists {
		fmt.Printf("Async-only template: %s\n", template.Template)
	}
}

// ExampleConfigPrinterUsage shows how to use ConfigPrinter
func ExampleConfigPrinterUsage() {
	fmt.Println("=== Example ConfigPrinter ===")

	// Load config
	source := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	config, err := goresponse.LoadConfig(source)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// Add manual template for demo
	manualTemplate := goresponse.NewMessageTemplateBuilder("custom_error").
		WithTemplate("Custom error occurred: $details").
		WithTranslations(map[string]string{
			"en": "Custom error occurred: $details",
			"id": "Error kustom terjadi: $details",
		}).
		WithCodeMappings(map[string]int{
			"http": 500,
			"grpc": 13,
		}).
		Build()

	config.AddMessageTemplate(manualTemplate)

	// Example 1: Print config with indent (default)
	fmt.Println("1. Print config with indent (default):")
	if err := config.PrintConfig(); err != nil {
		log.Printf("Error printing config: %v", err)
	}

	fmt.Println()

	// Example 2: Print config without indent
	fmt.Println("2. Print config without indent:")
	if err := config.PrintConfigWithIndent(false); err != nil {
		log.Printf("Error printing config: %v", err)
	}

	fmt.Println()

	// Example 3: Export config as string
	fmt.Println("3. Export config as string:")
	jsonStr, err := config.ExportConfig()
	if err != nil {
		log.Printf("Error exporting config: %v", err)
	} else {
		fmt.Printf("Exported config length: %d characters\n", len(jsonStr))
		fmt.Printf("First 100 chars: %s...\n", jsonStr[:100])
	}

	// Example 4: Export config to file
	fmt.Println("\n4. Export config to file:")
	if err := config.ExportConfigToFile("exported_config.json"); err != nil {
		log.Printf("Error exporting to file: %v", err)
	} else {
		fmt.Println("Config exported to exported_config.json")
	}

	// Example 5: Using method chaining
	fmt.Println("\n5. Using method chaining:")
	printer := config.Printer()

	// Print with indent
	if err := printer.WithIndent(true).Print(); err != nil {
		log.Printf("Error: %v", err)
	}

	// Export without indent
	compactJSON, err := printer.WithIndent(false).Export()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Compact JSON length: %d characters\n", len(compactJSON))
	}

	// Export to file with different name
	if err := printer.WithIndent(true).ExportToFile("formatted_config.json"); err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Formatted config exported to formatted_config.json")
	}
}

// ExampleAsyncConfigPrinterUsage shows ConfigPrinter with AsyncConfigManager
func ExampleAsyncConfigPrinterUsage() {
	fmt.Println("=== Example AsyncConfigManager Printer ===")

	// Setup async manager
	source := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	asyncManager := goresponse.NewAsyncConfigManager(source, 30*time.Second)

	if err := asyncManager.Start(); err != nil {
		log.Printf("Error starting async manager: %v", err)
		return
	}
	defer asyncManager.Stop()

	// Add manual template
	manualTemplate := goresponse.NewMessageTemplateBuilder("async_custom").
		WithTemplate("Async custom message: $message").
		WithTranslations(map[string]string{
			"en": "Async custom message: $message",
			"id": "Pesan kustom async: $message",
		}).
		WithCodeMappings(map[string]int{
			"http": 200,
			"grpc": 0,
		}).
		Build()

	asyncManager.AddMessageTemplate(manualTemplate)

	// Print config from async manager
	fmt.Println("Print config from AsyncConfigManager:")
	if err := asyncManager.PrintConfig(); err != nil {
		log.Printf("Error printing config: %v", err)
	}

	// Export config from async manager
	fmt.Println("\nExport config from AsyncConfigManager:")
	jsonStr, err := asyncManager.ExportConfig()
	if err != nil {
		log.Printf("Error exporting config: %v", err)
	} else {
		fmt.Printf("Exported config length: %d characters\n", len(jsonStr))
	}

	// Export to file
	if err := asyncManager.ExportConfigToFile("async_exported_config.json"); err != nil {
		log.Printf("Error exporting to file: %v", err)
	} else {
		fmt.Println("Async config exported to async_exported_config.json")
	}
}

// ExampleAdvancedMessageTemplateUsage shows advanced usage
func ExampleAdvancedMessageTemplateUsage() {
	fmt.Println("=== Example Advanced MessageTemplate Usage ===")

	// Create multiple templates with builder pattern
	templates := []*goresponse.MessageTemplate{
		goresponse.NewMessageTemplateBuilder("success").
			WithTemplate("Operation completed successfully").
			WithTranslation("en", "Operation completed successfully").
			WithTranslation("id", "Operasi berhasil diselesaikan").
			WithCodeMapping("http", 200).
			Build(),

		goresponse.NewMessageTemplateBuilder("error").
			WithTemplate("An error occurred: $details").
			WithTranslation("en", "An error occurred: $details").
			WithTranslation("id", "Terjadi kesalahan: $details").
			WithCodeMapping("http", 500).
			Build(),

		goresponse.NewMessageTemplateBuilder("not_found").
			WithTemplate("Resource not found: $resource").
			WithTranslation("en", "Resource not found: $resource").
			WithTranslation("id", "Resource tidak ditemukan: $resource").
			WithCodeMapping("http", 404).
			Build(),
	}

	// Load config and add templates
	source := goresponse.ConfigSource{Method: "file", Path: "config.json"}
	config, err := goresponse.LoadConfig(source)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	config.AddMessageTemplates(templates...)

	// Simulate usage in application
	fmt.Println("Simulating usage in application:")

	// Success case
	successMsg := config.GetMessageTemplateTranslationWithFallback("success", "en")
	fmt.Printf("Success message: %s\n", successMsg)

	// Error case with parameters
	errorMsg := config.GetMessageTemplateTranslationWithFallback("error", "id")
	fmt.Printf("Error message: %s\n", errorMsg)

	// Not found case
	notFoundMsg := config.GetMessageTemplateTranslationWithFallback("not_found", "en")
	fmt.Printf("Not found message: %s\n", notFoundMsg)

	// Test update template
	fmt.Println("\nUpdate template:")
	updatedTemplate := goresponse.NewMessageTemplateBuilder("success").
		WithTemplate("Operation completed successfully!").
		WithTranslation("en", "Operation completed successfully!").
		WithTranslation("id", "Operasi berhasil diselesaikan!").
		WithCodeMapping("http", 200).
		WithCodeMapping("grpc", 0).
		Build()

	config.UpdateMessageTemplate(updatedTemplate)

	// Test update result
	if template, exists := config.GetMessageTemplate("success"); exists {
		fmt.Printf("Updated template: %s\n", template.Template)
		fmt.Printf("GRPC code: %d\n", template.CodeMappings["grpc"])
	}

	// Test remove template
	fmt.Println("\nRemove template:")
	config.RemoveMessageTemplate("error")

	if _, exists := config.GetMessageTemplate("error"); !exists {
		fmt.Println("Template 'error' successfully removed")
	}
}

// ExampleTranslationSourceUsage shows how to use translation_source
func ExampleTranslationSourceUsage() {
	fmt.Println("=== Example Translation Source ===")

	// Example 1: Using config with translation_source
	fmt.Println("1. Load config with translation_source:")
	fileSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config_separated.json", // File with translation_source
	}

	config, err := goresponse.LoadConfig(fileSource)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	fmt.Printf("Default Language: %s\n", config.GetDefaultLanguage())
	fmt.Printf("Supported Languages: %v\n", config.GetSupportedLanguages())

	// Get translations from translation_source
	translation, exists := config.GetTranslation("en", "success")
	if exists {
		fmt.Printf("English 'success' from translation_source: %s\n", translation)
	}

	translation, exists = config.GetTranslation("id", "success")
	if exists {
		fmt.Printf("Indonesian 'success' from translation_source: %s\n", translation)
	}

	// Get message template
	template, exists := config.GetMessageTemplate("validation_failed")
	if exists {
		fmt.Printf("Validation failed template: %s\n", template.Template)
	}

	fmt.Println()

	// Example 2: Comparison inline vs translation_source
	fmt.Println("2. Comparison inline vs translation_source:")

	// Load config with inline translations
	inlineSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json", // File with inline translations
	}

	inlineConfig, err := goresponse.LoadConfig(inlineSource)
	if err != nil {
		log.Printf("Error loading inline config: %v", err)
		return
	}

	// Load config with translation_source
	separatedConfig, err := goresponse.LoadConfig(fileSource)
	if err != nil {
		log.Printf("Error loading separated config: %v", err)
		return
	}

	// Compare results
	inlineTrans, _ := inlineConfig.GetTranslation("en", "success")
	separatedTrans, _ := separatedConfig.GetTranslation("en", "success")

	fmt.Printf("Inline translation: %s\n", inlineTrans)
	fmt.Printf("Separated translation: %s\n", separatedTrans)
	fmt.Printf("Same result: %v\n", inlineTrans == separatedTrans)

	fmt.Println()

	// Example 3: Using async with translation_source
	fmt.Println("3. Async loading with translation_source:")

	asyncManager := goresponse.NewAsyncConfigManager(fileSource, 30*time.Second)

	// Callback for monitoring
	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		fmt.Printf("Config updated with translation_source at %s\n", time.Now().Format("15:04:05"))
		if newConfig != nil {
			trans, _ := newConfig.GetTranslation("en", "success")
			fmt.Printf("Latest English success: %s\n", trans)
		}
	})

	if err := asyncManager.Start(); err != nil {
		log.Printf("Error starting async manager: %v", err)
		return
	}

	// Simulate usage
	time.Sleep(2 * time.Second)

	translation = asyncManager.GetTranslationWithFallback("en", "success")
	fmt.Printf("Async translation: %s\n", translation)

	asyncManager.Stop()
}

// ExampleMixedTranslationUsage shows mixed inline and translation_source usage
func ExampleMixedTranslationUsage() {
	fmt.Println("=== Example Mixed Translation Usage ===")

	// Create config with inline translations
	config := &goresponse.ResponseConfig{
		DefaultLanguage: "en",
		Languages:       []string{"en", "id"},
		Translations: map[string]map[string]string{
			"en": {
				"success": "Success (inline)",
				"error":   "Error (inline)",
			},
			"id": {
				"success": "Berhasil (inline)",
				"error":   "Error (inline)",
			},
		},
		TranslationSources: map[string]goresponse.TranslationSource{
			"en": {
				Method: "file",
				Path:   "translations/en.json",
			},
			"id": {
				Method: "file",
				Path:   "translations/id.json",
			},
		},
	}

	// Load translation sources (using internal function)
	// Note: In real implementation, this will be called automatically by LoadConfig
	// Here we simulate for demo

	// Test merge result
	fmt.Println("Merge result inline + translation_source:")

	// Key that exists in both inline and translation_source
	translation, exists := config.GetTranslation("en", "success")
	if exists {
		fmt.Printf("English 'success': %s (translation_source override inline)\n", translation)
	}

	// Key that only exists in inline
	translation, exists = config.GetTranslation("en", "error")
	if exists {
		fmt.Printf("English 'error': %s (only from inline)\n", translation)
	}

	// Key that only exists in translation_source
	translation, exists = config.GetTranslation("en", "USER_NOT_FOUND")
	if exists {
		fmt.Printf("English 'USER_NOT_FOUND': %s (only from translation_source)\n", translation)
	}
}

// ExampleAsyncConfigManagerUsage shows how to use AsyncConfigManager
func ExampleAsyncConfigManagerUsage() {
	fmt.Println("=== Example AsyncConfigManager ===")

	// Create async config manager with refresh every 30 seconds
	fileSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	asyncManager := goresponse.NewAsyncConfigManager(fileSource, 30*time.Second)

	// Add callback for configuration changes
	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		fmt.Printf("Config updated! Old default lang: %s, New default lang: %s\n",
			oldConfig.GetDefaultLanguage(), newConfig.GetDefaultLanguage())
	})

	// Start async manager
	if err := asyncManager.Start(); err != nil {
		log.Printf("Error starting async manager: %v", err)
		return
	}

	fmt.Printf("Async manager started. Is running: %v\n", asyncManager.IsRunning())

	// Using configuration
	config := asyncManager.GetConfig()
	if config != nil {
		fmt.Printf("Current default language: %s\n", config.GetDefaultLanguage())

		// Get translation
		translation := asyncManager.GetTranslationWithFallback("en", "success")
		fmt.Printf("English success: %s\n", translation)
	}

	// Simulate running application for several seconds
	fmt.Println("Running for 2 minutes to demonstrate auto refresh...")
	time.Sleep(2 * time.Minute)

	// Stop async manager
	asyncManager.Stop()
	fmt.Printf("Async manager stopped. Is running: %v\n", asyncManager.IsRunning())
}

// ExampleAsyncConfigManagerWithURL shows async loading from URL
func ExampleAsyncConfigManagerWithURL() {
	fmt.Println("=== Example AsyncConfigManager with URL ===")

	// Create async config manager for URL with refresh every 1 minute
	urlSource := goresponse.ConfigSource{
		Method: "url",
		Path:   "https://api.example.com/config.json",
	}

	asyncManager := goresponse.NewAsyncConfigManager(urlSource, 1*time.Minute)

	// Callback for logging changes
	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		fmt.Printf("Config refreshed from URL at %s\n", time.Now().Format(time.RFC3339))
		if oldConfig != nil && newConfig != nil {
			fmt.Printf("Languages changed from %v to %v\n",
				oldConfig.GetSupportedLanguages(), newConfig.GetSupportedLanguages())
		}
	})

	// Start async manager
	if err := asyncManager.Start(); err != nil {
		log.Printf("Error starting async manager with URL: %v", err)
		return
	}

	// Simulate usage
	for i := 0; i < 5; i++ {
		time.Sleep(30 * time.Second)

		config := asyncManager.GetConfig()
		if config != nil {
			fmt.Printf("Iteration %d - Default language: %s\n", i+1, config.GetDefaultLanguage())
		}

		// Check for errors
		if err := asyncManager.GetLastError(); err != nil {
			fmt.Printf("Last error: %v\n", err)
		}
	}

	asyncManager.Stop()
}

// ExampleAsyncConfigManagerAdvanced shows advanced features
func ExampleAsyncConfigManagerAdvanced() {
	fmt.Println("=== Example AsyncConfigManager Advanced ===")

	fileSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	// Create manager with refresh every 10 seconds for demo
	asyncManager := goresponse.NewAsyncConfigManager(fileSource, 10*time.Second)

	// Multiple callbacks
	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		fmt.Printf("Callback 1: Config changed at %s\n", time.Now().Format("15:04:05"))
	})

	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		if oldConfig != nil && newConfig != nil {
			oldLang := oldConfig.GetDefaultLanguage()
			newLang := newConfig.GetDefaultLanguage()
			if oldLang != newLang {
				fmt.Printf("Callback 2: Language changed from %s to %s\n", oldLang, newLang)
			}
		}
	})

	// Start manager
	if err := asyncManager.Start(); err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Simulate interval change
	time.Sleep(15 * time.Second)
	fmt.Println("Changing refresh interval to 5 seconds...")
	asyncManager.UpdateInterval(5 * time.Second)

	// Simulate source change
	time.Sleep(10 * time.Second)
	fmt.Println("Updating source to different file...")
	newSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config_backup.json",
	}
	asyncManager.UpdateSource(newSource)

	// Force refresh
	time.Sleep(5 * time.Second)
	fmt.Println("Forcing refresh...")
	if err := asyncManager.ForceRefresh(); err != nil {
		fmt.Printf("Force refresh error: %v\n", err)
	}

	// Simulate running for 1 minute
	time.Sleep(1 * time.Minute)

	// Cleanup
	asyncManager.RemoveAllCallbacks()
	asyncManager.Stop()
	fmt.Println("Async manager stopped and callbacks removed")
}

// ExampleConfigComparison shows comparison between sync vs async
func ExampleConfigComparison() {
	fmt.Println("=== Comparison Sync vs Async Loading ===")

	fileSource := goresponse.ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	// 1. Sync Loading (One-time)
	fmt.Println("1. Sync Loading (One-time):")
	start := time.Now()
	config, err := goresponse.LoadConfig(fileSource)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("   Loaded in: %v\n", time.Since(start))
	fmt.Printf("   Default language: %s\n", config.GetDefaultLanguage())

	// 2. Async Loading (Auto refresh)
	fmt.Println("\n2. Async Loading (Auto refresh every 5 seconds):")
	asyncManager := goresponse.NewAsyncConfigManager(fileSource, 5*time.Second)

	// Callback for monitoring
	refreshCount := 0
	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		refreshCount++
		fmt.Printf("   Auto refresh #%d at %s\n", refreshCount, time.Now().Format("15:04:05"))
	})

	start = time.Now()
	if err := asyncManager.Start(); err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("   Started in: %v\n", time.Since(start))

	// Simulate usage for 20 seconds
	time.Sleep(20 * time.Second)

	asyncManager.Stop()
	fmt.Printf("   Total auto refreshes: %d\n", refreshCount)
	fmt.Printf("   Total runtime: %v\n", time.Since(start))
}
