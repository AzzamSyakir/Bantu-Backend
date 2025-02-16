package web

// import (
// 	"bantu-backend/src/gateway-service/entity"
// 	model_request "bantu-backend/src/gateway-service/model/request/controller"
// 	model_response "bantu-backend/src/gateway-service/model/response"
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"golang.org/x/crypto/bcrypt"

// 	"testing"

// 	"github.com/guregu/null"
// 	"github.com/stretchr/testify/assert"
// )

// type GatewayWeb struct {
// 	Test *testing.T
// 	Path string
// }

// func NewGatewayWeb(test *testing.T) *GatewayWeb {
// 	authWeb := &GatewayWeb{
// 		Test: test,
// 		Path: "auths",
// 	}
// 	return authWeb
// }
// func (authWeb *GatewayWeb) Start() {
// 	authWeb.Test.Run("GatewayWeb_Register_Succeed", authWeb.Register)
// 	authWeb.Test.Run("GatewayWeb_Login_Succeed", authWeb.Login)
// 	authWeb.Test.Run("GatewayWeb_Logout_Succeed", authWeb.Logout)
// 	authWeb.Test.Run("GatewayWeb_GetNewAccessToken_Succeed", authWeb.GetNewAccessToken)
// }

// func (authWeb *GatewayWeb) Register(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	defer testWeb.AllSeeder.Down()

// 	mockGateway := testWeb.AllSeeder.User.UserMock.Data[0]

// 	bodyRequest := &model_request.RegisterRequest{}
// 	bodyRequest.Name = null.NewString(mockGateway.Name.String, true)
// 	bodyRequest.Email = null.NewString(mockGateway.Email.String, true)
// 	bodyRequest.Password = null.NewString(mockGateway.Password.String, true)
// 	bodyRequest.Balance = null.NewInt(mockGateway.Balance.Int64, true)

// 	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
// 	if marshalErr != nil {
// 		t.Fatal(marshalErr)
// 	}
// 	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

// 	url := fmt.Sprintf("%s/%s/register", testWeb.Server.URL, authWeb.Path)
// 	request, newRequestErr := http.NewRequest(http.MethodPost, url, bodyRequestBuffer)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}

// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.User]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusCreated, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, mockGateway.Name.String, bodyResponse.Data.Name.String)
// 	assert.Equal(t, mockGateway.Email.String, bodyResponse.Data.Email.String)
// 	assert.Equal(t, mockGateway.Balance.Int64, bodyResponse.Data.Balance.Int64)
// 	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(mockGateway.Password.String)))

// 	newUserMock := bodyResponse.Data
// 	testWeb.AllSeeder.User.UserMock.Data = append(testWeb.AllSeeder.User.UserMock.Data, newUserMock)
// }

// func (authWeb *GatewayWeb) Login(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedUserMock := testWeb.AllSeeder.User.UserMock.Data[0]

// 	bodyRequest := &model_request.LoginRequest{}
// 	bodyRequest.Email = selectedUserMock.Email
// 	bodyRequest.Password = selectedUserMock.Password

// 	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
// 	if marshalErr != nil {
// 		t.Fatal(marshalErr)
// 	}
// 	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

// 	url := fmt.Sprintf("%s/%s/login", testWeb.Server.URL, authWeb.Path)
// 	request, newRequestErr := http.NewRequest(http.MethodPost, url, bodyRequestBuffer)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.Session]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}
// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedUserMock.Id, bodyResponse.Data.UserId)
// }

// func (authWeb *GatewayWeb) Logout(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	url := fmt.Sprintf("%s/%s/logout", testWeb.Server.URL, authWeb.Path)
// 	request, newRequestErr := http.NewRequest(http.MethodPost, url, nil)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}

// 	request.Header.Set("Gatewayorization", "Bearer "+selectedSessionMock.AccessToken.String)

// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// }
// func (authWeb *GatewayWeb) GetNewAccessToken(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	url := fmt.Sprintf("%s/%s/access-token", testWeb.Server.URL, authWeb.Path)
// 	request, newRequest := http.NewRequest(http.MethodPost, url, nil)
// 	if newRequest != nil {
// 		t.Fatal(newRequest)
// 	}

// 	request.Header.Set("Gatewayorization", "Bearer "+selectedSessionMock.RefreshToken.String)

// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	responseBody := &model_response.Response[*entity.Session]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(responseBody)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.True(t, selectedSessionMock.UpdatedAt.Time.Before(responseBody.Data.UpdatedAt.Time))
// }
