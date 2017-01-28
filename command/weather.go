package command

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cosban/lueshi/internal"
)

type WeatherData struct {
	Weather []struct {
		Id                      int
		Main, Description, Icon string
	}
	Main struct {
		Temp, Pressure, Humidity, Temp_min, Temp_max, Sea_level, Grnd_level float64
	}
	Wind struct {
		Speed, Deg float64
	}
	Name string
}

func Weather(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	r, err := unmarshalWeather(args)
	response := fmt.Sprintf("\u200B<@%s>: There is obviously no weather at that location, like, ever.", m.Author.ID)
	if err == nil && r != nil {
		response = fmt.Sprintf("\u200B<@%s>: %s: %s at %s wind at %s %s",
			m.Author.ID,
			r.Name,
			r.Weather[0].Main,
			tempString(r.Main.Temp),
			speedString(r.Wind.Speed),
			directionString(r.Wind.Deg),
		)
	}
	s.ChannelMessageSend(m.ChannelID, response)
}

func Temperature(args []string, s *discordgo.Session, m *discordgo.MessageCreate) {
	r, err := unmarshalWeather(args)
	response := fmt.Sprintf("\u200B<@%s>: 0.0 Kelvin. Seriously.", m.Author.ID)
	if err == nil && r != nil {
		response = fmt.Sprintf("\u200B<@%s>: %s: %s H:%s L:%s ",
			m.Author.ID,
			r.Name,
			tempString(r.Main.Temp),
			tempString(r.Main.Temp_max),
			tempString(r.Main.Temp_min))
	}
	s.ChannelMessageSend(m.ChannelID, response)
}

func unmarshalWeather(args []string) (*WeatherData, error) {
	var request string
	zip := parseZip(args)
	if len(zip) > 1 {
		q := url.QueryEscape(zip)
		request = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&zip=%s", weatherkey, q)
	} else {
		q := url.QueryEscape(strings.Join(args, " "))
		request = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s", weatherkey, q)
	}

	r := &WeatherData{}
	internal.GetJSON(request, r)

	if r.Main.Temp == 0 {
		r = nil
	}
	return r, nil
}

func parseZip(args []string) string {
	for _, element := range args {
		if strings.HasPrefix(element, "zip:") {
			return element[len("zip:"):]
		}
	}
	return ""
}

func tempString(k float64) string {
	return fmt.Sprintf("%.1fC (%.1fF)", kelvinToC(k), kelvinToF(k))
}

func speedString(k float64) string {
	return fmt.Sprintf("%.1f m/s", k)
}

func directionString(i float64) string {
	if i < 34 {
		return "NNE"
	} else if i < 56 {
		return "NE"
	} else if i < 79 {
		return "ENE"
	} else if i < 101 {
		return "E"
	} else if i < 124 {
		return "ESE"
	} else if i < 146 {
		return "SE"
	} else if i < 169 {
		return "SSE"
	} else if i < 191 {
		return "S"
	} else if i < 214 {
		return "SSW"
	} else if i < 236 {
		return "SW"
	} else if i < 259 {
		return "WSW"
	} else if i < 281 {
		return "W"
	} else if i < 304 {
		return "WNW"
	} else if i < 326 {
		return "NW"
	} else if i < 349 {
		return "NNW"
	}
	return "N"
}

func kelvinToC(k float64) float64 {
	return (k - 272.15)
}

func kelvinToF(k float64) float64 {
	return (k-273.15)*1.8 + 32
}
