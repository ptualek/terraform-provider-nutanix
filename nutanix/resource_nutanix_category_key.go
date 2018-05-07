package nutanix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixCategoryKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixCategoryKeyKeyCreateOrUpdate,
		Read:   resourceNutanixCategoryKeyKeyRead,
		Update: resourceNutanixCategoryKeyKeyCreateOrUpdate,
		Delete: resourceNutanixCategoryKeyKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCategoryKeySchema(),
	}
}

func resourceNutanixCategoryKeyKeyCreateOrUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating CategoryKey: %s", resourceData.Get("name").(string))

	conn := meta.(*NutanixClient).API

	request := &v3.CategoryKey{}

	name, nameOK := resourceData.GetOk("name")

	// Read Arguments and set request values
	if v, ok := resourceData.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}

	if desc, ok := resourceData.GetOk("description"); ok {
		request.Description = utils.String(desc.(string))
	}

	// validaste required fields
	if !nameOK {
		return fmt.Errorf("Please provide the required attribute name")
	}

	request.Name = utils.String(name.(string))

	//Make request to the API
	resp, err := conn.V3.CreateOrUpdateCategoryKey(request)

	if err != nil {
		return err
	}

	n := *resp.Name

	// set terraform state
	resourceData.SetId(n)

	return resourceNutanixCategoryKeyKeyRead(resourceData, meta)
}

func resourceNutanixCategoryKeyKeyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading CategoryKey: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*NutanixClient).API

	// Make request to the API
	resp, err := conn.V3.GetCategoryKey(d.Id())

	if err != nil {
		return err
	}

	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}

	if err := d.Set("name", utils.StringValue(resp.Name)); err != nil {
		return err
	}

	if err := d.Set("description", utils.StringValue(resp.Description)); err != nil {
		return err
	}

	if err := d.Set("system_defined", utils.BoolValue(resp.SystemDefined)); err != nil {
		return err
	}

	return nil
}

func resourceNutanixCategoryKeyKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API

	log.Printf("Destroying the category with the name %s", d.Id())
	fmt.Printf("Destroying the category with the name %s", d.Id())

	if err := conn.V3.DeleteCategoryKey(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func getCategoryKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"system_defined": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
