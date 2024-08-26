package hyperbaseclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Hyperbase struct {
	baseURL   string
	authToken string
}

type HyperbaseProject struct {
	h         *Hyperbase
	projectID uuid.UUID
}

type HyperbaseCollection struct {
	hp           *HyperbaseProject
	collectionID uuid.UUID
}

type hyperbaseRes[T any] struct {
	Data  T                 `json:"data"`
	Error hyperbaseErrorRes `json:"error"`
}

type hyperbaseErrorRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type hyperbaseAuthRes struct {
	Token string `json:"token"`
}

var (
	ErrDuplicate = errors.New("duplicate value violates unique constraint")
)

func New(baseURL string) *Hyperbase {
	return &Hyperbase{baseURL: baseURL}
}

func (h *Hyperbase) Authenticate(
	tokenID uuid.UUID,
	token string,
	collectionID uuid.UUID,
	authCredential map[string]any,
) error {
	url := h.baseURL + "/api/rest/auth/token-based"

	data := map[string]any{
		"token_id":      tokenID,
		"token":         token,
		"collection_id": collectionID,
		"data":          authCredential,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("not OK")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var authRes hyperbaseRes[hyperbaseAuthRes]
	if err := json.Unmarshal(body, &authRes); err != nil {
		return err
	}

	h.authToken = authRes.Data.Token
	return nil
}

func (h *Hyperbase) SetProject(projectID uuid.UUID) *HyperbaseProject {
	return &HyperbaseProject{h: h, projectID: projectID}
}

func (hp *HyperbaseProject) SetCollection(collectionID uuid.UUID) *HyperbaseCollection {
	return &HyperbaseCollection{hp: hp, collectionID: collectionID}
}

func (hc *HyperbaseCollection) InsertOne(data map[string]any) error {
	url := hc.hp.h.baseURL + "/api/rest/project/" + hc.hp.projectID.String() + "/collection/" + hc.collectionID.String() + "/record"

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+hc.hp.h.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resData hyperbaseRes[any]
	if err := json.Unmarshal(body, &resData); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		if strings.HasPrefix(resData.Error.Message, "Duplicate value") {
			return ErrDuplicate
		}
		return errors.New("not created")
	}

	return nil
}
