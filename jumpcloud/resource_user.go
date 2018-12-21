package jumpcloud

import (
	"context"
	//    "net/http"
	"fmt"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"xorgid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	configv2 := m.(*jcapiv2.Configuration)
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", configv2.DefaultHeader["x-api-key"])

	client := jcapiv1.NewAPIClient(configv1)

	var payload jcapiv1.Systemuserputpost
	payload.Username = d.Get("username").(string)
	payload.Email = d.Get("email").(string)
	payload.Firstname = d.Get("firstname").(string)
	payload.Lastname = d.Get("lastname").(string)
	payload.EnableUserPortalMultifactor = d.Get("enable_mfa").(bool)

	req := map[string]interface{}{
		"body":   payload,
		"xOrgId": d.Get("xorgid").(string),
	}
	returnstruc, _, err := client.SystemusersApi.SystemusersPost(context.TODO(),
		"", "", req)
	if err != nil {
		return err
	}
	d.SetId(returnstruc.Id)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	configv2 := m.(*jcapiv2.Configuration)
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", configv2.DefaultHeader["x-api-key"])

	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersGet(context.TODO(),
		d.Id(), "", "", nil)

	if err != nil {
		return err
	}

	// Need to understand why we need this
	/*
		if !ok {
			// not found
			d.SetId("")
			return nil
		}
	*/
	d.SetId(res.Id)

	if err := d.Set("username", res.Username); err != nil {
		return err
	}
	if err := d.Set("email", res.Email); err != nil {
		return err
	}
	if err := d.Set("firstname", res.Firstname); err != nil {
		return err
	}
	if err := d.Set("lastname", res.Lastname); err != nil {
		return err
	}

	if err := d.Set("enable_mfa", res.EnableUserPortalMultifactor); err != nil {
		return err
	}
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	configv2 := m.(*jcapiv2.Configuration)
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", configv2.DefaultHeader["x-api-key"])

	client := jcapiv1.NewAPIClient(configv1)

	var payload jcapiv1.Systemuserput
	payload.Username = d.Get("username").(string)
	payload.Email = d.Get("email").(string)
	payload.Firstname = d.Get("firstname").(string)
	payload.Lastname = d.Get("lastname").(string)
	payload.EnableUserPortalMultifactor = d.Get("enable_mfa").(bool)

	req := map[string]interface{}{
		"body":   payload,
		"xOrgId": d.Get("xorgid").(string),
	}
	_, _, err := client.SystemusersApi.SystemusersPut(context.TODO(),
		d.Id(), "", "", req)
	if err != nil {
		return err
	}
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	configv2 := m.(*jcapiv2.Configuration)
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", configv2.DefaultHeader["x-api-key"])

	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersDelete(context.TODO(),
		d.Id(), "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting user group:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}
