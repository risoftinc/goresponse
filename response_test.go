package goresponse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestResponseContextKey tests ResponseContextKey type
func TestResponseContextKey(t *testing.T) {
	tests := []struct {
		name     string
		key      ResponseContextKey
		expected string
	}{
		{
			name:     "LanguageKey",
			key:      LanguageKey,
			expected: "goresponse-language",
		},
		{
			name:     "ProtocolKey",
			key:      ProtocolKey,
			expected: "goresponse-protocol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.key) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.key))
			}
		})
	}
}

// TestNewResponseBuilder tests NewResponseBuilder function
func TestNewResponseBuilder(t *testing.T) {
	tests := []struct {
		name       string
		messageKey string
		expected   string
	}{
		{
			name:       "Valid message key",
			messageKey: "welcome",
			expected:   "welcome",
		},
		{
			name:       "Empty message key",
			messageKey: "",
			expected:   "",
		},
		{
			name:       "Complex message key",
			messageKey: "user.created.success",
			expected:   "user.created.success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewResponseBuilder(tt.messageKey)
			if builder == nil {
				t.Fatal("Expected builder to be created")
			}
			if builder.MessageKey != tt.expected {
				t.Errorf("Expected message key %s, got %s", tt.expected, builder.MessageKey)
			}
			if builder.Params != nil {
				t.Error("Expected Params to be nil initially")
			}
			if builder.Data != nil {
				t.Error("Expected Data to be nil initially")
			}
			if builder.Meta != nil {
				t.Error("Expected Meta to be nil initially")
			}
			if builder.Context != nil {
				t.Error("Expected Context to be nil initially")
			}
			if builder.Language != "" {
				t.Errorf("Expected empty language, got %s", builder.Language)
			}
			if builder.Protocol != "" {
				t.Errorf("Expected empty protocol, got %s", builder.Protocol)
			}
			if builder.ErrorData != nil {
				t.Error("Expected ErrorData to be nil initially")
			}
			if builder.IsBuiltError {
				t.Error("Expected IsBuiltError to be false initially")
			}
		})
	}
}

// TestWithProtocol tests WithProtocol function
func TestWithProtocol(t *testing.T) {
	ctx := context.Background()
	protocol := "http"

	newCtx := WithProtocol(ctx, protocol)

	// Check if protocol was added to context
	if newCtx.Value(ProtocolKey) != protocol {
		t.Errorf("Expected protocol %s in context, got %v", protocol, newCtx.Value(ProtocolKey))
	}

	// Test with different protocol
	grpcCtx := WithProtocol(ctx, "grpc")
	if grpcCtx.Value(ProtocolKey) != "grpc" {
		t.Errorf("Expected grpc protocol in context, got %v", grpcCtx.Value(ProtocolKey))
	}

	// Test with empty protocol
	emptyCtx := WithProtocol(ctx, "")
	if emptyCtx.Value(ProtocolKey) != "" {
		t.Errorf("Expected empty protocol in context, got %v", emptyCtx.Value(ProtocolKey))
	}
}

// TestWithLanguage tests WithLanguage function
func TestWithLanguage(t *testing.T) {
	ctx := context.Background()
	language := "en"

	newCtx := WithLanguage(ctx, language)

	// Check if language was added to context
	if newCtx.Value(LanguageKey) != language {
		t.Errorf("Expected language %s in context, got %v", language, newCtx.Value(LanguageKey))
	}

	// Test with different language
	idCtx := WithLanguage(ctx, "id")
	if idCtx.Value(LanguageKey) != "id" {
		t.Errorf("Expected id language in context, got %v", idCtx.Value(LanguageKey))
	}

	// Test with empty language
	emptyCtx := WithLanguage(ctx, "")
	if emptyCtx.Value(LanguageKey) != "" {
		t.Errorf("Expected empty language in context, got %v", emptyCtx.Value(LanguageKey))
	}
}

