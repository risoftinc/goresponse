package goresponse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strings"
)

// ResponseContextKey represents the type for context keys used in response building
type ResponseContextKey string

// ResponseBuilder provides a fluent interface for building standardized API responses
// It supports method chaining and parameter substitution in message templates
type ResponseBuilder struct {
	MessageKey   string          // Key to identify the message template
	Params       map[string]any  // Parameters for template substitution
	Data         map[string]any  // Response data payload
	Meta         map[string]any  // Response metadata
	Context      context.Context // Context containing language and protocol info
	Language     string          // Language code for message translation
	Protocol     string          // Protocol type (http, grpc, etc.)
	ErrorData    error           // Error information if this is an error response
	IsBuiltError bool            // Flag indicating if this builder represents an error
}

// Response represents the final standardized response structure
// It contains all necessary information for API responses across different protocols
type Response struct {
	Code     int            `json:"code"`           // Response code (HTTP status, gRPC code, etc.)
	Message  string         `json:"message"`        // Human-readable message
	Data     map[string]any `json:"data,omitempty"` // Response payload data
	Meta     map[string]any `json:"meta,omitempty"` // Additional metadata
	Error    error          `json:"-"`              // Error details if applicable (not serialized)
	Language string         `json:"-"`              // Language used (not serialized)
	Protocol string         `json:"-"`              // Protocol used (not serialized)
}

const (
	// LanguageKey is the context key for storing language information
	LanguageKey ResponseContextKey = "goresponse-language"
	// ProtocolKey is the context key for storing protocol information
	ProtocolKey ResponseContextKey = "goresponse-protocol"
)

// NewResponseBuilder creates a new ResponseBuilder instance with the given message key
// The message key is used to identify the appropriate message template from configuration
func NewResponseBuilder(messageKey string) *ResponseBuilder {
	return &ResponseBuilder{
		MessageKey: messageKey,
	}
}

// WithProtocol adds protocol information to the context
// This is useful for passing protocol type (http, grpc, etc.) through the request chain
func WithProtocol(ctx context.Context, protocol string) context.Context {
	return context.WithValue(ctx, ProtocolKey, protocol)
}

// WithLanguage adds language information to the context
// This is useful for passing language preference through the request chain
func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, LanguageKey, language)
}

// GetLanguageFromContext extracts language information from the context
// Returns the language string and a boolean indicating if the language was found
func GetLanguageFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	lang, ok := ctx.Value(LanguageKey).(string)
	return lang, ok
}

// GetProtocolFromContext extracts protocol information from the context
// Returns the protocol string and a boolean indicating if the protocol was found
func GetProtocolFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	protocol, ok := ctx.Value(ProtocolKey).(string)
	return protocol, ok
}

// GetLanguage extracts language information from the context
// Returns the language string, or empty string if not found or context is nil
func GetLanguage(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return ""
}

// GetProtocol extracts protocol information from the context
// Returns the protocol string, or empty string if not found or context is nil
func GetProtocol(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if protocol, ok := ctx.Value(ProtocolKey).(string); ok {
		return protocol
	}
	return ""
}

// WithContext sets the context and extracts language and protocol information if available
// This method allows the response builder to inherit language and protocol settings from the request context
func (rb *ResponseBuilder) WithContext(ctx context.Context) *ResponseBuilder {
	rb.Context = ctx

	if ctx != nil {
		// Extract language from context if available
		if lang, ok := ctx.Value(LanguageKey).(string); ok {
			rb.Language = lang
		}

		// Extract protocol from context if available
		if proto, ok := ctx.Value(ProtocolKey).(string); ok {
			rb.Protocol = proto
		}
	}

	return rb
}

// SetLanguage manually sets the language for message translation
// This overrides any language setting from context
func (rb *ResponseBuilder) SetLanguage(language string) *ResponseBuilder {
	rb.Language = language
	return rb
}

// SetProtocol manually sets the protocol type for response code mapping
// This overrides any protocol setting from context
func (rb *ResponseBuilder) SetProtocol(protocol string) *ResponseBuilder {
	rb.Protocol = protocol
	return rb
}

// SetError sets an error for this response builder and marks it as an error response
// This will cause the final response to be treated as an error
func (rb *ResponseBuilder) SetError(err error) *ResponseBuilder {
	rb.ErrorData = err
	rb.IsBuiltError = true
	return rb
}

