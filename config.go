package goresponse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// LoadConfig loads configuration based on ConfigSource (sync loading)
func LoadConfig(source ConfigSource) (*ResponseConfig, error) {
	var data []byte
	var err error

	switch strings.ToLower(source.Method) {
	case "file":
		data, err = loadFromFile(source.Path)
	case "url":
		data, err = loadFromURL(source.Path)
	default:
		return nil, fmt.Errorf("unsupported method: %s. Supported methods: file, url", source.Method)
	}

	if err != nil {
		return nil, err
	}

	var config ResponseConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load translations from translation_source if any
	if err := loadTranslationSources(&config); err != nil {
		return nil, fmt.Errorf("failed to load translation sources: %w", err)
	}

	return &config, nil
}

// loadFromFile loads data from file
func loadFromFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return data, nil
}

// loadFromURL loads data from URL
func loadFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch config, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// loadTranslationSources loads translations from translation_source
func loadTranslationSources(config *ResponseConfig) error {
	if len(config.TranslationSources) == 0 {
		return nil // No translation sources
	}

	// Initialize translations map if not exists
	if config.Translations == nil {
		config.Translations = make(map[string]map[string]string)
	}

	// Load translations for each language
	for lang, source := range config.TranslationSources {
		translations, err := loadTranslationFromSource(source)
		if err != nil {
			return fmt.Errorf("failed to load translations for language %s: %w", lang, err)
		}

		// Merge with existing translations (translation_source will override translations)
		if config.Translations[lang] == nil {
			config.Translations[lang] = make(map[string]string)
		}

		// Update/override with translations from source
		for key, value := range translations {
			config.Translations[lang][key] = value
		}
	}

	return nil
}

// loadTranslationFromSource loads translations from source (file or URL)
func loadTranslationFromSource(source TranslationSource) (map[string]string, error) {
	var data []byte
	var err error

	switch strings.ToLower(source.Method) {
	case "file":
		data, err = loadFromFile(source.Path)
	case "url":
		data, err = loadFromURL(source.Path)
	default:
		return nil, fmt.Errorf("unsupported translation source method: %s. Supported methods: file, url", source.Method)
	}

	if err != nil {
		return nil, err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal translations: %w", err)
	}

	return translations, nil
}

// GetTranslation gets translation based on language and key
func (c *ResponseConfig) GetTranslation(lang, key string) (string, bool) {
	if translations, exists := c.Translations[lang]; exists {
		if translation, exists := translations[key]; exists {
			return translation, true
		}
	}
	return "", false
}

// GetMessageTemplate gets message template based on key (manual priority > async)
func (c *ResponseConfig) GetMessageTemplate(key string) (*MessageTemplate, bool) {
	// Priority 1: Manual templates (added manually)
	if c.ManualMessageTemplates != nil {
		if template, exists := c.ManualMessageTemplates[key]; exists {
			return &template, true
		}
	}

	// Priority 2: Async templates (from config file/URL)
	if c.MessageTemplates != nil {
		if template, exists := c.MessageTemplates[key]; exists {
			return &template, true
		}
	}

	return nil, false
}

// GetSupportedLanguages returns list of supported languages
func (c *ResponseConfig) GetSupportedLanguages() []string {
	return c.Languages
}

// GetDefaultLanguage returns default language
func (c *ResponseConfig) GetDefaultLanguage() string {
	return c.DefaultLanguage
}

// GetTranslationWithFallback gets translation with fallback to default language
func (c *ResponseConfig) GetTranslationWithFallback(lang, key string) string {
	// Try requested language
	if translation, exists := c.GetTranslation(lang, key); exists {
		return translation
	}

	// Fallback to default language
	if translation, exists := c.GetTranslation(c.GetDefaultLanguage(), key); exists {
		return translation
	}

	// Fallback to key itself if not found
	return key
}

// ConfigManager for managing configuration synchronously
type ConfigManager struct {
	config *ResponseConfig
	source ConfigSource
}

// NewConfigManager creates ConfigManager instance
func NewConfigManager(source ConfigSource) *ConfigManager {
	return &ConfigManager{
		source: source,
	}
}

