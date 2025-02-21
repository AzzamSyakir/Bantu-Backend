package middleware_test

// import (
// 	"bantu-backend/src/internal/middlewares"
// 	"testing"
// )

// type TestMiddleware struct {
// 	Test       *testing.T
// 	Middleware *middlewares.Middleware
// }

// func NewTestMiddleware(test *testing.T) *TestMiddleware {
// 	return &TestMiddleware{
// 		Test:       test,
// 		Middleware: middlewares.NewMiddleware(),
// 	}
// }
// func (testMiddleware *TestMiddleware) Start() {
// 	testMiddleware.Test.Run("TestMiddleware_TestCors", testMiddleware.TestCors)
// 	testMiddleware.Test.Run("TestMiddleware_TestInputValidation", testMiddleware.TestCors)
// 	testMiddleware.Test.Run("TestMiddleware_TestRateLimit", testMiddleware.TestCors)
// }
// func (testMiddleware *TestMiddleware) TestCors(t *testing.T) {
// 	t.Parallel()
// 	testMiddleware.Middleware.InputValidationMiddleware()

// }
// func (testMiddleware *TestMiddleware) TestInputValidation(t *testing.T) {
// 	t.Parallel()

// }
// func (testMiddleware *TestMiddleware) TestRateLimit(t *testing.T) {
// 	t.Parallel()

// }
