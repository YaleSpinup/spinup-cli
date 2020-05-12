package spinup

// Image is a server image
type Image struct {
	Architecture string `json:"architecture,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	CreatedBy    string `json:"created_by,omitempty"`
	Description  string `json:"description,omitempty"`
	ID           string `json:"id"`
	Name         string `json:"name,omitempty"`
	ServerName   string `json:"server_name,omitempty"`
	State        string `json:"state,omitempty"`
	Status       string `json:"status,omitempty"`
	Volumes      map[string]struct {
		DeleteOnTermination bool   `json:"delete_on_termination"`
		Encrypted           bool   `json:"encrypted"`
		ID                  string `json:"snapshot_id"`
		Size                int    `json:"volume_size"`
		Type                string `json:"volume_type,omitempty"`
	} `json:"volumes,omitempty"`
	Offering *Offering `json:"offering,omitempty"`
}

// Images is a list of server images
type Images []*Image

// GetEndpoint gets the endpoint UR for an image list
func (i *Images) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["id"] + "/images"
}
