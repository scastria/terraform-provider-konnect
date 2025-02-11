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

func resourceConsumerACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerACLCreate,
		ReadContext:   resourceConsumerACLRead,
		UpdateContext: resourceConsumerACLUpdate,
		DeleteContext: resourceConsumerACLDelete,
		CustomizeDiff: resourceConsumerACLDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConsumerACLDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	c := m.(*client.Client)
	tags := []string{}
	tagsSet, ok := diff.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
	}
	allTags := unionArrays(tags, c.DefaultTags)
	diff.SetNew("all_tags", allTags)
	return nil
}

func fillConsumerACL(c *client.ConsumerACL, d *schema.ResourceData, defaultTags []string) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.ConsumerId = d.Get("consumer_id").(string)
	c.Group = d.Get("group").(string)
	tags := []string{}
	tagsSet, ok := d.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
		c.Tags = tags
	}
	c.AllTags = unionArrays(tags, defaultTags)
}

func fillResourceDataFromConsumerACL(c *client.ConsumerACL, d *schema.ResourceData, defaultTags []string) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("consumer_id", c.ConsumerId)
	d.Set("group", c.Group)
	d.Set("acl_id", c.Id)
	d.Set("all_tags", c.AllTags)
	d.Set("tags", subtractArrays(c.AllTags, defaultTags))
}

func resourceConsumerACLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumerACL := client.ConsumerACL{}
	fillConsumerACL(&newConsumerACL, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(newConsumerACL)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerACLPath, newConsumerACL.ControlPlaneId, newConsumerACL.ConsumerId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerACL{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newConsumerACL.ControlPlaneId
	retVal.ConsumerId = newConsumerACL.ConsumerId
	d.SetId(retVal.ConsumerACLEncodeId())
	fillResourceDataFromConsumerACL(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerACLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerACLDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerACLPathGet, controlPlaneId, consumerId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerACL{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerACL(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerACLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerACLDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumerACL := client.ConsumerACL{}
	fillConsumerACL(&upConsumerACL, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(upConsumerACL)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerACLPathGet, controlPlaneId, consumerId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerACL{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerACL(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerACLDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerACLPathGet, controlPlaneId, consumerId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
