package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/najibulloShapoatov/crud/pkg/customers"
)

//406.75
//4.0675p
//465p

//Server ...
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
}

//NewServer ... создает новый сервер
func NewServer(m *mux.Router, cSvc *customers.Service) *Server {
	return &Server{mux: m, customerSvc: cSvc}
}

// функция для запуска хендлеров через мукс
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init ... инициализация сервера
func (s *Server) Init() {
	//s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)

	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods("GET")
	s.mux.HandleFunc("/customers", s.handleSave).Methods("POST")
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods("GET")

	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods("GET")
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods("POST")
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods("DELETE")
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods("DELETE")
	

	/*
		http://127.0.0.1:9999/customers.save?id=0&name=Shahlo&phone=992928015290
		http://127.0.0.1:9999/customers.save?id=0&name=Najibullo&phone=992931441244
		http://127.0.0.1:9999/customers.getById?id=1
		http://127.0.0.1:9999/customers.getAll
		http://127.0.0.1:9999/customers.getAllActive
		http://127.0.0.1:9999/customers.blockById?id=1
		http://127.0.0.1:9999/customers.unblockById?id=1
		http://127.0.0.1:9999/customers.removeById?id=1

	*/
}

// хендлер метод для извлечения всех клиентов
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	//берем все клиенты
	items, err := s.customerSvc.All(r.Context())

	//если ест ошибка
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, items)
}

// хендлер метод для извлечения всех активных клиентов
func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {

	//берем все активные клиенты
	items, err := s.customerSvc.AllActive(r.Context())

	//если ест ошибка
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, items)
}

//хендлер который верет по айди
func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ByID(r.Context(), id)

	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

//хендлер для блокировки
func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//изменяем статус клиента на фалсе
	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)
	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

//хенндлер для разблокировки
func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)
	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	//idP := r.URL.Query().Get("id")
	idP := mux.Vars(r)["id"]

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//удаляем клиента из базу
	item, err := s.customerSvc.Delete(r.Context(), id)
	//если ошибка равно на notFound то вернем ошибку не найдено
	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	//передаем в функции respondJSON, ResponseWriter и данные (он отвечает клиенту)
	respondJSON(w, item)
}

//хендлер для сохранения и обновления
func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	/*
		//получаем данные из параметра запроса
		idP := r.FormValue("id")
		name := r.FormValue("name")
		phone := r.FormValue("phone")

		id, err := strconv.ParseInt(idP, 10, 64)
		//если получили ошибку то отвечаем с ошибкой
		if err != nil {
			//вызываем фукцию для ответа с ошибкой
			errorWriter(w, http.StatusBadRequest, err)
			return
		}
		//Здесь опционалная проверка то что если все данные приходит пустыми то вернем ошибку
		if name == "" && phone == "" {
			//вызываем фукцию для ответа с ошибкой
			errorWriter(w, http.StatusBadRequest, err)
			return
		}

		//обявляем новый клиент
		item := &customers.Customer{
			ID:    id,
			Name:  name,
			Phone: phone,
		}
	*/

	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//сохроняем или обновляем клиент
	customer, err := s.customerSvc.Save(r.Context(), item)

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	respondJSON(w, customer)
}

/*
+
+
+
+
+
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

/*
+
+
+
*/
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