// TestResponseBuilderWithContext tests WithContext method
func TestResponseBuilderWithContext(t *testing.T) {
	t.Run("Context with language and protocol", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithLanguage(ctx, "en")
		ctx = WithProtocol(ctx, "http")

		builder := NewResponseBuilder("test")
		result := builder.WithContext(ctx)

		if result != builder {
			t.Error("Expected method chaining to return same builder")
		}
		if builder.Context != ctx {
			t.Error("Expected context to be set")
		}
		if builder.Language != "en" {
			t.Errorf("Expected language 'en', got %s", builder.Language)
		}
		if builder.Protocol != "http" {
			t.Errorf("Expected protocol 'http', got %s", builder.Protocol)
		}
	})

	t.Run("Context with only language", func(t *testing.T) {
		ctx := WithLanguage(context.Background(), "id")

		builder := NewResponseBuilder("test")
		builder.WithContext(ctx)

		if builder.Language != "id" {
			t.Errorf("Expected language 'id', got %s", builder.Language)
		}
		if builder.Protocol != "" {
			t.Errorf("Expected empty protocol, got %s", builder.Protocol)
		}
	})

	t.Run("Context with only protocol", func(t *testing.T) {
		ctx := WithProtocol(context.Background(), "grpc")

		builder := NewResponseBuilder("test")
		builder.WithContext(ctx)

		if builder.Language != "" {
			t.Errorf("Expected empty language, got %s", builder.Language)
		}
		if builder.Protocol != "grpc" {
			t.Errorf("Expected protocol 'grpc', got %s", builder.Protocol)
		}
	})

	t.Run("Empty context", func(t *testing.T) {
		ctx := context.Background()

		builder := NewResponseBuilder("test")
		builder.WithContext(ctx)

		if builder.Language != "" {
			t.Errorf("Expected empty language, got %s", builder.Language)
		}
		if builder.Protocol != "" {
			t.Errorf("Expected empty protocol, got %s", builder.Protocol)
		}
	})

	t.Run("Context with wrong types", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), LanguageKey, 123)
		ctx = context.WithValue(ctx, ProtocolKey, true)

		builder := NewResponseBuilder("test")
		builder.WithContext(ctx)

		// Should not set language or protocol due to type assertion failure
		if builder.Language != "" {
			t.Errorf("Expected empty language, got %s", builder.Language)
		}
		if builder.Protocol != "" {
			t.Errorf("Expected empty protocol, got %s", builder.Protocol)
		}
	})
}

// TestResponseBuilderSetLanguage tests SetLanguage method
func TestResponseBuilderSetLanguage(t *testing.T) {
	builder := NewResponseBuilder("test")
	result := builder.SetLanguage("en")

	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Language != "en" {
		t.Errorf("Expected language 'en', got %s", builder.Language)
	}

	// Test overriding
	builder.SetLanguage("id")
	if builder.Language != "id" {
		t.Errorf("Expected language 'id', got %s", builder.Language)
	}

	// Test empty language
	builder.SetLanguage("")
	if builder.Language != "" {
		t.Errorf("Expected empty language, got %s", builder.Language)
	}
}

// TestResponseBuilderSetProtocol tests SetProtocol method
func TestResponseBuilderSetProtocol(t *testing.T) {
	builder := NewResponseBuilder("test")
	result := builder.SetProtocol("http")

	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Protocol != "http" {
		t.Errorf("Expected protocol 'http', got %s", builder.Protocol)
	}

	// Test overriding
	builder.SetProtocol("grpc")
	if builder.Protocol != "grpc" {
		t.Errorf("Expected protocol 'grpc', got %s", builder.Protocol)
	}

	// Test empty protocol
	builder.SetProtocol("")
	if builder.Protocol != "" {
		t.Errorf("Expected empty protocol, got %s", builder.Protocol)
	}
}

// TestResponseBuilderSetError tests SetError method
func TestResponseBuilderSetError(t *testing.T) {
	builder := NewResponseBuilder("test")
	testError := errors.New("test error")

	result := builder.SetError(testError)

	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.ErrorData != testError {
		t.Error("Expected error to be set")
	}
	if !builder.IsBuiltError {
		t.Error("Expected IsBuiltError to be true")
	}

	// Test with nil error
	builder.SetError(nil)
	if builder.ErrorData != nil {
		t.Error("Expected error to be nil")
	}
	if !builder.IsBuiltError {
		t.Error("Expected IsBuiltError to remain true")
	}
}

