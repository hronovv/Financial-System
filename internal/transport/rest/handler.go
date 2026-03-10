package rest

import (
	"financial_system/internal/service"
	"net/http"

	_ "financial_system/cmd/app/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services  *service.Services
	jwtSecret string
}

func NewHandler(services *service.Services, jwtSecret string) *Handler {
	return &Handler{
		services:  services,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// auth section
	auth := router.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
	auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost)

	//client section
	client := router.PathPrefix("/client").Subrouter()
	client.Use(h.authMiddleware("client"))

	client.HandleFunc("/banks", h.getBanks).Methods(http.MethodGet)
	client.HandleFunc("/enterprises", h.getEnterprises).Methods(http.MethodGet)

	client.HandleFunc("/accounts", h.openAccount).Methods(http.MethodPost)
	client.HandleFunc("/accounts/{id:[0-9]+}", h.closeAccount).Methods(http.MethodDelete)
	client.HandleFunc("/accounts/transfer", h.transferFromAccount).Methods(http.MethodPost)
	client.HandleFunc("/accounts/history", h.getAccountHistory).Methods(http.MethodGet)

	client.HandleFunc("/deposits", h.openDeposit).Methods(http.MethodPost)
	client.HandleFunc("/deposits/{id:[0-9]+}", h.closeDeposit).Methods(http.MethodDelete)
	client.HandleFunc("/deposits/transfer", h.transferFromDeposit).Methods(http.MethodPost)
	client.HandleFunc("/deposits/{id:[0-9]+}/accumulate", h.accumulateDeposit).Methods(http.MethodPost)

	client.HandleFunc("/salary-project/apply", h.applyForSalaryProject).Methods(http.MethodPost)
	client.HandleFunc("/salary-project/receive", h.receiveSalary).Methods(http.MethodPost)

	//manager section
	manager := router.PathPrefix("/manager").Subrouter()
	manager.Use(h.authMiddleware("manager"))

	manager.HandleFunc("/users/{id:[0-9]+}/approve", h.approveUser).Methods(http.MethodPost)
	manager.HandleFunc("/users/{id:[0-9]+}/history", h.getUserHistory).Methods(http.MethodGet)

	manager.HandleFunc("/accounts/{id:[0-9]+}/block", h.blockAccount).Methods(http.MethodPost)
	manager.HandleFunc("/accounts/{id:[0-9]+}/unblock", h.unblockAccount).Methods(http.MethodPost)
	manager.HandleFunc("/deposits/{id:[0-9]+}/block", h.blockDeposit).Methods(http.MethodPost)
	manager.HandleFunc("/deposits/{id:[0-9]+}/unblock", h.unblockDeposit).Methods(http.MethodPost)

	manager.HandleFunc("/enterprises", h.getEnterprisesWithEmployees).Methods(http.MethodGet)
	manager.HandleFunc("/enterprises/{id:[0-9]+}/employees", h.addEmployeeToEnterprise).Methods(http.MethodPost)
	manager.HandleFunc("/enterprises/{enterprise_id:[0-9]+}/employees/{user_id:[0-9]+}", h.removeEmployeeFromEnterprise).Methods(http.MethodDelete)
	manager.HandleFunc("/salary-project/applications/{id:[0-9]+}/approve", h.approveSalaryApplication).Methods(http.MethodPost)

	// admin section
	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(h.authMiddleware("admin"))

	admin.HandleFunc("/logs", h.getAllLogs).Methods(http.MethodGet)
	admin.HandleFunc("/logs/undo-all", h.undoAllActions).Methods(http.MethodPost)
	admin.HandleFunc("/logs/{id:[0-9]+}/undo", h.undoAction).Methods(http.MethodPost)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
