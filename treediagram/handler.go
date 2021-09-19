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

	return FormatWeatherResponse(weather)
}

func (h *handler) Forecast(request contract.Request) (*contract.Response, error) {
	weather, err := h.lookupWeather(request.Content)
	if err != nil {
		return nil, err
	}

	return FormatForecastResponse(weather)
}

func (h *handler) Plan(request contract.Request) (*contract.Response, error) {
	plan, err := h.lookupPlan(request)
	if err != nil {
		return FormatParseError(err)
	}

	return FormatPlanResponse(plan)
}

func (h *handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", h.makeLoggingHttpHandlerFunc("weather", h.Weather))
	mux.HandleFunc("/forecast", h.makeLoggingHttpHandlerFunc("forecast", h.Forecast))
	mux.HandleFunc("/plan", h.makeLoggingHttpHandlerFunc("plan", h.Plan))

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

func (h *handler) lookupPlan(request contract.Request) (*weatherpb.PlanReply, error) {
	req, err := ParsePlanRequest(request)
	if err != nil {
		return nil, err
	}

	weather, err := h.client.Plan(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (h *handler) makeLoggingHttpHandlerFunc(name string, f func(contract.Request) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeRequestHttpHandlerFunc(f)

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