// TestResponseBuilderSetParam tests SetParam method
func TestResponseBuilderSetParam(t *testing.T) {
	builder := NewResponseBuilder("test")

	// Test first parameter
	result := builder.SetParam("name", "John")
	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Params == nil {
		t.Error("Expected Params to be initialized")
	}
	if builder.Params["name"] != "John" {
		t.Errorf("Expected name 'John', got %v", builder.Params["name"])
	}

	// Test second parameter
	builder.SetParam("age", 25)
	if builder.Params["age"] != 25 {
		t.Errorf("Expected age 25, got %v", builder.Params["age"])
	}

	// Test overriding parameter
	builder.SetParam("name", "Jane")
	if builder.Params["name"] != "Jane" {
		t.Errorf("Expected name 'Jane', got %v", builder.Params["name"])
	}

	// Test with different types
	builder.SetParam("active", true)
	if builder.Params["active"] != true {
		t.Errorf("Expected active true, got %v", builder.Params["active"])
	}

	builder.SetParam("score", 98.5)
	if builder.Params["score"] != 98.5 {
		t.Errorf("Expected score 98.5, got %v", builder.Params["score"])
	}
}

// TestResponseBuilderSetParams tests SetParams method
func TestResponseBuilderSetParams(t *testing.T) {
	builder := NewResponseBuilder("test")

	// Test with empty map
	result := builder.SetParams(map[string]any{})
	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Params == nil {
		t.Error("Expected Params to be initialized")
	}
	if len(builder.Params) != 0 {
		t.Errorf("Expected 0 params, got %d", len(builder.Params))
	}

	// Test with single parameter
	params := map[string]any{
		"name": "John",
	}
	builder.SetParams(params)
	if len(builder.Params) != 1 {
		t.Errorf("Expected 1 param, got %d", len(builder.Params))
	}
	if builder.Params["name"] != "John" {
		t.Errorf("Expected name 'John', got %v", builder.Params["name"])
	}

	// Test with multiple parameters
	params = map[string]any{
		"name":  "Jane",
		"age":   30,
		"email": "jane@example.com",
	}
	builder.SetParams(params)
	if len(builder.Params) != 3 {
		t.Errorf("Expected 3 params, got %d", len(builder.Params))
	}
	if builder.Params["name"] != "Jane" {
		t.Errorf("Expected name 'Jane', got %v", builder.Params["name"])
	}
	if builder.Params["age"] != 30 {
		t.Errorf("Expected age 30, got %v", builder.Params["age"])
	}
	if builder.Params["email"] != "jane@example.com" {
		t.Errorf("Expected email 'jane@example.com', got %v", builder.Params["email"])
	}

	// Test with nil map
	builder.SetParams(nil)
	if builder.Params == nil {
		t.Error("Expected Params to remain initialized")
	}
}

// TestResponseBuilderSetData tests SetData method
func TestResponseBuilderSetData(t *testing.T) {
	builder := NewResponseBuilder("test")

	// Test first data field
	result := builder.SetData("user_id", 123)
	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Data == nil {
		t.Error("Expected Data to be initialized")
	}
	if builder.Data["user_id"] != 123 {
		t.Errorf("Expected user_id 123, got %v", builder.Data["user_id"])
	}

	// Test second data field
	builder.SetData("username", "john_doe")
	if builder.Data["username"] != "john_doe" {
		t.Errorf("Expected username 'john_doe', got %v", builder.Data["username"])
	}

	// Test overriding data field
	builder.SetData("user_id", 456)
	if builder.Data["user_id"] != 456 {
		t.Errorf("Expected user_id 456, got %v", builder.Data["user_id"])
	}

	// Test with complex data
	user := map[string]any{
		"id":   123,
		"name": "John Doe",
		"role": "admin",
	}
	builder.SetData("user", user)
	if builder.Data["user"] == nil {
		t.Error("Expected user data to be set")
	}
}

