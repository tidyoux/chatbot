package weather

import (
	"fmt"
	"net/url"

	"github.com/buger/jsonparser"
	"github.com/tidyoux/chatbot/utils"
)

const (
	apiKey = "3e43bb984fc141939b890ab7962d37a4"

	apiURL = "https://free-api.heweather.com/s6/weather?key=" + apiKey + "&location="

	success = "ok"
)

var lifestyleTypes = map[string]string{
	"comf":  "舒适度指数",
	"cw":    "洗车指数",
	"drsg":  "穿衣指数",
	"flu":   "感冒指数",
	"sport": "运动指数",
	"trav":  "旅游指数",
	"uv":    "紫外线指数",
	"air":   "空气污染扩散条件指数",
	"ac":    "空调开启指数",
	"ag":    "过敏指数",
	"gl":    "太阳镜指数",
	"mu":    "化妆指数",
	"airc":  "晾晒指数",
	"ptfc":  "交通指数",
	"fisin": "钓鱼指数",
	"spi":   "防晒指数",
}

type WeatherData []byte

func getWeatherData(location string) (WeatherData, error) {
	loc := url.QueryEscape(location)
	data, err := utils.GetBody(apiURL + loc)
	if err != nil {
		return nil, err
	}

	return newWeatherData(data)
}

func newWeatherData(data []byte) (WeatherData, error) {
	data, _, _, err := jsonparser.Get(data, "HeWeather6", "[0]")
	if err != nil {
		return nil, err
	}

	return WeatherData(data), nil
}

func (d WeatherData) status() string {
	status, _ := jsonparser.GetString(d, "status")
	return status
}

func (d WeatherData) basic() string {
	data, _, _, _ := jsonparser.Get(d, "basic")
	return formatFromJSON(data, "%s %s %s %s (%s, %s) %s",
		"cnty", "admin_area", "parent_city", "location", "lon", "lat", "tz")
}

func (d WeatherData) now() string {
	data, _, _, _ := jsonparser.Get(d, "now")
	return formatFromJSON(data, "%s, %s%s级, 气温%s℃, 体感温度%s℃, 相对湿度%s%%, 能见度%skm, 降水量%smm, 气压%shPa",
		"cond_txt", "wind_dir", "wind_sc", "tmp", "fl", "hum", "vis", "pcpn", "pres")
}

func (d WeatherData) dailyForecast() string {
	data, _, _, _ := jsonparser.Get(d, "daily_forecast")
	var msg string
	jsonparser.ArrayEach(data, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		msg += "\n"
		if equalFromJSON(data, "cond_txt_d", "cond_txt_n") {
			msg += formatFromJSON(value, "%s, %s, %s%s级, 气温%s℃ - %s℃, 相对湿度%s%%, 能见度%skm，降水概率%s%%, 降水量%smm, 气压%shPa, 日出日落%s-%s, 月升月落%s-%s",
				"date", "cond_txt_d", "wind_dir", "wind_sc", "tmp_max", "tmp_min", "hum", "vis", "pop", "pcpn", "pres", "sr", "ss", "mr", "ms")
		} else {
			msg += formatFromJSON(value, "%s, %s转%s, %s%s级, 气温%s℃ - %s℃, 相对湿度%s%%, 能见度%skm，降水概率%s%%, 降水量%smm, 气压%shPa, 日出日落%s-%s, 月升月落%s-%s",
				"date", "cond_txt_d", "cond_txt_n", "wind_dir", "wind_sc", "tmp_max", "tmp_min", "hum", "vis", "pop", "pcpn", "pres", "sr", "ss", "mr", "ms")
		}
	})
	return msg
}

func (d WeatherData) lifestyle() string {
	data, _, _, _ := jsonparser.Get(d, "lifestyle")
	var msg string
	jsonparser.ArrayEach(data, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		t, _ := jsonparser.GetString(value, "type")
		ts := lifestyleTypes[t]
		msg += "\n" + fmt.Sprintf("%s, %s", ts, formatFromJSON(value, "%s。%s", "brf", "txt"))
	})
	return msg
}

func (d WeatherData) update() string {
	update, _ := jsonparser.GetString(d, "update", "loc")
	return update
}

func (d WeatherData) format() string {
	status := d.status()
	if status != success {
		return status
	}

	return fmt.Sprintf("更新时间: %s\n", d.update()) +
		fmt.Sprintf("基本信息: %s\n", d.basic()) +
		fmt.Sprintf("实时天气: %s\n", d.now()) +
		fmt.Sprintf("\n天气预报: %s\n", d.dailyForecast()) +
		fmt.Sprintf("\n生活指数: %s\n", d.lifestyle())
}

func formatFromJSON(data []byte, format string, keys ...string) string {
	var values []interface{}
	for _, key := range keys {
		v, _ := jsonparser.GetString(data, key)
		values = append(values, v)
	}
	return fmt.Sprintf(format, values...)
}

func equalFromJSON(data []byte, key1, key2 string) bool {
	v1, _ := jsonparser.GetString(data, key1)
	v2, _ := jsonparser.GetString(data, key2)
	return v1 == v2
}
