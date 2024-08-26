package scraper

import (
	"bmkgscraper/model"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func GempaTerkini() (*model.DataModel, error) {
	gempaTerkiniURL := "https://data.bmkg.go.id/DataMKG/TEWS/gempaterkini.json"

	resp, err := http.Get(gempaTerkiniURL)
	if err != nil {
		return nil, errors.Join(errors.New(gempaTerkiniURL), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Join(errors.New(gempaTerkiniURL), errors.New("response status code is not ok"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New(gempaTerkiniURL), err)
	}

	var gempaTerkiniData model.DataModel
	if err := json.Unmarshal(body, &gempaTerkiniData); err != nil {
		return nil, errors.Join(errors.New(gempaTerkiniURL), err)
	}

	return &gempaTerkiniData, nil
}
