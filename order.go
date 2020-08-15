package gdac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

type (
	OrderService service
)

type OrderCreateOptions struct {
	Pair     string    `url:"pair"`
	Side     OrderSide `url:"side"`
	Price    string    `url:"price"`
	Quantity string    `url:"quantity"`
}

type OrderSide string

const (
	SideBuy  OrderSide = "B" // Buying, 매수
	SideSell OrderSide = "S" // Selling, 매도
)

type OrderListOptions struct {
	Pair   string `url:"pair,omitempty"`
	Offset string `url:"offset,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Status string `url:"status,omitempty"`
}

func (s *OrderService) CreateOrder(ctx context.Context, opt *OrderCreateOptions) (*Order, *http.Response, error) {
	qv, err := query.Values(opt)
	if err != nil {
		return nil, nil, err
	}

	qs := qv.Encode()
	u := fmt.Sprintf("v0.4/orders?%s", qs)
	req, err := s.client.NewRequest(http.MethodPost, u, nil)
	if err != nil {
		return nil, nil, err
	}

	err = s.client.withAuthentication(req)
	if err != nil {
		return nil, nil, err
	}

	o := Order{}
	resp, err := s.client.Do(ctx, req, &o)
	if err != nil {
		return nil, resp, err
	}

	return &o, resp, nil
}

func (s *OrderService) ListOrders(ctx context.Context, opt *OrderListOptions) ([]*Order, *http.Response, error) {
	qv, err := query.Values(opt)
	if err != nil {
		return nil, nil, err
	}

	qs := qv.Encode()
	u := fmt.Sprintf("v0.4/orders?%s", qs)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	err = s.client.withAuthentication(req)
	if err != nil {
		return nil, nil, err
	}

	o := []*Order{}
	resp, err := s.client.Do(ctx, req, &o)
	if err != nil {
		return nil, resp, err
	}

	return o, resp, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*Order, *http.Response, error) {
	u := fmt.Sprintf("v0.4/orders/%s", orderID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	err = s.client.withAuthentication(req)
	if err != nil {
		return nil, nil, err
	}

	o := &Order{}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, resp, err
	}

	return o, resp, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID string) (*Order, *http.Response, error) {
	u := fmt.Sprintf("v0.4/orders/%s", orderID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}

	err = s.client.withAuthentication(req)
	if err != nil {
		return nil, nil, err
	}

	o := &Order{}
	resp, err := s.client.Do(ctx, req, o)
	if err != nil {
		return nil, resp, err
	}

	return o, resp, nil
}