// TestResponseBuilderSetDatas tests SetDatas method
func TestResponseBuilderSetDatas(t *testing.T) {
	builder := NewResponseBuilder("test")

	// Test with empty map
	result := builder.SetDatas(map[string]any{})
	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Data == nil {
		t.Error("Expected Data to be initialized")
	}
	if len(builder.Data) != 0 {
		t.Errorf("Expected 0 data fields, got %d", len(builder.Data))
	}

	// Test with single data field
	data := map[string]any{
		"user_id": 123,
	}
	builder.SetDatas(data)
	if len(builder.Data) != 1 {
		t.Errorf("Expected 1 data field, got %d", len(builder.Data))
	}
	if builder.Data["user_id"] != 123 {
		t.Errorf("Expected user_id 123, got %v", builder.Data["user_id"])
	}

	// Test with multiple data fields
	data = map[string]any{
		"user_id":  456,
		"username": "jane_doe",
		"email":    "jane@example.com",
		"active":   true,
	}
	builder.SetDatas(data)
	if len(builder.Data) != 4 {
		t.Errorf("Expected 4 data fields, got %d", len(builder.Data))
	}
	if builder.Data["user_id"] != 456 {
		t.Errorf("Expected user_id 456, got %v", builder.Data["user_id"])
	}
	if builder.Data["username"] != "jane_doe" {
		t.Errorf("Expected username 'jane_doe', got %v", builder.Data["username"])
	}

	// Test with nil map
	builder.SetDatas(nil)
	if builder.Data == nil {
		t.Error("Expected Data to remain initialized")
	}
}

// TestResponseBuilderSetMeta tests SetMeta method
func TestResponseBuilderSetMeta(t *testing.T) {
	builder := NewResponseBuilder("test")

	// Test first meta field
	result := builder.SetMeta("request_id", "req-123")
	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Meta == nil {
		t.Error("Expected Meta to be initialized")
	}
	if builder.Meta["request_id"] != "req-123" {
		t.Errorf("Expected request_id 'req-123', got %v", builder.Meta["request_id"])
	}

	// Test second meta field
	builder.SetMeta("timestamp", 1640995200)
	if builder.Meta["timestamp"] != 1640995200 {
		t.Errorf("Expected timestamp 1640995200, got %v", builder.Meta["timestamp"])
	}

	// Test overriding meta field
	builder.SetMeta("request_id", "req-456")
	if builder.Meta["request_id"] != "req-456" {
		t.Errorf("Expected request_id 'req-456', got %v", builder.Meta["request_id"])
	}

	// Test with complex meta
	meta := map[string]any{
		"version": "1.0.0",
		"source":  "api",
	}
	builder.SetMeta("info", meta)
	if builder.Meta["info"] == nil {
		t.Error("Expected info meta to be set")
	}
}

// TestResponseBuilderSetMetas tests SetMetas method
func TestResponseBuilderSetMetas(t *testing.T) {
	builder := NewResponseBuilder("test")

	// Test with empty map
	result := builder.SetMetas(map[string]any{})
	if result != builder {
		t.Error("Expected method chaining to return same builder")
	}
	if builder.Meta == nil {
		t.Error("Expected Meta to be initialized")
	}
	if len(builder.Meta) != 0 {
		t.Errorf("Expected 0 meta fields, got %d", len(builder.Meta))
	}

	// Test with single meta field
	meta := map[string]any{
		"request_id": "req-123",
	}
	builder.SetMetas(meta)
	if len(builder.Meta) != 1 {
		t.Errorf("Expected 1 meta field, got %d", len(builder.Meta))
	}
	if builder.Meta["request_id"] != "req-123" {
		t.Errorf("Expected request_id 'req-123', got %v", builder.Meta["request_id"])
	}

	// Test with multiple meta fields
	meta = map[string]any{
		"request_id": "req-456",
		"timestamp":  1640995200,
		"version":    "1.0.0",
		"source":     "api",
	}
	builder.SetMetas(meta)
	if len(builder.Meta) != 4 {
		t.Errorf("Expected 4 meta fields, got %d", len(builder.Meta))
	}
	if builder.Meta["request_id"] != "req-456" {
		t.Errorf("Expected request_id 'req-456', got %v", builder.Meta["request_id"])
	}
	if builder.Meta["timestamp"] != 1640995200 {
		t.Errorf("Expected timestamp 1640995200, got %v", builder.Meta["timestamp"])
	}

	// Test with nil map
	builder.SetMetas(nil)
	if builder.Meta == nil {
		t.Error("Expected Meta to remain initialized")
	}
}

