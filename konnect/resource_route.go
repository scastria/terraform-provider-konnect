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

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouteCreate,
		ReadContext:   resourceRouteRead,
		UpdateContext: resourceRouteUpdate,
		DeleteContext: resourceRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
				},
			},
			"methods": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"GET", "PUT", "POST", "PATCH", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE"}, false),
				},
			},
			"hosts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"paths": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"https_redirect_status_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{426, 301, 302, 307, 308}),
				Default:      426,
			},
			"regex_priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"strip_path": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"path_handling": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"v0", "v1"}, false),
				Default:      "v0",
			},
			"preserve_host": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"request_buffering": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"response_buffering": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"header": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
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

func fillRoute(c *client.Route, d *schema.ResourceData) {
	c.ControlPlaneId = d.Get("control_plane_id").(string)
	c.StripPath = d.Get("strip_path").(bool)
	c.PreserveHost = d.Get("preserve_host").(bool)
	c.RequestBuffering = d.Get("request_buffering").(bool)
	c.ResponseBuffering = d.Get("response_buffering").(bool)
	name, ok := d.GetOk("name")
	if ok {
		c.Name = name.(string)
	}
	protocols, ok := d.GetOk("protocols")
	if ok {
		c.Protocols = convertSetToArray(protocols.(*schema.Set))
	}
	methods, ok := d.GetOk("methods")
	if ok {
		c.Methods = convertSetToArray(methods.(*schema.Set))
	}
	hosts, ok := d.GetOk("hosts")
	if ok {
		c.Hosts = convertSetToArray(hosts.(*schema.Set))
	}
	paths, ok := d.GetOk("paths")
	if ok {
		c.Paths = convertSetToArray(paths.(*schema.Set))
	}
	httpsRedirectStatusCode, ok := d.GetOk("https_redirect_status_code")
	if ok {
		c.HTTPSRedirectStatusCode = httpsRedirectStatusCode.(int)
	}
	regexPriority, ok := d.GetOk("regex_priority")
	if ok {
		c.RegexPriority = regexPriority.(int)
	}
	pathHandling, ok := d.GetOk("path_handling")
	if ok {
		c.PathHandling = pathHandling.(string)
	}
	service, ok := d.GetOk("service_id")
	if ok {
		c.Service = &client.EntityId{}
		c.Service.Id = service.(string)
	}
	konnectHeaders, ok := d.GetOk("header")
	if ok {
		c.Headers = map[string][]string{}
		for _, item := range konnectHeaders.(*schema.Set).List() {
			itemMap := item.(map[string]interface{})
			itemName := itemMap["name"].(string)
			itemValuesList := itemMap["values"].(*schema.Set).List()
			var itemValues []string
			for _, value := range itemValuesList {
				itemValues = append(itemValues, value.(string))
			}
			c.Headers[itemName] = itemValues
		}
	} else {
		c.Headers = nil
	}
}

func fillResourceDataFromRoute(c *client.Route, d *schema.ResourceData) {
	d.Set("control_plane_id", c.ControlPlaneId)
	d.Set("name", c.Name)
	d.Set("protocols", c.Protocols)
	d.Set("methods", c.Methods)
	d.Set("hosts", c.Hosts)
	d.Set("paths", c.Paths)
	d.Set("https_redirect_status_code", c.HTTPSRedirectStatusCode)
	d.Set("regex_priority", c.RegexPriority)
	d.Set("strip_path", c.StripPath)
	d.Set("path_handling", c.PathHandling)
	d.Set("preserve_host", c.PreserveHost)
	d.Set("request_buffering", c.RequestBuffering)
	d.Set("response_buffering", c.ResponseBuffering)
	serviceId := ""
	if c.Service != nil {
		serviceId = c.Service.Id
	}
	d.Set("service_id", serviceId)
	d.Set("route_id", c.Id)
	var konnectHeaders []map[string]interface{}
	konnectHeaders = nil
	if c.Headers != nil {
		for name, values := range c.Headers {
			itemMap := map[string]interface{}{}
			itemMap["name"] = name
			itemMap["values"] = values
			konnectHeaders = append(konnectHeaders, itemMap)
		}
	}
	d.Set("header", konnectHeaders)
}

func resourceRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	newRoute := client.Route{}
	fillRoute(&newRoute, d)
	err := json.NewEncoder(&buf).Encode(newRoute)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RoutePath, newRoute.ControlPlaneId)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPost, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal := &client.Route{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = newRoute.ControlPlaneId
	d.SetId(retVal.RouteEncodeId())
	fillResourceDataFromRoute(retVal, d)
	return diags
}

func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.RouteDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RoutePathGet, controlPlaneId, id)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re := err.(*client.RequestError)
		if re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.Route{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	fillResourceDataFromRoute(retVal, d)
	return diags
}

func resourceRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.RouteDecodeId(d.Id())
	c := m.(*client.Client)
	buf := bytes.Buffer{}
	upRoute := client.Route{}
	fillRoute(&upRoute, d)
	// Hide non-updateable fields
	//upTeam.IsPredefined = false
	err := json.NewEncoder(&buf).Encode(upRoute)
	if err != nil {
		return diag.FromErr(err)
	}
	requestPath := fmt.Sprintf(client.RoutePathGet, controlPlaneId, id)
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
	}
	body, err := c.HttpRequest(ctx, true, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.Route{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal.ControlPlaneId = controlPlaneId
	fillResourceDataFromRoute(retVal, d)
	return diags
}

func resourceRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	controlPlaneId, id := client.RouteDecodeId(d.Id())
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.RoutePathGet, controlPlaneId, id)
	_, err := c.HttpRequest(ctx, true, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
