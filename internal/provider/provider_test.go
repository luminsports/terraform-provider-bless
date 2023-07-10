package provider_test

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/luminsports/terraform-provider-bless/internal/aws"
	"github.com/luminsports/terraform-provider-bless/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type KMSMock struct {
	kmsiface.KMSAPI
	mock.Mock
}

func (k *KMSMock) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	args := k.Called(input)
	output := args.Get(0).(*kms.EncryptOutput)
	return output, args.Error(1)
}

func (k *KMSMock) GetPublicKey(input *kms.GetPublicKeyInput) (*kms.GetPublicKeyOutput, error) {
	args := k.Called(input)
	output := args.Get(0).(*kms.GetPublicKeyOutput)
	return output, args.Error(1)
}

func getTestProviders() (map[string]func() (*schema.Provider, error), *KMSMock) {
	kmsMock := &KMSMock{}
	providers := map[string]func() (*schema.Provider, error){
		"bless": func() (*schema.Provider, error) {
			ca := provider.Provider()
			ca.ConfigureContextFunc = func(_ context.Context, s *schema.ResourceData) (interface{}, diag.Diagnostics) {
				client := &aws.Client{
					KMS: aws.KMS{Svc: kmsMock},
				}
				return client, nil
			}
			return ca, nil
		},
	}

	return providers, kmsMock
}

func TestProvider(t *testing.T) {
	assert := assert.New(t)
	p := provider.Provider()
	err := p.InternalValidate()
	assert.Nil(err)
}
