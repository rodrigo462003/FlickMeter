package movieAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rodrigo462003/FlickMeter/model"
)

type MovieGetter interface {
	GetMovie(id string) (movie *model.Movie, err error)
	GetVideos(id string) (videos []model.Video, err error)
	MoviesIndex() (movieIndex []model.MovieIndex)
}

type movieGet struct {
	auth   string
	bkTree []model.MovieIndex
}

func NewMovieGet(token string, filePath string) *movieGet {
    _ = (filePath)
	return &movieGet{fmt.Sprintf("Bearer %s", token), []model.MovieIndex{}}
}

func (m *movieGet) MoviesIndex() []model.MovieIndex {
    return m.bkTree
}

func (m *movieGet) GetMovie(id string) (*model.Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", m.auth)

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

func (m *movieGet) GetVideos(id string) ([]model.Video, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s/videos", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", m.auth)

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

	response := struct {
		Results []model.Video `json:"results"`
	}{}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

func movieTree(filePath string) (movieIdx []model.MovieIndex) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var results []model.MovieIndex
	if err := json.Unmarshal(bytes, &results); err != nil {
		panic(err)
	}

	return results
}
