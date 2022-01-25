package sentry

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryKeyCreate,
		Read:   resourceSentryKeyRead,
		Update: resourceSentryKeyUpdate,
		Delete: resourceSentryKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the key should be created for",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project the key should be created for",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the key",
			},
			"public": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rate_limit_window": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"rate_limit_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dsn_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsn_public": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsn_csp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSentryKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	proj, resp, err := client.Projects.Get(org, project)
	logging.LogHttpResponse(resp, proj, logging.TraceLevel)
	if found, err := checkClientGet(resp, err, d); !found {
		return fmt.Errorf("project not found \"%v\": %w", project, err)
	}

	params := &sentry.CreateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	logging.Debugf("Creating Sentry key named %s in org %s for project %s", params.Name, org, project)
	key, resp, err := client.ProjectKeys.Create(org, project, params)
	logging.LogHttpResponse(resp, key, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Created Sentry key named %s in org %s for project %s", key.Name, org, project)
	d.SetId(key.ID)

	return resourceSentryKeyRead(d, meta)
}

func resourceSentryKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Reading rule with ID %v in org %v for project %v", id, org, project)
	keys, resp, err := client.ProjectKeys.List(org, project)
	logging.LogHttpResponse(resp, keys, logging.TraceLevel)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}

	found := false

	for _, key := range keys {
		if key.ID == id {
			logging.Debugf("Read the sentry key with ID %s", id)
			d.SetId(key.ID)
			d.Set("name", key.Name)
			d.Set("public", key.Public)
			d.Set("secret", key.Secret)
			d.Set("project_id", key.ProjectID)
			d.Set("is_active", key.IsActive)

			if key.RateLimit != nil {
				d.Set("rate_limit_window", key.RateLimit.Window)
				d.Set("rate_limit_count", key.RateLimit.Count)
			}

			d.Set("dsn_secret", key.DSN.Secret)
			d.Set("dsn_public", key.DSN.Public)
			d.Set("dsn_csp", key.DSN.CSP)

			found = true

			break
		}
	}

	if !found {
		logging.Debugf("Sentry key with ID %s could not be found...", id)
		d.SetId("")
	}

	return nil
}

func resourceSentryKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	params := &sentry.UpdateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	logging.Debugf("Updating Sentry key with ID %s", id)
	key, resp, err := client.ProjectKeys.Update(org, project, id, params)
	logging.LogHttpResponse(resp, key, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Updated Sentry key ID to %s", key.ID)

	d.SetId(key.ID)
	return resourceSentryKeyRead(d, meta)
}

func resourceSentryKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Deleting Sentry key with ID %s", id)
	resp, err := client.ProjectKeys.Delete(org, project, id)
	logging.LogHttpResponse(resp, nil, logging.TraceLevel)
	logging.Debugf("Deleted Sentry key with ID %s", id)
	return err
}
