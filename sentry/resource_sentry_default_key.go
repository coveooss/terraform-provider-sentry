package sentry

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryDefaultKey() *schema.Resource {
	// reuse read and update operations
	dKey := resourceSentryKey()
	dKey.Create = resourceSentryDefaultKeyCreate
	dKey.Delete = resourceAwsDefaultVpcDelete

	// Key name is a computed resource for default key
	dKey.Schema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    true,
		Description: "The name of the key",
	}

	return dKey
}

func resourceSentryDefaultKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	keys, resp, err := client.ProjectKeys.List(org, project)
	logging.LogHttpResponse(resp, keys, logging.TraceLevel)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}

	if len(keys) < 1 {
		return fmt.Errorf("Default key not found on the project")
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].DateCreated.Before(keys[j].DateCreated)
	})

	id := keys[0].ID
	params := &sentry.UpdateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	logging.Debugf("Creating Sentry default key in org %s for project %s with ID %s", org, project, id)
	key, resp, err := client.ProjectKeys.Update(org, project, id, params)
	logging.LogHttpResponse(resp, key, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Created Sentry default key in org %s for project %s with ID %s", org, project, id)

	d.SetId(id)
	return resourceSentryKeyRead(d, meta)
}

func resourceAwsDefaultVpcDelete(d *schema.ResourceData, meta interface{}) error {
	logging.Warning("Cannot destroy Default Key. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
