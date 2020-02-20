package spinup

// Resource is a specific resource in the database, it represents an actual instance, container, s3 bucket, etc
type Resource struct {
	Admin     string    `json:"admin,omitempty"`
	CreatedAt string    `json:"created_at"`
	DeletedAt string    `json:"deleted_at,omitempty"`
	Flavor    string    `json:"flavor"`
	ID        *FlexInt  `json:"id"`
	IP        string    `json:"ip,omitempty"`
	Is        string    `json:"is_a"`
	Name      string    `json:"name"`
	ServerID  string    `json:"server_id,omitempty"`
	SizeID    *FlexInt  `json:"size_id"`
	SpaceID   *FlexInt  `json:"space_id"`
	Status    string    `json:"status"`
	TypeID    *FlexInt  `json:"type_id"`
	Task      string    `json:"task,omitempty"`
	Type      *Offering `json:"type,omitempty"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

// GetEndpoint returns the URL to get a resource
func (r *Resource) GetEndpoint(id string) string {
	return BaseURL + ResourceURI + "/" + id
}
