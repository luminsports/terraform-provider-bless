package provider_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestKMSPublicKey(t *testing.T) {
	r := require.New(t)
	providers, kmsMock := getTestProviders()

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	r.NoError(err)

	derBytes, err := x509.MarshalPKIXPublicKey(priv.Public())
	r.NoError(err)
	output := &kms.GetPublicKeyOutput{
		PublicKey: derBytes,
		KeyId:     aws.String("key id"),
	}

	kmsMock.On("GetPublicKey", mock.Anything).Return(output, nil)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providers,
		Steps: []resource.TestStep{
			{
				Config: `
              provider "bless" {
                region = "us-east-1"
              }

              data "bless_kms_public_key" "bless" {
                kms_key_id = "testo"
              }

              output "public_key" {
                value = "${data.bless_kms_public_key.bless.public_key}"
              }
            `,
				Check: func(s *terraform.State) error {
					publicSSHUntyped := s.RootModule().Outputs["public_key"].Value
					publicSSH, ok := publicSSHUntyped.(string)
					r.True(ok)
					r.Regexp(
						regexp.MustCompile("^ecdsa-sha2-nistp384 "),
						publicSSH)
					return nil
				},
			},
		},
	})
}
