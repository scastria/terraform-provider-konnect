package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
)

func resourceConsumerJWT() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerJWTCreate,
		ReadContext:   resourceConsumerJWTRead,
		UpdateContext: resourceConsumerJWTUpdate,
		DeleteContext: resourceConsumerJWTDelete,
		CustomizeDiff: resourceConsumerJWTDiff,
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
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HS256",
				ValidateFunc: validation.StringInSlice([]string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512", "ES256", "ES384"}, false),
			},
			"rsa_public_key": {
				Type:     schema.TypeString,
				Optional: true,
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
			"jwt_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConsumerJWTDiff(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
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

func fillConsumerJWT(c *client.ConsumerJWT, d *schema.ResourceData, defaultTags []string) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.ConsumerId = d.Get("consumer_id").(string)
	c.Algorithm = d.Get("algorithm").(string)
	key, ok := d.GetOk("key")
	if ok {
		c.Key = key.(string)
	}
	secret, ok := d.GetOk("secret")
	if ok {
		c.Secret = secret.(string)
	}
	rsaPublicKey, ok := d.GetOk("rsa_public_key")
	if ok {
		c.RSAPublicKey = rsaPublicKey.(string)
	}
	tags := []string{}
	tagsSet, ok := d.GetOk("tags")
	if ok {
		tags = convertSetToArray(tagsSet.(*schema.Set))
		c.Tags = tags
	}
	c.AllTags = unionArrays(tags, defaultTags)
}

func fillResourceDataFromConsumerJWT(c *client.ConsumerJWT, d *schema.ResourceData, defaultTags []string) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("consumer_id", c.ConsumerId)
	d.Set("algorithm", c.Algorithm)
	d.Set("key", c.Key)
	d.Set("secret", c.Secret)
	d.Set("rsa_public_key", c.RSAPublicKey)
	d.Set("jwt_id", c.Id)
	d.Set("all_tags", c.AllTags)
	d.Set("tags", subtractArrays(c.AllTags, defaultTags))
}

func resourceConsumerJWTCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newConsumerJWT := client.ConsumerJWT{}
	fillConsumerJWT(&newConsumerJWT, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(newConsumerJWT)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerJWTPath, newConsumerJWT.ControlPlaneId, newConsumerJWT.ConsumerId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerJWT{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newConsumerJWT.ControlPlaneId
	retVal.ConsumerId = newConsumerJWT.ConsumerId
	d.SetId(retVal.ConsumerJWTEncodeId())
	fillResourceDataFromConsumerJWT(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerJWTRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerJWTDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerJWTPathGet, controlPlaneId, consumerId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerJWT{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerJWT(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerJWTUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerJWTDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upConsumerJWT := client.ConsumerJWT{}
	fillConsumerJWT(&upConsumerJWT, d, c.DefaultTags)
	err := json.NewEncoder(&buf).Encode(upConsumerJWT)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.ConsumerJWTPathGet, controlPlaneId, consumerId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.ConsumerJWT{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	retVal.ConsumerId = consumerId
	fillResourceDataFromConsumerJWT(retVal, d, c.DefaultTags)
	return diags
}

func resourceConsumerJWTDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, consumerId, id := client.ConsumerJWTDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.ConsumerJWTPathGet, controlPlaneId, consumerId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
