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

func (r *Route) Register() {
	r.Router.HandleFunc("/register", r.AuthController.Register).Methods("POST")
	r.Router.HandleFunc("/login", r.AuthController.Login).Methods("POST")

	r.Router.HandleFunc("/admin/register", r.AuthController.AdminRegister).Methods("POST")
	r.Router.HandleFunc("/admin/login", r.AuthController.AdminLogin).Methods("POST")

	r.Router.HandleFunc("/jobs", r.JobController.GetJobs).Methods("GET")
	r.Router.HandleFunc("/jobs", r.JobController.CreateJob).Methods("POST")
	r.Router.HandleFunc("/jobs/{id}", r.JobController.GetJobByID).Methods("GET")
	r.Router.HandleFunc("/jobs/{id}", r.JobController.UpdateJob).Methods("PUT")
	r.Router.HandleFunc("/jobs/{id}", r.JobController.DeleteJob).Methods("DELETE")

	r.Router.HandleFunc("/jobs/{id}/proposals", r.ProposalController.GetProposals).Methods("GET")
	r.Router.HandleFunc("/jobs/{id}/proposal", r.ProposalController.CreateProposal).Methods("POST")
	r.Router.HandleFunc("/jobs/{id}/proposal/{proposalId}", r.ProposalController.UpdateProposal).Methods("PUT")
	r.Router.HandleFunc("/jobs/{id}/proposal/{proposalId}/accept", r.ProposalController.AcceptProposal).Methods("PUT")
	// r.Router.HandleFunc("/jobs/{id}/payment", r.PaymentController.ProcessPayment).Methods("POST")
}
