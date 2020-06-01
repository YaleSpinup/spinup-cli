package spinup

type Secret struct {
	ARN              string
	Name             string
	Description      string
	KeyId            string
	Type             string
	LastModifiedDate string
	Version          int64
}

type SecretName string
type Secrets []SecretName

// GetEndpoint returns the endpoint to get details about a secret
func (s *Secret) GetEndpoint(params map[string]string) string {
	return BaseURL + SecretsURI + "/" + params["id"] + "/secrets/" + params["secretId"]
}

// GetEndpoint returns the endpoint to get a list of secrets in a space
func (s *Secrets) GetEndpoint(params map[string]string) string {
	return BaseURL + SecretsURI + "/" + params["id"] + "/secrets"
}
