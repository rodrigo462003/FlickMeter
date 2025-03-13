package movieAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rodrigo462003/FlickMeter/fuzzy"
	"github.com/rodrigo462003/FlickMeter/model"
)

type MovieGetter interface {
	GetMovie(id string) (movie *model.Movie, err error)
	GetVideos(id string) (videos []model.Video, err error)
	Tree() (tree fuzzy.Tree)
}

type movieGet struct {
	auth string
	tree fuzzy.Tree
}

func NewMovieGet(token string, filePath string) *movieGet {
	movies := GetMovieIndex(filePath)
	stringers := make([]fuzzy.Stringer, len(movies))
	for i := range stringers {
		stringers[i] = movies[i]
	}
	tree := fuzzy.NewTree(stringers)

	return &movieGet{fmt.Sprintf("Bearer %s", token), *tree}
}

func (m *movieGet) Tree() fuzzy.Tree {
	return m.tree
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

func GetMovieIndex(filePath string) (movies []model.MovieIndex) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&movies); err != nil {
		panic(err)
	}

	return movies
}
