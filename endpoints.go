package osu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"time"
)

// BeatmapsIncludeConverted includes converted beatmaps (only valid when mode is not Osu)
func BeatmapsIncludeConverted() BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprint("&a=1")
	}
}

// BeatmapsLimit specifies how many beatmaps to return (default and maximum value is 500)
func BeatmapsLimit(limit int) BeatmapOption {
	if limit < 1 {
		limit = 1
	} else if limit > 500 {
		limit = 500
	}
	return func(s string) string {
		return s + fmt.Sprintf("&limit=%d", limit)
	}
}

// BeatmapsWithMode specifies which mode to confine results to
func BeatmapsWithMode(mode mode) BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprintf("&m=%d", mode)
	}
}

// BeatmapsWithHash confines the search to beatmaps with a specific BeatmapID
func BeatmapsWithHash(ID string) BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprintf("&h=%v", url.QueryEscape(ID))
	}
}

// BeatmapsWithID confines the search to beatmaps with a specific BeatmapID
func BeatmapsWithID(ID string) BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprintf("&b=%v", url.QueryEscape(ID))
	}
}

// BeatmapsWithSetID confines the search to beatmaps with a specific BeatmapsetID
func BeatmapsWithSetID(ID string) BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprintf("&s=%v", url.QueryEscape(ID))
	}
}

// BeatmapsByCreator confines the search to beatmaps created by a specific user
func BeatmapsByCreator(ID string, IDType usernameType) BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprintf("&u=%v&type=%v", url.QueryEscape(ID), url.QueryEscape(string(IDType)))
	}
}

// BeatmapsSince confines the search to beatmaps ranked or loved since date
func BeatmapsSince(date time.Time) BeatmapOption {
	return func(s string) string {
		return s + fmt.Sprintf("&since=%v", date.Format("2006-01-02"))
	}
}

// Beatmaps fetches a list of beatmaps
func (client *Client) Beatmaps(opts ...BeatmapOption) ([]*Beatmap, error) {
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
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.Beatmaps: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	maps := make([]*Beatmap, 0)
	err = json.Unmarshal(body, &maps)
	return maps, err
}

// UserMode specifies which mode to show info for in the User struct (default is Osu)
func UserMode(mode mode) UserOption {
	return func(s string) string {
		return s + fmt.Sprintf("&m=%d", mode)
	}
}

// UserEventsSince specifies the max number of days between now and last event date. Range of 1-31 (default is 1)
func UserEventsSince(days int) UserOption {
	if days < 1 {
		days = 1
	} else if days > 31 {
		days = 31
	}
	return func(s string) string {
		return s + fmt.Sprintf("&event_days=%d", days)
	}
}

// User fetches information for a specific user
func (client *Client) User(ID string, IDType usernameType, opts ...UserOption) (*User, error) {
	query := apiURL + "get_user?k=" + client.key
	query = BeatmapsByCreator(ID, IDType)(query)
	for _, opt := range opts {
		query = opt(query)
	}
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.User: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.User: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.User: " + err.Error())
	}
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.User: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	user := make([]*User, 0)
	err = json.Unmarshal(body, &user)
	if len(user) == 1 {
		return user[0], err
	}
	return nil, err
}

// ScoresWithMode confines results to those with the specified mode
func ScoresWithMode(m mode) ScoresOption {
	return func(s string) string {
		return s + fmt.Sprintf("&m=%d", m)
	}
}

// ScoresByUser confines results to just scores from the specific user
func ScoresByUser(ID string, IDType usernameType) ScoresOption {
	return func(s string) string {
		return s + fmt.Sprintf("&u=%v&type=%v", url.QueryEscape(ID), url.QueryEscape(string(IDType)))
	}
}

// ScoresWithMods confines results to those that have the given mods
func ScoresWithMods(mods ...Mod) ScoresOption {
	m := 0
	for _, v := range mods {
		m |= int(v)
	}
	return func(s string) string {
		return s + fmt.Sprintf("&mods=%d", m)
	}
}

// ScoresLimit specifies how many beatmaps to return (default 50, max 100)
func ScoresLimit(limit int) ScoresOption {
	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}
	return func(s string) string {
		return s + fmt.Sprintf("&limit=%d", limit)
	}
}

