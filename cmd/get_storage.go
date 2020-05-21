package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
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
	params := map[string]string{"id": id}
	resource := &spinup.Resource{}
	if err := SpinupClient.GetResource(params, resource); err != nil {
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
	if err := SpinupClient.GetResource(params, resource); err != nil {
		return []byte{}, err
	}

	state := "populated"
	if info.Empty {
		state = "empty"
	}

	return json.MarshalIndent(newResourceSummary(resource, size, state), "", "  ")
}

func storageDetails(id string) ([]byte, error) {
	params := map[string]string{"id": id}
	resource := &spinup.Resource{}
	if err := SpinupClient.GetResource(params, resource); err != nil {
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
	if err := SpinupClient.GetResource(params, resource); err != nil {
		return []byte{}, err
	}

	users := spinup.S3StorageUsers{}
	if err := SpinupClient.GetResource(params, &users); err != nil {
		return []byte{}, err
	}

	state := "populated"
	if info.Empty {
		state = "empty"
	}

	type User struct {
		Username  string   `json:"username"`
		CreatedAt string   `json:"created_at"`
		LastUsed  string   `json:"last_used"`
		Keys      []string `json:"key_id"`
	}

	userList := []*User{}
	for _, u := range users {
		params["name"] = u.Username
		user := spinup.S3StorageUser{}
		if err = SpinupClient.GetResource(params, &user); err != nil {
			return []byte{}, err
		}

		keys := make([]string, 0, len(user.AccessKeys))
		for _, k := range user.AccessKeys {
			keys = append(keys, k.AccessKeyId)
		}

		userList = append(userList, &User{
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			LastUsed:  u.LastUsed,
			Keys:      keys,
		})
	}

	output := struct {
		*ResourceSummary
		Empty bool    `json:"empty"`
		Users []*User `json:"users"`
	}{
		newResourceSummary(resource, size, state),
		info.Empty,
		userList,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
