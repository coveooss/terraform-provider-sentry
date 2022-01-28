package sentry

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryOrganizationCreate,
		Read:   resourceSentryOrganizationRead,
		Update: resourceSentryOrganizationUpdate,
		Delete: resourceSentryOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The human readable name for the organization",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique URL slug for this organization",
				Computed:    true,
			},
			"agree_terms": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "You agree to the applicable terms of service and privacy policy",
			},
		},
	}
}

func resourceSentryOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	params := &sentry.CreateOrganizationParams{
		Name:       d.Get("name").(string),
		Slug:       d.Get("slug").(string),
		AgreeTerms: sentry.Bool(d.Get("agree_terms").(bool)),
	}

	logging.Debugf("Creating Sentry organization %s", params.Name)
	org, _, err := client.Organizations.Create(params)
	if err != nil {
		return err
	}
	logging.Debugf("Created Sentry organization %s", org.Name)

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(d, meta)
}

func resourceSentryOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()

	logging.Debugf("Reading Sentry organization %s", slug)
	org, resp, err := client.Organizations.Get(slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}
	logging.Debugf("Read Sentry organization %s", org.Slug)

	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)
	return nil
}

func resourceSentryOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	params := &sentry.UpdateOrganizationParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	logging.Debugf("Updating Sentry organization %s", slug)
	org, _, err := client.Organizations.Update(slug, params)
	if err != nil {
		return err
	}
	logging.Debugf("Updated Sentry organization %s", org.Slug)

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(d, meta)
}

func resourceSentryOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()

	logging.Debugf("Deleting Sentry organization %s", slug)
	_, err := client.Organizations.Delete(slug)
	logging.Debugf("Deleted Sentry organization %s", slug)

	return err
}
