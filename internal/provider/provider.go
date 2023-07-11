package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/luminsports/terraform-provider-bless/internal/aws"
)

// Provider is a provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				InputDefault: "us-east-1",
			},
			"profile": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"role_arn"},
			},
			"role_arn": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"profile"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bless_ca":         CA(),
			"bless_ecdsa_ca":   ECDSACA(),
			"bless_ed25519_ca": ED25519CA(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bless_lambda":         Lambda(),
			"bless_kms_public_key": KMSPublicKey(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, s *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client, err := aws.NewClient(s)

	if err != nil {
		return client, diag.FromErr(err)
	}

	return client, nil
}
