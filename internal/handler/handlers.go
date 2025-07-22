package handler

import "net/http"

type Handler struct {
	mux *http.ServeMux // Указатель на стандартный роутер, который связывает url с обработчиками http-запросов (роутер для обработки запросов)
}

// NewHandler - функция конструктор для структуры Handler, создает экземпляр, новый роутер

func NewHandler() *Handler {
	return &Handler{
		mux: http.NewServeMux()}
}
