package api

import (
	"FirstWeb/orders/pkg/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура API-приложения
type API struct {
	r  *mux.Router //Маршрутизатор запросов
	db *db.DB      //База данных
}

// Конструктор для API
func New(db *db.DB) *API {
	api := API{}
	api.r = mux.NewRouter()
	api.db = db
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов, т.к переменная r не экспортируема
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/orders", api.ordersHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/orders", api.newOrderHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/orders/{id}", api.updateOrderHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/orders/{id}", api.deleteOrderHandler).Methods(http.MethodDelete)
}

// ordersHandler возвращает все заказы.
func (api *API) ordersHandler(w http.ResponseWriter, r *http.Request) {

	//Получения данных из БД
	orders := api.db.Orders()

	//Отпрвака данных клиентуу в формате JSON
	json.NewEncoder(w).Encode(orders)
}

// newOrderehandler создает новый заказ
func (api *API) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	var o db.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	id := api.db.NewOrder(o)
	w.Write([]byte(strconv.Itoa(id)))
}

// updateOrderHandler обновляет данные заказа по ID.
func (api *API) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра {id} из пути запроса.
	// Например, /orders/45.
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Декодирование в переменную тела запроса,
	// которое должно содержать JSON-представление
	// обновляемого объекта.
	var o db.Order
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o.ID = id
	// Обновление данных в БД.
	api.db.UpdateOrder(o)
	// Отправка клиенту статуса успешного выполнения запроса
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	api.db.DeleteOrder(id)
	w.WriteHeader(http.StatusOK)
}
