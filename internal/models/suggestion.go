package models

import "encoding/json"

type SuggestionItem struct {
	Label string `json:"label" redis:"label"`
	Value string `json:"value" redis:"value"`
}

func (bbi *SuggestionItem) MarshalBinary() ([]byte, error) {
	return json.Marshal(bbi)
}

func (bbi *SuggestionItem) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &bbi); err != nil {
		return err
	}

	return nil
}
