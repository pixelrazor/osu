package osu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"
)

// WithMode specifies which mode to confine results to in searches
func WithMode(mode Mode) Option {
	return func(s string) string {
		return s + fmt.Sprintf("&m=%d", mode)
	}
}

// WithID confines the search to beatmaps with a specific BeatmapID
func WithID(ID string) Option {
	return func(s string) string {
		return s + fmt.Sprintf("&b=%v", ID)
	}
}

// WithSetID confines the search to beatmaps with a specific BeatmapsetID
func WithSetID(ID string) Option {
	return func(s string) string {
		return s + fmt.Sprintf("&s=%v", ID)
	}
}

// Since confines the search to beatmaps ranked or loved since date
func Since(date time.Time) Option {
	return func(s string) string {
		return s + fmt.Sprintf("&since=%v", date.Format("2006-01-02"))
	}
}

// Beatmaps fetches a list of beatmaps
func (client *Client) Beatmaps(opts ...Option) ([]*Beatmap, error) {
	query := apiURL + "get_beatmaps?k=" + client.key
	for _, opt := range opts {
		query = opt(query)
	}
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.Beatmaps: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.Beatmaps: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.Beatmaps: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	maps := make([]*Beatmap, 0)
	err = json.Unmarshal(body, &maps)
	return maps, err
}
