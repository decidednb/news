package postgres

import (
	"fmt"
	"math/rand"
	"news/pkg/storage"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		t.Fatal(fmt.Errorf("Переменная окружения не получена"))
	}
	// conn - строка подключения к базе данных
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	s, err := New(pgConn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
}

func TestStorage_StoreNews(t *testing.T) {
	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		t.Fatal(fmt.Errorf("Переменная окружения не получена"))
	}
	// conn - строка подключения к базе данных
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	s, err := New(pgConn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	rand.Seed(time.Now().UnixNano())
	randString := strconv.Itoa(rand.Intn(1_000_000_000))
	news := []storage.Post{
		{
			Title:   "Test post title " + randString,
			Link:    "https://testtest.ru/news/testpost/" + randString,
			PubTime: time.Now().Unix(),
			Content: "Test post content " + randString,
		},
	}

	if err := s.StoreNews(news); err != nil {
		t.Fatalf("Storage.StoreNews() error = %v", err)
	}

}

func TestStorage_LastNews(t *testing.T) {
	type args struct {
		n int
	}

	pgPass := os.Getenv("pgPass")
	if pgPass == "" {
		t.Fatal(fmt.Errorf("Переменная окружения не получена"))
	}
	// conn - строка подключения к базе данных
	pgConn := "postgres://postgres:" + pgPass + "@localhost:5432/news"

	s, err := New(pgConn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Кейс 1 (n = 0) - ожидаем количество новостей по умолчанию",
			args: args{
				n: 0,
			},
			want:    10,
			wantErr: false,
		}, {
			name: "Кейс 2 - (n = 20)",
			args: args{
				n: 20,
			},
			want:    20,
			wantErr: false,
		}, {
			name: "Кейс 3 - (n = 5)",
			args: args{
				n: 5,
			},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.LastNews(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.LastNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("Storage.LastNews() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
