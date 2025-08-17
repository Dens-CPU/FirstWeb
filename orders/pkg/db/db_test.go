package db

import (
	"reflect"
	"testing"
	"time"
)

// Тест всех CRUD-операций с заказами.
func TestDB_Order(t *testing.T) {
	// Создаем БД.
	db := New()
	o := Order{
		IsOpen:       true,
		DeliveryTime: time.Now().Unix(),
	}

	// Тест создания записи в БД.
	o.ID = db.NewOrder(o)
	// Проверка.
	ord := db.Orders()
	if !reflect.DeepEqual(ord[0], o) {
		t.Errorf("не найден созданный заказ")
	}

	// Тест обновления записи в БД.
	o.IsOpen = false
	o.DeliveryAddress = "Адрес доставки"
	db.UpdateOrder(o)
	// Проверка.
	ord = db.Orders()
	if !reflect.DeepEqual(ord[0], o) {
		t.Errorf("не найден обновленный заказ")
	}

	// Тест удаления записи из БД.
	db.DeleteOrder(o.ID)
	// Проверка.
	ord = db.Orders()
	if len(ord) != 0 {
		t.Errorf("заказ не был удален")
	}
}
