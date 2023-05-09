package numbersapi

import (
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type Api struct {
	baseUrl string
}

func NewApi() *Api {
	baseUrl := viper.GetString("BASE_URL")

	return &Api{
		baseUrl: baseUrl,
	}
}

func (a *Api) GetFact(param, category string) (string, error) {
	req, err := http.NewRequest("GET", a.baseUrl+param+"/"+category, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
