package goresponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ConfigSource struct to specify configuration source
type ConfigSource struct {
	Method string `json:"method"` // "file" or "url"
	Path   string `json:"path"`   // file path or URL
}

// TranslationSource struct to specify translation source per language
type TranslationSource struct {
	Method string `json:"method"` // "file" or "url"
	Path   string `json:"path"`   // file path or URL
}

// ResponseConfig struct to store response configuration
type ResponseConfig struct {
	MessageTemplates       map[string]MessageTemplate   `json:"message_templates"`
	ManualMessageTemplates map[string]MessageTemplate   `json:"-"` // Manual templates (high priority)
	DefaultLanguage        string                       `json:"default_language"`
	Languages              []string                     `json:"languages"`
	Translations           map[string]map[string]string `json:"translations"`       // Inline translations
	TranslationSources     map[string]TranslationSource `json:"translation_source"` // Separate translation sources
}

// MessageTemplate struct for message template
type MessageTemplate struct {
	Key          string            `json:"key"`
	Template     string            `json:"template"`
	CodeMappings map[string]int    `json:"code_mappings"`
	Translations map[string]string `json:"translations,omitempty"` // Translations per language
}

// ConfigChangeCallback is function type for callback when config changes
type ConfigChangeCallback func(oldConfig, newConfig *ResponseConfig)

// MessageTemplateBuilder for method chaining in creating message template
type MessageTemplateBuilder struct {
	template *MessageTemplate
}

// NewMessageTemplateBuilder creates MessageTemplateBuilder instance
func NewMessageTemplateBuilder(key string) *MessageTemplateBuilder {
	return &MessageTemplateBuilder{
		template: &MessageTemplate{
			Key:          key,
			Template:     "",
			CodeMappings: make(map[string]int),
		},
	}
}

// WithTemplate adds template string
func (mtb *MessageTemplateBuilder) WithTemplate(template string) *MessageTemplateBuilder {
	mtb.template.Template = template
	return mtb
}

// WithTranslation adds translation for specific language
func (mtb *MessageTemplateBuilder) WithTranslation(lang, translation string) *MessageTemplateBuilder {
	if mtb.template.Translations == nil {
		mtb.template.Translations = make(map[string]string)
	}
	mtb.template.Translations[lang] = translation
	return mtb
}

// WithTranslations adds multiple translations at once
func (mtb *MessageTemplateBuilder) WithTranslations(translations map[string]string) *MessageTemplateBuilder {
	if mtb.template.Translations == nil {
		mtb.template.Translations = make(map[string]string)
	}
	for lang, translation := range translations {
		mtb.template.Translations[lang] = translation
	}
	return mtb
}

// WithCodeMapping adds code mapping (HTTP status, etc.)
func (mtb *MessageTemplateBuilder) WithCodeMapping(mappingType string, code int) *MessageTemplateBuilder {
	mtb.template.CodeMappings[mappingType] = code
	return mtb
}

// WithCodeMappings adds multiple code mappings at once
func (mtb *MessageTemplateBuilder) WithCodeMappings(codeMappings map[string]int) *MessageTemplateBuilder {
	for mappingType, code := range codeMappings {
		mtb.template.CodeMappings[mappingType] = code
	}
	return mtb
}

// Build returns completed MessageTemplate
func (mtb *MessageTemplateBuilder) Build() *MessageTemplate {
	return mtb.template
}

// ConfigPrinter for printing and exporting config with method chaining
type ConfigPrinter struct {
	config *ResponseConfig
	indent bool
}

// NewConfigPrinter creates ConfigPrinter instance
func NewConfigPrinter(config *ResponseConfig) *ConfigPrinter {
	return &ConfigPrinter{
		config: config,
		indent: true, // Default with indent
	}
}

// WithIndent sets whether to use indent or not
func (cp *ConfigPrinter) WithIndent(useIndent bool) *ConfigPrinter {
	cp.indent = useIndent
	return cp
}

// Print prints config to console
func (cp *ConfigPrinter) Print() error {
	jsonData, err := cp.toJSON()
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

// Export returns config as JSON string
func (cp *ConfigPrinter) Export() (string, error) {
	jsonData, err := cp.toJSON()
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// ExportToFile saves config to file
func (cp *ConfigPrinter) ExportToFile(filename string) error {
	jsonData, err := cp.toJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// toJSON converts config to JSON
func (cp *ConfigPrinter) toJSON() ([]byte, error) {
	// Handle nil config
	if cp.config == nil {
		return nil, errors.New("config is nil")
	}

	// Create struct for export that merges manual and async templates
	exportConfig := struct {
		MessageTemplates   map[string]MessageTemplate   `json:"message_templates"`
		DefaultLanguage    string                       `json:"default_language"`
		Languages          []string                     `json:"languages"`
		Translations       map[string]map[string]string `json:"translations"`
		TranslationSources map[string]TranslationSource `json:"translation_source,omitempty"`
	}{
		MessageTemplates:   make(map[string]MessageTemplate),
		DefaultLanguage:    cp.config.DefaultLanguage,
		Languages:          cp.config.Languages,
		Translations:       cp.config.Translations,
		TranslationSources: cp.config.TranslationSources,
	}

	// Merge async templates first
	for key, template := range cp.config.MessageTemplates {
		exportConfig.MessageTemplates[key] = template
	}

	// Override with manual templates if any (manual priority)
	if cp.config.ManualMessageTemplates != nil {
		for key, template := range cp.config.ManualMessageTemplates {
			exportConfig.MessageTemplates[key] = template
		}
	}

	// Marshal to JSON
	if cp.indent {
		return json.MarshalIndent(exportConfig, "", "  ")
	}
	return json.Marshal(exportConfig)
}
