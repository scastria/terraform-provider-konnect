package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
)

func resourceTeamMappings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamMappingsCreate,
		ReadContext:   resourceTeamMappingsRead,
		UpdateContext: resourceTeamMappingsUpdate,
		DeleteContext: resourceTeamMappingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"mapping": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Required: true,
						},
						"team_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func fillTeamMappings(c *client.TeamMappings, d *schema.ResourceData) {
	konnectMappings, ok := d.GetOk("mapping")
	if ok {
		c.MappingsRead = []client.TeamMapping{}
		for _, item := range konnectMappings.(*schema.Set).List() {
			itemMap := item.(map[string]interface{})
			itemGroup := itemMap["group"].(string)
			itemTeamIdsList := itemMap["team_ids"].(*schema.Set).List()
			var itemTeamIds []string
			for _, value := range itemTeamIdsList {
				itemTeamIds = append(itemTeamIds, value.(string))
			}
			mapping := client.TeamMapping{
				Group:   itemGroup,
				TeamIds: itemTeamIds,
			}
			c.MappingsRead = append(c.MappingsRead, mapping)
		}
	} else {
		c.MappingsRead = nil
	}
	c.MappingsWrite = c.MappingsRead
}

func fillResourceDataFromTeamMappings(c *client.TeamMappings, d *schema.ResourceData) {
	var konnectMappings []map[string]interface{}
	konnectMappings = nil
	if c.MappingsWrite != nil {
		c.MappingsRead = c.MappingsWrite
	}
	if c.MappingsRead != nil {
		for _, mapping := range c.MappingsRead {
			itemMap := map[string]interface{}{}
			itemMap["group"] = mapping.Group
			itemMap["team_ids"] = mapping.TeamIds
			konnectMappings = append(konnectMappings, itemMap)
		}
	}
	d.Set("mapping", konnectMappings)
}

func resourceTeamMappingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := resourceTeamMappingsCreateUpdate(ctx, d, m)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(client.TeamMappingsId)
	return diags
}

func resourceTeamMappingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if d.Id() != client.TeamMappingsId {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Id must be equal to %s", client.TeamMappingsId))
	}
	requestPath := fmt.Sprintf(client.TeamMappingsPath)
	body, err := c.HttpRequest(ctx, false, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.TeamMappings{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromTeamMappings(retVal, d)
	return diags
}

func resourceTeamMappingsCreateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upTeamMappings := client.TeamMappings{}
	fillTeamMappings(&upTeamMappings, d)
	//When writing, Konnect API uses mappings from MappingsWrite
	upTeamMappings.MappingsRead = nil
	err := json.NewEncoder(&buf).Encode(upTeamMappings)
	if err != nil {
		return err
	}
	requestPath := fmt.Sprintf(client.TeamMappingsPath)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, false, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return err
	}
	retVal := &client.TeamMappings{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return err
	}
	fillResourceDataFromTeamMappings(retVal, d)
	return nil
}

func resourceTeamMappingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	err := resourceTeamMappingsCreateUpdate(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceTeamMappingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}
