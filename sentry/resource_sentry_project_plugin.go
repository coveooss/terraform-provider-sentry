package sentry

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryPluginCreate,
		Read:   resourceSentryPluginRead,
		Update: resourceSentryPluginUpdate,
		Delete: resourceSentryPluginDelete,
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

func resourceSentryPluginCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	plugin := d.Get("plugin").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Creating plugin %v in org %v for project %v", plugin, org, project)
	resp, err := client.ProjectPlugins.Enable(org, project, plugin)
	logging.LogHttpResponse(resp, nil, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Created plugin %v in org %v for project %v", plugin, org, project)

	d.SetId(plugin)

	params := d.Get("config").(map[string]interface{})
	pluginObj, resp, err := client.ProjectPlugins.Update(org, project, plugin, params)
	logging.LogHttpResponse(resp, pluginObj, logging.TraceLevel)
	if err != nil {
		return err
	}

	return resourceSentryPluginRead(d, meta)
}

func resourceSentryPluginRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Reading plugin with ID %v in org %v for project %v", id, org, project)
	plugin, resp, err := client.ProjectPlugins.Get(org, project, id)
	logging.LogHttpResponse(resp, plugin, logging.TraceLevel)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
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

func resourceSentryPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Updating plugin with ID %v in org %v for project %v", id, org, project)
	params := d.Get("config").(map[string]interface{})
	plugin, resp, err := client.ProjectPlugins.Update(org, project, id, params)
	logging.LogHttpResponse(resp, plugin, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Updated plugin with ID %v in org %v for project %v", id, org, project)

	return resourceSentryPluginRead(d, meta)
}

func resourceSentryPluginDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Deleting plugin with ID %v in org %v for project %v", id, org, project)
	resp, err := client.ProjectPlugins.Disable(org, project, id)
	logging.LogHttpResponse(resp, nil, logging.TraceLevel)
	logging.Debugf("Deleted plugin with ID %v in org %v for project %v", id, org, project)

	return err
}
