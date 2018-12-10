package osu

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ulikunitz/xz/lzma"
)

// Mod is used to represent which modifiers are applied to beatmaps
type Mod int64

// All mods
const (
	None   Mod = 0
	NoFail Mod = 1 << iota
	Easy
	TouchDevice
	Hidden
	HardRock
	SuddenDeath
	DoubleTime
	Relax
	HalfTime
	Nightcore
	Flashlight
	Autoplay
	SpunOut
	Relax2
	Perfect
	Key4
	Key5
	Key6
	Key7
	Key8
	FadeIn
	Random
	Cinema
	Target
	Key9
	KeyCoop
	Key1
	Key3
	Key2
	ScoreV2
	LastMod
	KeyMod            = Key1 | Key2 | Key3 | Key4 | Key5 | Key6 | Key7 | Key8 | Key9 | KeyCoop
	FreeModAllowed    = NoFail | Easy | Hidden | HardRock | SuddenDeath | Flashlight | FadeIn | Relax | Relax2 | SpunOut | KeyMod
	ScoreIncreaseMods = Hidden | HardRock | DoubleTime | Flashlight | FadeIn
)

func (m Mod) String() string {
	switch m {
	case None:
		return "None"
	case NoFail:
		return "NoFail"
	case Easy:
		return "Easy"
	case TouchDevice:
		return "TouchDevice"
	case Hidden:
		return "Hidden"
	case HardRock:
		return "HardRock"
	case SuddenDeath:
		return "SuddenDeath"
	case DoubleTime:
		return "DoubleTime"
	case Relax:
		return "Relax"
	case HalfTime:
		return "HalfTime"
	case Nightcore:
		return "Nightcore"
	case Flashlight:
		return "Flashlight"
	case Autoplay:
		return "Autoplay"
	case SpunOut:
		return "SpunOut"
	case Relax2:
		return "Relax2"
	case Perfect:
		return "Perfect"
	case Key4:
		return "Key4"
	case Key5:
		return "Key5"
	case Key6:
		return "Key6"
	case Key7:
		return "Key7"
	case Key8:
		return "Key8"
	case FadeIn:
		return "FadeIn"
	case Random:
		return "Random"
	case Cinema:
		return "Cinema"
	case Target:
		return "Target"
	case Key9:
		return "Key9"
	case KeyCoop:
		return "KeyCoop"
	case Key1:
		return "Key1"
	case Key3:
		return "Key3"
	case Key2:
		return "Key2"
	case ScoreV2:
		return "ScoreV2"
	case LastMod:
		return "LastMod"
	default:
		return "Invalid Mod"
	}
}

// Mods represents a list of mods
type Mods int

func (mods Mods) String() string {
	l := mods.List()
	if len(l) == 0 {
		return "[]"
	}
	ret := ""
	for _, v := range l {
		ret += ", " + v.String()
	}
	return "[" + ret[2:] + "]"
}

// List returns a slice containing the individual mods contained by a field
func (mods Mods) List() []Mod {
	num := int64(mods)
	list := make([]Mod, 0)
	i := uint(0)
	for num != 0 {
		if num&1 == 1 {
			list = append(list, Mod(1<<i))
		}
		i++
		num >>= 1
	}
	return list
}

type status int

// Status holds possible values of Beatmap.Approved
var Status = struct {
	Loved, Qualified, Approved, Ranked, Pending, WIP, Graveyard status
}{4, 3, 2, 1, 0, -1, -2}

type usernameType string

// UsernameType holds the two different ways a user can be represented in a query
var UsernameType = struct {
	// ID is for referring to the ID number of a user
	ID usernameType
	// Name is for referring to the name of a user
	Name usernameType
}{"int", "string"}

type mode int

// Mode holds all possible game modes
var Mode = struct {
	Osu, Taiko, Ctb, Mania mode
}{0, 1, 2, 3}

type language int

// Language holds all languages of a Beatmap
var Language = struct {
	Any, Other, English, Japanese, Chinese, Instrumental, Korean, French, German, Swedish, Spanish, Italian language
}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

type genre int

// Genre holds the genres of Beatmaps
var Genre = struct {
	Any, Unspecified, VideoGame, Anime, Rock, Pop, OtherGenre, Novelty, HipHop, Electronic genre
}{0, 1, 2, 3, 4, 5, 6, 7, 9, 10}

