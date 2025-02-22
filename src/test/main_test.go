package test

import (
	"bantu-backend/src/test/middleware_test"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test(t *testing.T) {
	chdirErr := os.Chdir("../../.")
	if chdirErr != nil {
		t.Fatal(chdirErr)
	}
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	fmt.Println("Test started.")
	middlewareTest := middleware_test.NewTestMiddleware(t)
	middlewareTest.Start()
	fmt.Println("Test finished.")
}
