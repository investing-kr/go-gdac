package gdac

type Trade struct {
	Pair        string `json:"pair"`
	TradeID     string `json:"trade_id"`
	OrderID     string `json:"order_id"`
	Side        string `json:"side"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
	TradeDtime  string `json:"trade_dtime"`
	Fee         string `json:"fee"`
	FeeCurrency string `json:"fee_currency"`
	Type        string `json:"type"`
}

// O:Open, P:PartiallyFilled, F:Filled, C:Canceled
const (
	OrderStatusOpen            = "O"
	OrderStatusPartiallyFilled = "P"
	OrderStatusFilled          = "F"
	OrderStatusCanceled        = "C"
)

type Order struct {
	Pair         string `json:"pair"`
	OrderID      string `json:"order_id"`
	Status       string `json:"status"`
	Price        string `json:"price"`
	Quantity     string `json:"quantity"`
	OpenQuantity string `json:"open_quantity"`
	CreatedDtime string `json:"created_dtime"`
}

type Balance struct {
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
}
