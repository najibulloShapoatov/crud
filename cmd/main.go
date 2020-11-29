package main

import (
	"github.com/najibulloShapoatov/crud/cmd/app"
	"github.com/najibulloShapoatov/crud/pkg/customers"
	"net/http"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

func main() {
	//это хост
	host :="0.0.0.0"
	//это порт
	port := "9999"
	//это строка подключения к бд
	dbConnectionString :="postgres://app:pass@localhost:5432/db"
	//запускаем функцию execute c проверкой на err
	if err := execute(host, port, dbConnectionString); err != nil{
		//если получили ошибку то закрываем приложения
		log.Print(err)
		os.Exit(1)
	}
}

//функция запуска сервера
func execute(host, port, dbConnectionString string) (err error){
	//поключаемся к бд
	db, err := sql.Open("pgx", dbConnectionString)
	//если получили ошибку то вернем его
	if err !=nil{
		return err
	}
	//в конце закрываем подключения к бд
	defer db.Close()

	//обьявляем новый мукс
	mux := http.NewServeMux()
	//обьявляем новый сервис с бд
	customerService := customers.NewService(db)
	//обьявляем новый сервер с мукс и сервисами
	server := app.NewServer(mux, customerService)
	//Иницализируем наш сервер регистрируем роуты
	server.Init()

	//создаем новый Http server
	httpServer := &http.Server{
		Addr:host+":"+port,
		Handler: server,
	}

	//запускаем сервер
	return httpServer.ListenAndServe()
}