// TestResponseBuilderError tests Error method
func TestResponseBuilderError(t *testing.T) {
	builder := NewResponseBuilder("test")
	builder.SetLanguage("en")
	builder.SetProtocol("http")
	builder.SetParam("name", "John")
	builder.SetData("user_id", 123)
	builder.SetMeta("request_id", "req-123")
	builder.SetError(errors.New("test error"))

	errorStr := builder.Error()
	if errorStr == "" {
		t.Error("Expected non-empty error string")
	}

	// Verify it's valid JSON
	var errorData map[string]interface{}
	err := json.Unmarshal([]byte(errorStr), &errorData)
	if err != nil {
		t.Errorf("Error string is not valid JSON: %v", err)
	}

	// Verify it contains expected fields
	if errorData["MessageKey"] != "test" {
		t.Errorf("Expected MessageKey 'test', got %v", errorData["MessageKey"])
	}
	if errorData["Language"] != "en" {
		t.Errorf("Expected Language 'en', got %v", errorData["Language"])
	}
	if errorData["Protocol"] != "http" {
		t.Errorf("Expected Protocol 'http', got %v", errorData["Protocol"])
	}
}

// TestResponseBuilderToError tests ToError method
func TestResponseBuilderToError(t *testing.T) {
	builder := NewResponseBuilder("test")
	err := builder.ToError()

	if err == nil {
		t.Error("Expected error to be returned")
	}

	// Verify it's the same builder
	if err != builder {
		t.Error("Expected ToError to return the same builder")
	}
}

// TestParseResponseBuilderError tests ParseResponseBuilderError function
func TestParseResponseBuilderError(t *testing.T) {
	t.Run("Valid response builder error", func(t *testing.T) {
		builder := NewResponseBuilder("test")
		builder.SetLanguage("en")
		builder.SetProtocol("http")

		err := builder.ToError()
		parsedBuilder, ok := ParseResponseBuilderError(err)

		if !ok {
			t.Error("Expected to successfully parse response builder error")
		}
		if parsedBuilder == nil {
			t.Fatal("Expected parsed builder to not be nil")
		}
		if parsedBuilder.MessageKey != "test" {
			t.Errorf("Expected MessageKey 'test', got %s", parsedBuilder.MessageKey)
		}
		if parsedBuilder.Language != "en" {
			t.Errorf("Expected Language 'en', got %s", parsedBuilder.Language)
		}
		if parsedBuilder.Protocol != "http" {
			t.Errorf("Expected Protocol 'http', got %s", parsedBuilder.Protocol)
		}
	})

	t.Run("Nil error", func(t *testing.T) {
		parsedBuilder, ok := ParseResponseBuilderError(nil)
		if ok {
			t.Error("Expected to not parse nil error")
		}
		if parsedBuilder != nil {
			t.Error("Expected parsed builder to be nil")
		}
	})

	t.Run("Regular error", func(t *testing.T) {
		regularErr := errors.New("regular error")
		parsedBuilder, ok := ParseResponseBuilderError(regularErr)
		if ok {
			t.Error("Expected to not parse regular error")
		}
		if parsedBuilder != nil {
			t.Error("Expected parsed builder to be nil")
		}
	})

	t.Run("Wrapped response builder error", func(t *testing.T) {
		builder := NewResponseBuilder("wrapped")
		builder.SetLanguage("id")

		wrappedErr := fmt.Errorf("wrapped: %w", builder.ToError())
		parsedBuilder, ok := ParseResponseBuilderError(wrappedErr)
		if !ok {
			t.Error("Expected to successfully parse wrapped response builder error")
		}
		if parsedBuilder == nil {
			t.Fatal("Expected parsed builder to not be nil")
		}
		if parsedBuilder.MessageKey != "wrapped" {
			t.Errorf("Expected MessageKey 'wrapped', got %s", parsedBuilder.MessageKey)
		}
		if parsedBuilder.Language != "id" {
			t.Errorf("Expected Language 'id', got %s", parsedBuilder.Language)
		}
	})
}

// TestResponseStruct tests Response struct
func TestResponseStruct(t *testing.T) {
	response := Response{
		Code:     200,
		Message:  "Success",
		Data:     map[string]any{"user_id": 123},
		Meta:     map[string]any{"request_id": "req-123"},
		Error:    errors.New("test error"),
		Language: "en",
		Protocol: "http",
	}

	// Test all fields
	if response.Code != 200 {
		t.Errorf("Expected code 200, got %d", response.Code)
	}
	if response.Message != "Success" {
		t.Errorf("Expected message 'Success', got %s", response.Message)
	}
	if response.Data == nil {
		t.Error("Expected data to be set")
	}
	if response.Meta == nil {
		t.Error("Expected meta to be set")
	}
	if response.Error == nil {
		t.Error("Expected error to be set")
	}
	if response.Language != "en" {
		t.Errorf("Expected language 'en', got %s", response.Language)
	}
	if response.Protocol != "http" {
		t.Errorf("Expected protocol 'http', got %s", response.Protocol)
	}
}

