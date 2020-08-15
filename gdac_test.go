package gdac_test

import (
	"context"
	"os"
	"testing"

	"github.com/investing-kr/go-gdac"
)

var c *gdac.Client

func TestMain(m *testing.M) {
	var err error
	opts := gdac.ClientOptionsFromEnv()
	c, err = gdac.NewClient(nil, opts)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestBalance(t *testing.T) {
	ctx := context.Background()
	balances, _, err := c.Accounts.Balance(ctx, "USDT")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range balances {
		t.Logf("%+v", *v)
	}
}

func TestOrder(t *testing.T) {
	ctx := context.Background()
	order, _, err := c.Orders.CreateOrder(ctx, &gdac.OrderCreateOptions{
		Pair:     "KLAY/KRW",
		Side:     gdac.SideBuy,
		Price:    "100",
		Quantity: "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", *order)

	_, _, err = c.Orders.GetOrder(ctx, order.OrderID)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.Orders.CancelOrder(ctx, order.OrderID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Trades(t *testing.T) {
	ctx := context.Background()
	trades, _, err := c.Trades.ListTrades(ctx, &gdac.TradeListOptions{
		Pair: "KLAY/KRW",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range trades {
		t.Logf("%+v", trade)
	}
}
