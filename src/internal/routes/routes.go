package routes

import (
	"bantu-backend/src/internal/controllers"

	"github.com/gorilla/mux"
)

type Route struct {
	Router                *mux.Router
	AuthController        *controllers.AuthController
	ChatController        *controllers.ChatController
	JobController         *controllers.JobController
	ProposalController    *controllers.ProposalController
	TransactionController *controllers.TransactionController
	UserController        *controllers.UserController
}

func NewRoute(router *mux.Router, authController *controllers.AuthController, chatController *controllers.ChatController, jobController *controllers.JobController, proposalController *controllers.ProposalController, transactionController *controllers.TransactionController) *Route {
	subRouter := router.PathPrefix("/api").Subrouter()
	return &Route{
		Router:                subRouter,
		AuthController:        authController,
		ChatController:        chatController,
		JobController:         jobController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
	}
}

func (r *Route) Register() {
	// eg
	// r.Router.HandleFunc("/jobs", r.JobController.CreateJob).Methods("POST")
	// r.Router.HandleFunc("/jobs", r.JobController.GetJobs).Methods("GET")
	// r.Router.HandleFunc("/jobs/{id}", r.JobController.GetJobByID).Methods("GET")
	// r.Router.HandleFunc("/jobs/{id}", r.JobController.UpdateJob).Methods("PUT")
	// r.Router.HandleFunc("/jobs/{id}", r.JobController.DeleteJob).Methods("DELETE")
	// r.Router.HandleFunc("/jobs/{id}/apply", r.JobController.ApplyJob).Methods("POST")
	// r.Router.HandleFunc("/jobs/{id}/proposals", r.ProposalController.CreateProposal).Methods("POST")
	// r.Router.HandleFunc("/jobs/{id}/proposals", r.ProposalController.GetProposals).Methods("GET")
	// r.Router.HandleFunc("/jobs/{id}/proposals/{proposalId}", r.ProposalController.UpdateProposal).Methods("PUT")
	// r.Router.HandleFunc("/jobs/{id}/payment", r.PaymentController.ProcessPayment).Methods("POST")
}
