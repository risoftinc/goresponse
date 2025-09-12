package goresponse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ExampleResponseUsage demonstrates various usage patterns of ResponseBuilder and ResponseManager
func ExampleResponseUsage() {
	// Load configuration first
	source := ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	config, err := LoadConfig(source)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	fmt.Println("=== ResponseBuilder and ResponseManager Examples ===")

	// Example 1: Basic Success Response
	exampleBasicSuccess(config)

	// Example 2: Error Response with Parameters
	exampleErrorResponse(config)

	// Example 3: Context Integration
	exampleContextIntegration(config)

	// Example 4: Data Payload Response
	exampleDataPayload(config)

	// Example 4.5: Metadata Response
	exampleMetadataResponse(config)

	// Example 5: Service Layer Integration
	exampleServiceLayer(config)

	// Example 6: Handler Layer Integration
	exampleHandlerLayer(config)

	// Example 7: Multiple Language Support
	exampleMultipleLanguages(config)

	// Example 8: Protocol-Specific Responses
	exampleProtocolSpecific(config)

	// Example 9: Error Recovery and Parsing
	exampleErrorRecovery(config)

	// Example 10: Complex Business Logic
	exampleComplexBusinessLogic(config)
}

// Example 1: Basic Success Response
func exampleBasicSuccess(config *ResponseConfig) {
	fmt.Println("1. Basic Success Response:")
	fmt.Println("------------------------")

	// Create response builder
	builder := NewResponseBuilder("user_created")
	builder.SetParam("name", "John Doe")
	builder.SetParam("email", "john@example.com")

	// Build response
	response, err := config.BuildResponse(builder)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Print response
	printResponse(response)
	fmt.Println()
}

// Example 2: Error Response with Parameters
func exampleErrorResponse(config *ResponseConfig) {
	fmt.Println("2. Error Response with Parameters:")
	fmt.Println("----------------------------------")

	// Create error response
	builder := NewResponseBuilder("validation_failed")
	builder.SetParam("field", "email")
	builder.SetParam("value", "invalid-email")
	builder.SetParam("rule", "must be valid email format")
	builder.SetError(fmt.Errorf("email validation failed"))

	// Build response
	response, err := config.BuildResponse(builder)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Print response
	printResponse(response)
	fmt.Println()
}

// Example 3: Context Integration
func exampleContextIntegration(config *ResponseConfig) {
	fmt.Println("3. Context Integration:")
	fmt.Println("----------------------")

	// Create context with language and protocol
	ctx := context.Background()
	ctx = WithLanguage(ctx, "id")
	ctx = WithProtocol(ctx, "http")

	// Create response builder with context
	builder := NewResponseBuilder("welcome_message")
	builder.SetParam("name", "Ahmad")
	builder.WithContext(ctx)

	// Build response
	response, err := config.BuildResponse(builder)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Print response
	printResponse(response)
	fmt.Println()
}

// Example 4: Data Payload Response
func exampleDataPayload(config *ResponseConfig) {
	fmt.Println("4. Data Payload Response:")
	fmt.Println("-------------------------")

	// Sample user data
	users := []map[string]any{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
	}

	// Create response with data
	builder := NewResponseBuilder("users_retrieved")
	builder.SetParam("count", len(users))
	builder.SetData("users", users)
	builder.SetData("pagination", map[string]any{
		"page":     1,
		"per_page": 10,
		"total":    len(users),
	})

	// Build response
	response, err := config.BuildResponse(builder)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Print response
	printResponse(response)
	fmt.Println()
}

// Example 4.5: Metadata Response
func exampleMetadataResponse(config *ResponseConfig) {
	fmt.Println("4.5. Metadata Response:")
	fmt.Println("----------------------")

	// Response with metadata
	builder := NewResponseBuilder("users_retrieved")
	builder.SetParam("count", 3)
	builder.SetData("users", []map[string]any{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
	})

	// Add metadata
	builder.SetMeta("request_id", "req-12345")
	builder.SetMeta("processing_time", "150ms")
	builder.SetMeta("cache_hit", true)
	builder.SetMeta("api_version", "v1.2.0")
	builder.SetMeta("rate_limit", map[string]any{
		"remaining":  95,
		"reset_time": "2024-01-01T12:00:00Z",
	})

	response, err := config.BuildResponse(builder)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	printResponse(response)
	fmt.Println()

	// Example with multiple metadata at once
	fmt.Println("Multiple Metadata at Once:")
	builder2 := NewResponseBuilder("data_processed")
	builder2.SetParam("operation", "bulk_update")
	builder2.SetData("updated_count", 25)
	builder2.SetMetas(map[string]any{
		"batch_id":     "batch-789",
		"start_time":   "2024-01-01T10:00:00Z",
		"end_time":     "2024-01-01T10:05:00Z",
		"success_rate": 0.96,
		"errors":       []string{"row_15", "row_23"},
	})

	response2, err := config.BuildResponse(builder2)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	printResponse(response2)
	fmt.Println()
}

