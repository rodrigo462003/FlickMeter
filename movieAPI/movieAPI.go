package movieAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rodrigo462003/FlickMeter/model"
)

type MovieGetter interface {
	GetMovie(id string) (movie *model.Movie, err error)
}

type movieGet struct {
	Auth string
}

func NewMovieGet(token string) *movieGet {
	return &movieGet{fmt.Sprintf("Bearer %s", token)}
}

func (m *movieGet) GetMovie(id string) (*model.Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", m.Auth)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	movie := new(model.Movie)
	if err := json.Unmarshal(body, movie); err != nil {
		return nil, err
	}

	return movie, nil
}
