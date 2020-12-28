package summary

import (
	"context"

	"github.com/jonfriesen/finance-go"
	"github.com/jonfriesen/finance-go/form"
)

// Client is used to invoke chart APIs.
type Client struct {
	B finance.Backend
}

func getC() Client {
	return Client{finance.GetBackend(finance.YFinBackend)}
}

// Get returns a quote summary.
func Get(ctx context.Context, symbol string) (*finance.QuoteSummary, error) {
	return getC().Get(ctx, symbol)
}

// Get returns a quote summary.
func (c Client) Get(ctx context.Context, symbol string) (*finance.QuoteSummary, error) {
	if symbol == "" {
		return nil, finance.CreateArgumentError()
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	// Build request.
	body := &form.Values{}
	// form.AppendTo(body, params)
	// Set request meta data.
	body.Set("region", "US")
	body.Set("corsDomain", "com.finance.yahoo")
	body.Set("formatted", "false")
	body.Set("lang", "en-US")
	body.Set("region", "US")

	// Note: If these modules are expanded the QuoteSummary struct will also need to be
	// expanded to match.
	body.Set("modules", "summaryProfile,financialData")

	resp := response{}
	err := c.B.Call("v10/finance/quoteSummary/"+symbol, body, &ctx, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Inner.Error != nil {
		return nil, resp.Inner.Error
	}

	result := resp.Inner.Results[0]
	if result == nil {
		return nil, finance.CreateRemoteErrorS("no quote summary in response")
	}

	return result, nil
}

// response is a yfin chart response.
type response struct {
	Inner struct {
		Results []*finance.QuoteSummary `json:"result"`
		Error   *finance.YfinError      `json:"error"`
	} `json:"quoteSummary"`
}