// Example 5: Service Layer Integration
func exampleServiceLayer(config *ResponseConfig) {
	fmt.Println("5. Service Layer Integration:")
	fmt.Println("-----------------------------")

	// Simulate service method
	userService := &UserService{config: config}

	// Test cases
	testCases := []struct {
		name  string
		email string
		age   int
	}{
		{"Valid User", "valid@example.com", 25},
		{"Invalid Email", "invalid-email", 30},
		{"Underage", "young@example.com", 15},
	}

	for _, tc := range testCases {
		fmt.Printf("Testing: %s\n", tc.name)

		err := userService.CreateUser(context.Background(), tc.name, tc.email, tc.age)
		if err != nil {
			// Check if it's a ResponseBuilder error
			if builder, ok := ParseResponseBuilderError(err); ok {
				response, _ := config.BuildResponse(builder)
				printResponse(response)
			} else {
				fmt.Printf("Regular error: %v\n", err)
			}
		} else {
			fmt.Println("User created successfully")
		}
		fmt.Println()
	}
}

// Example 6: Handler Layer Integration
func exampleHandlerLayer(config *ResponseConfig) {
	fmt.Println("6. Handler Layer Integration:")
	fmt.Println("-----------------------------")

	// Simulate HTTP handler
	handler := &UserHandler{config: config}

	// Test different scenarios
	scenarios := []struct {
		name    string
		request *CreateUserRequest
	}{
		{
			name: "Valid Request",
			request: &CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
		},
		{
			name: "Invalid Email",
			request: &CreateUserRequest{
				Name:  "Jane Doe",
				Email: "invalid-email",
				Age:   30,
			},
		},
		{
			name: "Underage",
			request: &CreateUserRequest{
				Name:  "Young User",
				Email: "young@example.com",
				Age:   15,
			},
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("Scenario: %s\n", scenario.name)
		response := handler.CreateUser(scenario.request)
		printResponse(response)
		fmt.Println()
	}
}

// Example 7: Multiple Language Support
func exampleMultipleLanguages(config *ResponseConfig) {
	fmt.Println("7. Multiple Language Support:")
	fmt.Println("-----------------------------")

	languages := []string{"en", "id", "es", "fr"}
	messageKey := "welcome_message"

	for _, lang := range languages {
		fmt.Printf("Language: %s\n", lang)

		builder := NewResponseBuilder(messageKey)
		builder.SetLanguage(lang)
		builder.SetProtocol("http")
		builder.SetParam("name", "User")

		response, err := config.BuildResponse(builder)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		fmt.Printf("Message: %s\n", response.Message)
		fmt.Printf("Code: %d\n", response.Code)
		fmt.Println()
	}
}

// Example 8: Protocol-Specific Responses
func exampleProtocolSpecific(config *ResponseConfig) {
	fmt.Println("8. Protocol-Specific Responses:")
	fmt.Println("-------------------------------")

	protocols := []string{"http", "grpc", "rest", "graphql"}
	messageKey := "validation_failed"

	for _, protocol := range protocols {
		fmt.Printf("Protocol: %s\n", protocol)

		builder := NewResponseBuilder(messageKey)
		builder.SetProtocol(protocol)
		builder.SetParam("field", "email")
		builder.SetParam("value", "invalid")

		response, err := config.BuildResponse(builder)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		fmt.Printf("Message: %s\n", response.Message)
		fmt.Printf("Code: %d\n", response.Code)
		fmt.Println()
	}
}

// Example 9: Error Recovery and Parsing
func exampleErrorRecovery(config *ResponseConfig) {
	fmt.Println("9. Error Recovery and Parsing:")
	fmt.Println("-------------------------------")

	// Create an error from ResponseBuilder
	builder := NewResponseBuilder("user_not_found")
	builder.SetParam("id", "123")
	builder.SetError(fmt.Errorf("user with id 123 not found"))

	// Convert to error
	err := builder.ToError()
	fmt.Printf("Original error: %v\n", err)

	// Parse the error back to ResponseBuilder
	parsedBuilder, ok := ParseResponseBuilderError(err)
	if ok {
		fmt.Println("Successfully parsed ResponseBuilder error")

		// Build response from parsed builder
		response, buildErr := config.BuildResponse(parsedBuilder)
		if buildErr != nil {
			log.Printf("Error building response: %v", buildErr)
		} else {
			printResponse(response)
		}
	} else {
		fmt.Println("Failed to parse ResponseBuilder error")
	}
	fmt.Println()
}

// Example 10: Complex Business Logic
func exampleComplexBusinessLogic(config *ResponseConfig) {
	fmt.Println("10. Complex Business Logic:")
	fmt.Println("---------------------------")

	// Simulate complex business logic
	orderService := &OrderService{config: config}

	// Test order processing
	order := &Order{
		ID:       "ORD-001",
		UserID:   "user-123",
		Items:    []string{"laptop", "mouse", "keyboard"},
		Total:    1500.00,
		Currency: "USD",
	}

	// Process order
	response := orderService.ProcessOrder(context.Background(), order)
	printResponse(response)
	fmt.Println()
}

// Helper function to print response
func printResponse(response *Response) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}
	fmt.Printf("Response: %s\n", string(jsonData))
}

