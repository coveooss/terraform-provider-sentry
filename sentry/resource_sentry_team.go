package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryTeamCreate,
		ReadContext:   resourceSentryTeamRead,
		UpdateContext: resourceSentryTeamUpdate,
		DeleteContext: resourceSentryTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSentryTeamImport,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the team should be created for",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional slug for this team",
				Computed:    true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_pending": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_member": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSentryTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	params := &sentry.CreateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}
	logging.Debugf("Creating Sentry team %s in org %s", params.Name, org)

	team, _, err := client.Teams.Create(org, params)
	if err != nil {
		return diag.FromErr(err)
	}
	logging.Debugf("Created Sentry team %s in org %s", team.Name, org)

	d.SetId(team.Slug)
	return resourceSentryTeamRead(ctx, d, meta)
}

func resourceSentryTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	logging.Debugf("Reading Sentry team %s in org %s", slug, org)
	team, resp, err := client.Teams.Get(org, slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	logging.Debugf("Read Sentry team %s in org %s", team.Slug, org)

	d.SetId(team.Slug)
	d.Set("team_id", team.ID)
	d.Set("name", team.Name)
	d.Set("slug", team.Slug)
	d.Set("organization", org)
	d.Set("has_access", team.HasAccess)
	d.Set("is_pending", team.IsPending)
	d.Set("is_member", team.IsMember)
	return nil
}

func resourceSentryTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	params := &sentry.UpdateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	logging.Debugf("Updating Sentry team %s in org %s", slug, org)
	team, _, err := client.Teams.Update(org, slug, params)
	if err != nil {
		return diag.FromErr(err)
	}
	logging.Debugf("Updated Sentry team %s in org %s", team.Slug, org)

	d.SetId(team.Slug)
	return resourceSentryTeamRead(ctx, d, meta)
}

func resourceSentryTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)

	logging.Debugf("Deleting Sentry team %s in org %s", slug, org)
	_, err := client.Teams.Delete(org, slug)
	logging.Debugf("Deleted Sentry team %s in org %s", slug, org)

	return diag.FromErr(err)
}
