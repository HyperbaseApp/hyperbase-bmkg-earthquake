package collector

import (
	"bmkgearthquakecollector/model"
	"bytes"
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

	cleanedBody := bytes.ReplaceAll(body, []byte("\n"), []byte(" "))

	var gempaTerkiniData model.DataModel
	if err := json.Unmarshal(cleanedBody, &gempaTerkiniData); err != nil {
		return nil, errors.Join(errors.New(gempaDirasakanURL), err)
	}

	return &gempaTerkiniData, nil
}
