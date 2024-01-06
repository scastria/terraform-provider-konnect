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

func fillConsumerBasic(c *client.ConsumerBasic, d *schema.ResourceData) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.ConsumerId = d.Get("consumer_id").(string)
	c.Username = d.Get("username").(string)
	c.Password = d.Get("password").(string)
}

func fillResourceDataFromConsumerBasic(c *client.ConsumerBasic, d *schema.ResourceData) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("consumer_id", c.ConsumerId)
	d.Set("username", c.Username)
	d.Set("password_hash", c.Password)
	d.Set("basic_id", c.Id)
}

func resourceConsumerBasicCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumerBasic := client.ConsumerBasic{}
	fillConsumerBasic(&newConsumerBasic, d)
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
	fillResourceDataFromConsumerBasic(retVal, d)
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
	fillResourceDataFromConsumerBasic(retVal, d)
	return diags
}

func resourceConsumerBasicUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerBasicDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumerBasic := client.ConsumerBasic{}
	fillConsumerBasic(&upConsumerBasic, d)
	// Hide non-updateable fields
	//upTeam.IsPredefined = false
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
	fillResourceDataFromConsumerBasic(retVal, d)
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
