package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/najibulloShapoatov/crud/pkg/managers"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/najibulloShapoatov/crud/cmd/app"
	"github.com/najibulloShapoatov/crud/pkg/customers"
	"go.uber.org/dig"
)

func main() {
	//это хост
	host := "0.0.0.0"
	//это порт
	port := "9999"
	//это строка подключения к бд
	dbConnectionString := "postgres://app:pass@localhost:5432/db"
	//запускаем функцию execute c проверкой на err
	if err := execute(host, port, dbConnectionString); err != nil {
		//если получили ошибку то закрываем приложения
		log.Print(err)
		os.Exit(1)
	}
}

//функция запуска сервера
func execute(host, port, dbConnectionString string) (err error) {

	//здес обявляем слайс с зависимостями тоест добавляем все сервисы и конструкторы
	dependencies := []interface{}{
		app.NewServer, //это сервер
		mux.NewRouter, //это роутер
		func() (*pgxpool.Pool, error) { //это фукция конструктор который принимает *pgxpool.Pool, error
			connCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(connCtx, dbConnectionString)
		},
		customers.NewService, //это сервис клиентов
		managers.NewService,  //это сервис менеджеров
		func(server *app.Server) *http.Server { //это фукция конструктор который принимает *app.Server и вернет *http.Server
			return &http.Server{
				Addr:    host + ":" + port,
				Handler: server,
			}
		},
	}

	//обявляем новый контейнер
	container := dig.New()
	//в цикле регистрируем все зависимостив контейнер
	for _, v := range dependencies {
		err = container.Provide(v)
		if err != nil {
			return err
		}
	}

	/*вызываем метод Invoke позволяет вызвать на контейнере функøия, при этом подставит нам в
	параметры тот объект, который нужно "собрать" (именно в ÿтот момент все
	зависимости будут собраны, либо мы полуùим ощибку)*/
	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	//если получили ошибку то вернем его
	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
