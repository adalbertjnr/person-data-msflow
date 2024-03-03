package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adalbertjnr/ws-person/types"
)

type HTTPClientEndpoint struct {
	endpoint string
}

func NewHTTPClientEndpoint(endpoint string) *HTTPClientEndpoint {
	return &HTTPClientEndpoint{
		endpoint: endpoint,
	}
}

func (e *HTTPClientEndpoint) Aggregate(ctx context.Context, data *types.AggregatePerson) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error serializing data. http client error %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, e.endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error creating new post request to %s %w", e.endpoint, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending the http post request %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("not expected http post request response. got %d", resp.StatusCode)
	}
	return nil
}

func (e *HTTPClientEndpoint) GetPersonById(ctx context.Context, id int) (*types.Person, error) {
	b, err := json.Marshal(id)
	if err != nil {
		return nil, fmt.Errorf("error serializing user id to create the request for gateway %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, e.endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("not expected status code from the gateway server %d", resp.StatusCode)
	}
	p := types.Person{}
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
