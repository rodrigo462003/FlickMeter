package movieAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rodrigo462003/FlickMeter/model"
)

type MovieFetcher interface {
	Get(id uint) (movie *model.Movie, err error)
	Search(query string) (movies []model.Movie, err error)
}

type movieClient struct {
	http *http.Client
}

type movieTransport struct {
	transport   http.RoundTripper
	bearerToken string
}

func (t *movieTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.bearerToken))
	req.Header.Add("Accept", "application/json")

	return t.transport.RoundTrip(req)
}

func newTransport(token string) *movieTransport {
	return &movieTransport{
		transport:   http.DefaultTransport,
		bearerToken: token,
	}
}

func NewMovieGet(token string) *movieClient {
	return &movieClient{
		http: &http.Client{
			Transport: newTransport(token),
		},
	}
}

func (c *movieClient) Get(id uint) (*model.Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
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

	videos, err := c.Videos(movie.ID)
	if err != nil {
		return nil, err
	}

	movie.Videos = videos
	return movie, nil
}

func (c *movieClient) Videos(id int) ([]model.Video, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/videos", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
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

func (c *movieClient) Search(query string) ([]model.Movie, error) {
	query = url.PathEscape(query)
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?query=%s", query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := struct {
		Results []model.Movie `json:"results"`
	}{}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}
