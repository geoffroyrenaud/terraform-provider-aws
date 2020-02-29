package aws

import (
	//"encoding/base64"
	"log"
	"time"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAwsKmsCiphertextDecrypt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsKmsCiphertextDecryptRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAwsKmsCiphertextDecryptRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).kmsconn

	d.SetId(time.Now().UTC().String())

	req := &kms.DecryptInput{
		// KeyId:          aws.String(d.Get("key_id").(string)),
		CiphertextBlob: []byte(d.Get("ciphertext_blob").(string)),
	}

	if ec := d.Get("context"); ec != nil {
		req.EncryptionContext = stringMapToPointers(ec.(map[string]interface{}))
	}

	log.Printf("[DEBUG] KMS decrypt for key: %s", d.Get("key_id").(string))
	resp, err := conn.Decrypt(req)
	if err != nil {
		return err
	}

	//d.Set("plaintext", base64.StdEncoding.DecodeString(resp.Plaintext))
	d.Set("plaintext", resp.Plaintext)

	return nil
}
