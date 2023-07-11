package provider

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/luminsports/terraform-provider-bless/internal/aws"
	"github.com/luminsports/terraform-provider-bless/internal/util"
	"github.com/pkg/errors"
)

// ED25519CA is an ED25519 CA resource.
func ED25519CA() *schema.Resource {
	ca := newResourceED25519CA()
	return &schema.Resource{
		Create: ca.Create,
		Read:   ca.Read,
		Delete: ca.Delete,

		Schema: map[string]*schema.Schema{
			schemaKmsKeyID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The kms key with which we should encrypt the CA password.",
				ForceNew:    true,
			},

			// computed
			schemaEncryptedPrivateKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the base64 encoded CA encrypted private key.",
			},
			schemaPublicKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the plaintext CA public key in openssh format.",
			},
			schemaEncryptedPassword: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the kms encrypted password.",
			},
		},
	}
}

// resourceCA is a namespace.
type resourceED25519CA struct{}

func newResourceED25519CA() *resourceED25519CA {
	return &resourceED25519CA{}
}

// Create creates a CA.
func (ca *resourceED25519CA) Create(d *schema.ResourceData, meta interface{}) error {
	awsClient, ok := meta.(*aws.Client)
	if !ok {
		return errors.New("meta is not of type *aws.Client")
	}

	kmsKeyID := d.Get(schemaKmsKeyID).(string)
	keyPair, err := ca.createKeypair()
	if err != nil {
		return err
	}

	encryptedPassword, err := awsClient.KMS.EncryptBytes(keyPair.Password, kmsKeyID)
	if err != nil {
		return err
	}

	d.Set(schemaEncryptedPrivateKey, keyPair.B64EncryptedPrivateKey) // nolint
	d.Set(schemaPublicKey, keyPair.PublicKey)                        // nolint
	d.Set(schemaEncryptedPassword, encryptedPassword)                // nolint
	d.SetId(util.HashForState(keyPair.PublicKey))
	return nil
}

// Read reads the ca.
func (ca *resourceED25519CA) Read(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// Delete deletes the ca.
func (ca *resourceED25519CA) Delete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

// ------------ helpers ------------------.
func (ca *resourceED25519CA) createKeypair() (*util.CA, error) {
	// generate private key
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "Private key generation failed")
	}
	return util.NewCA(privateKey, publicKey, caPasswordBytes)
}
