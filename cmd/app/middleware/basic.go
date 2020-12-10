package middleware

import (
	"log"
	"errors"
	"encoding/base64"
	"strings"
	"net/http"
)
//Basic ...
func Basic(checkAuth func(string, string)bool) func(handler http.Handler)http.Handler{

	// вернем функцию который принимает хендлер и возврашает тоже хендлер
	return func(handler http.Handler)http.Handler{

		//вернем агтсцию хендлер 
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){

			//здес извелекаем логин и парол
			login, pass, err := getLoginPass(r)
			//если при извлечении получили ошибке отвечаем с ошибкой 401 и печатаем ошибку в лог
			if err!=nil{
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusUnauthorized),http.StatusUnauthorized)
				return
			}
			//если данные нет в базу отвечаем с ошибкой 401 
			if !checkAuth(login, pass) {
				http.Error(w, http.StatusText(http.StatusUnauthorized),http.StatusUnauthorized)
				return
			}
			//если все хорошо вызываем уже нужный хендлер
			handler.ServeHTTP(w,r)
		})
	}
}

// это функция который извелекаеть данные из запроса и вернет логин и парол и ошибку если есть ошибка
func getLoginPass(r *http.Request) (string, string, error){
		
		//здес мы берем из хедера заначению Authorization потом разделим его на две части
        auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		//здесь проверим точно там через пробел два значения и первы из них должен быть слова Basic
        if len(auth) != 2 || auth[0] != "Basic" {
			//если нет то вернем ошибку что такой метод авторизации унас нету
            return "", "", errors.New("invalid auth method") 
        }

		//берем 2 значению тоест после слова Basic то что было и со стандартными библиотекой декодируем тесть извелекаем
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		//здесь разделим по ":" и скадываем их в массив
		pair := strings.SplitN(string(payload), ":", 2)
		//если длина массива не равно два тоест ползовател отправил нам не достоверные данные 
		if len(pair) != 2 {
			//и вернем ошибку что не не достоверные данные 
			return "", "", errors.New("invalid auth data") 
		}
		// если все хорошо то вернем лоиг и парол и нил тоест нет ошибки
		return pair[0], pair[1], nil
}