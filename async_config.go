package goresponse

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AsyncConfigManager for managing configuration asynchronously with auto refresh
type AsyncConfigManager struct {
	source    ConfigSource
	config    *ResponseConfig
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	interval  time.Duration
	callbacks []ConfigChangeCallback
	isRunning bool
	lastError error
}

// NewAsyncConfigManager creates AsyncConfigManager instance
func NewAsyncConfigManager(source ConfigSource, refreshInterval time.Duration) *AsyncConfigManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &AsyncConfigManager{
		source:    source,
		ctx:       ctx,
		cancel:    cancel,
		interval:  refreshInterval,
		callbacks: make([]ConfigChangeCallback, 0),
		isRunning: false,
	}
}

// Start starts auto refresh configuration
func (acm *AsyncConfigManager) Start() error {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if acm.isRunning {
		return fmt.Errorf("async config manager is already running")
	}

	// Load configuration for the first time
	config, err := LoadConfig(acm.source)
	if err != nil {
		return fmt.Errorf("failed to load initial config: %w", err)
	}

	acm.config = config
	acm.isRunning = true

	// Start goroutine for auto refresh
	go acm.refreshLoop()

	return nil
}

// Stop stops auto refresh
func (acm *AsyncConfigManager) Stop() {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if !acm.isRunning {
		return
	}

	acm.cancel()
	acm.isRunning = false
}

// refreshLoop runs loop for auto refresh
func (acm *AsyncConfigManager) refreshLoop() {
	ticker := time.NewTicker(acm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-acm.ctx.Done():
			return
		case <-ticker.C:
			acm.refreshConfig()
		}
	}
}

// refreshConfig performs configuration refresh
func (acm *AsyncConfigManager) refreshConfig() {
	newConfig, err := LoadConfig(acm.source)
	if err != nil {
		acm.mu.Lock()
		acm.lastError = err
		acm.mu.Unlock()
		return
	}

	acm.mu.Lock()
	oldConfig := acm.config
	newConfig.ManualMessageTemplates = oldConfig.ManualMessageTemplates
	acm.config = newConfig
	acm.lastError = nil
	callbacks := acm.callbacks
	acm.mu.Unlock()

	// Call all callbacks
	for _, callback := range callbacks {
		callback(oldConfig, newConfig)
	}
}

// GetConfig returns current configuration
func (acm *AsyncConfigManager) GetConfig() *ResponseConfig {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	return acm.config
}

// GetTranslation gets translation with thread safety
func (acm *AsyncConfigManager) GetTranslation(lang, key string) (string, bool) {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return "", false
	}
	return acm.config.GetTranslation(lang, key)
}

// GetTranslationWithFallback gets translation with fallback and thread safety
func (acm *AsyncConfigManager) GetTranslationWithFallback(lang, key string) string {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return key
	}
	return acm.config.GetTranslationWithFallback(lang, key)
}

// GetMessageTemplate gets message template with thread safety
func (acm *AsyncConfigManager) GetMessageTemplate(key string) (*MessageTemplate, bool) {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return nil, false
	}
	return acm.config.GetMessageTemplate(key)
}

// GetSupportedLanguages returns list of supported languages with thread safety
func (acm *AsyncConfigManager) GetSupportedLanguages() []string {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return []string{}
	}
	return acm.config.GetSupportedLanguages()
}

// GetDefaultLanguage returns default language with thread safety
func (acm *AsyncConfigManager) GetDefaultLanguage() string {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return ""
	}
	return acm.config.GetDefaultLanguage()
}

// AddCallback adds callback for configuration changes
func (acm *AsyncConfigManager) AddCallback(callback ConfigChangeCallback) {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	acm.callbacks = append(acm.callbacks, callback)
}

// RemoveAllCallbacks removes all callbacks
func (acm *AsyncConfigManager) RemoveAllCallbacks() {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	acm.callbacks = make([]ConfigChangeCallback, 0)
}

// IsRunning returns status whether manager is running
func (acm *AsyncConfigManager) IsRunning() bool {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	return acm.isRunning
}

