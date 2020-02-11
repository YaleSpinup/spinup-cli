package cmd

import (
	"encoding/json"
	"strconv"

	"github.com/YaleSpinup/spinup/pkg/spinup"
)

func resourceSummary(resource *spinup.Resource, size spinup.Size, state string) ([]byte, error) {
	beta := false
	if b, err := strconv.Atoi(resource.Type.Beta); err != nil && b != 0 {
		beta = true
	}

	tryit := false
	if size.GetPrice() == "tryit" {
		tryit = true
	}

	output := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Status   string `json:"status"`
		Type     string `json:"type"`
		Flavor   string `json:"flavor"`
		Security string `json:"security"`
		SpaceID  string `json:"space_id"`
		Beta     bool   `json:"beta"`
		Size     string `json:"size"`
		TryIT    bool   `json:"tryit"`
		State    string `json:"state,omitempty"`
	}{
		ID:       resource.ID.String(),
		Name:     resource.Name,
		Status:   resource.Status,
		Type:     resource.Type.Name,
		Flavor:   resource.Type.Flavor,
		Security: resource.Type.Security,
		Size:     size.GetName(),
		SpaceID:  resource.SpaceID.String(),
		Beta:     beta,
		TryIT:    tryit,
		State:    state,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
