package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MXLange/desafio-pos-weather-api/env"
	"github.com/MXLange/desafio-pos-weather-api/internal/domain/app"
	viacep "github.com/MXLange/desafio-pos-weather-api/internal/domain/via_cep"
	weatherapi "github.com/MXLange/desafio-pos-weather-api/internal/domain/weather_api"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e, err := env.New(ctx)
	if err != nil {
		panic(err)
	}

	vc, err := viacep.NewService(e.APICEPBaseURL)
	if err != nil {
		panic(err)
	}

	wa, err := weatherapi.NewService(e.WeatherBaseURL, e.WeatherAPIKey)
	if err != nil {
		panic(err)
	}

	as := app.NewService(vc, wa)
	h := app.NewHandler(as)

	mux := http.NewServeMux()
	mux.HandleFunc("/temperature", h.GetTemperature)

	server := &http.Server{
		Addr:    ":" + e.ApiPort,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	stop()

	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}
