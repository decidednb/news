package rss

import (
	"testing"
)

func TestParseFeed(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Кейс 1 - rss лента stoneforest.ru",
			args: args{
				url: "https://stoneforest.ru/feed/",
			},
			wantErr: false,
		}, {
			name: "Кейс 2 - rss лента darkside.ru",
			args: args{
				url: "https://www.darkside.ru/rss/",
			},
			wantErr: false,
		}, {
			name: "Кейс 3 - rss лента habr.com/ru/rss/best/daily",
			args: args{
				url: "https://habr.com/ru/rss/best/daily/?fl=ru",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFeed(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFeed() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(got) == 0 {
				t.Errorf("Данные не распарсились или лента пуста")
			}

			t.Logf("Получено %v новостей\n%v", len(got), got)
		})
	}
}