// GetLastError returns last error that occurred
func (acm *AsyncConfigManager) GetLastError() error {
	acm.mu.RLock()
	defer acm.mu.RUnlock()
	return acm.lastError
}

// ForceRefresh forces configuration refresh manually
func (acm *AsyncConfigManager) ForceRefresh() error {
	acm.refreshConfig()
	return acm.GetLastError()
}

// UpdateSource changes configuration source
func (acm *AsyncConfigManager) UpdateSource(newSource ConfigSource) {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	acm.source = newSource
}

// UpdateInterval changes refresh interval
func (acm *AsyncConfigManager) UpdateInterval(newInterval time.Duration) {
	acm.mu.Lock()
	defer acm.mu.Unlock()
	acm.interval = newInterval
}

// AddMessageTemplate adds message template (thread-safe, manual priority)
func (acm *AsyncConfigManager) AddMessageTemplate(template *MessageTemplate) {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if acm.config == nil {
		return
	}
	acm.config.AddMessageTemplate(template)
}

// AddMessageTemplates adds multiple message templates (thread-safe)
func (acm *AsyncConfigManager) AddMessageTemplates(templates ...*MessageTemplate) {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if acm.config == nil {
		return
	}
	acm.config.AddMessageTemplates(templates...)
}

// RemoveMessageTemplate removes message template (thread-safe)
func (acm *AsyncConfigManager) RemoveMessageTemplate(key string) {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if acm.config == nil {
		return
	}
	acm.config.RemoveMessageTemplate(key)
}

// UpdateMessageTemplate updates message template (thread-safe)
func (acm *AsyncConfigManager) UpdateMessageTemplate(template *MessageTemplate) {
	acm.mu.Lock()
	defer acm.mu.Unlock()

	if acm.config == nil {
		return
	}
	acm.config.UpdateMessageTemplate(template)
}

// GetMessageTemplateTranslation gets translation from message template (thread-safe)
func (acm *AsyncConfigManager) GetMessageTemplateTranslation(templateKey, lang string) (string, bool) {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return "", false
	}
	return acm.config.GetMessageTemplateTranslation(templateKey, lang)
}

// GetMessageTemplateTranslationWithFallback gets translation with fallback (thread-safe)
func (acm *AsyncConfigManager) GetMessageTemplateTranslationWithFallback(templateKey, lang string) string {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return templateKey
	}
	return acm.config.GetMessageTemplateTranslationWithFallback(templateKey, lang)
}

// Printer returns ConfigPrinter for print/export config (thread-safe)
func (acm *AsyncConfigManager) Printer() *ConfigPrinter {
	acm.mu.RLock()
	defer acm.mu.RUnlock()

	if acm.config == nil {
		return nil
	}
	return NewConfigPrinter(acm.config)
}

// PrintConfig prints config to console (thread-safe shortcut method)
func (acm *AsyncConfigManager) PrintConfig() error {
	printer := acm.Printer()
	if printer == nil {
		return fmt.Errorf("config is not loaded")
	}
	return printer.Print()
}

// PrintConfigWithIndent prints config with indent (thread-safe shortcut method)
func (acm *AsyncConfigManager) PrintConfigWithIndent(useIndent bool) error {
	printer := acm.Printer()
	if printer == nil {
		return fmt.Errorf("config is not loaded")
	}
	return printer.WithIndent(useIndent).Print()
}

// ExportConfig returns config as JSON string (thread-safe shortcut method)
func (acm *AsyncConfigManager) ExportConfig() (string, error) {
	printer := acm.Printer()
	if printer == nil {
		return "", fmt.Errorf("config is not loaded")
	}
	return printer.Export()
}

// ExportConfigToFile saves config to file (thread-safe shortcut method)
func (acm *AsyncConfigManager) ExportConfigToFile(filename string) error {
	printer := acm.Printer()
	if printer == nil {
		return fmt.Errorf("config is not loaded")
	}
	return printer.ExportToFile(filename)
}