// Beatmap contains all data relating to an individual beatmap
type Beatmap struct {
	Approved         status    `json:"approved,string"`
	ApprovedDate     time.Time `json:"approved_date,string"`
	LastUpdate       time.Time `json:"last_update,string"`
	Artist           string    `json:"artist"`
	BeatmapID        string    `json:"beatmap_id"`
	BeatmapsetID     string    `json:"beatmapset_id"`
	Bpm              float64   `json:"bpm,string"`
	Creator          string    `json:"creator"`
	CreatorID        string    `json:"creator_id"`
	Difficultyrating float64   `json:"difficultyrating,string"`
	DiffSize         float64   `json:"diff_size,string"`
	DiffOverall      float64   `json:"diff_overall,string"`
	DiffApproach     float64   `json:"diff_approach,string"`
	DiffDrain        float64   `json:"diff_drain,string"`
	HitLength        int       `json:"hit_length,string"`
	Source           string    `json:"source"`
	GenreID          genre     `json:"genre_id,string"`
	LanguageID       language  `json:"language_id,string"`
	Title            string    `json:"title"`
	TotalLength      int       `json:"total_length,string"`
	Version          string    `json:"version"`
	FileMd5          string    `json:"file_md5"`
	Mode             mode      `json:"mode,string"`
	Tags             string    `json:"tags"`
	FavouriteCount   int64     `json:"favourite_count,string"`
	Playcount        int64     `json:"playcount,string"`
	Passcount        int64     `json:"passcount,string"`
	MaxCombo         int64     `json:"max_combo,string"`
}

// User holds all information relating to a user
type User struct {
	UserID             string       `json:"user_id"`
	Username           string       `json:"username"`
	JoinDate           time.Time    `json:"join_date,string"`
	Count300           int64        `json:"count300,string"`
	Count100           int64        `json:"count100,string"`
	Count50            int64        `json:"count50,string"`
	Playcount          int64        `json:"playcount,string"`
	RankedScore        int64        `json:"ranked_score,string"`
	TotalScore         int64        `json:"total_score,string"`
	PpRank             int64        `json:"pp_rank,string"`
	Level              float64      `json:"level,string"`
	PpRaw              float64      `json:"pp_raw,string"`
	Accuracy           float64      `json:"accuracy,string"`
	CountRankSs        int64        `json:"count_rank_ss,string"`
	CountRankSSH       int64        `json:"count_rank_ssh,string"`
	CountRankS         int64        `json:"count_rank_s,string"`
	CountRankSh        int64        `json:"count_rank_sh,string"`
	CountRankA         int64        `json:"count_rank_a,string"`
	Country            string       `json:"country"`
	TotalSecondsPlayed int64        `json:"total_seconds_played,string"`
	PpCountryRank      int64        `json:"pp_country_rank,string"`
	Events             []*UserEvent `json:"events"`
}

// UserEvent holds information about recent events for a user
type UserEvent struct {
	DisplayHTML  string    `json:"display_html"`
	BeatmapID    string    `json:"beatmap_id"`
	BeatmapsetID string    `json:"beatmapset_id"`
	Date         time.Time `json:"date,string"`
	Epicfactor   int64     `json:"epicfactor,string"`
}

// Score holds iformation about a score for a specific beatmap
type Score struct {
	ScoreID         string    `json:"score_id"`
	Score           int64     `json:"score,string"`
	Username        string    `json:"username"`
	Count300        int64     `json:"count300,string"`
	Count100        int64     `json:"count100,string"`
	Count50         int64     `json:"count50,string"`
	Countmiss       int64     `json:"countmiss,string"`
	Maxcombo        int64     `json:"maxcombo,string"`
	Countkatu       int64     `json:"countkatu,string"`
	Countgeki       int64     `json:"countgeki,string"`
	Perfect         string    `json:"perfect"`
	EnabledMods     Mods      `json:"enabled_mods,string"`
	UserID          string    `json:"user_id"`
	Date            time.Time `json:"date,string"`
	Rank            string    `json:"rank"`
	Pp              float64   `json:"pp,string"`
	ReplayAvailable string    `json:"replay_available"`
}

// BestScore holds the information on the top scores for a user
type BestScore struct {
	BeatmapID   string    `json:"beatmap_id"`
	Score       int64     `json:"score,string"`
	Maxcombo    int64     `json:"maxcombo,string"`
	Count300    int64     `json:"count300,string"`
	Count100    int64     `json:"count100,string"`
	Count50     int64     `json:"count50,string"`
	Countmiss   int64     `json:"countmiss,string"`
	Countkatu   int64     `json:"countkatu,string"`
	Countgeki   int64     `json:"countgeki,string"`
	Perfect     string    `json:"perfect"`
	EnabledMods Mods      `json:"enabled_mods,string"`
	UserID      string    `json:"user_id"`
	Date        time.Time `json:"date,string"`
	Rank        string    `json:"rank"`
	Pp          float64   `json:"pp,string"`
}

// RecentScore holds the information on the top scores for a user
type RecentScore struct {
	BeatmapID   string    `json:"beatmap_id"`
	Score       int64     `json:"score,string"`
	Maxcombo    int64     `json:"maxcombo,string"`
	Count300    int64     `json:"count300,string"`
	Count100    int64     `json:"count100,string"`
	Count50     int64     `json:"count50,string"`
	Countmiss   int64     `json:"countmiss,string"`
	Countkatu   int64     `json:"countkatu,string"`
	Countgeki   int64     `json:"countgeki,string"`
	Perfect     string    `json:"perfect"`
	EnabledMods Mods      `json:"enabled_mods,string"`
	UserID      string    `json:"user_id"`
	Date        time.Time `json:"date,string"`
	Rank        string    `json:"rank"`
}

