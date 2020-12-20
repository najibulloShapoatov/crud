package app

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/najibulloShapoatov/crud/pkg/customers"
	"encoding/json"
	"net/http"
)


func (s *Server) handleCustomerRegistration(w http.ResponseWriter, r *http.Request) {
	
	//обявляем структура клиента для запраса
	var item *customers.Customer

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//Генерируем bcrypt хеш от реалного пароля
	hashed, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//и поставляем хеш в поле парол
	item.Password = string(hashed)

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






func (s *Server) handleCustomerGetToken(w http.ResponseWriter, r *http.Request) {
	//обявляем структуру для запроса
	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	//извелекаем данные из запраса
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}
	//взываем из сервиса  securitySvc метод AuthenticateCustomer
	token, err := s.customerSvc.Token(r.Context(), item.Login, item.Password)

	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//вызываем функцию для ответа в формате JSON
	respondJSON(w, map[string]interface{}{"status": "ok", "token": token})

}



func (s *Server) handleCustomerGetProducts(w http.ResponseWriter, r *http.Request) {

	items, err := s.customerSvc.Products(r.Context())
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	respondJSON(w, items)

}

