# GoResponse Configuration

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Package for managing flexible response configuration with support for multiple sources (files and URLs).

## Table of Contents

- [Installation](#installation)
  - [Requirements](#requirements)
  - [Install](#install)
  - [Import](#import)
- [Quick Start](#quick-start)
  - [1. Create Configuration File](#1-create-configuration-file)
  - [2. Load Configuration](#2-load-configuration)
  - [3. Async Configuration (Optional)](#3-async-configuration-optional)
- [Testing](#testing)
  - [Test Coverage](#test-coverage)
- [File Explanation](#file-explanation)
  - [types.go](#typesgo)
  - [config.go](#configgo)
  - [async_config.go](#asyncconfiggo)
  - [response.go](#responsego)
  - [example_usage.go](#exampleusagego)
  - [example_response.go](#exampleresponsego)
- [Configuration Structure](#configuration-structure)
- [Usage Examples](#usage-examples)
  - [Sync Loading](#sync-loading)
  - [Async Loading](#async-loading)
  - [Translation Source](#translation-source)
  - [MessageTemplateBuilder](#messagetemplatebuilder)
  - [ConfigPrinter](#configprinter)
  - [ResponseBuilder](#responsebuilder)
  - [ResponseManager](#responsemanager)
- [How Translation Source Works](#how-translation-source-works)
- [Message Templates Priority](#message-templates-priority)
- [ConfigPrinter](#configprinter-1)
- [ResponseManager - Automatic Response Generation](#responsemanager---automatic-response-generation)
- [Available Methods](#available-methods)
  - [ConfigManager Methods](#configmanager-methods)
  - [AsyncConfigManager Methods](#asyncconfigmanager-methods)
  - [ResponseBuilder Methods](#responsebuilder-methods)
  - [ResponseBuilder Helper Functions](#responsebuilder-helper-functions)
  - [ResponseConfig Methods (Response Building)](#responseconfig-methods-response-building)
- [Echo Framework Integration](#echo-framework-integration)
- [Service Layer with Error Return Pattern](#service-layer-with-error-return-pattern)
- [ResponseBuilder Examples](#responsebuilder-examples)
- [Error Handling](#error-handling)
- [Extensibility](#extensibility)
- [Sync vs Async Loading Comparison](#sync-vs-async-loading-comparison)
- [License](#license)
- [Contributing](#contributing)
  - [Development Setup](#development-setup)

## File Structure

```
goresponse/
├── types.go              # Shared types (ConfigSource, ResponseConfig, MessageTemplate, ConfigChangeCallback)
├── config.go             # Sync loading functions and ConfigManager
├── async_config.go       # Async loading with AsyncConfigManager
├── response.go           # ResponseBuilder and ResponseManager for standardized API responses
├── example_usage.go      # Usage examples for sync and async
├── example_response.go   # ResponseBuilder and ResponseManager usage examples
├── config.json           # Example configuration file
└── README.md             # Documentation
```

## Features

- ✅ Load configuration from JSON files
- ✅ Load configuration from URLs (HTTP/HTTPS)
- ✅ **Sync Loading** - Load configuration once
- ✅ **Async Loading** - Auto refresh configuration periodically
- ✅ **Translation Source** - Load translations from separate files/URLs
- ✅ **MessageTemplateBuilder** - Method chaining for creating message templates
- ✅ **ConfigPrinter** - Print and export config with method chaining
- ✅ **ResponseBuilder** - Fluent interface for building standardized API responses
- ✅ **ResponseManager** - Automatic response code and message generation
- ✅ Simple and flexible configuration structure
- ✅ Multi-language support
- ✅ Message templates with parameters
- ✅ Fallback translation to default language
- ✅ Callback system for configuration change notifications
- ✅ Thread-safe operations
- ✅ Good error handling

## Installation

### Requirements

- Go 1.21 or higher
- No external dependencies (uses only standard library)

### Install

```bash
go get github.com/risoftinc/goresponse
```

### Import

```go
import "github.com/risoftinc/goresponse"
```

## Quick Start

### 1. Create Configuration File

Create a `config.json` file:

```json
{
  "default_language": "en",
  "languages": ["en", "id"],
  "message_templates": {
    "welcome": {
      "key": "welcome",
      "template": "Welcome $name",
      "code_mappings": {
        "http": 200,
        "grpc": 0
      }
    }
  },
  "translations": {
    "en": {
      "welcome": "Welcome $name"
    },
    "id": {
      "welcome": "Selamat datang $name"
    }
  }
}
```

### 2. Load Configuration

```go
package main

import (
    "fmt"
    "github.com/risoftinc/goresponse"
)

func main() {
    // Load configuration from file
    config, err := goresponse.LoadConfig(goresponse.ConfigSource{
        Method: "file",
        Path:   "config.json",
    })
    if err != nil {
        panic(err)
    }

    // Create response builder
    builder := goresponse.NewResponseBuilder("welcome").
        SetLanguage("en").
        SetProtocol("http").
        SetParam("name", "John").
        SetData("user_id", 123)

    // Build response
    response, err := config.BuildResponse(builder)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Response: %+v\n", response)
    // Output: Response: &{Code:200 Message:Welcome John Data:map[user_id:123] Meta:<nil> Error:<nil> Language:en Protocol:http}
}
```

### 3. Async Configuration (Optional)

```go
// Create async manager with 30-second refresh interval
manager := goresponse.NewAsyncConfigManager(
    goresponse.ConfigSource{
        Method: "file",
        Path:   "config.json",
    },
    30*time.Second,
)

// Start auto-refresh
err := manager.Start()
if err != nil {
    panic(err)
}
defer manager.Stop()

// Use manager
config := manager.GetConfig()
// ... use config as above
```

## Testing

The package includes comprehensive unit tests for all components:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestLoadConfig

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

### Test Coverage

- **config.go** - 27KB, 1051 lines of tests
- **async_config.go** - 28KB, 1212 lines of tests  
- **types.go** - 27KB, 973 lines of tests
- **response.go** - 33KB, 1206 lines of tests

Total: **115KB, 4442 lines** of comprehensive test coverage including:
- Unit tests for all functions and methods
- Edge cases and error conditions
- Thread safety tests
- Performance benchmarks
- JSON marshaling/unmarshaling tests
- Method chaining tests

## File Explanation

### `types.go`
Contains all shared types used by sync and async loading:
- `ConfigSource` - Structure for specifying configuration source
- `ResponseConfig` - Main configuration structure
- `MessageTemplate` - Template for messages
- `TranslationSource` - Structure for separate translation sources
- `ConfigChangeCallback` - Callback type for async loading

### `config.go`
Contains functions and structs for **sync loading** (one-time):
- `LoadConfig()` - Main function for loading configuration
- `ConfigManager` - Manager for sync loading
- Helper functions for file and URL loading

### `async_config.go`
Contains structs and functions for **async loading** (auto refresh):
- `AsyncConfigManager` - Manager for async loading
- Auto refresh with configurable interval
- Thread-safe operations
- Callback system for monitoring

### `response.go`
Contains ResponseBuilder and ResponseManager for standardized API responses:
- `responseBuilder` - Fluent interface for building responses with method chaining
- `Response` - Final standardized response structure
- Context integration for language and protocol detection
- Parameter substitution in message templates
- Error handling and recovery
- `BuildResponse()` - Automatic response generation with proper codes and messages

### `example_usage.go`
Contains usage examples for both types of loading:
- Sync loading examples
- Async loading examples
- Advanced usage patterns
- Comparison examples

### `example_response.go`
Contains comprehensive examples for ResponseBuilder and ResponseManager:
- Basic success and error response examples
- Context integration examples
- Data payload response examples
- Service and handler layer integration
- Multiple language support examples
- Protocol-specific response examples
- Error recovery and parsing examples
- Complex business logic examples
- Async configuration usage examples
- Method chaining patterns

## Configuration Structure

### ConfigSource
```go
type ConfigSource struct {
    Method string `json:"method"` // "file" or "url"
    Path   string `json:"path"`   // file path or URL
}
```

### ResponseConfig
```go
type ResponseConfig struct {
    MessageTemplates   map[string]MessageTemplate   `json:"message_templates"`
    DefaultLanguage    string                       `json:"default_language"`
    Languages          []string                     `json:"languages"`
    Translations       map[string]map[string]string `json:"translations"`        // Inline translations
    TranslationSources map[string]TranslationSource `json:"translation_source"`  // Separate translation sources
}
```

### TranslationSource
```go
type TranslationSource struct {
    Method string `json:"method"` // "file" or "url"
    Path   string `json:"path"`   // file path or URL
}
```

## Usage

### 1. Load from File

```go
package main

import (
    "fmt"
    "log"
    "your-project/utils/goresponse"
)

func main() {
    // Create ConfigSource for file
    source := goresponse.ConfigSource{
        Method: "file",
        Path:   "config.json",
    }
    
    // Load configuration
    config, err := goresponse.LoadConfig(source)
    if err != nil {
        log.Fatal(err)
    }
    
    // Using configuration
    fmt.Printf("Default Language: %s\n", config.GetDefaultLanguage())
    
    // Get translation
    translation, exists := config.GetTranslation("en", "success")
    if exists {
        fmt.Printf("Success message: %s\n", translation)
    }
}
```

### 2. Load from URL

```go
// Create ConfigSource for URL
source := goresponse.ConfigSource{
    Method: "url",
    Path:   "https://api.example.com/config.json",
}

config, err := goresponse.LoadConfig(source)
if err != nil {
    log.Fatal(err)
}
```

### 3. Async Loading with Auto Refresh

```go
// Create AsyncConfigManager with refresh every 5 minutes
source := goresponse.ConfigSource{
    Method: "file",
    Path:   "config.json",
}

asyncManager := goresponse.NewAsyncConfigManager(source, 5*time.Minute)

// Add callback for configuration changes
asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
    fmt.Printf("Config updated! New default language: %s\n", newConfig.GetDefaultLanguage())
})

// Start async manager
if err := asyncManager.Start(); err != nil {
    log.Fatal(err)
}

// Using configuration (thread-safe)
translation := asyncManager.GetTranslationWithFallback("en", "success")
fmt.Printf("Translation: %s\n", translation)

// Stop when done
defer asyncManager.Stop()
```

### 4. Using ConfigManager (Sync)

```go
// Create ConfigManager
source := goresponse.ConfigSource{
    Method: "file",
    Path:   "config.json",
}

manager := goresponse.NewConfigManager(source)

// Load configuration
if err := manager.Load(); err != nil {
    log.Fatal(err)
}

// Using with fallback
translation := manager.GetTranslationWithFallback("fr", "success") // Fallback to default language
fmt.Printf("Translation: %s\n", translation)
```

### 5. Translation Source (Separate Translations)

```go
// Config with translation_source
source := ConfigSource{
    Method: "file",
    Path:   "config_separated.json", // File with translation_source
}

config, err := LoadConfig(source)
if err != nil {
    log.Fatal(err)
}

// Translations will be loaded automatically from separate files
translation := config.GetTranslationWithFallback("en", "success")
fmt.Printf("Translation: %s\n", translation)
```

### 6. MessageTemplateBuilder (Method Chaining)

```go
// Create message template with method chaining
template := goresponse.NewMessageTemplateBuilder("user_created").
    WithTemplate("User $name has been created successfully").
    WithTranslation("en", "User $name has been created successfully").
    WithTranslation("id", "User $name berhasil dibuat").
    WithCodeMapping("http", 201).
    WithCodeMapping("grpc", 0).
    Build()

// Using variadic functions for multiple values
translations := map[string]string{
    "en": "User $name has been created successfully",
    "id": "User $name berhasil dibuat",
    "es": "Usuario $name ha sido creado exitosamente",
    "fr": "Utilisateur $name a été créé avec succès",
}

codeMappings := map[string]int{
    "http": 201,
    "grpc": 0,
    "rest": 201,
    "graphql": 200,
}

template2 := goresponse.NewMessageTemplateBuilder("user_created_v2").
    WithTemplate("User $name has been created successfully").
    WithTranslations(translations).        // Multiple translations at once
    WithCodeMappings(codeMappings).        // Multiple code mappings at once
    Build()

// Add to config (manual priority)
config.AddMessageTemplate(template2)

// Or add multiple templates
config.AddMessageTemplates(template1, template2, template3)

// Using template translation
translation := config.GetMessageTemplateTranslationWithFallback("user_created", "id")
fmt.Printf("Translation: %s\n", translation)
```

### 7. ConfigPrinter (Print & Export)

```go
// Print config to console
config.PrintConfig()                    // With indent (default)
config.PrintConfigWithIndent(false)    // Without indent

// Export config as string
jsonStr, err := config.ExportConfig()

// Export config to file
err := config.ExportConfigToFile("exported_config.json")

// Using method chaining
printer := config.Printer()
printer.WithIndent(true).Print()                           // Print with indent
printer.WithIndent(false).Export()                         // Export without indent
printer.WithIndent(true).ExportToFile("formatted.json")    // Export to file
```

### 8. Load from JSON String

```go
import "encoding/json"

// JSON string
jsonStr := `{"method": "file", "path": "config.json"}`

var source goresponse.ConfigSource
json.Unmarshal([]byte(jsonStr), &source)

config, err := goresponse.LoadConfig(source)
```

## Available Methods

### ResponseConfig Methods

- `GetTranslation(lang, key string) (string, bool)` - Get translation
- `GetMessageTemplate(key string) (*MessageTemplate, bool)` - Get message template
- `GetSupportedLanguages() []string` - List of supported languages
- `GetDefaultLanguage() string` - Default language
- `GetTranslationWithFallback(lang, key string) string` - Translation with fallback
- `AddMessageTemplate(template *MessageTemplate)` - Add message template
- `AddMessageTemplates(templates ...*MessageTemplate)` - Add multiple templates
- `RemoveMessageTemplate(key string)` - Remove message template
- `UpdateMessageTemplate(template *MessageTemplate)` - Update message template
- `GetMessageTemplateTranslation(templateKey, lang string) (string, bool)` - Get translation from template
- `GetMessageTemplateTranslationWithFallback(templateKey, lang string) string` - Template translation with fallback
- `Printer() *ConfigPrinter` - Get ConfigPrinter for print/export
- `PrintConfig() error` - Print config to console (shortcut)
- `PrintConfigWithIndent(useIndent bool) error` - Print config with indent option (shortcut)
- `ExportConfig() (string, error)` - Export config as JSON string (shortcut)
- `ExportConfigToFile(filename string) error` - Export config to file (shortcut)

### ConfigPrinter Methods

- `NewConfigPrinter(config *ResponseConfig) *ConfigPrinter` - Create new printer
- `WithIndent(useIndent bool) *ConfigPrinter` - Set indent option
- `Print() error` - Print config to console
- `Export() (string, error)` - Export config as JSON string
- `ExportToFile(filename string) error` - Export config to file

### MessageTemplateBuilder Methods

- `NewMessageTemplateBuilder(key string) *MessageTemplateBuilder` - Create new builder
- `WithTemplate(template string) *MessageTemplateBuilder` - Set template string
- `WithTranslation(lang, translation string) *MessageTemplateBuilder` - Add translation
- `WithTranslations(translations map[string]string) *MessageTemplateBuilder` - Add multiple translations
- `WithCodeMapping(mappingType string, code int) *MessageTemplateBuilder` - Add code mapping
- `WithCodeMappings(codeMappings map[string]int) *MessageTemplateBuilder` - Add multiple code mappings
- `Build() *MessageTemplate` - Build final template

### ConfigManager Methods (Sync)

- `Load() error` - Load configuration
- `GetConfig() *ResponseConfig` - Get configuration
- `Reload() error` - Reload configuration
- `GetTranslationWithFallback(lang, key string) string` - Translation with fallback

### ResponseBuilder Methods

- `NewResponseBuilder(messageKey string) *responseBuilder` - Create new response builder
- `WithContext(ctx context.Context) *responseBuilder` - Set context and extract language/protocol
- `SetLanguage(language string) *responseBuilder` - Set language manually
- `SetProtocol(protocol string) *responseBuilder` - Set protocol manually
- `SetError(err error) *responseBuilder` - Set error and mark as error response
- `SetParam(key string, value any) *responseBuilder` - Add single parameter
- `SetParams(params map[string]any) *responseBuilder` - Add multiple parameters
- `SetData(key string, value any) *responseBuilder` - Add single data field
- `SetDatas(data map[string]any) *responseBuilder` - Add multiple data fields
- `SetMeta(key string, value any) *responseBuilder` - Add single metadata field
- `SetMetas(meta map[string]any) *responseBuilder` - Add multiple metadata fields
- `Error() string` - Get JSON representation of builder
- `ToError() error` - Convert builder to error type

### ResponseBuilder Helper Functions

- `WithProtocol(ctx context.Context, protocol string) context.Context` - Add protocol to context
- `WithLanguage(ctx context.Context, language string) context.Context` - Add language to context
- `ParseResponseBuilderError(err error) (*responseBuilder, bool)` - Extract builder from error

### ResponseConfig Methods (Response Building)

- `BuildResponse(rb *responseBuilder) (*Response, error)` - Build final response from builder

### AsyncConfigManager Methods

- `Start() error` - Start auto refresh
- `Stop()` - Stop auto refresh
- `GetConfig() *ResponseConfig` - Get current configuration (thread-safe)
- `GetTranslation(lang, key string) (string, bool)` - Get translation (thread-safe)
- `GetTranslationWithFallback(lang, key string) string` - Translation with fallback (thread-safe)
- `GetMessageTemplate(key string) (*MessageTemplate, bool)` - Get message template (thread-safe)
- `GetSupportedLanguages() []string` - List of supported languages (thread-safe)
- `GetDefaultLanguage() string` - Default language (thread-safe)
- `AddCallback(callback ConfigChangeCallback)` - Add callback for changes
- `RemoveAllCallbacks()` - Remove all callbacks
- `IsRunning() bool` - Status whether manager is running
- `GetLastError() error` - Last error that occurred
- `ForceRefresh() error` - Force refresh configuration
- `UpdateSource(newSource ConfigSource)` - Change configuration source
- `UpdateInterval(newInterval time.Duration)` - Change refresh interval
- `AddMessageTemplate(template *MessageTemplate)` - Add message template (thread-safe)
- `AddMessageTemplates(templates ...*MessageTemplate)` - Add multiple templates (thread-safe)
- `RemoveMessageTemplate(key string)` - Remove message template (thread-safe)
- `UpdateMessageTemplate(template *MessageTemplate)` - Update message template (thread-safe)
- `GetMessageTemplateTranslation(templateKey, lang string) (string, bool)` - Get translation from template (thread-safe)
- `GetMessageTemplateTranslationWithFallback(templateKey, lang string) string` - Template translation with fallback (thread-safe)
- `Printer() *ConfigPrinter` - Get ConfigPrinter for print/export (thread-safe)
- `PrintConfig() error` - Print config to console (thread-safe shortcut)
- `PrintConfigWithIndent(useIndent bool) error` - Print config with indent option (thread-safe shortcut)
- `ExportConfig() (string, error)` - Export config as JSON string (thread-safe shortcut)
- `ExportConfigToFile(filename string) error` - Export config to file (thread-safe shortcut)

## Configuration File Examples

### 1. Configuration with Inline Translations

```json
{
  "message_templates": {
    "validation_failed": {
      "key": "validation_failed",
      "template": "Validation failed for field: $field",
      "code_mappings": {
        "http": 422
      }
    }
  },
  "default_language": "en",
  "languages": ["en", "id"],
  "translations": {
    "en": {
      "success": "Success",
      "not_found": "$resource not found",
      "validation_failed": "Validation failed for field: $field"
    },
    "id": {
      "success": "Berhasil",
      "not_found": "$resource tidak ditemukan",
      "validation_failed": "Validasi gagal untuk field: $field"
    }
  }
}
```

### 2. Configuration with Translation Source

```json
{
  "message_templates": {
    "validation_failed": {
      "key": "validation_failed",
      "template": "Validation failed for field: $field",
      "code_mappings": {
        "http": 422
      }
    }
  },
  "default_language": "en",
  "languages": ["en", "id"],
  "translation_source": {
    "en": {
      "method": "file",
      "path": "utils/goresponse/translations/en.json"
    },
    "id": {
      "method": "file",
      "path": "utils/goresponse/translations/id.json"
    }
  }
}
```

### 3. Separate Translation File (en.json)

```json
{
  "success": "Success",
  "not_found": "$resource not found",
  "validation_failed": "Validation failed for field: $field",
  "USER_NOT_FOUND": "User not found",
  "INVALID_CREDENTIALS": "Invalid credentials provided"
}
```

## How Translation Source Works

1. **Load Configuration** - Load main configuration file
2. **Detect Translation Source** - Check if `translation_source` exists
3. **Load Translations** - Load translations from separate files/URLs for each language
4. **Merge** - Combine with existing `translations` (translation_source will override)
5. **Ready to Use** - Configuration ready to use with complete translations

### Merge Priority:
- `translation_source` **overrides** `translations` (if same key exists)
- `translations` still used for keys not in `translation_source`

## Message Templates Priority

### Manual Templates (High Priority)
Templates added manually have the highest priority and won't be overridden by async reload:

```go
// Manual template (high priority)
manualTemplate := NewMessageTemplateBuilder("validation_failed").
    WithTemplate("Custom validation failed (MANUAL)").
    WithTranslations(map[string]string{
        "en": "Custom validation failed (MANUAL)",
        "id": "Validasi kustom gagal (MANUAL)",
    }).
    WithCodeMappings(map[string]int{
        "http": 400,
        "grpc": 3,
    }).
    Build()

config.AddMessageTemplate(manualTemplate)
```

### Async Templates (Low Priority)
Templates from config file/URL will be overridden by manual templates:

```go
// Template from config file (low priority)
// Will be overridden if manual template with same key exists
```

### Priority Order:
1. **Manual Templates** - Templates added via `AddMessageTemplate()`
2. **Async Templates** - Templates from config file/URL
3. **Fallback** - Template string or key

### Advantages of Manual Priority:
- ✅ **Override Protection** - Manual templates won't be overridden by async reload
- ✅ **Runtime Customization** - Can add/update templates at runtime
- ✅ **Hot Fix** - Can fix templates without restarting application
- ✅ **A/B Testing** - Can test different templates dynamically

## ConfigPrinter - Print & Export

### ConfigPrinter Features:
- ✅ **Print to Console** - Print config to console with JSON format
- ✅ **Export to String** - Get config as JSON string
- ✅ **Export to File** - Save config to JSON file
- ✅ **Method Chaining** - With `WithIndent()` for indent control
- ✅ **Thread-Safe** - Safe to use with AsyncConfigManager
- ✅ **Manual Priority** - Export merges manual and async templates

### Usage Examples:

```go
// Print config to console
config.PrintConfig()                    // With indent (default)
config.PrintConfigWithIndent(false)    // Without indent

// Export config
jsonStr, err := config.ExportConfig()
err := config.ExportConfigToFile("config.json")

// Method chaining
printer := config.Printer()
printer.WithIndent(true).Print()                           // Print with indent
printer.WithIndent(false).Export()                         // Export without indent
printer.WithIndent(true).ExportToFile("formatted.json")    // Export to file

// With AsyncConfigManager (thread-safe)
asyncManager.PrintConfig()
asyncManager.ExportConfigToFile("async_config.json")
```

### Output Format:
ConfigPrinter will export config in the same format as the original config file, with:
- Manual templates have priority (override async templates)
- All translations and code mappings are preserved
- Consistent and readable JSON format

## ResponseManager - Automatic Response Generation

### ResponseManager Features:
- ✅ **Automatic Response Code** - Automatically sets appropriate response codes based on protocol
- ✅ **Language-Aware Messages** - Generates messages in the correct language
- ✅ **Template-Based** - Uses message templates from configuration
- ✅ **Parameter Substitution** - Supports dynamic parameter replacement in messages
- ✅ **Context Integration** - Automatically detects language and protocol from context
- ✅ **Protocol Support** - Works with HTTP, gRPC, REST, GraphQL, and other protocols
- ✅ **Error Handling** - Built-in error response generation
- ✅ **Method Chaining** - Fluent interface for easy response building

### ResponseBuilder Usage:

```go
// Basic usage with automatic code and message generation
func handleCreateUser(ctx context.Context, name, email string) (*goresponse.Response, error) {
    // Create response builder
    builder := goresponse.NewResponseBuilder("user_created")
    
    // Set parameters for template substitution
    builder.SetParam("name", name)
    builder.SetParam("email", email)
    
    // Set context (automatically extracts language and protocol)
    builder.WithContext(ctx)
    
    // Build response with automatic code and message
    return config.BuildResponse(builder)
}
```

### Advanced Usage with Manual Settings:

```go
// Manual language and protocol setting
func handleValidationError(field, value string) (*goresponse.Response, error) {
    builder := goresponse.NewResponseBuilder("validation_failed")
    
    // Set parameters
    builder.SetParam("field", field)
    builder.SetParam("value", value)
    
    // Manual settings (overrides context)
    builder.SetLanguage("id")
    builder.SetProtocol("http")
    
    // Build response
    return config.BuildResponse(builder)
}
```

### Context Integration:

```go
// Create context with language and protocol
ctx := context.Background()
ctx = goresponse.WithLanguage(ctx, "en")
ctx = goresponse.WithProtocol(ctx, "http")

// Response builder will automatically use these settings
builder := goresponse.NewResponseBuilder("success")
builder.WithContext(ctx)
```

### Error Response Generation:

```go
// Error response with automatic error handling
func handleError(err error) (*goresponse.Response, error) {
    builder := goresponse.NewResponseBuilder("internal_error")
    
    // Set error information
    builder.SetError(err)
    builder.SetParam("service", "user_service")
    
    // Build error response
    return config.BuildResponse(builder)
}
```

### Data Payload:

```go
// Response with data payload
func handleUserList(users []User) (*goresponse.Response, error) {
    builder := goresponse.NewResponseBuilder("users_retrieved")
    
    // Set data payload
    builder.SetData("users", users)
    builder.SetData("count", len(users))
    
    // Set parameters
    builder.SetParam("count", len(users))
    
    return config.BuildResponse(builder)
}
```

### Response Structure:

The generated response follows this structure:

```go
type Response struct {
    Code     int         `json:"code"`              // Automatic response code
    Message  string      `json:"message"`           // Translated message
    Data     interface{} `json:"data,omitempty"`    // Response payload
    Meta     interface{} `json:"meta,omitempty"`    // Additional metadata
    Error    error       `json:"error,omitempty"`   // Error details (if error)
    Language string      `json:"-"`                 // Language used
    Protocol string      `json:"-"`                 // Protocol used
}
```

### How ResponseManager Works:

1. **Template Resolution** - Finds message template by key from configuration
2. **Language Detection** - Determines language from context or manual setting
3. **Message Translation** - Gets translated message for the language
4. **Parameter Substitution** - Replaces `$param` placeholders with actual values
5. **Code Mapping** - Maps response code based on protocol (HTTP: 200, gRPC: 0, etc.)
6. **Response Building** - Constructs final Response with all information

### Protocol Code Mapping:

```go
// Automatic code mapping based on protocol
switch protocol {
case "http":
    code = template.CodeMappings["http"]  // e.g., 200, 201, 400, 500
case "grpc":
    code = template.CodeMappings["grpc"]  // e.g., 0, 3, 5, 13
case "rest":
    code = template.CodeMappings["rest"]  // e.g., 200, 201, 400, 500
case "graphql":
    code = template.CodeMappings["graphql"] // e.g., 200, 400, 500
}
```

### Service-Level Integration:

```go
// Service layer - returns error
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // Business logic...
    if err := s.validateUser(req); err != nil {
        // Return ResponseBuilder as error
        builder := goresponse.NewResponseBuilder("validation_failed")
        builder.SetParam("field", "email")
        builder.SetParam("value", req.Email)
        return builder.ToError()
    }
    
    // Success case
    return nil
}

// Handler layer - converts error to response
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Call service
    err := h.userService.CreateUser(r.Context(), req)
    
    if err != nil {
        // Check if it's a ResponseBuilder error
        if builder, ok := goresponse.ParseResponseBuilderError(err); ok {
            // Build proper response
            response, _ := h.config.BuildResponse(builder)
            h.writeResponse(w, response)
            return
        }
        
        // Handle other errors
        h.writeError(w, err)
        return
    }
    
    // Success response
    builder := goresponse.NewResponseBuilder("user_created")
    builder.SetParam("name", req.Name)
    response, _ := h.config.BuildResponse(builder)
    h.writeResponse(w, response)
}
```

### Advantages of ResponseManager:

- ✅ **Consistency** - All responses follow the same structure
- ✅ **Localization** - Automatic language detection and translation
- ✅ **Protocol Agnostic** - Works with any protocol
- ✅ **Template-Driven** - Centralized message management
- ✅ **Parameter Support** - Dynamic message customization
- ✅ **Error Integration** - Seamless error handling
- ✅ **Context Aware** - Automatic language/protocol detection
- ✅ **Type Safe** - Compile-time safety with Go types

### ResponseBuilder Examples

For comprehensive examples of ResponseBuilder and ResponseManager usage, see `example_response.go`:

```go
// Basic usage
func ExampleResponseUsage() {
    // Load configuration
    config, err := LoadConfig(source)
    if err != nil {
        log.Fatal(err)
    }

    // Basic success response
    builder := NewResponseBuilder("user_created")
    builder.SetParam("name", "John Doe")
    builder.SetParam("email", "john@example.com")
    
    response, err := config.BuildResponse(builder)
    // response.Code = 201, response.Message = "User John Doe has been created successfully"
}

// Service layer integration
func (s *UserService) CreateUser(ctx context.Context, name, email string) error {
    if !isValidEmail(email) {
        builder := NewResponseBuilder("validation_failed")
        builder.SetParam("field", "email")
        builder.SetParam("value", email)
        return builder.ToError()
    }
    return nil
}

// Handler layer integration
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    err := h.userService.CreateUser(r.Context(), req.Name, req.Email)
    if err != nil {
        if builder, ok := ParseResponseBuilderError(err); ok {
            response, _ := h.config.BuildResponse(builder)
            h.writeResponse(w, response)
            return
        }
    }
    // Handle success...
}

// Method chaining with metadata
response, err := config.BuildResponse(
    NewResponseBuilder("data_retrieved").
        SetParams(map[string]any{
            "table": "users",
            "count": 100,
        }).
        SetDatas(map[string]any{
            "users": []User{...},
        }).
        SetMetas(map[string]any{
            "request_id": "req-123",
            "processing_time": "150ms",
            "cache_hit": true,
        }).
        SetLanguage("en").
        SetProtocol("http"),
)
```

### Available Example Functions

- `ExampleResponseUsage()` - Basic ResponseBuilder examples
- `ExampleAsyncResponseUsage()` - Async configuration examples
- `ExampleResponseBuilderChaining()` - Method chaining examples

### Echo Framework Integration

For Echo framework integration with automatic language detection from headers:

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "your-project/utils/goresponse"
)

// EchoMiddleware creates middleware for automatic context setup
func EchoMiddleware(config *goresponse.ResponseConfig) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Extract language from Content-Language header
            language := c.Request().Header.Get("Content-Language")
            if language == "" {
                language = "en" // Default language
            }
            
            // Extract protocol from request
            protocol := "http"
            if c.Request().TLS != nil {
                protocol = "https"
            }
            
            // Create context with language and protocol
            ctx := c.Request().Context()
            ctx = goresponse.WithLanguage(ctx, language)
            ctx = goresponse.WithProtocol(ctx, protocol)
            
            // Update request context
            c.SetRequest(c.Request().WithContext(ctx))
            
            return next(c)
        }
    }
}

// ResponseHelper provides helper methods for Echo handlers
type ResponseHelper struct {
    config *goresponse.ResponseConfig
}

func NewResponseHelper(config *goresponse.ResponseConfig) *ResponseHelper {
    return &ResponseHelper{config: config}
}

// Success sends a success response
func (h *ResponseHelper) Success(c echo.Context, messageKey string, params map[string]any, data map[string]any) error {
    builder := goresponse.NewResponseBuilder(messageKey)
    builder.WithContext(c.Request().Context())
    
    if params != nil {
        builder.SetParams(params)
    }
    
    if data != nil {
        builder.SetDatas(data)
    }
    
    response, err := h.config.BuildResponse(builder)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to build response",
        })
    }
    
    return c.JSON(response.Code, response)
}

// Error sends an error response
func (h *ResponseHelper) Error(c echo.Context, messageKey string, params map[string]any, err error) error {
    builder := goresponse.NewResponseBuilder(messageKey)
    builder.WithContext(c.Request().Context())
    
    if params != nil {
        builder.SetParams(params)
    }
    
    if err != nil {
        builder.SetError(err)
    }
    
    response, buildErr := h.config.BuildResponse(builder)
    if buildErr != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to build response",
        })
    }
    
    return c.JSON(response.Code, response)
}

// Handler example
type UserHandler struct {
    responseHelper *ResponseHelper
}

func NewUserHandler(config *goresponse.ResponseConfig) *UserHandler {
    return &UserHandler{
        responseHelper: NewResponseHelper(config),
    }
}

// CreateUser handler
func (h *UserHandler) CreateUser(c echo.Context) error {
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }
    
    if err := c.Bind(&req); err != nil {
        return h.responseHelper.Error(c, "validation_failed", map[string]any{
            "field": "request_body",
            "value": "invalid",
        }, err)
    }
    
    // Validation
    if req.Email == "" {
        return h.responseHelper.Error(c, "validation_failed", map[string]any{
            "field": "email",
            "value": req.Email,
            "rule":  "required",
        }, nil)
    }
    
    if req.Age < 18 {
        return h.responseHelper.Error(c, "validation_failed", map[string]any{
            "field": "age",
            "value": req.Age,
            "rule":  "must be at least 18 years old",
        }, nil)
    }
    
    // Business logic here...
    
    // Success response
    return h.responseHelper.Success(c, "user_created", map[string]any{
        "name":  req.Name,
        "email": req.Email,
    }, map[string]any{
        "user_id": "12345",
        "created_at": "2024-01-01T00:00:00Z",
    })
}

// GetUsers handler
func (h *UserHandler) GetUsers(c echo.Context) error {
    // Simulate user data
    users := []map[string]any{
        {"id": 1, "name": "Alice", "email": "alice@example.com"},
        {"id": 2, "name": "Bob", "email": "bob@example.com"},
    }
    
    return h.responseHelper.Success(c, "users_retrieved", map[string]any{
        "count": len(users),
    }, map[string]any{
        "users": users,
        "pagination": map[string]any{
            "page": 1,
            "size": 10,
            "total": len(users),
        },
    })
}

// Main application setup
func main() {
    // Load configuration
    source := goresponse.ConfigSource{
        Method: "file",
        Path:   "config.json",
    }
    
    config, err := goresponse.LoadConfig(source)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create Echo instance
    e := echo.New()
    
    // Add middleware
    e.Use(EchoMiddleware(config))
    
    // Create handlers
    userHandler := NewUserHandler(config)
    
    // Routes
    e.POST("/users", userHandler.CreateUser)
    e.GET("/users", userHandler.GetUsers)
    
    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}
```

### Advanced Echo Integration with Custom Headers

For more advanced header handling:

```go
// AdvancedEchoMiddleware with custom header handling
func AdvancedEchoMiddleware(config *goresponse.ResponseConfig) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Extract language from various headers
            language := extractLanguage(c)
            
            // Extract protocol and version
            protocol := extractProtocol(c)
            
            // Extract additional context
            ctx := c.Request().Context()
            ctx = goresponse.WithLanguage(ctx, language)
            ctx = goresponse.WithProtocol(ctx, protocol)
            
            // Add custom context values
            ctx = context.WithValue(ctx, "request_id", c.Request().Header.Get("X-Request-ID"))
            ctx = context.WithValue(ctx, "user_agent", c.Request().UserAgent())
            ctx = context.WithValue(ctx, "ip_address", c.RealIP())
            
            c.SetRequest(c.Request().WithContext(ctx))
            
            return next(c)
        }
    }
}

func extractLanguage(c echo.Context) string {
    // Priority order for language detection
    if lang := c.Request().Header.Get("Content-Language"); lang != "" {
        return lang
    }
    if lang := c.Request().Header.Get("Accept-Language"); lang != "" {
        // Parse Accept-Language header (e.g., "en-US,en;q=0.9,id;q=0.8")
        if parts := strings.Split(lang, ","); len(parts) > 0 {
            if langCode := strings.Split(parts[0], ";")[0]; langCode != "" {
                return strings.Split(langCode, "-")[0] // Extract "en" from "en-US"
            }
        }
    }
    if lang := c.QueryParam("lang"); lang != "" {
        return lang
    }
    return "en" // Default
}

func extractProtocol(c echo.Context) string {
    if c.Request().TLS != nil {
        return "https"
    }
    if proto := c.Request().Header.Get("X-Forwarded-Proto"); proto != "" {
        return proto
    }
    return "http"
}
```

### Error Handling in Echo

```go
// Custom error handler for Echo
func CustomErrorHandler(config *goresponse.ResponseConfig) echo.HTTPErrorHandler {
    return func(err error, c echo.Context) {
        // Check if it's a ResponseBuilder error
        if builder, ok := goresponse.ParseResponseBuilderError(err); ok {
            response, buildErr := config.BuildResponse(builder)
            if buildErr != nil {
                c.JSON(http.StatusInternalServerError, map[string]string{
                    "error": "Failed to build response",
                })
                return
            }
            c.JSON(response.Code, response)
            return
        }
        
        // Handle Echo errors
        if he, ok := err.(*echo.HTTPError); ok {
            c.JSON(he.Code, map[string]interface{}{
                "code":    he.Code,
                "message": he.Message,
            })
            return
        }
        
        // Default error
        c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Internal server error",
        })
    }
}

// Usage in main
func main() {
    e := echo.New()
    
    // Set custom error handler
    e.HTTPErrorHandler = CustomErrorHandler(config)
    
    // ... rest of setup
}
```

### Service Layer with Error Return Pattern

For a more clean separation of concerns, you can set parameters at service level and return errors:

```go
// Service layer - sets parameters and returns ResponseBuilder as error
type UserService struct {
    config *goresponse.ResponseConfig
}

func NewUserService(config *goresponse.ResponseConfig) *UserService {
    return &UserService{config: config}
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // Validation with detailed parameters
    if req.Email == "" {
        builder := goresponse.NewResponseBuilder("validation_failed")
        builder.SetParam("field", "email")
        builder.SetParam("value", req.Email)
        builder.SetParam("rule", "required")
        builder.SetParam("message", "Email is required")
        return builder.ToError()
    }
    
    if !isValidEmail(req.Email) {
        builder := goresponse.NewResponseBuilder("validation_failed")
        builder.SetParam("field", "email")
        builder.SetParam("value", req.Email)
        builder.SetParam("rule", "email_format")
        builder.SetParam("message", "Invalid email format")
        return builder.ToError()
    }
    
    if req.Age < 18 {
        builder := goresponse.NewResponseBuilder("validation_failed")
        builder.SetParam("field", "age")
        builder.SetParam("value", req.Age)
        builder.SetParam("rule", "minimum_age")
        builder.SetParam("message", "Must be at least 18 years old")
        return builder.ToError()
    }
    
    // Business logic validation
    if s.isEmailExists(req.Email) {
        builder := goresponse.NewResponseBuilder("email_already_exists")
        builder.SetParam("email", req.Email)
        builder.SetParam("message", "Email already registered")
        return builder.ToError()
    }
    
    // Simulate database error
    if err := s.saveUser(req); err != nil {
        builder := goresponse.NewResponseBuilder("database_error")
        builder.SetParam("operation", "create_user")
        builder.SetParam("table", "users")
        builder.SetError(err)
        return builder.ToError()
    }
    
    // Success case - no error
    return nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) error {
    user, err := s.findUserByID(userID)
    if err != nil {
        builder := goresponse.NewResponseBuilder("user_not_found")
        builder.SetParam("user_id", userID)
        builder.SetParam("message", "User not found")
        builder.SetError(err)
        return builder.ToError()
    }
    
    // Store user data in context for handler to use
    ctx = context.WithValue(ctx, "user_data", user)
    return nil
}

// Handler layer - converts service errors to responses
type UserHandler struct {
    userService *UserService
    config      *goresponse.ResponseConfig
}

func NewUserHandler(config *goresponse.ResponseConfig) *UserHandler {
    return &UserHandler{
        userService: NewUserService(config),
        config:      config,
    }
}

func (h *UserHandler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        // Handle binding error
        builder := goresponse.NewResponseBuilder("validation_failed")
        builder.SetParam("field", "request_body")
        builder.SetParam("value", "invalid_json")
        builder.SetParam("rule", "valid_json")
        builder.SetError(err)
        
        response, _ := h.config.BuildResponse(builder)
        return c.JSON(response.Code, response)
    }
    
    // Call service
    err := h.userService.CreateUser(c.Request().Context(), &req)
    if err != nil {
        // Check if it's a ResponseBuilder error
        if builder, ok := goresponse.ParseResponseBuilderError(err); ok {
            // Add context information to builder
            builder.WithContext(c.Request().Context())
            
            response, buildErr := h.config.BuildResponse(builder)
            if buildErr != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{
                    "error": "Failed to build response",
                })
            }
            
            return c.JSON(response.Code, response)
        }
        
        // Handle other errors
        builder := goresponse.NewResponseBuilder("internal_error")
        builder.WithContext(c.Request().Context())
        builder.SetError(err)
        
        response, _ := h.config.BuildResponse(builder)
        return c.JSON(response.Code, response)
    }
    
    // Success response
    builder := goresponse.NewResponseBuilder("user_created")
    builder.WithContext(c.Request().Context())
    builder.SetParam("name", req.Name)
    builder.SetParam("email", req.Email)
    builder.SetData("user_id", "12345")
    builder.SetData("created_at", time.Now().Format(time.RFC3339))
    builder.SetMeta("request_id", c.Request().Header.Get("X-Request-ID"))
    builder.SetMeta("processing_time", "50ms")
    
    response, err := h.config.BuildResponse(builder)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to build response",
        })
    }
    
    return c.JSON(response.Code, response)
}

func (h *UserHandler) GetUser(c echo.Context) error {
    userID := c.Param("id")
    
    // Call service
    err := h.userService.GetUser(c.Request().Context(), userID)
    if err != nil {
        // Check if it's a ResponseBuilder error
        if builder, ok := goresponse.ParseResponseBuilderError(err); ok {
            builder.WithContext(c.Request().Context())
            
            response, buildErr := h.config.BuildResponse(builder)
            if buildErr != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{
                    "error": "Failed to build response",
                })
            }
            
            return c.JSON(response.Code, response)
        }
        
        // Handle other errors
        builder := goresponse.NewResponseBuilder("internal_error")
        builder.WithContext(c.Request().Context())
        builder.SetError(err)
        
        response, _ := h.config.BuildResponse(builder)
        return c.JSON(response.Code, response)
    }
    
    // Get user data from context
    userData := c.Request().Context().Value("user_data")
    
    // Success response
    builder := goresponse.NewResponseBuilder("user_retrieved")
    builder.WithContext(c.Request().Context())
    builder.SetParam("user_id", userID)
    builder.SetData("user", userData)
    builder.SetMeta("cache_hit", false)
    builder.SetMeta("query_time", "25ms")
    
    response, err := h.config.BuildResponse(builder)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to build response",
        })
    }
    
    return c.JSON(response.Code, response)
}

// Helper function for error conversion in Echo
func ConvertServiceErrorToResponse(err error, c echo.Context, config *goresponse.ResponseConfig) error {
    if err == nil {
        return nil
    }
    
    // Check if it's a ResponseBuilder error
    if builder, ok := goresponse.ParseResponseBuilderError(err); ok {
        builder.WithContext(c.Request().Context())
        
        response, buildErr := config.BuildResponse(builder)
        if buildErr != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "Failed to build response",
            })
        }
        
        return c.JSON(response.Code, response)
    }
    
    // Handle other errors
    builder := goresponse.NewResponseBuilder("internal_error")
    builder.WithContext(c.Request().Context())
    builder.SetError(err)
    
    response, _ := config.BuildResponse(builder)
    return c.JSON(response.Code, response)
}

// Simplified handler using helper
func (h *UserHandler) CreateUserSimplified(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        builder := goresponse.NewResponseBuilder("validation_failed")
        builder.SetParam("field", "request_body")
        builder.SetError(err)
        return ConvertServiceErrorToResponse(builder.ToError(), c, h.config)
    }
    
    // Call service
    err := h.userService.CreateUser(c.Request().Context(), &req)
    if err != nil {
        return ConvertServiceErrorToResponse(err, c, h.config)
    }
    
    // Success response
    builder := goresponse.NewResponseBuilder("user_created")
    builder.WithContext(c.Request().Context())
    builder.SetParam("name", req.Name)
    builder.SetParam("email", req.Email)
    builder.SetMeta("version", "v1.0")
    builder.SetMeta("endpoint", "POST /users")
    
    response, err := h.config.BuildResponse(builder)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to build response",
        })
    }
    
    return c.JSON(response.Code, response)
}

// Mock functions for service
func (s *UserService) isEmailExists(email string) bool {
    // Mock implementation
    return email == "existing@example.com"
}

func (s *UserService) saveUser(req *CreateUserRequest) error {
    // Mock implementation
    if req.Email == "error@example.com" {
        return fmt.Errorf("database connection failed")
    }
    return nil
}

func (s *UserService) findUserByID(userID string) (map[string]any, error) {
    // Mock implementation
    if userID == "notfound" {
        return nil, fmt.Errorf("user not found")
    }
    return map[string]any{
        "id":    userID,
        "name":  "John Doe",
        "email": "john@example.com",
    }, nil
}

func isValidEmail(email string) bool {
    return len(email) > 0 && 
           email != "invalid-email" && 
           email != "error@example.com"
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}
```

### Testing Echo Integration

```go
func TestEchoIntegration(t *testing.T) {
    // Load test configuration
    config, err := goresponse.LoadConfig(goresponse.ConfigSource{
        Method: "file",
        Path:   "test_config.json",
    })
    assert.NoError(t, err)
    
    // Create Echo instance
    e := echo.New()
    e.Use(EchoMiddleware(config))
    
    // Create handler
    handler := NewUserHandler(config)
    e.POST("/users", handler.CreateUser)
    
    // Test request with Content-Language header
    req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{
        "name": "John Doe",
        "email": "john@example.com",
        "age": 25
    }`))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Content-Language", "id") // Indonesian
    
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    
    // Execute
    err = handler.CreateUser(c)
    assert.NoError(t, err)
    
    // Verify response
    assert.Equal(t, http.StatusCreated, rec.Code)
    
    var response goresponse.Response
    err = json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "id", response.Language)
    assert.Contains(t, response.Message, "John Doe")
}
```

## Error Handling

All functions return errors that can be handled properly:

```go
config, err := goresponse.LoadConfig(source)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "failed to open config file"):
        // Handle file not found
    case strings.Contains(err.Error(), "failed to fetch config from URL"):
        // Handle network error
    case strings.Contains(err.Error(), "failed to unmarshal config"):
        // Handle JSON parsing error
    default:
        // Handle other errors
    }
}
```

## Extensibility

This structure is easy to extend:

1. **Database Source**: Add "database" method in ConfigSource
2. **Environment Variables**: Add "env" method
3. **Remote Config**: Add "remote" method with authentication
4. **Caching**: Add caching layer in ConfigManager

## Sync vs Async Loading Comparison

### Sync Loading (One-time)
- ✅ **Simple** - Load once, use continuously
- ✅ **Lightweight** - No background process overhead
- ✅ **Suitable for** - Static configuration, simple applications
- ❌ **Not real-time** - Need restart to update configuration

### Async Loading (Auto Refresh)
- ✅ **Real-time** - Configuration updates automatically
- ✅ **Flexible** - Refresh interval can be configured
- ✅ **Monitoring** - Callback for tracking changes
- ✅ **Thread-safe** - Safe to use in concurrent environment
- ✅ **Suitable for** - Microservices, production applications
- ❌ **More complex** - Has background process overhead
- ❌ **Resource usage** - Uses goroutines and memory

### When to Use Which?

**Use Sync Loading if:**
- Configuration rarely changes
- Simple application
- Want maximum performance
- Don't need real-time updates

**Use Async Loading if:**
- Configuration changes frequently
- Production/microservices application
- Need monitoring of configuration changes
- Want zero-downtime updates

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/risoftinc/goresponse.git
cd goresponse

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Check for linting issues
golangci-lint run
```

Copyright (c) 2025 Risoftinc.