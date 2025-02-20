package routes

import (
	"bantu-backend/src/internal/controllers"
	"bantu-backend/src/internal/middlewares"

	"github.com/gorilla/mux"
)

type Route struct {
	Middleware            *middlewares.Middleware
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
	middleware *middlewares.Middleware,
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
	// eg
	r.Router.HandleFunc("/jobs", r.JobController.GetJobs).Methods("GET")
	r.Router.HandleFunc("/jobs", r.JobController.CreateJob).Methods("POST")
	r.Router.HandleFunc("/jobs/{id}", r.JobController.GetJobByID).Methods("GET")
	r.Router.HandleFunc("/jobs/{id}", r.JobController.UpdateJob).Methods("PUT")
	r.Router.HandleFunc("/jobs/{id}", r.JobController.DeleteJob).Methods("DELETE")
	r.Router.HandleFunc("/jobs/{id}/apply", r.JobController.ApplyJob).Methods("POST")
	// r.Router.HandleFunc("/jobs/{id}/proposals", r.ProposalController.CreateProposal).Methods("POST")
	// r.Router.HandleFunc("/jobs/{id}/proposals", r.ProposalController.GetProposals).Methods("GET")
	// r.Router.HandleFunc("/jobs/{id}/proposals/{proposalId}", r.ProposalController.UpdateProposal).Methods("PUT")
	// r.Router.HandleFunc("/jobs/{id}/payment", r.PaymentController.ProcessPayment).Methods("POST")
}
