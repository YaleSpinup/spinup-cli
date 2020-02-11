package spinup

// Offering is the Spinup representation of an offering (or type)
type Offering struct {
	Beta              string   `json:"beta,omitempty"`
	CreatedAt         string   `json:"created_at,omitempty"`
	DefaultDiskSize   string   `json:"default_disk_size,omitempty"`
	DeletedAt         string   `json:"deleted_at,omitempty"`
	Details           string   `json:"details,omitempty"`
	Flavor            string   `json:"flavor"`
	ID                *FlexInt `json:"id"`
	Image             string   `json:"image,omitempty"`
	ImageOwner        string   `json:"image_owner,omitempty"`
	ImageSearchString string   `json:"search_string,omitempty"`
	Logo              string   `json:"logo,omitempty"`
	MinDiskSize       string   `json:"min_disk_size,omitempty"`
	Name              string   `json:"name"`
	NameBroker        string   `json:"name_broker,omitempty"`
	Options           string   `json:"options,omitempty"`
	ProductCode       string   `json:"product_code,omitempty"`
	Provider          string   `json:"provider,omitempty"`
	ProviderBroker    string   `json:"provider_name,omitempty"`
	Security          string   `json:"security,omitempty"`
	SecurityGroup     string   `json:"sgs,omitempty"`
	Subnet            string   `json:"subnet,omitempty"`
	Type              string   `json:"type"`
	UpdatedAt         string   `json:"updated_at,omitempty"`
	UserData          string   `json:"userdata,omitempty"`
}
