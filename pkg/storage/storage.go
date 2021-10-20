package storage

// Модель данных Публикация
type Post struct {
	ID      int    `bson:"id" json:"ID"`            // идентификатор новости
	Title   string `bson:"title" json:"Title"`      // заголовок публикации
	Content string `bson:"content" json:"Content"`  // содержание публикации
	PubTime int64  `bson:"pub_time" json:"PubTime"` // время публикации
	Link    string `bson:"link" json:"Link"`        // ссылка на источник
}

// Interface задаёт контракт на работу с базой данных.
type DBInterface interface {
	LastNews(n int) ([]Post, error) // получение n-количества последних новостей
	StoreNews([]Post) error         // создание массива публикации
	Close()                         // закрывает соединение с БД
}