// Mock structures for examples

type UserService struct {
	config *ResponseConfig
}

func (s *UserService) CreateUser(ctx context.Context, name, email string, age int) error {
	// Validation logic
	if email == "" || !isValidEmail(email) {
		builder := NewResponseBuilder("validation_failed")
		builder.SetParam("field", "email")
		builder.SetParam("value", email)
		builder.SetParam("rule", "must be valid email format")
		return builder.ToError()
	}

	if age < 18 {
		builder := NewResponseBuilder("validation_failed")
		builder.SetParam("field", "age")
		builder.SetParam("value", age)
		builder.SetParam("rule", "must be at least 18 years old")
		return builder.ToError()
	}

	// Simulate database error
	if email == "error@example.com" {
		builder := NewResponseBuilder("internal_error")
		builder.SetParam("service", "user_service")
		builder.SetError(fmt.Errorf("database connection failed"))
		return builder.ToError()
	}

	// Success case
	return nil
}

type UserHandler struct {
	config *ResponseConfig
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func (h *UserHandler) CreateUser(req *CreateUserRequest) *Response {
	// Simulate service call
	userService := &UserService{config: h.config}
	err := userService.CreateUser(context.Background(), req.Name, req.Email, req.Age)

	if err != nil {
		// Check if it's a ResponseBuilder error
		if builder, ok := ParseResponseBuilderError(err); ok {
			response, _ := h.config.BuildResponse(builder)
			return response
		}

		// Handle other errors
		builder := NewResponseBuilder("internal_error")
		builder.SetError(err)
		response, _ := h.config.BuildResponse(builder)
		return response
	}

	// Success response
	builder := NewResponseBuilder("user_created")
	builder.SetParam("name", req.Name)
	builder.SetParam("email", req.Email)
	response, _ := h.config.BuildResponse(builder)
	return response
}

type OrderService struct {
	config *ResponseConfig
}

type Order struct {
	ID       string   `json:"id"`
	UserID   string   `json:"user_id"`
	Items    []string `json:"items"`
	Total    float64  `json:"total"`
	Currency string   `json:"currency"`
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *Order) *Response {
	// Complex business logic simulation
	if order.Total <= 0 {
		builder := NewResponseBuilder("validation_failed")
		builder.SetParam("field", "total")
		builder.SetParam("value", order.Total)
		builder.SetParam("rule", "must be greater than 0")
		response, _ := s.config.BuildResponse(builder)
		return response
	}

	if len(order.Items) == 0 {
		builder := NewResponseBuilder("validation_failed")
		builder.SetParam("field", "items")
		builder.SetParam("value", "empty")
		builder.SetParam("rule", "must have at least one item")
		response, _ := s.config.BuildResponse(builder)
		return response
	}

	// Simulate processing
	time.Sleep(100 * time.Millisecond)

	// Success response with order data
	builder := NewResponseBuilder("order_processed")
	builder.SetParam("order_id", order.ID)
	builder.SetParam("total", order.Total)
	builder.SetParam("currency", order.Currency)
	builder.SetData("order", order)
	builder.SetData("processed_at", time.Now().Format(time.RFC3339))
	builder.SetData("status", "completed")

	response, _ := s.config.BuildResponse(builder)
	return response
}

// Helper function to validate email (simple validation)
func isValidEmail(email string) bool {
	return len(email) > 0 &&
		email != "invalid-email" &&
		email != "error@example.com"
}

// ExampleAsyncResponseUsage demonstrates ResponseBuilder with AsyncConfigManager
func ExampleAsyncResponseUsage() {
	fmt.Println("=== Async ResponseBuilder Examples ===")

	// Create async manager
	source := ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	asyncManager := NewAsyncConfigManager(source, 5*time.Minute)

	// Add callback for configuration changes
	asyncManager.AddCallback(func(oldConfig, newConfig *ResponseConfig) {
		fmt.Printf("Configuration updated! New default language: %s\n", newConfig.GetDefaultLanguage())
	})

	// Start async manager
	if err := asyncManager.Start(); err != nil {
		log.Printf("Error starting async manager: %v", err)
		return
	}
	defer asyncManager.Stop()

	// Wait for initial load
	time.Sleep(100 * time.Millisecond)

	// Example 1: Basic response with async config
	fmt.Println("1. Basic Response with Async Config:")
	fmt.Println("------------------------------------")

	builder := NewResponseBuilder("success")
	builder.SetParam("service", "async_service")
	builder.SetLanguage("en")
	builder.SetProtocol("http")

	// Get config from async manager (thread-safe)
	config := asyncManager.GetConfig()
	response, err := config.BuildResponse(builder)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	printResponse(response)
	fmt.Println()

	// Example 2: Multiple concurrent responses
	fmt.Println("2. Multiple Concurrent Responses:")
	fmt.Println("---------------------------------")

	// Simulate concurrent requests
	for i := 0; i < 5; i++ {
		go func(id int) {
			builder := NewResponseBuilder("concurrent_request")
			builder.SetParam("request_id", fmt.Sprintf("req-%d", id))
			builder.SetParam("timestamp", time.Now().Format(time.RFC3339))

			config := asyncManager.GetConfig()
			response, err := config.BuildResponse(builder)
			if err != nil {
				log.Printf("Error in goroutine %d: %v", id, err)
				return
			}

			fmt.Printf("Response %d: %s\n", id, response.Message)
		}(i)
	}

	// Wait for goroutines to complete
	time.Sleep(200 * time.Millisecond)
	fmt.Println()
}

// ExampleResponseBuilderChaining demonstrates method chaining patterns
func ExampleResponseBuilderChaining() {
	fmt.Println("=== ResponseBuilder Method Chaining Examples ===")

	// Load configuration
	source := ConfigSource{
		Method: "file",
		Path:   "config.json",
	}

	config, err := LoadConfig(source)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// Example 1: Basic chaining
	fmt.Println("1. Basic Method Chaining:")
	fmt.Println("-------------------------")

	response, err := config.BuildResponse(
		NewResponseBuilder("user_created").
			SetParam("name", "Alice").
			SetParam("email", "alice@example.com").
			SetLanguage("en").
			SetProtocol("http"),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	printResponse(response)
	fmt.Println()

	// Example 2: Complex chaining with data
	fmt.Println("2. Complex Chaining with Data:")
	fmt.Println("------------------------------")

	response, err = config.BuildResponse(
		NewResponseBuilder("data_retrieved").
			SetParams(map[string]any{
				"table":  "users",
				"count":  100,
				"filter": "active",
			}).
			SetDatas(map[string]any{
				"users": []map[string]any{
					{"id": 1, "name": "User 1"},
					{"id": 2, "name": "User 2"},
				},
				"pagination": map[string]any{
					"page": 1,
					"size": 10,
				},
			}).
			SetLanguage("id").
			SetProtocol("rest"),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	printResponse(response)
	fmt.Println()

	// Example 3: Error response chaining
	fmt.Println("3. Error Response Chaining:")
	fmt.Println("---------------------------")

	response, err = config.BuildResponse(
		NewResponseBuilder("validation_failed").
			SetParam("field", "password").
			SetParam("rule", "minimum 8 characters").
			SetError(fmt.Errorf("password too short")).
			SetMeta("validation_rule", "password_strength").
			SetMeta("attempt_count", 3).
			SetLanguage("en").
			SetProtocol("http"),
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	printResponse(response)
	fmt.Println()
}
