package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jianyuan/terraform-provider-sentry/sentry"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sentry.Provider,
		Logger:       hclog.Default(), // might not be mandatory, should test removing this
	})
}
