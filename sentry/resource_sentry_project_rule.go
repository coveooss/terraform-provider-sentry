package sentry

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
	"github.com/mitchellh/mapstructure"
)

const (
	defaultActionMatch = "any"
	defaultFilterMatch = "any"
	defaultFrequency   = 30
)

func resourceSentryRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryRuleCreate,
		Read:   resourceSentryRuleRead,
		Update: resourceSentryRuleUpdate,
		Delete: resourceSentryRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryRuleImporter,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule name",
			},
			"action_match": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter_match": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Perform actions at most once every X minutes",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Perform rule in a specific environment",
			},
		},
	}
}

func resourceSentryRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	name := d.Get("name").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	environment := d.Get("environment").(string)
	actionMatch := d.Get("action_match").(string)
	filterMatch := d.Get("filter_match").(string)
	inputConditions := d.Get("conditions").([]interface{})
	inputActions := d.Get("actions").([]interface{})
	inputFilters := d.Get("filters").([]interface{})
	frequency := d.Get("frequency").(int)

	if actionMatch == "" {
		actionMatch = defaultActionMatch
	}
	if filterMatch == "" {
		filterMatch = defaultFilterMatch
	}
	if frequency == 0 {
		frequency = defaultFrequency
	}

	conditions := make([]sentry.ConditionType, len(inputConditions))
	for i, ic := range inputConditions {
		var condition sentry.ConditionType
		mapstructure.WeakDecode(ic, &condition)
		conditions[i] = condition
	}
	actions := make([]sentry.ActionType, len(inputActions))
	for i, ia := range inputActions {
		var action sentry.ActionType
		mapstructure.WeakDecode(ia, &action)
		actions[i] = action
	}
	filters := make([]sentry.FilterType, len(inputFilters))
	for i, ia := range inputFilters {
		var filter sentry.FilterType
		mapstructure.WeakDecode(ia, &filter)
		filters[i] = filter
	}

	params := &sentry.CreateRuleParams{
		ActionMatch: actionMatch,
		FilterMatch: filterMatch,
		Environment: environment,
		Frequency:   frequency,
		Name:        name,
		Conditions:  conditions,
		Actions:     actions,
		Filters:     filters,
	}

	if environment != "" {
		params.Environment = environment
	}

	logging.Debugf("Creating rule with name %v in org %v for project %v", name, org, project)
	rule, resp, err := client.Rules.Create(org, project, params)
	logging.LogHttpResponse(resp, rule, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Created rule with name %v and ID %s in org %v for project %v", rule.Name, rule.ID, org, project)

	d.SetId(rule.ID)

	return resourceSentryRuleRead(d, meta)
}

func resourceSentryRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	id := d.Id()

	logging.Debugf("Reading rule with ID %v in org %v for project %v", id, org, project)
	rules, resp, err := client.Rules.List(org, project)
	logging.LogHttpResponse(resp, rules, logging.TraceLevel)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}

	var rule *sentry.Rule
	for _, r := range rules {
		if r.ID == id {
			rule = &r
			break
		}
	}

	if rule == nil {
		if id == "" {
			logging.Error("The rule ID was never set (ID was '')")
		}
		return errors.New("Could not find rule with ID " + id)
	}
	logging.Debugf("Read rule with ID %v in org %v for project %v", rule.ID, org, project)

	// workaround for
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/62
	// as the data sent by Sentry is integer
	for _, f := range rule.Actions {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	for _, f := range rule.Conditions {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	for _, f := range rule.Filters {
		for k, v := range f {
			switch vv := v.(type) {
			case float64:
				f[k] = fmt.Sprintf("%.0f", vv)
			}
		}
	}

	d.SetId(rule.ID)
	d.Set("name", rule.Name)
	d.Set("frequency", rule.Frequency)
	d.Set("environment", rule.Environment)
	d.Set("filters", rule.Filters)
	d.Set("actions", rule.Actions)
	d.Set("conditions", rule.Conditions)
	d.Set("action_match", rule.ActionMatch)
	d.Set("filter_match", rule.FilterMatch)

	return nil
}

func resourceSentryRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	name := d.Get("name").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	environment := d.Get("environment").(string)
	actionMatch := d.Get("action_match").(string)
	filterMatch := d.Get("filter_match").(string)
	inputConditions := d.Get("conditions").([]interface{})
	inputActions := d.Get("actions").([]interface{})
	inputFilters := d.Get("filters").([]interface{})
	frequency := d.Get("frequency").(int)

	if actionMatch == "" {
		actionMatch = defaultActionMatch
	}
	if filterMatch == "" {
		filterMatch = defaultFilterMatch
	}
	if frequency == 0 {
		frequency = defaultFrequency
	}

	conditions := make([]sentry.ConditionType, len(inputConditions))
	for i, ic := range inputConditions {
		var condition sentry.ConditionType
		mapstructure.WeakDecode(ic, &condition)
		conditions[i] = condition
	}
	actions := make([]sentry.ActionType, len(inputActions))
	for i, ia := range inputActions {
		var action sentry.ActionType
		mapstructure.WeakDecode(ia, &action)
		actions[i] = action
	}
	filters := make([]sentry.FilterType, len(inputFilters))
	for i, ia := range inputFilters {
		var filter sentry.FilterType
		mapstructure.WeakDecode(ia, &filter)
		filters[i] = filter
	}

	params := &sentry.Rule{
		ID:          id,
		ActionMatch: actionMatch,
		FilterMatch: filterMatch,
		Frequency:   frequency,
		Name:        name,
		Conditions:  conditions,
		Actions:     actions,
		Filters:     filters,
	}

	if environment != "" {
		params.Environment = &environment
	}

	logging.Debugf("Updating rule with ID %v in org %v for project %v", id, org, project)
	rule, resp, err := client.Rules.Update(org, project, id, params)
	logging.LogHttpResponse(resp, rule, logging.TraceLevel)
	if err != nil {
		return err
	}
	logging.Debugf("Updated rule with ID %v in org %v for project %v", id, org, project)

	return resourceSentryRuleRead(d, meta)
}

func resourceSentryRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	logging.Debugf("Deleting rule with ID %v in org %v for project %v", id, org, project)
	resp, err := client.Rules.Delete(org, project, id)
	logging.LogHttpResponse(resp, nil, logging.TraceLevel)
	logging.Debugf("Deleted rule with ID %v in org %v for project %v", id, org, project)

	return err
}
