package treediagram

import (
	"bytes"
	"flag"
	"net/url"
	"strconv"

	"github.com/jukeizu/contract"
	"github.com/jukeizu/weather/api/protobuf-spec/weatherpb"
	shellwords "github.com/mattn/go-shellwords"
)

func ParsePlanRequest(request contract.Request) (*weatherpb.PlanRequest, error) {
	values := url.Values(request.QueryParams)
	args, err := shellwords.Parse(request.Content)
	if err != nil {
		return nil, err
	}

	outputBuffer := bytes.NewBuffer([]byte{})

	parser := flag.NewFlagSet("", flag.ContinueOnError)
	parser.SetOutput(outputBuffer)

	paramDaylight, _ := strconv.ParseBool(values.Get("daylight"))
	paramUnits := values.Get("units")
	if paramUnits == "" {
		paramUnits = "imperial"
	}

	defaultWind := float64(0)
	paramWind := values.Get("windmax")
	if paramUnits != "" {
		w, err := strconv.ParseFloat(paramWind, 64)
		if err == nil {
			defaultWind = w
		}
	}

	defaultWindGust := float64(0)
	paramWindGust := values.Get("windgustmax")
	if paramUnits != "" {
		w, err := strconv.ParseFloat(paramWindGust, 64)
		if err == nil {
			defaultWindGust = w
		}
	}

	location := parser.String("l", values.Get("location"), "The location.")
	daylight := parser.Bool("d", paramDaylight, "Only include daylight hours.")
	wind := parser.Float64("ws", defaultWind, "The max wind speed.")
	windGust := parser.Float64("wgs", defaultWindGust, "The max wind gust speed.")
	units := parser.String("u", paramUnits, "The units to use.")

	err = parser.Parse(args[1:])
	if err != nil {
		return nil, ParseError{Message: outputBuffer.String()}
	}

	req := &weatherpb.PlanRequest{
		Location: *location,
		Daylight: *daylight,
		Units:    *units,
		Wind: &weatherpb.DataRange{
			Max: *wind,
		},
		WindGust: &weatherpb.DataRange{
			Max: *windGust,
		},
	}

	return req, nil
}
