package middleware_test

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/middleware"
	"bantu-backend/src/internal/rabbitmq/producer"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type TestMiddleware struct {
	Test       *testing.T
	Middleware *middleware.Middleware
}

func NewTestMiddleware(test *testing.T) *TestMiddleware {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
	}
	fmt.Println("Current directory:", dir)
	envConfig := configs.NewEnvConfig()
	servicesProducer := producer.CreateNewServicesProducer(envConfig.RabbitMq)
	rabbitmqConfig := configs.NewRabbitMqConfig(envConfig)
	return &TestMiddleware{
		Test:       test,
		Middleware: middleware.NewMiddleware(rabbitmqConfig, servicesProducer, envConfig),
	}
}
func (testMiddleware *TestMiddleware) Start() {
	testMiddleware.Test.Run("TestMiddleware_TestCors", testMiddleware.TestCorsMiddleware)
	testMiddleware.Test.Run("TestMiddleware_TestInputValidation", testMiddleware.TestInputValidationMiddleware)
	testMiddleware.Test.Run("TestMiddleware_TestRateLimit", testMiddleware.TestRateLimitMiddleware)
	// testMiddleware.Test.Run("TestMiddleware_TestApplyMiddleware", testMiddleware.TestApplyMiddleware)
}
func (testMiddleware *TestMiddleware) TestCorsMiddleware(t *testing.T) {
	t.Parallel()
	t.Log("Starting TestCorsMiddleware")

	// Simple final handler returning "OK"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Apply CORS middleware
	handler := testMiddleware.Middleware.CorsMiddleware(finalHandler)
	t.Log("CORS middleware applied")

	// Create GET request with Origin header
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Origin", "http://example.com")
	t.Log("Created request with Origin header:", req.Header.Get("Origin"))

	// Capture response
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Check if Access-Control-Allow-Origin is set to "*"
	t.Logf("Response Headers: %v", rec.Header())
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected Access-Control-Allow-Origin '*' but got '%s'", rec.Header().Get("Access-Control-Allow-Origin"))
	} else {
		t.Log("Access-Control-Allow-Origin correctly set to '*'")
	}

	t.Log("TestCorsMiddleware completed")
}

func (testMiddleware *TestMiddleware) TestInputValidationMiddleware(t *testing.T) {
	t.Parallel()
	t.Log("Starting TestInputValidationMiddleware")

	// Simple final handler returning "OK"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Apply Input Validation middleware
	handler := testMiddleware.Middleware.InputValidationMiddleware(finalHandler)
	t.Log("InputValidation middleware applied")

	// Create POST request with valid JSON body and Content-Type header
	req := httptest.NewRequest("POST", "http://example.com", strings.NewReader(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")
	t.Log("Created request with Content-Type:", req.Header.Get("Content-Type"))

	// Capture response
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	t.Logf("Response status: %d", rec.Code)

	// Check if status is OK (200)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rec.Code)
	} else {
		t.Log("InputValidationMiddleware passed with valid Content-Type")
	}

	t.Log("TestInputValidationMiddleware completed")
}

func (testMiddleware *TestMiddleware) TestRateLimitMiddleware(t *testing.T) {
	t.Parallel()
	t.Log("Starting TestRateLimitMiddleware")

	// Simple final handler returning "OK"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Apply Rate Limit middleware
	handler := testMiddleware.Middleware.RateLimitMiddleware(finalHandler)
	t.Log("RateLimit middleware applied")

	allowedCount := 0
	rejectedCount := 0
	totalRequests := 110

	// Send 110 GET requests to test rate limiting
	for i := 0; i < totalRequests; i++ {
		req := httptest.NewRequest("GET", "http://example.com", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code == http.StatusOK {
			allowedCount++
			t.Logf("Request %d allowed (status code %d)", i+1, rec.Code)
		} else if rec.Code == http.StatusTooManyRequests {
			rejectedCount++
			t.Logf("Request %d rejected (status code %d)", i+1, rec.Code)
		}
	}

	t.Logf("Total allowed requests: %d", allowedCount)
	t.Logf("Total rejected requests: %d", rejectedCount)

	// Expect exactly 100 allowed and 10 rejected requests
	if allowedCount != 100 {
		t.Errorf("expected allowed requests to be 100, but got %d", allowedCount)
	}
	if rejectedCount != 10 {
		t.Errorf("expected rejected requests to be 10, but got %d", rejectedCount)
	}

	t.Log("TestRateLimitMiddleware completed")
}

func (testMiddleware *TestMiddleware) TestApplyMiddleware(t *testing.T) {
	t.Parallel()
	t.Log("Starting TestApplyMiddleware")

	// Simple final handler returning "OK"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Apply all middlewares together
	handler := testMiddleware.Middleware.ApplyMiddleware(finalHandler)
	t.Log("Combined middlewares applied")

	// Create a valid POST request with proper headers and JSON body
	req := httptest.NewRequest("POST", "http://example.com", strings.NewReader(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://example.com")
	t.Log("Created valid request with proper headers")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	t.Logf("Valid request response: status code %d, body: %s", rec.Code, rec.Body.String())

	// Check if valid request returns status OK and body "OK"
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("expected body 'OK' but got '%s'", rec.Body.String())
	} else {
		t.Log("Valid request passed successfully")
	}

	// Create an invalid request with wrong Content-Type header
	reqInvalid := httptest.NewRequest("POST", "http://example.com", strings.NewReader(`{"key": "value"}`))
	reqInvalid.Header.Set("Content-Type", "text/plain")
	t.Log("Created invalid request with wrong Content-Type")
	recInvalid := httptest.NewRecorder()
	handler.ServeHTTP(recInvalid, reqInvalid)
	t.Logf("Invalid request response: status code %d", recInvalid.Code)

	// Check if invalid request returns Unsupported Media Type (415)
	if recInvalid.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected status %d but got %d", http.StatusUnsupportedMediaType, recInvalid.Code)
	} else {
		t.Log("Invalid request correctly rejected with Unsupported Media Type")
	}

	t.Log("TestApplyMiddleware completed")
}
