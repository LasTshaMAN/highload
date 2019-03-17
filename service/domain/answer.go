package domain

import (
	"encoding/json"
	"io"
)

type Answer struct {
	Value int `json:"value"`
}

func parse(respBody io.Reader) (Answer, error) {
	var result Answer
	if err := json.NewDecoder(respBody).Decode(&result); err != nil {
		return Answer{}, err
	}
	return result, nil
}
