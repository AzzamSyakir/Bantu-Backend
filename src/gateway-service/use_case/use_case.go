package use_case

import (
	"bantu-backend/grpc/pb"
	"bantu-backend/src/gateway-service/config"
	"bantu-backend/src/gateway-service/delivery/grpc/client"
	"bantu-backend/src/gateway-service/entity"
	model_request "bantu-backend/src/gateway-service/model/request/controller"
	model_response "bantu-backend/src/gateway-service/model/response"
	"bantu-backend/src/gateway-service/repository"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type GatewayUseCase struct {
	pb.UnimplementedGatewayServiceServer
	DatabaseConfig    *config.DatabaseConfig
	GatewayRepository *repository.GatewayRepository
	Env               *config.EnvConfig
	userClient        *client.UserServiceClient
}

func NewGatewayUseCase(
	databaseConfig *config.DatabaseConfig,
	authRepository *repository.GatewayRepository,
	env *config.EnvConfig,
	initUserClient *client.UserServiceClient,
) *GatewayUseCase {
	authUseCase := &GatewayUseCase{
		UnimplementedGatewayServiceServer: pb.UnimplementedGatewayServiceServer{},
		userClient:                        initUserClient,
		DatabaseConfig:                    databaseConfig,
		GatewayRepository:                 authRepository,
		Env:                               env,
	}
	return authUseCase
}

func (authUseCase *GatewayUseCase) Login(request *model_request.LoginRequest) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.GatewayDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "GatewayUseCase Login failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundUser, err := authUseCase.userClient.GetUserByEmail(request.Email.String)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: foundUser.Message,
			Data:    nil,
		}
		return result, rollback
	}
	if foundUser.Data == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: foundUser.Message,
			Data:    nil,
		}
		return result, rollback
	}

	comparePasswordErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Data.Password), []byte(request.Password.String))
	if comparePasswordErr != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "GatewayUseCase Login is failed, password is not match.",
			Data:    nil,
		}
		return result, rollback
	}

	accessToken := null.NewString(uuid.NewString(), true)
	refreshToken := null.NewString(uuid.NewString(), true)
	currentTime := null.NewTime(time.Now(), true)
	accessTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Minute*10), true)
	refreshTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Hour*24*2), true)

	foundSession, err := authUseCase.GatewayRepository.GetOneByUserId(begin, foundUser.Data.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase Login failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession != nil {

		foundSession.AccessToken = accessToken
		foundSession.RefreshToken = refreshToken
		foundSession.AccessTokenExpiredAt = accessTokenExpiredAt
		foundSession.RefreshTokenExpiredAt = refreshTokenExpiredAt
		foundSession.UpdatedAt = currentTime
		patchedSession, err := authUseCase.GatewayRepository.PatchOneById(begin, foundSession.Id.String, foundSession)
		if err != nil {
			rollback := begin.Rollback()
			result = &model_response.Response[*entity.Session]{
				Code:    http.StatusBadRequest,
				Message: "GatewayUseCase Login failed, query updateSession  fail, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		commit := begin.Commit()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusOK,
			Message: "GatewayUseCase Login is succeed",
			Data:    patchedSession,
		}
		return result, commit
	}

	newSession := &entity.Session{
		Id:                    null.NewString(uuid.NewString(), true),
		UserId:                null.NewString(foundUser.Data.Id, true),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  accessTokenExpiredAt,
		RefreshTokenExpiredAt: refreshTokenExpiredAt,
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
	}

	createdSession, err := authUseCase.GatewayRepository.CreateSession(begin, newSession)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase Login failed, query createSession fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "GatewayUseCase Login is succeed",
		Data:    createdSession,
	}
	return result, commit
}

func (authUseCase *GatewayUseCase) Logout(accessToken string) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.GatewayDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "GatewayUseCase Logout failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundSession, err := authUseCase.GatewayRepository.FindOneByAccToken(begin, accessToken)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase Logout failed, Invalid token, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if foundSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase Logout is failed, session is not found by access token.",
			Data:    nil,
		}
		return result, rollback
	}
	deletedSession, err := authUseCase.GatewayRepository.DeleteOneById(begin, foundSession.Id.String)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase Logout failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if deletedSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase Logout failed, delete session failed",
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "GatewayUseCase Logout is succeed.",
		Data:    deletedSession,
	}
	return result, commit
}

func (authUseCase *GatewayUseCase) LogoutWithUserId(context context.Context, id *pb.ByUserId) (empty *pb.Empty, err error) {
	begin, err := authUseCase.DatabaseConfig.GatewayDB.Connection.Begin()
	if err != nil {
		begin.Rollback()
		return &pb.Empty{}, err
	}

	foundSession, err := authUseCase.GatewayRepository.GetOneByUserId(begin, id.Id)
	if err != nil {
		begin.Rollback()
		return &pb.Empty{}, err
	}
	if foundSession == nil {
		begin.Rollback()
		return &pb.Empty{}, err
	}
	_, err = authUseCase.GatewayRepository.DeleteOneByUserId(begin, foundSession.UserId.String)
	if err != nil {
		begin.Rollback()
		return &pb.Empty{}, err
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		return &pb.Empty{}, commitErr
	}

	return &pb.Empty{}, nil
}

func (authUseCase *GatewayUseCase) GetNewAccessToken(refreshToken string) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.GatewayDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "GatewayUseCase GetNewAccesToken failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	foundSession, err := authUseCase.GatewayRepository.FindOneByRefToken(begin, refreshToken)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase GetNewAccesToken failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase GetNewAccesToken  failed, session is not found by refresh token.",
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession.RefreshTokenExpiredAt.Time.Before(time.Now()) {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "GatewayUseCase GetNewAccessToken is failed, refresh token is expired.",
			Data:    nil,
		}
		return result, rollback
	}

	foundSession.AccessToken = null.NewString(uuid.NewString(), true)
	foundSession.UpdatedAt = null.NewTime(time.Now(), true)
	patchedSession, err := authUseCase.GatewayRepository.PatchOneById(begin, foundSession.Id.String, foundSession)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "GatewayUseCase GetNewAccesToken  failed, query to db fail," + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "GatewayUseCase GetNewAccessToken is succeed.",
		Data:    patchedSession,
	}
	return result, commit

}
