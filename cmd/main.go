package main

import (
	"go.uber.org/dig"
	"time"
	"context"
	"github.com/najibulloShapoatov/crud/cmd/app"
	"github.com/najibulloShapoatov/crud/pkg/customers"
	"net/http"
	"github.com/jackc/pgx/v4/pgxpool"
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
	

	dependencies := []interface{}{
		app.NewServer,
		http.NewServeMux,
		func() (*pgxpool.Pool, error){
			connCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(connCtx, dbConnectionString)
		},
		customers.NewService,
		func(server *app.Server)*http.Server{
			return &http.Server{
				Addr:host+":"+port,
				Handler: server,
			}
		},
	}


	container := dig.New()
	for _, v := range dependencies {
		err = container.Provide(v)
		if err !=nil{
			return err
		}
	}

	err = container.Invoke(func(server *app.Server){
		server.Init()
	})
	if err != nil{
		return err
	}

	return container.Invoke(func(server *http.Server) error{
		return server.ListenAndServe()
	})
}