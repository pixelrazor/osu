package osu

import (
	"net/http"
	"time"
)

// Mod is used to represent which modifiers are applied to beatmaps
type Mod int

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
	num := int(mods)
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

// Possible values of Beatmap.Approved
const (
	Loved     = 4
	Qualified = 3
	Approved  = 2
	Ranked    = 1
	Pending   = 0
	WIP       = -1
	Graveyard = -2
)

// UsernameType represents the two different ways a user can be represented in a query
type UsernameType string

// All username types
const (
	// UserID is for referring to the ID number of a user
	UserID UsernameType = "int"
	// Username is for referring to the name of a user
	Username UsernameType = "string"
)

// Mode represents the game mode for a Beatmap
type Mode int

// All modes
const (
	Osu Mode = iota
	Taiko
	Ctb
	Mania
)

// Language represents the language of a Beatmap
type Language int

// All languages
const (
	AnyLanguage Language = iota
	OtherLangauge
	English
	Japanese
	Chinese
	Instrumental
	Korean
	French
	German
	Swedish
	Spanish
	Italian
)

// Genre represents the genre of a Beatmap
type Genre int

// All genres
const (
	AnyGenre Genre = iota
	UnspecifiedGenre
	VideoGame
	Anime
	Rock
	Pop
	OtherGenre
	Novelty
	_
	HipHop
	Electronic
)

// Beatmap contains all data relating to an individual beatmap
type Beatmap struct {
	Approved         int       `json:"approved,string"`
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
	GenreID          Genre     `json:"genre_id,string"`
	LanguageID       Language  `json:"language_id,string"`
	Title            string    `json:"title"`
	TotalLength      int       `json:"total_length,string"`
	Version          string    `json:"version"`
	FileMd5          string    `json:"file_md5"`
	Mode             Mode      `json:"mode,string"`
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

// BeatmapOption is used to add optional queries to Client.Beatmaps
type BeatmapOption func(string) string

// UserOption is used to add optional queries to Client.User
type UserOption func(string) string

// ScoresOption is used to add optional queries to Client.Scores
type ScoresOption func(string) string

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

const apiURL = "https://osu.ppy.sh/api/"
