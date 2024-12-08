package rest

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/Utro-tvar/medods-test/internal/email"
	mwlog "github.com/Utro-tvar/medods-test/internal/middleware/logger"
	"github.com/Utro-tvar/medods-test/internal/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Service interface {
	Generate(user models.User) models.TokensPair
	Refresh(tokens models.TokensPair, ip net.IP) models.TokensPair
	GetUser(access string) models.User
}

type App struct {
	service Service
	logger  *slog.Logger
	router  chi.Router
}

func New(logger *slog.Logger, service Service) *App {
	app := App{logger: logger, service: service}

	app.router = chi.NewRouter()

	app.router.Use(middleware.Logger)
	app.router.Use(mwlog.New(logger))
	app.router.Use(middleware.Recoverer)
	app.router.Use(middleware.URLFormat)

	app.router.Route("/generate/{guid}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			guid := chi.URLParam(r, "guid")
			ipstr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				logger.Error("Cannot get user ip", slog.Any("error", err))
				render.Data(w, r, []byte("Cannot parse your ip"))
				return
			}
			ip := net.ParseIP(ipstr)
			render.JSON(w, r, service.Generate(models.User{GUID: guid, IP: ip}))
		})
	})

	app.router.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		tokens := models.TokensPair{}
		render.DecodeJSON(r.Body, &tokens)

		ipstr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Error("Cannot get user ip", slog.Any("error", err))
			render.Data(w, r, []byte("Cannot parse your ip"))
			return
		}
		ip := net.ParseIP(ipstr)
		user := service.GetUser(tokens.Access)
		if !net.IP.Equal(ip, user.IP) {
			email.Send("mock@email.com", []byte("Your IP has been changed"))
		}

		new := service.Refresh(tokens, ip)
		render.JSON(w, r, new)
	})

	return &app
}

func (a *App) MustRun(addr string) {
	err := http.ListenAndServe(addr, a.router)
	if err != nil {
		panic(err)
	}
}
