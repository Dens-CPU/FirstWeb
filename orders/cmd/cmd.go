package main

import (
	"FirstWeb/orders/pkg/api"
	"FirstWeb/orders/pkg/db"
	"log"
	"net/http"
	"time"
)

func main() {
	//Инициализаци БД в памяти.
	dbase := db.New()
	//Создание объекта API, используующего БД в памяти.
	api := api.New(dbase)
	p := []db.Product{
		{
			Name:  "Iphone",
			Price: 100.0,
		},
		{
			Name:  "Samsung",
			Price: 101.0,
		},
	}
	o := db.Order{
		IsOpen:       true,
		DeliveryTime: time.Now().Unix(),
		Products:     p,
	}
	dbase.NewOrder(o)
	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 80.

	err := http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
