package managers

import (
	"errors"
	"strings"

	"github.com/uretgec/go-sonic/sonic"
)

type SearchManager struct {
	Rs   *sonic.Client
	lang string
}

func NewSearchManager(rs *sonic.Client) *SearchManager {
	return &SearchManager{
		Rs:   rs,
		lang: sonic.LangEng,
	}
}

func (sm *SearchManager) UpdateLang(lang string) {
	sm.lang = lang
}

func (sm *SearchManager) HealthChecked() error {
	return sm.Rs.Ping(sm.Rs.Context()).Err()
}

func (sm *SearchManager) Close() {
	sm.Rs.Close()
}

// Only Use Search Mode
func (sm *SearchManager) Search(collection, bucket, terms string, limit, offset int) ([]string, error) {
	if sm.Rs.Options().ChannelMode != sonic.ChannelSearch {
		return nil, errors.New("not use this function this mode")
	}

	results, err := sm.Rs.Query(sm.Rs.Context(), collection, bucket, terms, limit, offset, sm.lang).Slice()
	if err != nil {
		return nil, err
	}

	prefix := "post:"

	if len(results) > 0 {
		for in, uid := range results {
			results[in] = strings.Replace(uid, prefix, "", 1)
		}
	}

	return results, nil
}

func (sm *SearchManager) Suggest(collection, bucket, query string, limit int) ([]string, error) {
	if sm.Rs.Options().ChannelMode != sonic.ChannelSearch {
		return nil, errors.New("not use this function this mode")
	}

	results, err := sm.Rs.Suggest(sm.Rs.Context(), collection, bucket, query, limit).Slice()
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Only Use IngestMode
func (sm *SearchManager) Set(collection, bucket, object, text string) error {
	if sm.Rs.Options().ChannelMode != sonic.ChannelIngest {
		return errors.New("not use this function this mode")
	}

	var chunks []string

	if sm.Rs.IsPushContentReady(text) {
		chunks = sm.Rs.SplitPushContent(text)
	} else {
		chunks = append(chunks, text)
	}

	for _, text := range chunks {
		if err := sm.Rs.Push(sm.Rs.Context(), collection, bucket, object, text, sm.lang).Err(); err != nil {
			return err
		}
	}

	chunks = nil
	return nil
}

func (sm *SearchManager) Del(collection, bucket, object string) error {
	if sm.Rs.Options().ChannelMode != sonic.ChannelIngest {
		return errors.New("not use this function this mode")
	}

	return sm.Rs.FlushObject(sm.Rs.Context(), collection, bucket, object).Err()
}

func (sm *SearchManager) SetSuggest(collection, bucket, object, text string) error {
	if sm.Rs.Options().ChannelMode != sonic.ChannelIngest {
		return errors.New("not use this function this mode")
	}

	var chunks []string

	if sm.Rs.IsPushContentReady(text) {
		chunks = sm.Rs.SplitPushContent(text)
	} else {
		chunks = append(chunks, text)
	}

	for _, text := range chunks {
		if err := sm.Rs.Push(sm.Rs.Context(), collection, bucket, object, text, sm.lang).Err(); err != nil {
			return err
		}
	}

	chunks = nil
	return nil
}

func (sm *SearchManager) DelSuggest(collection, bucket, object string) error {
	if sm.Rs.Options().ChannelMode != sonic.ChannelIngest {
		return errors.New("not use this function this mode")
	}

	return sm.Rs.FlushObject(sm.Rs.Context(), collection, bucket, object).Err()
}