// TestResponseJSON tests Response JSON marshaling
func TestResponseJSON(t *testing.T) {
	response := Response{
		Code:     200,
		Message:  "Success",
		Data:     map[string]any{"user_id": 123},
		Meta:     map[string]any{"request_id": "req-123"},
		Error:    errors.New("test error"),
		Language: "en",
		Protocol: "http",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal Response: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal Response: %v", err)
	}

	// Verify unmarshaled data
	if unmarshaled.Code != response.Code {
		t.Errorf("Expected code %d, got %d", response.Code, unmarshaled.Code)
	}
	if unmarshaled.Message != response.Message {
		t.Errorf("Expected message %s, got %s", response.Message, unmarshaled.Message)
	}
	// Note: Language and Protocol are not serialized (json:"-")
	if unmarshaled.Language != "" {
		t.Errorf("Expected empty language (not serialized), got %s", unmarshaled.Language)
	}
	if unmarshaled.Protocol != "" {
		t.Errorf("Expected empty protocol (not serialized), got %s", unmarshaled.Protocol)
	}
}

// TestBuildResponse tests BuildResponse method
func TestBuildResponse(t *testing.T) {
	// Create test config
	config := &ResponseConfig{
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
			"error": {
				Key:      "error",
				Template: "Error: $message",
				CodeMappings: map[string]int{
					"http": 400,
					"grpc": 3,
				},
				Translations: map[string]string{
					"en": "Error: $message",
					"id": "Kesalahan: $message",
				},
			},
		},
		DefaultLanguage: "en",
		Languages:       []string{"en", "id"},
	}

	t.Run("Successful response with translation", func(t *testing.T) {
		builder := NewResponseBuilder("welcome")
		builder.SetLanguage("en")
		builder.SetProtocol("http")
		builder.SetParam("name", "John")
		builder.SetData("user_id", 123)
		builder.SetMeta("request_id", "req-123")

		response, err := config.BuildResponse(builder)
		if err != nil {
			t.Errorf("BuildResponse failed: %v", err)
		}

		if response.Code != 200 {
			t.Errorf("Expected code 200, got %d", response.Code)
		}
		if response.Message != "Welcome John" {
			t.Errorf("Expected message 'Welcome John', got %s", response.Message)
		}
		if response.Data == nil {
			t.Error("Expected data to be set")
		}
		if response.Meta == nil {
			t.Error("Expected meta to be set")
		}
		if response.Language != "en" {
			t.Errorf("Expected language 'en', got %s", response.Language)
		}
		if response.Protocol != "http" {
			t.Errorf("Expected protocol 'http', got %s", response.Protocol)
		}
	})

	t.Run("Response with fallback to template", func(t *testing.T) {
		builder := NewResponseBuilder("welcome")
		builder.SetLanguage("es") // Spanish not available, should fallback to template
		builder.SetProtocol("grpc")
		builder.SetParam("name", "Jane")

		response, err := config.BuildResponse(builder)
		if err != nil {
			t.Errorf("BuildResponse failed: %v", err)
		}

		if response.Code != 0 {
			t.Errorf("Expected code 0, got %d", response.Code)
		}
		if response.Message != "Welcome Jane" {
			t.Errorf("Expected message 'Welcome Jane', got %s", response.Message)
		}
	})

	t.Run("Error response", func(t *testing.T) {
		builder := NewResponseBuilder("error")
		builder.SetLanguage("id")
		builder.SetProtocol("http")
		builder.SetParam("message", "Invalid input")
		builder.SetError(errors.New("validation failed"))

		response, err := config.BuildResponse(builder)
		if err != nil {
			t.Errorf("BuildResponse failed: %v", err)
		}

		if response.Code != 400 {
			t.Errorf("Expected code 400, got %d", response.Code)
		}
		if response.Message != "Kesalahan: Invalid input" {
			t.Errorf("Expected message 'Kesalahan: Invalid input', got %s", response.Message)
		}
		if response.Error == nil {
			t.Error("Expected error to be set")
		}
	})

	t.Run("Nil response builder", func(t *testing.T) {
		response, err := config.BuildResponse(nil)
		if err == nil {
			t.Error("Expected error for nil response builder")
		}
		if response != nil {
			t.Error("Expected nil response for nil builder")
		}
		if !strings.Contains(err.Error(), "response builder is nil") {
			t.Errorf("Expected error about nil builder, got: %v", err)
		}
	})

	t.Run("Message template not found", func(t *testing.T) {
		builder := NewResponseBuilder("nonexistent")
		builder.SetLanguage("en")
		builder.SetProtocol("http")

		response, err := config.BuildResponse(builder)
		if err == nil {
			t.Error("Expected error for nonexistent message template")
		}
		if response != nil {
			t.Error("Expected nil response for nonexistent template")
		}
		if !strings.Contains(err.Error(), "message template not found") {
			t.Errorf("Expected error about template not found, got: %v", err)
		}
	})

	t.Run("Empty parameters", func(t *testing.T) {
		builder := NewResponseBuilder("welcome")
		builder.SetLanguage("en")
		builder.SetProtocol("http")
		// No parameters set

		response, err := config.BuildResponse(builder)
		if err != nil {
			t.Errorf("BuildResponse failed: %v", err)
		}

		if response.Message != "Welcome $name" {
			t.Errorf("Expected message 'Welcome $name', got %s", response.Message)
		}
	})

	t.Run("Multiple parameters", func(t *testing.T) {
		builder := NewResponseBuilder("welcome")
		builder.SetLanguage("en")
		builder.SetProtocol("http")
		builder.SetParam("name", "Alice")
		builder.SetParam("place", "World")

		response, err := config.BuildResponse(builder)
		if err != nil {
			t.Errorf("BuildResponse failed: %v", err)
		}

		if response.Message != "Welcome Alice" {
			t.Errorf("Expected message 'Welcome Alice', got %s", response.Message)
		}
	})
}

