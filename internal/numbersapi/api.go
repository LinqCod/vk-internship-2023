package numbersapi

import (
	"io"
	"net/http"
	"os"
)

type Api struct {
	baseUrl string
}

func NewApi() *Api {
	baseUrl := os.Getenv("BASE_URL")

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
