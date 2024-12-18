package collector

import (
	"bmkgearthquakecollector/model"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func AutoGempa() (*model.AutoGempaModel, error) {
	autoGempaURL := "https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json"

	resp, err := http.Get(autoGempaURL)
	if err != nil {
		return nil, errors.Join(errors.New(autoGempaURL), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Join(errors.New(autoGempaURL), errors.New("response status code is not ok"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Join(errors.New(autoGempaURL), err)
	}

	cleanedBody := bytes.ReplaceAll(body, []byte("\n"), []byte(" "))

	var autoGempaData model.AutoGempaModel
	if err := json.Unmarshal(cleanedBody, &autoGempaData); err != nil {
		return nil, errors.Join(errors.New(autoGempaURL), err)
	}

	return &autoGempaData, nil
}
