package gdac

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type (
	AccountService service
)

func (s *AccountService) Balance(ctx context.Context, currency string) ([]*Balance, *http.Response, error) {
	u := fmt.Sprintf("v0.4/balance")

	if currency != "" {
		params := url.Values{}
		params.Add("currency", currency)
		qs := params.Encode()
		u = fmt.Sprintf("%s?%s", u, qs)
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	err = s.client.withAuthentication(req)
	if err != nil {
		return nil, nil, err
	}

	o := []*Balance{}
	resp, err := s.client.Do(ctx, req, &o)
	if err != nil {
		return nil, resp, err
	}

	return o, resp, nil
}
