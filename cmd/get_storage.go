package cmd

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
)

func getStorage(params map[string]string, resource *spinup.Resource) error {
	var j []byte
	var err error
	status := resource.Status

	if status != "created" && status != "creating" && status != "deleting" {
		j, err = ingStatus(resource)
		if err != nil {
			return err
		}
	} else {
		if detailedGetCmd {
			if j, err = storageDetails(params, resource); err != nil {
				return err
			}
		} else {
			if j, err = storage(params, resource); err != nil {
				return err
			}
		}
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write(j)

	return nil
}

func storage(params map[string]string, resource *spinup.Resource) ([]byte, error) {
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

func storageDetails(params map[string]string, resource *spinup.Resource) ([]byte, error) {
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
		params["username"] = u.Username
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
