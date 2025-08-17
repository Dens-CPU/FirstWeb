package api

import (
	"FirstWeb/orders/pkg/db"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_ordersHandler(t *testing.T) {
	// Создание чистого объекта API для теста.
	dbase := db.New()
	dbase.NewOrder(db.Order{})
	api := New(dbase)

	//Создание HTTp-запроса.
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)

	//Создание объекта для записи ответа обработчика.
	rr := httptest.NewRecorder()

	//Вызов маршрутизатора. Маршруутизатор для пути и метода запроса вызовет обработчик.
	//Обработчик запишет ответ в созданный объект.

	api.r.ServeHTTP(rr, req)

	//Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хоели %d", rr.Code, http.StatusOK)
	}
	//Читаем тело ответа.
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	//Раскодирование JSON в массив заказов.
	var data []db.Order
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	//Проверка, что в массиве ровно один элемент.
	const wantLen = 1
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
	// Также можно проверить совпадение заказов в результате
	// с добавленными в БД для теста.
}

func TestAPI_newOrderHandler(t *testing.T) {
	// Создание чисто объекта API.
	dbase := db.New()
	api := New(dbase)
	p := []db.Product{
		{
			Name:  "Phone",
			Price: 1,
		},
	}
	o := db.Order{
		IsOpen:       true,
		DeliveryTime: time.Now().Unix(),
		Products:     p,
	}
	b, err := json.Marshal(o)
	if err != nil {
		t.Errorf("Ошибка кодирования данных")
	}
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()

	//Вызов маршрутизатора.Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("получили %d,ожидалось %d", rr.Code, http.StatusOK)
	}
}

func TestAPI_updateOrderHandler(t *testing.T) {
	dbase := db.New()
	dbase.NewOrder(db.Order{})
	api := New(dbase)

	o := db.Order{}
	b, _ := json.Marshal(o)
	req := httptest.NewRequest(http.MethodPatch, "/orders/45", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("получили %d,ожидалось %d", rr.Code, http.StatusOK)
	}
}

func TestAPI_deleteOrderHandler(t *testing.T) {
	dbase := db.New()
	dbase.NewOrder(db.Order{})
	api := New(dbase)

	req := httptest.NewRequest(http.MethodDelete, "/orders/45", nil)
	rr := httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("получили %d,ожидалось %d", rr.Code, http.StatusOK)
	}
}
