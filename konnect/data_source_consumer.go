package konnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-konnect/konnect/client"
	"net/http"
	"strings"
)

func dataSourceConsumer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConsumerRead,
		Schema: map[string]*schema.Schema{
			"control_plane_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"search_custom_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConsumerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	//consumers do not support searching and filtering so do it manually after reading all consumers
	controlPlaneId := d.Get("control_plane_id").(string)
	requestPath := fmt.Sprintf(client.ConsumerPath, controlPlaneId)
	body, err := c.HttpRequest(ctx, true, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	retVals := &client.ConsumerCollection{}
	err = json.NewDecoder(body).Decode(retVals)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	//Check for a quick exit
	if len(retVals.Consumers) == 0 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No consumer exists"))
	}
	//Do manual searching
	filteredList := []client.Consumer{}
	searchUsername, ok := d.GetOk("search_username")
	if ok {
		searchUsernameLower := strings.ToLower(searchUsername.(string))
		for _, c := range retVals.Consumers {
			if strings.Contains(strings.ToLower(c.Username), searchUsernameLower) {
				filteredList = append(filteredList, c)
			}
		}
		retVals.Consumers = filteredList
		filteredList = []client.Consumer{}
	}
	searchCustomId, ok := d.GetOk("search_custom_id")
	if ok {
		searchCustomIdLower := strings.ToLower(searchCustomId.(string))
		for _, c := range retVals.Consumers {
			if strings.Contains(strings.ToLower(c.CustomId), searchCustomIdLower) {
				filteredList = append(filteredList, c)
			}
		}
		retVals.Consumers = filteredList
		filteredList = []client.Consumer{}
	}
	username, ok := d.GetOk("username")
	if ok {
		usernameLower := strings.ToLower(username.(string))
		for _, c := range retVals.Consumers {
			if strings.ToLower(c.Username) == usernameLower {
				filteredList = append(filteredList, c)
			}
		}
		retVals.Consumers = filteredList
		filteredList = []client.Consumer{}
	}
	customId, ok := d.GetOk("custom_id")
	if ok {
		customIdLower := strings.ToLower(customId.(string))
		for _, c := range retVals.Consumers {
			if strings.ToLower(c.CustomId) == customIdLower {
				filteredList = append(filteredList, c)
			}
		}
		retVals.Consumers = filteredList
		filteredList = []client.Consumer{}
	}
	numConsumers := len(retVals.Consumers)
	if numConsumers > 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single consumer"))
	} else if numConsumers != 1 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No consumer exists with that filter criteria"))
	}
	retVal := retVals.Consumers[0]
	retVal.ControlPlaneId = controlPlaneId
	d.Set("control_plane_id", retVal.ControlPlaneId)
	d.Set("username", retVal.Username)
	d.Set("custom_id", retVal.CustomId)
	d.Set("consumer_id", retVal.Id)
	d.SetId(retVal.ConsumerEncodeId())
	return diags
}
