package routes

import (
	"bantu-backend/src/internal/controllers"
	"bantu-backend/src/internal/middleware"

	"github.com/gorilla/mux"
)

type Route struct {
	Middleware            *middleware.Middleware
	Router                *mux.Router
	AuthController        *controllers.AuthController
	ChatController        *controllers.ChatController
	JobController         *controllers.JobController
	ProposalController    *controllers.ProposalController
	TransactionController *controllers.TransactionController
	UserController        *controllers.UserController
}

func NewRoute(
	router *mux.Router,
	middleware *middleware.Middleware,
	authController *controllers.AuthController,
	chatController *controllers.ChatController,
	jobController *controllers.JobController,
	proposalController *controllers.ProposalController,
	transactionController *controllers.TransactionController,
) *Route {
	subRouter := router.PathPrefix("/api").Subrouter()
	return &Route{
		Middleware:            middleware,
		Router:                subRouter,
		AuthController:        authController,
		ChatController:        chatController,
		JobController:         jobController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
	}
}

func (route *Route) Register() {
	route.Router.HandleFunc("/register", route.AuthController.Register).Methods("POST")
	route.Router.HandleFunc("/login", route.AuthController.Login).Methods("POST")

	route.Router.HandleFunc("/admin/register", route.AuthController.AdminRegister).Methods("POST")
	route.Router.HandleFunc("/admin/login", route.AuthController.AdminLogin).Methods("POST")

	route.Router.HandleFunc("/jobs", route.JobController.GetJobs).Methods("GET")
	route.Router.HandleFunc("/jobs", route.JobController.CreateJob).Methods("POST")
	route.Router.HandleFunc("/jobs/{id}", route.JobController.GetJobByID).Methods("GET")
	route.Router.HandleFunc("/jobs/{id}", route.JobController.UpdateJob).Methods("PUT")
	route.Router.HandleFunc("/jobs/{id}", route.JobController.DeleteJob).Methods("DELETE")

	route.Router.HandleFunc("/jobs/{id}/proposals", route.ProposalController.GetProposals).Methods("GET")
	route.Router.HandleFunc("/jobs/{id}/proposal", route.ProposalController.CreateProposal).Methods("POST")
	route.Router.HandleFunc("/jobs/{id}/proposal/{proposalId}", route.ProposalController.UpdateProposal).Methods("PUT")
	route.Router.HandleFunc("/jobs/{id}/proposal/{proposalId}/accept", route.ProposalController.AcceptProposal).Methods("PUT")

	route.Router.HandleFunc("/transaction/wallet/topup", route.TransactionController.TopUpBalance).Methods("POST")
	route.Router.HandleFunc("/transaction/wallet/pay/{proposalId}", route.TransactionController.PayFreelancer).Methods("POST")
}
