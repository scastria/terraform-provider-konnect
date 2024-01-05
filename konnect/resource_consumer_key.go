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

func resourceConsumerKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerKeyCreate,
		ReadContext:   resourceConsumerKeyRead,
		UpdateContext: resourceConsumerKeyUpdate,
		DeleteContext: resourceConsumerKeyDelete,
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
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func fillConsumerKey(c *client.ConsumerKey, d *schema.ResourceData) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.ConsumerId = d.Get("consumer_id").(string)
	key, ok := d.GetOk("key")
	if ok {
		c.Key = key.(string)
	}
}

func fillResourceDataFromConsumerKey(c *client.ConsumerKey, d *schema.ResourceData) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("consumer_id", c.ConsumerId)
	d.Set("key", c.Key)
	d.Set("key_id", c.Id)
}

func resourceConsumerKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumerKey := client.ConsumerKey{}
	fillConsumerKey(&newConsumerKey, d)
	err := json.NewEncoder(&buf).Encode(newConsumerKey)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerKeyPath, newConsumerKey.ControlPlaneId, newConsumerKey.ConsumerId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerKey{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newConsumerKey.ControlPlaneId
	retVal.ConsumerId = newConsumerKey.ConsumerId
	d.SetId(retVal.ConsumerKeyEncodeId())
	fillResourceDataFromConsumerKey(retVal, d)
	return diags
}

func resourceConsumerKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerKeyDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerKeyPathGet, controlPlaneId, consumerId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerKey{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerKey(retVal, d)
	return diags
}

func resourceConsumerKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerKeyDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumerKey := client.ConsumerKey{}
	fillConsumerKey(&upConsumerKey, d)
	// Hide non-updateable fields
	//upTeam.IsPredefined = false
	err := json.NewEncoder(&buf).Encode(upConsumerKey)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerKeyPathGet, controlPlaneId, consumerId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerKey{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerKey(retVal, d)
	return diags
}

func resourceConsumerKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerKeyDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerKeyPathGet, controlPlaneId, consumerId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
