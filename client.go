package gdac

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	defaultBaseURL = "https://partner.gdac.com/"
)

type ClientOptions struct {
	APIKey    string
	SecretKey string
	ServerURL string
}

func ClientOptionsFromEnv() *ClientOptions {
	opt := &ClientOptions{
		SecretKey: os.Getenv("GDAC_OPEN_API_SECRET_KEY"),
		ServerURL: os.Getenv("GDAC_OPEN_API_SERVER_URL"),
		APIKey:    os.Getenv("GDAC_OPEN_API_API_KEY"),
	}
	return opt
}

type service struct {
	client *Client
}

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	common     service

	apiKey    string
	secretKey string

	Orders   *OrderService
	Trades   *TradeService
	Accounts *AccountService
}

func NewClient(httpClient *http.Client, opt *ClientOptions) (*Client, error) {
	serverURL := defaultBaseURL
	if opt.ServerURL != "" {
		serverURL = opt.ServerURL
	}

	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		apiKey:     opt.APIKey,
		secretKey:  opt.SecretKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	c.common.client = c
	c.Orders = (*OrderService)(&c.common)
	c.Trades = (*TradeService)(&c.common)
	c.Accounts = (*AccountService)(&c.common)
	return c, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		ctx = context.TODO()
	}

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

var ErrNotImplemented = fmt.Errorf("upbit: NotImplemented")

type ErrResponse struct {
	Code string `json:"code"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("gdac: %s", e.Code)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	errResp := &ErrResponse{}
	err = json.Unmarshal(body, errResp)

	if err == nil {
		return errResp
	}

	return errors.New(string(body))
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) withAuthentication(req *http.Request) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"api_key": "api:ak:" + c.apiKey,
		"nonce":   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	})

	tokenString, err := token.SignedString([]byte(c.secretKey))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	return nil
}
