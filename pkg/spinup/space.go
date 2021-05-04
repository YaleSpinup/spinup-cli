package spinup

// Space holds details about a spinup space
type Space struct {
	Id             *FlexInt `json:"id"`
	Name           string   `json:"name,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	Department     string   `json:"department,omitempty"`
	Contact        string   `json:"contact,omitempty"`
	QuestionaireID string   `json:"questid,omitempty"`
	SecurityGroup  string   `json:"sg,omitempty"`
	Security       string   `json:"security,omitempty"`
	DataTypes      []struct {
		Id   *FlexInt
		Name string
	} `json:"data_types,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	DeletedAt string      `json:"deleted_at,omitempty"`
	Mine      bool        `json:"mine,omitempty"`
	Resources []*Resource `json:"resources,omitempty"`
	Cost      *SpaceCost  `json:"cost,omitempty"`
}

// GetSpace is a space returned from a wonky endpoint
type GetSpace struct {
	Space *Space `json:"space"`
}

// Spaces is a list of spaces
type Spaces struct {
	Spaces []*Space `json:"spaces"`
}

// SoaceCost is the cost estimate for a space
type SpaceCost struct {
	Amount string
	Unit   string
	End    string
	Start  string
}

// GetEndpoint returns the endpoint to get the list of spaces
func (s *Spaces) GetEndpoint(_ map[string]string) string {
	return BaseURL + SpaceURI
}

// GetEndpoint returns the endpoint to get details about a space
func (s *Space) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["id"]
}

// GetEndpoint returns the endpoint to get details about a space
func (s *GetSpace) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["id"]
}

// GetEndpoint returns the endpoint to get cost of a space
func (s *SpaceCost) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["id"] + "/cost"
}
