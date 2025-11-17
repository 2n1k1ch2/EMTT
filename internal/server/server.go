package handler

import (
	_ "EMTT/docs"
	"EMTT/internal/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
)

type Server struct {
	ser *http.Server
	mux *http.ServeMux

	log slog.Logger
}

func NewServer(addr string, log *slog.Logger) *Server {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	return &Server{
		ser: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		mux: mux,
		log: *log,
	}
}

func (s *Server) Start(h *handlers.Handler) {
	s.mux.HandleFunc("/subscriptions/list", h.ListSubscriptions)
	s.mux.HandleFunc("/subscriptions/create", h.CreateSubscription)
	s.mux.HandleFunc("/subscriptions/get/", h.GetSubscription)
	s.mux.HandleFunc("/subscriptions/update/", h.UpdateSubscription)
	s.mux.HandleFunc("/subscriptions/delete/", h.DeleteSubscription)
	s.mux.HandleFunc("/subscriptions/total", h.TotalCost)

	s.log.Info("server starting on " + s.ser.Addr)

	if err := s.ser.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Error("cannot start server", slog.Any("error", err))
	}
}
