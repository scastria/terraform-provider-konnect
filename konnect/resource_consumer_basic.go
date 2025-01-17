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

func resourceConsumerBasic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerBasicCreate,
		ReadContext:   resourceConsumerBasicRead,
		UpdateContext: resourceConsumerBasicUpdate,
		DeleteContext: resourceConsumerBasicDelete,
		CustomizeDiff: resourceConsumerBasicDiff,
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
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
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
			"password_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"basic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConsumerBasicDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
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

func fillConsumerBasic(c *client.ConsumerBasic, d *schema.ResourceData, defaultTags []string) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.ConsumerId = d.Get("consumer_id").(string)
	c.Username = d.Get("username").(string)
	c.Password = d.Get("password").(string)
	tags := []string{}
	tagsSet, ok := d.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
		c.Tags = tags
	}
	c.AllTags = unionArrays(tags, defaultTags)
}

func fillResourceDataFromConsumerBasic(c *client.ConsumerBasic, d *schema.ResourceData, defaultTags []string) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("consumer_id", c.ConsumerId)
	d.Set("username", c.Username)
	d.Set("password_hash", c.Password)
	d.Set("basic_id", c.Id)
	d.Set("all_tags", c.AllTags)
	d.Set("tags", subtractArrays(c.AllTags, defaultTags))
}

func resourceConsumerBasicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumerBasic := client.ConsumerBasic{}
	fillConsumerBasic(&newConsumerBasic, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(newConsumerBasic)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerBasicPath, newConsumerBasic.ControlPlaneId, newConsumerBasic.ConsumerId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerBasic{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newConsumerBasic.ControlPlaneId
	retVal.ConsumerId = newConsumerBasic.ConsumerId
	d.SetId(retVal.ConsumerBasicEncodeId())
	fillResourceDataFromConsumerBasic(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerBasicRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerBasicDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerBasicPathGet, controlPlaneId, consumerId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerBasic{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerBasic(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerBasicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerBasicDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumerBasic := client.ConsumerBasic{}
	fillConsumerBasic(&upConsumerBasic, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(upConsumerBasic)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerBasicPathGet, controlPlaneId, consumerId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerBasic{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerBasic(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerBasicDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerBasicDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerBasicPathGet, controlPlaneId, consumerId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
