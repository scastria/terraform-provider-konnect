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

func resourceConsumerHMAC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerHMACCreate,
		ReadContext:   resourceConsumerHMACRead,
		UpdateContext: resourceConsumerHMACUpdate,
		DeleteContext: resourceConsumerHMACDelete,
		CustomizeDiff: resourceConsumerHMACDiff,
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
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"hmac_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConsumerHMACDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
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

func fillConsumerHMAC(c *client.ConsumerHMAC, d *schema.ResourceData, defaultTags []string) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.ConsumerId = d.Get("consumer_id").(string)
	c.Username = d.Get("username").(string)
	secret, ok := d.GetOk("secret")
	if ok {
		c.Secret = secret.(string)
	}
	tags := []string{}
	tagsSet, ok := d.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
		c.Tags = tags
	}
	c.AllTags = unionArrays(tags, defaultTags)
	d.Set("all_tags", c.AllTags)
	d.Set("tags", subtractArrays(c.AllTags, defaultTags))
}

func fillResourceDataFromConsumerHMAC(c *client.ConsumerHMAC, d *schema.ResourceData, defaultTags []string) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("consumer_id", c.ConsumerId)
	d.Set("username", c.Username)
	d.Set("secret", c.Secret)
	d.Set("hmac_id", c.Id)
}

func resourceConsumerHMACCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumerHMAC := client.ConsumerHMAC{}
	fillConsumerHMAC(&newConsumerHMAC, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(newConsumerHMAC)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerHMACPath, newConsumerHMAC.ControlPlaneId, newConsumerHMAC.ConsumerId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerHMAC{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newConsumerHMAC.ControlPlaneId
	retVal.ConsumerId = newConsumerHMAC.ConsumerId
	d.SetId(retVal.ConsumerHMACEncodeId())
	fillResourceDataFromConsumerHMAC(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerHMACRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerHMACDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerHMACPathGet, controlPlaneId, consumerId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerHMAC{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerHMAC(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerHMACUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerHMACDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumerHMAC := client.ConsumerHMAC{}
	fillConsumerHMAC(&upConsumerHMAC, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(upConsumerHMAC)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerHMACPathGet, controlPlaneId, consumerId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerHMAC{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerHMAC(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerHMACDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerHMACDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerHMACPathGet, controlPlaneId, consumerId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
