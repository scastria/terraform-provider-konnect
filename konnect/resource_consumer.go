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

func resourceConsumer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerCreate,
		ReadContext:   resourceConsumerRead,
		UpdateContext: resourceConsumerUpdate,
		DeleteContext: resourceConsumerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"custom_id"},
			},
			"custom_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"username"},
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fillConsumer(c *client.Consumer, d *schema.ResourceData) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	username, ok := d.GetOk("username")
	if ok {
		c.Username = username.(string)
	}
	customId, ok := d.GetOk("custom_id")
	if ok {
		c.CustomId = customId.(string)
	}
}

func fillResourceDataFromConsumer(c *client.Consumer, d *schema.ResourceData) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("username", c.Username)
	d.Set("custom_id", c.CustomId)
	d.Set("consumer_id", c.Id)
}

func resourceConsumerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumer := client.Consumer{}
	fillConsumer(&newConsumer, d)
	err := json.NewEncoder(&buf).Encode(newConsumer)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerPath, newConsumer.ControlPlaneId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Consumer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newConsumer.ControlPlaneId
	d.SetId(retVal.ConsumerEncodeId())
	fillResourceDataFromConsumer(retVal, d)
	return diags
}

func resourceConsumerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.ConsumerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerPathGet, controlPlaneId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Consumer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	fillResourceDataFromConsumer(retVal, d)
	return diags
}

func resourceConsumerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.ConsumerDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumer := client.Consumer{}
	fillConsumer(&upConsumer, d)
	err := json.NewEncoder(&buf).Encode(upConsumer)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerPathGet, controlPlaneId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Consumer{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	fillResourceDataFromConsumer(retVal, d)
	return diags
}

func resourceConsumerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.ConsumerDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerPathGet, controlPlaneId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
