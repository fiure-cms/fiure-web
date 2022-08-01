package models

import (
	"encoding/json"
)

type Page struct {
	Name    string `json:"name" redis:"name"`
	Slug    string `json:"slug" redis:"slug"`
	Updated int64  `json:"updated" redis:"updated"`
	Content string `json:"content" redis:"content"`
	Status  string `json:"status" redis:"status"`
}

func (bbi *Page) MarshalBinary() ([]byte, error) {
	return json.Marshal(bbi)
}

func (bbi *Page) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &bbi); err != nil {
		return err
	}

	return nil
}