// Load loads configuration using source
func (cm *ConfigManager) Load() error {
	config, err := LoadConfig(cm.source)
	if err != nil {
		return err
	}
	cm.config = config
	return nil
}

// GetConfig returns loaded configuration
func (cm *ConfigManager) GetConfig() *ResponseConfig {
	return cm.config
}

// Reload reloads configuration
func (cm *ConfigManager) Reload() error {
	return cm.Load()
}

// GetTranslationWithFallback gets translation with fallback to default language
func (cm *ConfigManager) GetTranslationWithFallback(lang, key string) string {
	if cm.config == nil {
		return ""
	}
	return cm.config.GetTranslationWithFallback(lang, key)
}

// AddMessageTemplate adds message template to ResponseConfig (manual priority)
func (c *ResponseConfig) AddMessageTemplate(template *MessageTemplate) {
	if c.ManualMessageTemplates == nil {
		c.ManualMessageTemplates = make(map[string]MessageTemplate)
	}
	c.ManualMessageTemplates[template.Key] = *template
}

// AddMessageTemplates adds multiple message templates
func (c *ResponseConfig) AddMessageTemplates(templates ...*MessageTemplate) {
	for _, template := range templates {
		c.AddMessageTemplate(template)
	}
}

// RemoveMessageTemplate removes message template (from manual templates)
func (c *ResponseConfig) RemoveMessageTemplate(key string) {
	if c.ManualMessageTemplates != nil {
		delete(c.ManualMessageTemplates, key)
	}
}

// UpdateMessageTemplate updates existing message template (manual priority)
func (c *ResponseConfig) UpdateMessageTemplate(template *MessageTemplate) {
	c.AddMessageTemplate(template) // Add will override if key already exists
}

// GetMessageTemplateTranslation gets translation from message template (manual priority > async)
func (c *ResponseConfig) GetMessageTemplateTranslation(templateKey, lang string) (string, bool) {
	// Priority 1: Manual templates
	if c.ManualMessageTemplates != nil {
		if template, exists := c.ManualMessageTemplates[templateKey]; exists {
			if translations := template.Translations; translations != nil {
				if translation, exists := translations[lang]; exists {
					return translation, true
				}
			}
		}
	}

	// Priority 2: Async templates
	if c.MessageTemplates != nil {
		if template, exists := c.MessageTemplates[templateKey]; exists {
			if translations := template.Translations; translations != nil {
				if translation, exists := translations[lang]; exists {
					return translation, true
				}
			}
		}
	}

	// Fallback to default language
	if lang != c.GetDefaultLanguage() {
		return c.GetMessageTemplateTranslation(templateKey, c.GetDefaultLanguage())
	}

	// Fallback to template string
	if template, exists := c.GetMessageTemplate(templateKey); exists {
		return template.Template, true
	}

	return "", false
}

// GetMessageTemplateTranslationWithFallback gets translation with fallback
func (c *ResponseConfig) GetMessageTemplateTranslationWithFallback(templateKey, lang string) string {
	// Try requested language
	if translation, exists := c.GetMessageTemplateTranslation(templateKey, lang); exists {
		return translation
	}

	// Fallback to default language
	if translation, exists := c.GetMessageTemplateTranslation(templateKey, c.GetDefaultLanguage()); exists {
		return translation
	}

	// Fallback to template string
	if template, exists := c.GetMessageTemplate(templateKey); exists {
		return template.Template
	}

	// Fallback to key
	return templateKey
}

// Printer returns ConfigPrinter for print/export config
func (c *ResponseConfig) Printer() *ConfigPrinter {
	return NewConfigPrinter(c)
}

// PrintConfig prints config to console (shortcut method)
func (c *ResponseConfig) PrintConfig() error {
	return c.Printer().Print()
}

// PrintConfigWithIndent prints config with indent (shortcut method)
func (c *ResponseConfig) PrintConfigWithIndent(useIndent bool) error {
	return c.Printer().WithIndent(useIndent).Print()
}

// ExportConfig returns config as JSON string (shortcut method)
func (c *ResponseConfig) ExportConfig() (string, error) {
	return c.Printer().Export()
}

// ExportConfigToFile saves config to file (shortcut method)
func (c *ResponseConfig) ExportConfigToFile(filename string) error {
	return c.Printer().ExportToFile(filename)
}
