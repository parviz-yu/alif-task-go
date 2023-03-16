package server

type Config struct {
	// Адрес запуска вебсервера
	BindAddr string

	// Конфиг базы данных — путь
	DatabaseURL string

	CacheAddr string
}
