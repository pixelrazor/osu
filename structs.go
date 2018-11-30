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
	UserID   UsernameType = "int"
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
	FavouriteCount   int       `json:"favourite_count,string"`
	Playcount        int       `json:"playcount,string"`
	Passcount        int       `json:"passcount,string"`
	MaxCombo         int       `json:"max_combo,string"`
}

// Option is used to add optional queries to the API calls
type Option func(string) string

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
