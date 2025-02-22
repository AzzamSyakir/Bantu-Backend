package test

import (
	"bantu-backend/src/test/middleware_test"
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	chdirErr := os.Chdir("../../.")
	if chdirErr != nil {
		t.Fatal(chdirErr)
	}
	fmt.Println("Test started.")
	middlewareTest := middleware_test.NewTestMiddleware(t)
	middlewareTest.Start()
	fmt.Println("Test finished.")
}