// TestSubstituteParams tests substituteParams function
func TestSubstituteParams(t *testing.T) {
	tests := []struct {
		name     string
		template string
		params   map[string]any
		expected string
	}{
		{
			name:     "Single parameter",
			template: "Hello $name",
			params:   map[string]any{"name": "John"},
			expected: "Hello John",
		},
		{
			name:     "Multiple parameters",
			template: "Hello $name, welcome to $place",
			params:   map[string]any{"name": "Jane", "place": "World"},
			expected: "Hello Jane, welcome to World",
		},
		{
			name:     "No parameters",
			template: "Hello world",
			params:   map[string]any{},
			expected: "Hello world",
		},
		{
			name:     "Nil parameters",
			template: "Hello $name",
			params:   nil,
			expected: "Hello $name",
		},
		{
			name:     "Parameter not found",
			template: "Hello $name",
			params:   map[string]any{"other": "value"},
			expected: "Hello $name",
		},
		{
			name:     "Empty parameter value",
			template: "Hello $name",
			params:   map[string]any{"name": ""},
			expected: "Hello ",
		},
		{
			name:     "Numeric parameter",
			template: "User ID: $id",
			params:   map[string]any{"id": 123},
			expected: "User ID: 123",
		},
		{
			name:     "Boolean parameter",
			template: "Status: $active",
			params:   map[string]any{"active": true},
			expected: "Status: true",
		},
		{
			name:     "Float parameter",
			template: "Score: $score",
			params:   map[string]any{"score": 98.5},
			expected: "Score: 98.5",
		},
		{
			name:     "Repeated parameter",
			template: "Hello $name, $name is here",
			params:   map[string]any{"name": "Alice"},
			expected: "Hello Alice, Alice is here",
		},
		{
			name:     "Complex template",
			template: "User $name (ID: $id) has $count items in $category",
			params: map[string]any{
				"name":     "Bob",
				"id":       456,
				"count":    5,
				"category": "electronics",
			},
			expected: "User Bob (ID: 456) has 5 items in electronics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteParams(tt.template, tt.params)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestResponseBuilderMethodChaining tests method chaining
func TestResponseBuilderMethodChaining(t *testing.T) {
	builder := NewResponseBuilder("welcome").
		WithContext(WithLanguage(WithProtocol(context.Background(), "http"), "en")).
		SetLanguage("id").
		SetProtocol("grpc").
		SetParam("name", "John").
		SetParams(map[string]any{"age": 25, "city": "Jakarta"}).
		SetData("user_id", 123).
		SetDatas(map[string]any{"role": "admin", "permissions": []string{"read", "write"}}).
		SetMeta("request_id", "req-123").
		SetMetas(map[string]any{"timestamp": 1640995200, "version": "1.0.0"}).
		SetError(errors.New("test error"))

	// Verify all settings
	if builder.MessageKey != "welcome" {
		t.Errorf("Expected MessageKey 'welcome', got %s", builder.MessageKey)
	}
	if builder.Language != "id" {
		t.Errorf("Expected Language 'id', got %s", builder.Language)
	}
	if builder.Protocol != "grpc" {
		t.Errorf("Expected Protocol 'grpc', got %s", builder.Protocol)
	}
	if builder.Params["name"] != "John" {
		t.Errorf("Expected name 'John', got %v", builder.Params["name"])
	}
	if builder.Params["age"] != 25 {
		t.Errorf("Expected age 25, got %v", builder.Params["age"])
	}
	if builder.Data["user_id"] != 123 {
		t.Errorf("Expected user_id 123, got %v", builder.Data["user_id"])
	}
	if builder.Meta["request_id"] != "req-123" {
		t.Errorf("Expected request_id 'req-123', got %v", builder.Meta["request_id"])
	}
	if builder.ErrorData == nil {
		t.Error("Expected error to be set")
	}
	if !builder.IsBuiltError {
		t.Error("Expected IsBuiltError to be true")
	}
}

// TestResponseBuilderEdgeCases tests edge cases
func TestResponseBuilderEdgeCases(t *testing.T) {
	t.Run("SetParam with nil Params", func(t *testing.T) {
		builder := &responseBuilder{MessageKey: "test"}
		builder.SetParam("key", "value")
		if builder.Params == nil {
			t.Error("Expected Params to be initialized")
		}
		if builder.Params["key"] != "value" {
			t.Errorf("Expected key 'value', got %v", builder.Params["key"])
		}
	})

	t.Run("SetData with nil Data", func(t *testing.T) {
		builder := &responseBuilder{MessageKey: "test"}
		builder.SetData("key", "value")
		if builder.Data == nil {
			t.Error("Expected Data to be initialized")
		}
		if builder.Data["key"] != "value" {
			t.Errorf("Expected key 'value', got %v", builder.Data["key"])
		}
	})

	t.Run("SetMeta with nil Meta", func(t *testing.T) {
		builder := &responseBuilder{MessageKey: "test"}
		builder.SetMeta("key", "value")
		if builder.Meta == nil {
			t.Error("Expected Meta to be initialized")
		}
		if builder.Meta["key"] != "value" {
			t.Errorf("Expected key 'value', got %v", builder.Meta["key"])
		}
	})

	t.Run("WithContext with nil context", func(t *testing.T) {
		builder := NewResponseBuilder("test")
		builder.WithContext(nil)
		// Should not panic
	})

	t.Run("Error with empty builder", func(t *testing.T) {
		builder := &responseBuilder{MessageKey: "test"}
		errorStr := builder.Error()
		if errorStr == "" {
			t.Error("Expected non-empty error string")
		}
	})
}

// BenchmarkResponseBuilderCreation benchmarks response builder creation
func BenchmarkResponseBuilderCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewResponseBuilder("test")
	}
}

// BenchmarkResponseBuilderMethodChaining benchmarks method chaining
func BenchmarkResponseBuilderMethodChaining(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewResponseBuilder("test").
			SetLanguage("en").
			SetProtocol("http").
			SetParam("name", "John").
			SetData("user_id", 123).
			SetMeta("request_id", "req-123")
	}
}

// BenchmarkSubstituteParams benchmarks parameter substitution
func BenchmarkSubstituteParams(b *testing.B) {
	template := "Hello $name, welcome to $place. Your ID is $id and you have $count items."
	params := map[string]any{
		"name":  "John",
		"place": "World",
		"id":    123,
		"count": 5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = substituteParams(template, params)
	}
}

// BenchmarkBuildResponse benchmarks BuildResponse method
func BenchmarkBuildResponse(b *testing.B) {
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
				},
			},
		},
		DefaultLanguage: "en",
		Languages:       []string{"en"},
	}

	builder := NewResponseBuilder("welcome")
	builder.SetLanguage("en")
	builder.SetProtocol("http")
	builder.SetParam("name", "John")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = config.BuildResponse(builder)
	}
}
