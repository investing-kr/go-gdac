package gdac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type (
	TradeService service
)

type TradeListOptions struct {
	Pair   string `url:"pair,omitempty"`
	Offset string `url:"offset,omitempty"`
	Limit  int    `url:"limit,omitempty"`
}

func (s *TradeService) ListTrades(ctx context.Context, listOpt *TradeListOptions) ([]*Trade, *http.Response, error) {
	qv, err := query.Values(listOpt)
	if err != nil {
		return nil, nil, err
	}

	qs := qv.Encode()
	u := fmt.Sprintf("v0.4/public/trades?%s", qs)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	err = s.client.withAuthentication(req)
	if err != nil {
		return nil, nil, err
	}

	lst := []*Trade{}
	resp, err := s.client.Do(ctx, req, &lst)
	if err != nil {
		return nil, resp, err
	}

	return lst, resp, nil
}
