package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryPluginCreate,
		ReadContext:   resourceSentryPluginRead,
		UpdateContext: resourceSentryPluginUpdate,
		DeleteContext: resourceSentryPluginDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryPluginImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the project belongs to",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project to create the plugin for",
			},
			"plugin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plugin ID",
			},
			"config": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Plugin config",
			},
		},
	}
}

func resourceSentryPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	plugin := d.Get("plugin").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Creating plugin %v in org %v for project %v", plugin, org, project)
	_, err := client.ProjectPlugins.Enable(org, project, plugin)
	if err != nil {
		return diag.FromErr(err)
	}
	logging.Debugf("Created plugin %v in org %v for project %v", plugin, org, project)

	d.SetId(plugin)

	params := d.Get("config").(map[string]interface{})
	if _, _, err := client.ProjectPlugins.Update(org, project, plugin, params); err != nil {
		return diag.FromErr(err)
	}

	return resourceSentryPluginRead(ctx, d, meta)
}

func resourceSentryPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Reading plugin with ID %v in org %v for project %v", id, org, project)
	plugin, resp, err := client.ProjectPlugins.Get(org, project, id)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	logging.Debugf("Read plugin with ID %v in org %v for project %v", plugin.ID, org, project)

	d.SetId(plugin.ID)

	pluginConfig := make(map[string]string)
	for _, entry := range plugin.Config {
		if v, ok := entry.Value.(string); ok {
			pluginConfig[entry.Name] = v
		}
	}

	config := make(map[string]string)
	for k := range d.Get("config").(map[string]interface{}) {
		config[k] = pluginConfig[k]
	}

	d.Set("config", config)

	return nil
}

func resourceSentryPluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Updating plugin with ID %v in org %v for project %v", id, org, project)
	params := d.Get("config").(map[string]interface{})
	_, _, err := client.ProjectPlugins.Update(org, project, id, params)
	if err != nil {
		return diag.FromErr(err)
	}
	logging.Debugf("Updated plugin with ID %v in org %v for project %v", id, org, project)

	return resourceSentryPluginRead(ctx, d, meta)
}

func resourceSentryPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Deleting plugin with ID %v in org %v for project %v", id, org, project)
	_, err := client.ProjectPlugins.Disable(org, project, id)
	logging.Debugf("Deleted plugin with ID %v in org %v for project %v", id, org, project)

	return diag.FromErr(err)
}
