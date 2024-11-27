package collector

import (
	"bmkgearthquakecollector/model"
	"bytes"
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

	cleanedBody := bytes.ReplaceAll(body, []byte("\n"), []byte(" "))

	var gempaTerkiniData model.DataModel
	if err := json.Unmarshal(cleanedBody, &gempaTerkiniData); err != nil {
		return nil, errors.Join(errors.New(gempaTerkiniURL), err)
	}

	return &gempaTerkiniData, nil
}
