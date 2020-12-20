package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/najibulloShapoatov/crud/cmd/app/middleware"

	"github.com/gorilla/mux"

	"github.com/najibulloShapoatov/crud/pkg/customers"
	"github.com/najibulloShapoatov/crud/pkg/managers"
)

//Server ...
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
	managerSvc  *managers.Service
}

//NewServer ... создает новый сервер
func NewServer(m *mux.Router, cSvc *customers.Service, mSvc *managers.Service) *Server {
	return &Server{
		mux:         m,
		customerSvc: cSvc,
		managerSvc:  mSvc,
	}
}

// функция для запуска хендлеров через мукс
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init ... инициализация сервера
func (s *Server) Init() {

	customersAuthenticateMd := middleware.Authenticate(s.customerSvc.IDByToken)
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customersAuthenticateMd)

	customersSubrouter.HandleFunc("", s.handleCustomerRegistration).Methods("POST")
	customersSubrouter.HandleFunc("/token", s.handleCustomerGetToken).Methods("POST")
	customersSubrouter.HandleFunc("/products", s.handleCustomerGetProducts).Methods("GET")

	managersAuthenticateMd := middleware.Authenticate(s.managerSvc.IDByToken)
	managersSubRouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubRouter.Use(managersAuthenticateMd)
	managersSubRouter.HandleFunc("", s.handleManagerRegistration).Methods("POST")
	managersSubRouter.HandleFunc("/token", s.handleManagerGetToken).Methods("POST")
	managersSubRouter.HandleFunc("/sales", s.handleManagerGetSales).Methods("GET")
	managersSubRouter.HandleFunc("/sales", s.handleManagerMakeSales).Methods("POST")
	managersSubRouter.HandleFunc("/products", s.handleManagerGetProducts).Methods("GET")
	managersSubRouter.HandleFunc("/products", s.handleManagerChangeProducts).Methods("POST")
	managersSubRouter.HandleFunc("/products/{id:[0-9]+}", s.handleManagerRemoveProductByID).Methods("DELETE")
	managersSubRouter.HandleFunc("/customers", s.handleManagerGetCustomers).Methods("GET")
	managersSubRouter.HandleFunc("/customers", s.handleManagerChangeCustomer).Methods("POST")
	managersSubRouter.HandleFunc("/customers/{id:[0-9]+}", s.handleManagerRemoveCustomerByID).Methods("DELETE")

}

/*
+
+
*/
//это фукция для записывание ошибки в responseWriter или просто для ответа с ошиками
func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	//печатаем ошибку
	log.Print(err)
	//отвечаем ошибку с помошю библиотеке net/http
	http.Error(w, http.StatusText(httpSts), httpSts)
}

//это функция для ответа в формате JSON (он принимает интерфейс по этому мы можем в нем передат все что захочется)
func respondJSON(w http.ResponseWriter, iData interface{}) {

	//преобразуем данные в JSON
	data, err := json.Marshal(iData)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//поставить хедер "Content-Type: application/json" в ответе
	w.Header().Set("Content-Type", "application/json")
	//пишем ответ
	_, err = w.Write(data)
	//если получили ошибку
	if err != nil {
		//печатаем ошибку
		log.Print(err)
	}
}
