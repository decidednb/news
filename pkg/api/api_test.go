package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"news/pkg/storage"
	"news/pkg/storage/memdb"
	"strconv"
	"testing"
)

// Количество новостей в фикстуре
const wantLen = 2

func TestAPI_lastNewsHandler(t *testing.T) {

	// Объект API для теста
	dbase := memdb.New()
	api := New(dbase)

	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/"+strconv.Itoa(wantLen), nil)

	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()

	// Вызываем маршрутизатор.
	api.router.ServeHTTP(rr, req)

	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем ответ.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	// Раскодируем JSON в массив заказов.
	var news []storage.Post
	err = json.Unmarshal(b, &news)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	// Проверяем, что количество элементов в массиве соответсвует
	// количеству новостей в фикстуре.
	if len(news) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(news), wantLen)
	}

}
