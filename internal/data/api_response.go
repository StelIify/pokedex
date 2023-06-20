package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://pokeapi.co/api/v2"

type LocationAreasResp struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}
type LocationAreaDetailsResp struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}
type PokemonResp struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Species        struct {
		Name string `json:"name"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type Client struct {
	cache  *Cache
	client http.Client
}

func NewClient(timeout time.Duration, cache *Cache) Client {
	return Client{
		client: http.Client{Timeout: timeout},
		cache:  cache,
	}
}

func (c *Client) GetLocationAreas(next *string) (*LocationAreasResp, error) {
	endpoint := "/location-area"
	fullPath := baseUrl + endpoint

	if next != nil {
		fullPath = *next
	}

	body, ok := c.cache.Get(fullPath)
	if ok {
		var response LocationAreasResp
		err := json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("error during response unmarshal: %v", err)
		}
		return &response, nil
	}
	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during get request to the next batch of data: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return nil, fmt.Errorf("bad status code in the response: %v", response.StatusCode)
	}

	var apiResponse LocationAreasResp
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error during response unmarshal: %v", err)
	}
	c.cache.Put(fullPath, body)

	return &apiResponse, nil
}

func (c *Client) GetLocationAreaDetails(location string) (*LocationAreaDetailsResp, error) {
	endpoint := fmt.Sprintf("/location-area/%s", location)
	fullPath := baseUrl + endpoint

	body, ok := c.cache.Get(fullPath)
	if ok {
		var response LocationAreaDetailsResp
		err := json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("error during response unmarshal: %v", err)
		}
		return &response, nil
	}
	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during get request to location-area: %s, %v", location, err)
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return nil, fmt.Errorf("bad status code in the response: %v", response.StatusCode)
	}

	var apiResponse LocationAreaDetailsResp
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error during response unmarshal: %v", err)
	}
	c.cache.Put(fullPath, body)
	return &apiResponse, nil
}

func (c *Client) GetPokemonDetails(name string) (*PokemonResp, error) {
	endpoint := fmt.Sprintf("/pokemon/%s", name)
	fullPath := baseUrl + endpoint

	body, ok := c.cache.Get(fullPath)
	if ok {
		var response PokemonResp
		err := json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("error during response unmarshal: %v", err)
		}
		return &response, nil
	}
	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during get request to get pokemon info: %s, %v", name, err)
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return nil, fmt.Errorf("bad status code in the response: %v", response.StatusCode)
	}

	var apiResponse PokemonResp
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error during response unmarshal: %v", err)
	}
	c.cache.Put(fullPath, body)
	return &apiResponse, nil
}
