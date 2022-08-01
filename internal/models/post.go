package models

import (
	"encoding/json"
	"strings"

	"github.com/fiure-cms/fiure-web/internal/fcore"
)

type Post struct {
	Collection string `json:"collection" redis:"collection"`
	Bucket     string `json:"bucket" redis:"bucket"`

	UID     string `json:"uid" redis:"uid"`
	Status  string `json:"status" redis:"status"`
	Updated int64  `json:"updated" redis:"updated"`
	Score   int    `json:"score" redis:"score"`

	ID      string `json:"id" redis:"id"`
	License string `json:"license" redis:"license"`
	Name    string `json:"name" redis:"name"`
	Slug    string `json:"slug" redis:"slug"`
	Content string `json:"content" redis:"content"`

	Url     string `json:"url" redis:"url"`
	Website string `json:"website" redis:"website"`

	Cats       []PostTerm      `json:"cats" redis:"cats"`
	Tags       []PostTerm      `json:"tags" redis:"tags"`
	Images     []PostImageItem `json:"images" redis:"images"`
	Thumbnails []PostImageItem `json:"thumbs" redis:"thumbs"`
}

func (bbi *Post) MarshalBinary() ([]byte, error) {
	return json.Marshal(bbi)
}

func (bbi *Post) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &bbi); err != nil {
		return err
	}

	return nil
}

func (bbi *Post) Convert2SearchItem() (kv KV) {

	kv.Key = bbi.UID

	fullText := []string{}
	fullText = append(fullText, bbi.Name)

	// Remove multiline from content: Sonic doesnt like multiline
	singleLineContent := fcore.MakeSingleLineString(bbi.Content)
	singleLineContent = fcore.SanitizeString(singleLineContent)
	fullText = append(fullText, singleLineContent)

	fullText = append(fullText, bbi.Url)
	fullText = append(fullText, bbi.Website)

	var catList []string
	if len(bbi.Cats) > 0 {
		for _, term := range bbi.Cats {
			catList = append(catList, term.Name)
		}

		fullText = append(fullText, strings.Join(catList, ", "))
	}

	var tagList []string
	if len(bbi.Tags) > 0 {
		for _, term := range bbi.Tags {
			tagList = append(tagList, term.Name)
		}

		fullText = append(fullText, strings.Join(tagList, ", "))
	}

	// Make single line string
	kv.Value = strings.Join(fullText, " | ")

	return kv
}

func (bbi *Post) Convert2SuggestItem() (kv KV) {

	kv.Key = "sug:" + bbi.UID

	fullText := []string{}
	fullText = append(fullText, bbi.Name)

	var catList []string
	if len(bbi.Cats) > 0 {
		for _, term := range bbi.Cats {
			catList = append(catList, term.Name)
		}

		fullText = append(fullText, strings.Join(catList, ", "))
	}

	var tagList []string
	if len(bbi.Tags) > 0 {
		for _, term := range bbi.Tags {
			tagList = append(tagList, term.Name)
		}

		fullText = append(fullText, strings.Join(tagList, ", "))
	}

	kv.Value = strings.Join(fullText, " | ")

	return kv
}

type PostTerm struct {
	Name string `json:"name" redis:"name"`
	Slug string `json:"slug" redis:"slug"`
}

type PostImageItem struct {
	Name string `json:"name" redis:"name"`
	Path string `json:"path" redis:"path"`
}
