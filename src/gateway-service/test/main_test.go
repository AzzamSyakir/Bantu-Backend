package test

// import (
// 	"bantu-backend/src/gateway-service/test/web"
// 	"fmt"
// 	"os"
// 	"testing"
// )

// func Test(t *testing.T) {
// 	chdirErr := os.Chdir("../../../.")
// 	if chdirErr != nil {
// 		t.Fatal(chdirErr)
// 	}
// 	fmt.Println("TestWeb started.")
// 	authWeb := web.NewGatewayWeb(t)
// 	authWeb.Start()

// 	userWeb := web.NewUserWeb(t)
// 	userWeb.Start()

// 	productWeb := web.NewProductWeb(t)
// 	productWeb.Start()

// 	categoryWeb := web.NewCategoryWeb(t)
// 	categoryWeb.Start()

// 	orderWeb := web.NewOrderWeb(t)
// 	orderWeb.Start()
// 	fmt.Println("TestWeb finished.")
// }
