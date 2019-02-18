package treediagram

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jukeizu/contacts/api/protobuf-spec/contactspb"
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
	args, err := shellwords.Parse(request.Content)
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
			Text: "ツリーダイアグラム最終決定",
		},
	}

	if len(weather.Alerts) > 0 {
		buffer := bytes.Buffer{}

		for _, alert := range weather.Alerts {
			buffer.WriteString(fmt.Sprintf("[%s](%s)\n", alert.Message, alert.Uri))
		}

		alerts := contract.EmbedField{
			Name:  "Alerts",
			Value: buffer.String(),
		}

		embed.Fields = append(embed.Fields, &alerts)
	}

	message := contract.Message{
		Embed: &embed,
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}

func (h *handler) Forecast(request contract.Request) (*contract.Response, error) {
	return &contract.Response{Messages: []*contract.Message{}}, nil
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

func parseNameValue(command string, content string) (string, string) {
	input := strings.SplitAfterN(content, command, 2)[1]
	split := strings.SplitN(input, "'", 3)
	name, value := split[1], strings.TrimSpace(split[2])

	return name, value
}

func formatContact(contact *contactspb.Contact) string {
	if contact == nil {
		return ""
	}

	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("**%s**\n", contact.Name))
	buffer.WriteString(fmt.Sprintf(":house: %s", contact.Address))
	buffer.WriteString(fmt.Sprintf("\n\n:iphone: %s", contact.Phone))
	buffer.WriteString("\n\n\n")

	return buffer.String()
}