// Scores fetches a list of Scores for a specific beatmap
func (client *Client) Scores(ID string, opts ...ScoresOption) ([]*Score, error) {
	query := apiURL + "get_scores?k=" + client.key
	query = BeatmapsWithID(ID)(query)
	for _, opt := range opts {
		query = opt(query)
	}
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.Scores: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.Scores: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.Scores: " + err.Error())
	}
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.Scores: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	scores := make([]*Score, 0)
	err = json.Unmarshal(body, &scores)
	return scores, err
}

// UserBestLimit specifies how many beatmaps to return (default 10, max 100)
func UserBestLimit(limit int) UserBestOption {
	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}
	return func(s string) string {
		return s + fmt.Sprintf("&limit=%d", limit)
	}
}

// UserBestWithMode confines results to those with the specified mode
func UserBestWithMode(m mode) UserBestOption {
	return func(s string) string {
		return s + fmt.Sprintf("&m=%d", m)
	}
}

// UserBest returns a list of the top scores for the specified user
func (client *Client) UserBest(ID string, IDType usernameType, opts ...UserBestOption) ([]*BestScore, error) {
	query := apiURL + "get_user_best?k=" + client.key
	query = BeatmapsByCreator(ID, IDType)(query)
	for _, opt := range opts {
		query = opt(query)
	}
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.UserBest: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.UserBest: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.UserBest: " + err.Error())
	}
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.UserBest: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	scores := make([]*BestScore, 0)
	err = json.Unmarshal(body, &scores)
	return scores, err
}

// UserRecentLimit specifies how many beatmaps to return (default 10, max 50)
func UserRecentLimit(limit int) UserRecentOption {
	if limit < 1 {
		limit = 1
	} else if limit > 50 {
		limit = 100
	}
	return func(s string) string {
		return s + fmt.Sprintf("&limit=%d", limit)
	}
}

// UserRecentWithMode confines results to those with the specified mode
func UserRecentWithMode(m mode) UserRecentOption {
	return func(s string) string {
		return s + fmt.Sprintf("&m=%d", m)
	}
}

// UserRecent returns a list of the top scores for the specified user
func (client *Client) UserRecent(ID string, IDType usernameType, opts ...UserRecentOption) ([]*RecentScore, error) {
	query := apiURL + "get_user_recent?k=" + client.key
	query = BeatmapsByCreator(ID, IDType)(query)
	for _, opt := range opts {
		query = opt(query)
	}
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.UserRecent: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.UserRecent: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.UserRecent: " + err.Error())
	}
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.UserRecent: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	scores := make([]*RecentScore, 0)
	err = json.Unmarshal(body, &scores)
	return scores, err
}

// ReplayWithMods confines the result to a replay with the given mods
func ReplayWithMods(mods ...Mod) ReplayOption {
	m := 0
	for _, v := range mods {
		m |= int(v)
	}
	return func(s string) string {
		return s + fmt.Sprintf("&mods=%d", m)
	}
}

// Replay returns the data for the given beatmap, played by the specified user in the specified mode
func (client *Client) Replay(mode mode, beatmapID string, userID string, opts ...ReplayOption) (*Replay, error) {
	query := apiURL + "get_replay?k=" + client.key
	query = BeatmapsWithMode(mode)(query)
	query = BeatmapsWithID(beatmapID)(query)
	query += "&u=" + userID
	for _, opt := range opts {
		query = opt(query)
	}
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.Replay: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.Replay: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.Replay: " + err.Error())
	}
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.Replay: " + err.Error())
	}
	var regx = regexp.MustCompile(`(date"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	var replay Replay
	err = json.Unmarshal(body, &replay)
	return &replay, err
}

// Match fetches a multiplayer match with the given ID
func (client *Client) Match(matchID string) (*Match, error) {
	query := apiURL + "get_match?k=" + client.key + "&mp=" + matchID
	resp, err := client.c.Get(query)
	if err != nil {
		return nil, errors.New("osu.Client.Match: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("osu.Client.Match: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("osu.Client.Match: " + err.Error())
	}
	err = errorCheck(body)
	if err != nil {
		return nil, errors.New("osu.Client.Match: " + err.Error())
	}
	var regx = regexp.MustCompile(`(time"[[:space:]]*:[[:space:]]*"[0-9]{4}-[0-9]{2}-[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2})"`)
	body = regx.ReplaceAll(body, []byte(`${1}T${2}-00:00"`))
	match := make([]*Match, 0)
	err = json.Unmarshal(body, &match)
	if len(match) == 1 {
		return match[0], err
	}
	return nil, err
}
