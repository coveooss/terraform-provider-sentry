package sentry

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jianyuan/go-sentry/sentry"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	Token   string
	BaseURL string
}

// Client to connect to Sentry.
func (c *Config) Client(ctx context.Context) (interface{}, error) {
	var baseURL *url.URL
	var err error

	if c.BaseURL != "" {
		baseURL, err = url.Parse(c.BaseURL)
		logging.Debugf("Parsing base url %s", c.BaseURL)
		tflog.Debug(ctx, "Parsing base url", "BaseUrl", c.BaseURL)
		if err != nil {
			return nil, err
		}
	} else {
		tflog.Warn(ctx, "No base URL was set for the Sentry client")
	}

	tflog.Info(ctx, "Instantiating Sentry client...")
	cl := sentry.NewClient(nil, baseURL, c.Token)

	return cl, nil
}
