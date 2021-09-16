package treediagram

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/rs/zerolog"
)

type Handler interface {
	Weather(contract.Request) (*contract.Response, error)
	Forecast(contract.Request) (*contract.Response, error)
	Start() error
	Stop() error
}

type handler struct {
	logger     zerolog.Logger
	client     weatherpb.WeatherClient
	httpServer *http.Server
}

func NewHandler(logger zerolog.Logger, client weatherpb.WeatherClient, addr string) Handler {
	logger = logger.With().Str("component", "intent.endpoint.weather").Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return &handler{logger, client, &httpServer}
}

func (h *handler) Weather(request contract.Request) (*contract.Response, error) {
	weather, err := h.lookupWeather(request.Content)
	if err != nil {
		return nil, err
	}

	embed := contract.Embed{
		Title: weather.Location,
		Color: 6139372,
		Fields: []*contract.EmbedField{
			&contract.EmbedField{
				Name:   getEmojiForIcon(weather.Currently.Data.Icon) + weather.Currently.Data.Temperature + " " + weather.Currently.Description,
				Value:  weather.Currently.Summary + "\n",
				Inline: true,
			},
			&contract.EmbedField{
				Name:   "Probability",
				Value:  generateProbabilitySummary(weather.Currently.Data),
				Inline: true,
			},
			&contract.EmbedField{
				Name:   "Data",
				Value:  generateDataSummary(weather.Currently.Data),
				Inline: true,
			},
		},
		Footer: &contract.EmbedFooter{
			Text: finalDecisionFooter,
		},
	}

	if len(weather.Alerts) > 0 {
		alerts := contract.EmbedField{
			Name:  "Alerts",
			Value: generateAlertsSummary(weather.Alerts),
		}

		embed.Fields = append(embed.Fields, &alerts)
	}

	message := contract.Message{
		Embed: &embed,
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) Forecast(request contract.Request) (*contract.Response, error) {
	weather, err := h.lookupWeather(request.Content)
	if err != nil {
		return nil, err
	}

	embed := contract.Embed{
		Title: weather.Location,
		Color: 6139372,
		Footer: &contract.EmbedFooter{
			Text: finalDecisionFooter,
		},
	}

	for _, day := range weather.Forecast.Days {
		embedDay := contract.EmbedField{
			Name:  generateDayTitle(day),
			Value: day.Summary,
		}

		embed.Fields = append(embed.Fields, &embedDay)
	}

	message := contract.Message{
		Embed: &embed,
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", h.makeLoggingHttpHandlerFunc("weather", h.Weather))
	mux.HandleFunc("/forecast", h.makeLoggingHttpHandlerFunc("forecast", h.Forecast))

	h.httpServer.Handler = mux

	return h.httpServer.ListenAndServe()
}

func (h *handler) Stop() error {
	h.logger.Info().Msg("stopping")

	return h.httpServer.Shutdown(context.Background())
}

func (h *handler) lookupWeather(content string) (*weatherpb.WeatherReply, error) {
	args, err := shellwords.Parse(content)
	if err != nil {
		return nil, err
	}

	weatherRequest := weatherpb.WeatherRequest{
		Location: strings.Join(args[1:], " "),
	}

	weather, err := h.client.Weather(context.Background(), &weatherRequest)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (h *handler) makeLoggingHttpHandlerFunc(name string, f func(contract.Request) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeHttpHandlerFunc(f)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			h.logger.Info().
				Str("intent", name).
				Str("took", time.Since(begin).String()).
				Msg("called")
		}(time.Now())

		contractHandlerFunc.ServeHTTP(w, r)
	}
}
