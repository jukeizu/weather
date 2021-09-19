package treediagram

import (
	"fmt"
	"time"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
)

func FormatParseError(err error) (*contract.Response, error) {
	switch err.(type) {
	case ParseError:
		return contract.StringResponse(err.Error()), nil
	}

	return nil, err
}

func FormatWeatherResponse(weather *weatherpb.WeatherReply) (*contract.Response, error) {
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

func FormatForecastResponse(weather *weatherpb.WeatherReply) (*contract.Response, error) {
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

func FormatPlanResponse(plan *weatherpb.PlanReply) (*contract.Response, error) {
	embed := contract.Embed{
		Title: plan.Location,
		Color: 6139372,
		Footer: &contract.EmbedFooter{
			Text: finalDecisionFooter,
		},
	}

	days := map[string][]*weatherpb.Hour{}

	loc, _ := time.LoadLocation("America/Denver")

	for _, hour := range plan.Hours {
		key := time.Unix(hour.Data.Timestamp, 0).In(loc).Format("Mon 01/02")
		days[key] = append(days[key], hour)
	}

	for day, hours := range days {
		embedDay := contract.EmbedField{
			Name: day,
		}
		for _, hour := range hours {
			ts := time.Unix(hour.Data.Timestamp, 0).In(loc).Format("3:04 PM")
			embedDay.Value += fmt.Sprintf("\n%s - %s", ts, hour.Data.Wind)
		}

		embed.Fields = append(embed.Fields, &embedDay)
	}

	message := contract.Message{
		Embed: &embed,
	}

	return &contract.Response{Messages: []*contract.Message{&message}}, nil
}