// SetParam adds a single parameter for template substitution
// Parameters are used to replace placeholders like $name in message templates
func (rb *ResponseBuilder) SetParam(key string, value any) *ResponseBuilder {
	if rb.Params == nil {
		rb.Params = make(map[string]any)
	}
	rb.Params[key] = value
	return rb
}

// SetParams adds multiple parameters for template substitution
// Existing parameters with the same key will be replaced
func (rb *ResponseBuilder) SetParams(params map[string]any) *ResponseBuilder {
	if rb.Params == nil {
		rb.Params = make(map[string]any)
	}

	// Copy all parameters, replacing existing ones if they exist
	maps.Copy(rb.Params, params)

	return rb
}

// SetData adds a single data field to the response payload
// Data fields are included in the final response as the data payload
func (rb *ResponseBuilder) SetData(key string, value any) *ResponseBuilder {
	if rb.Data == nil {
		rb.Data = make(map[string]any)
	}
	rb.Data[key] = value
	return rb
}

// SetDatas adds multiple data fields to the response payload
// Existing data fields with the same key will be replaced
func (rb *ResponseBuilder) SetDatas(data map[string]any) *ResponseBuilder {
	if rb.Data == nil {
		rb.Data = make(map[string]any)
	}

	// Copy all data fields, replacing existing ones if they exist
	maps.Copy(rb.Data, data)

	return rb
}

// SetMeta adds a single metadata field to the response
// Metadata fields are included in the final response as additional information
func (rb *ResponseBuilder) SetMeta(key string, value any) *ResponseBuilder {
	if rb.Meta == nil {
		rb.Meta = make(map[string]any)
	}
	rb.Meta[key] = value
	return rb
}

// SetMetas adds multiple metadata fields to the response
// Existing metadata fields with the same key will be replaced
func (rb *ResponseBuilder) SetMetas(meta map[string]any) *ResponseBuilder {
	if rb.Meta == nil {
		rb.Meta = make(map[string]any)
	}

	// Copy all metadata fields, replacing existing ones if they exist
	maps.Copy(rb.Meta, meta)

	return rb
}

// Error implements the error interface by returning a JSON representation of the response builder
// This allows the response builder to be used as an error type
func (rb *ResponseBuilder) Error() string {
	b, _ := json.Marshal(rb)
	return string(b)
}

// ToError converts the response builder to an error type
// This allows the response builder to be returned as an error from functions
func (rb *ResponseBuilder) ToError() error {
	return rb
}

// ParseResponseBuilderError attempts to extract a ResponseBuilder from an error
// This is useful for recovering response builder information from error chains
func ParseResponseBuilderError(err error) (*ResponseBuilder, bool) {
	if err == nil {
		return nil, false
	}

	var rb *ResponseBuilder
	if errors.As(err, &rb) {
		return rb, true
	}
	return nil, false
}

// BuildResponse constructs the final Response from a ResponseBuilder using the configuration
// This method handles message template resolution, parameter substitution, and code mapping
// Note: This method may experience data inconsistency if called during async configuration reload
func (c *ResponseConfig) BuildResponse(rb *ResponseBuilder) (*Response, error) {
	if rb == nil {
		return nil, errors.New("response builder is nil")
	}

	r := &Response{}

	if rb.Language == "" {
		r.Language = c.GetDefaultLanguage()
	}

	// Get the message template for the specified message key
	template, exists := c.GetMessageTemplate(rb.MessageKey)
	if !exists {
		return nil, errors.New("message template not found")
	}

	// Determine the message text based on language preference
	// Fallback to template default if translation is not available
	if template.Translations[rb.Language] == "" {
		if translation, exists := c.GetTranslation(rb.Language, rb.MessageKey); exists {
			r.Message = translation
		} else {
			r.Message = template.Template
		}
	} else {
		r.Message = template.Translations[rb.Language]
	}

	// Substitute parameters in the message template
	r.Message = substituteParams(r.Message, rb.Params)

	// Map the response code based on protocol
	r.Code = template.CodeMappings[rb.Protocol]

	// Set response data and metadata
	r.Data = rb.Data
	r.Meta = rb.Meta
	r.Error = rb.ErrorData
	r.Language = rb.Language
	r.Protocol = rb.Protocol

	return r, nil
}

// substituteParams replaces parameter placeholders in a template string with actual values
// Placeholders are in the format $paramName and are replaced with string representations of values
func substituteParams(template string, params map[string]any) string {
	result := template
	for key, value := range params {
		placeholder := fmt.Sprintf("$%s", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}
