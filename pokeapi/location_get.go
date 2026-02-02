package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) GetLocation(locationName string) (Location, error) {
	url := BaseURL + LocationEndpoint + locationName
	
	if val, ok := c.cache.Get(url); ok {
		locResp := Location{}
		err := json.Unmarshal(val, &locResp)
		if err != nil {
			return Location{}, err
		}
		return locResp, nil
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}
	
	locResp := Location{}
	err = json.Unmarshal(data, &locResp)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, data)
	
	return locResp, nil
}