package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/YaleSpinup/spinup/pkg/spinup"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getStorageCmd)
}

var getStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Get details about a storage resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 storage service id is required")
		}

		var j []byte
		if detailedGetCmd {
			var err error
			if j, err = storageDetails(args[0]); err != nil {
				return err
			}
		} else {
			var err error
			if j, err = storage(args[0]); err != nil {
				return err
			}
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func storage(id string) ([]byte, error) {
	resource := &spinup.Resource{}
	err := SpinupClient.GetResource(id, resource)
	if err != nil {
		return []byte{}, err
	}

	status := resource.Status
	if status != "created" && status != "creating" && status != "deleting" {
		return json.MarshalIndent(struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			SpaceID string `json:"space_id"`
		}{
			ID:      resource.ID.String(),
			Name:    resource.Name,
			Status:  resource.Status,
			SpaceID: resource.SpaceID.String(),
		}, "", "  ")
	}

	size, err := SpinupClient.S3StorageSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.S3StorageInfo{}
	err = SpinupClient.GetResource(id, resource)
	if err != nil {
		return []byte{}, err
	}

	state := "populated"
	if info.Empty {
		state = "empty"
	}

	return resourceSummary(resource, size, state)
}

func storageDetails(id string) ([]byte, error) {
	resource := &spinup.Resource{}
	err := SpinupClient.GetResource(id, resource)
	if err != nil {
		return []byte{}, err
	}

	status := resource.Status
	if status != "created" && status != "creating" && status != "deleting" {
		return json.MarshalIndent(struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			SpaceID string `json:"space_id"`
		}{
			ID:      resource.ID.String(),
			Name:    resource.Name,
			Status:  resource.Status,
			SpaceID: resource.SpaceID.String(),
		}, "", "  ")
	}

	size, err := SpinupClient.S3StorageSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.S3StorageInfo{}
	err = SpinupClient.GetResource(id, resource)
	if err != nil {
		return []byte{}, err
	}

	users := spinup.S3StorageUsers{}
	err = SpinupClient.GetResource(id, &users)
	if err != nil {
		return []byte{}, err
	}

	beta := false
	if b, err := strconv.Atoi(resource.Type.Beta); err != nil && b != 0 {
		beta = true
	}

	tryit := false
	if size.Price == "tryit" {
		tryit = true
	}

	type User struct {
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
		LastUsed  string `json:"last_used"`
		UserId    string `json:"user_id"`
	}

	userList := []*User{}
	for _, u := range users {
		userList = append(userList, &User{
			UserId:    u.UserId,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			LastUsed:  u.LastUsed,
		})
	}

	output := struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Status   string  `json:"status"`
		Type     string  `json:"type"`
		Flavor   string  `json:"flavor"`
		Security string  `json:"security"`
		SpaceID  string  `json:"space_id"`
		Beta     bool    `json:"beta"`
		Size     string  `json:"size"`
		State    string  `json:"state"`
		Empty    bool    `json:"empty"`
		TryIT    bool    `json:"tryit"`
		Users    []*User `json:"users"`
	}{
		ID:       resource.ID.String(),
		Name:     resource.Name,
		Status:   resource.Status,
		Type:     resource.Type.Name,
		Flavor:   resource.Type.Flavor,
		Security: resource.Type.Security,
		Size:     size.Name,
		SpaceID:  resource.SpaceID.String(),
		Beta:     beta,
		TryIT:    tryit,
		// State:    info.Status,
		Empty: info.Empty,
		Users: userList,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
