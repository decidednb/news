-- Схема базы данных сайта - новостного агрегатора
DROP TABLE IF EXISTS news;

-- Таблица Новости
CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE
);