package swift

// type representing a clouds.yaml file
type CloudsYAML struct {
	Clouds map[string]Cloud `yaml:"clouds"`
}

// type representing a cloud in a clouds.yaml file
type Cloud struct {
	Auth               CloudAuth `yaml:"auth"`
	RegionName         string    `yaml:"region_name,omitempty"`
	Interface          string    `yaml:"interface,omitempty"`
	IdentityAPIVersion int       `yaml:"identity_api_version,omitempty"`
	VolumeAPIVersion   int       `yaml:"volume_api_version,omitempty"`
}

// type representing an authentication method in a cloud
type CloudAuth struct {
	AuthURL                     string `yaml:"auth_url"`
	Username                    string `yaml:"username,omitempty"`
	Password                    string `yaml:"password,omitempty"`
	ProjectName                 string `yaml:"project_name,omitempty"`
	ProjectID                   string `yaml:"project_id,omitempty"`
	UserDomainName              string `yaml:"user_domain_name,omitempty"`
	UserDomainID                string `yaml:"user_domain_id,omitempty"`
	ProjectDomainName           string `yaml:"project_domain_name,omitempty"`
	ProjectDomainID             string `yaml:"project_domain_id,omitempty"`
	TenantName                  string `yaml:"tenant_name,omitempty"`
	TenantID                    string `yaml:"tenant_id,omitempty"`
	ApplicationCredentialID     string `yaml:"application_credential_id,omitempty"`
	ApplicationCredentialSecret string `yaml:"application_credential_secret,omitempty"`
}

type AuthMethod int

const (
	AuthTempauth AuthMethod = iota
	AuthKeystoneV2
	AuthKeystoneV3Password
)
