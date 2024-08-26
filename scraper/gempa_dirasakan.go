package scraper

import (
	"bmkgscraper/model"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func GempaDirasakan() (*model.DataModel, error) {
	gempaDirasakanURL := "https://data.bmkg.go.id/DataMKG/TEWS/gempadirasakan.json"

	resp, err := http.Get(gempaDirasakanURL)
	if err != nil {
		return nil, errors.Join(errors.New(gempaDirasakanURL), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Join(errors.New(gempaDirasakanURL), errors.New("response status code is not ok"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New(gempaDirasakanURL), err)
	}

	var gempaTerkiniData model.DataModel
	if err := json.Unmarshal(body, &gempaTerkiniData); err != nil {
		return nil, errors.Join(errors.New(gempaDirasakanURL), err)
	}

	return &gempaTerkiniData, nil
}
