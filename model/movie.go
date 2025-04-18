package model

import (
	"slices"
	"time"
)

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	ISO6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

type Video struct {
	ISO639_1    string    `json:"iso_639_1"`
	ISO3166_1   string    `json:"iso_3166_1"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Site        string    `json:"site"`
	Size        int       `json:"size"`
	Type        string    `json:"type"`
	Official    bool      `json:"official"`
	PublishedAt time.Time `json:"published_at"`
	ID          string    `json:"id"`
}

type TopMovies struct {
	HotDay  []Movie
	HotWeek []Movie
	AllTime []Movie
}

type Movie struct {
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection any                 `json:"belongs_to_collection"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	ID                  uint                 `json:"id"`
	IMDBID              string              `json:"imdb_id"`
	OriginCountry       []string            `json:"origin_country"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []map[string]string `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
	Videos              []Video             `json:"videos"`
	Reviews             []Review
	Related             []Movie
}

func (m *Movie) Trailer() *Video {
	//Seems like api usually returns trailers at the end.
	for _, video := range slices.Backward(m.Videos) {
		if video.Type == "Trailer" {
			return &video
		}
	}

	return &Video{}
}
