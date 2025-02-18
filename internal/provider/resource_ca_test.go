package provider_test

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	a := assert.New(t)
	providers, kmsMock := getTestProviders()

	ciphertext := make([]byte, 10)
	_, err := rand.Read(ciphertext)
	a.NoError(err)
	output := &kms.EncryptOutput{
		CiphertextBlob: ciphertext,
	}
	kmsMock.On("Encrypt", mock.Anything).Return(output, nil)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providers,
		Steps: []resource.TestStep{
			{
				Config: `
				provider "bless" {
					region = "us-east-1"
				}

				resource "bless_ca" "bless" {
					kms_key_id = "testo"
				}

				output "private_key" {
					value = "${bless_ca.bless.encrypted_ca}"
				}
				output "public_key" {
					value = "${bless_ca.bless.public_key}"
				}
				output "password" {
					value = "${bless_ca.bless.encrypted_password}"
				}
			`,
				Check: func(s *terraform.State) error {
					privateUntyped := s.RootModule().Outputs["private_key"].Value
					private, ok := privateUntyped.(string)
					a.True(ok)
					bytesPrivate, err := base64.StdEncoding.DecodeString(private)
					a.Nil(err)
					a.Regexp(
						regexp.MustCompile("^-----BEGIN RSA PRIVATE KEY-----"),
						string(bytesPrivate))
					a.Regexp(
						regexp.MustCompile(`AES-256-CBC`),
						string(bytesPrivate))

					publicSSHUntyped := s.RootModule().Outputs["public_key"].Value
					publicSSH, ok := publicSSHUntyped.(string)
					a.True(ok)
					a.Regexp(
						regexp.MustCompile("^ssh-rsa "),
						publicSSH)
					return nil
				},
			},
		},
	})
}