// Match contains the information for a multiplayer match
type Match struct {
	Match MatchInfo    `json:"match"`
	Games []*MatchGame `json:"games"`
}

// MatchInfo contains information about the multiplayer room
type MatchInfo struct {
	MatchID   string     `json:"match_id"`
	Name      string     `json:"name"`
	StartTime time.Time  `json:"start_time,string"`
	EndTime   *time.Time `json:"end_time,string"`
}

// ScoringType represents the scoring method in a multiplayer match
var ScoringType = struct {
	Score, Accuracy, Combo, ScoreV2 scoringType
}{0, 1, 2, 3}

type scoringType int

//MatchGame contains information about beatmaps that have been played in a multiplayer match
type MatchGame struct {
	GameID      string        `json:"game_id"`
	StartTime   time.Time     `json:"start_time,string"`
	EndTime     time.Time     `json:"end_time,string"`
	BeatmapID   string        `json:"beatmap_id"`
	PlayMode    mode          `json:"play_mode,string"`
	MatchType   int64         `json:"match_type,string"`
	ScoringType scoringType   `json:"scoring_type,string"`
	TeamType    string        `json:"team_type"`
	Mods        Mods          `json:"mods,string"`
	Scores      []*MatchScore `json:"scores,string"`
}

// MatchScore contains the information about the score of an individual user in a multiplayer match
type MatchScore struct {
	Slot      int64  `json:"slot,string"`
	Team      int64  `json:"team,string"`
	UserID    string `json:"user_id"`
	Score     string `json:"score"`
	Maxcombo  string `json:"maxcombo"`
	Rank      string `json:"rank"`
	Count50   int64  `json:"count50,string"`
	Count100  int64  `json:"count100,string"`
	Count300  int64  `json:"count300,string"`
	Countmiss int64  `json:"countmiss,string"`
	Countgeki int64  `json:"countgeki,string"`
	Countkatu int64  `json:"countkatu,string"`
	Perfect   string `json:"perfect"`
	Pass      string `json:"pass"`
}

// ReplayPoint holds a piece of replay data
type ReplayPoint struct {
	// Time since the last action
	TimeSinceLast time.Duration
	// x-coordinate of the cursor from 0 - 512
	X float64
	// y-coordinate of the cursor from 0 - 384
	Y float64
	// Bitwise combination of keys/mouse buttons pressed
	Keys int
}

// ReplayContent just wraps []*ReplayPoint so it can implement json's Unmarhsaler interface
type ReplayContent []*ReplayPoint

// Replay holds a list of points that represent the different states in a replay
type Replay struct {
	Content ReplayContent `json:"content,string"`
}

// UnmarshalJSON satisfies the Unmarshaler interface
func (rc *ReplayContent) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Replace(str, `"`, "", -1)
	str = strings.Replace(str, `\`, "", -1)
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	reader, err := lzma.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	outRaw, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	groups := strings.Split(string(outRaw), ",")
	out := make([]*ReplayPoint, 0)
	for i := range groups {
		var (
			w    int64
			x, y float64
			z    int
		)
		if strings.Count(groups[i], "|") != 3 {
			continue
		}
		fmt.Sscanf(groups[i], "%v|%v|%v|%v", &w, &x, &y, &z)
		if w < 0 {
			continue
		}
		out = append(out, &ReplayPoint{time.Duration(w) * time.Millisecond, x, y, z})
	}
	*rc = out
	return nil
}

// BeatmapOption is used to add optional queries to Client.Beatmaps
type BeatmapOption func(string) string

// UserOption is used to add optional queries to Client.User
type UserOption func(string) string

// ScoresOption is used to add optional queries to Client.Scores
type ScoresOption func(string) string

// UserBestOption is used to add optional queries to Client.UserBest
type UserBestOption func(string) string

// UserRecentOption is used to add optional queries to Client.UserRecent
type UserRecentOption func(string) string

// ReplayOption is used to add optional queries to Client.Replay
type ReplayOption func(string) string

// Client executes requests to the endpoints
type Client struct {
	key string
	c   http.Client
}

// NewClient creates a Client with the given key
func NewClient(key string) *Client {
	client := new(Client)
	client.key = key
	return client
}
func errorCheck(data []byte) error {
	regex := regexp.MustCompile(`"error"[[:space:]]*:[[:space:]]*"(.*)"`)
	if regex.Match(data) {
		return errors.New(string(regex.ExpandString(nil, "$1", string(data), regex.FindSubmatchIndex(data))))
	}
	return nil
}

const apiURL = "https://osu.ppy.sh/api/"
