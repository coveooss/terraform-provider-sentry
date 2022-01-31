package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func dataSourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSentryOrganizationRead,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},

			"internal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Get("slug").(string)

	logging.Debugf("Reading Sentry org named with ID: %s", slug)
	org, _, err := client.Organizations.Get(slug)
	if err != nil {
		return diag.FromErr(err)
	}
	logging.Debugf("Read Sentry org named %s with ID: %s", org.Name, org.Slug)

	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)

	return nil
}
