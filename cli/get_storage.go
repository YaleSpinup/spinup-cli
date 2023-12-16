package cli

import (
	"encoding/json"
	"fmt"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getStorageCmd)
}

var getStorageCmd = &cobra.Command{
	Use:     "storage [space]/[resource]",
	Short:   "Get a storage service",
	PreRunE: getCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("get storage: %+v", args)

		status := getResource.Status
		if status != "created" && status != "creating" && status != "deleting" {
			return ingStatus(getResource)
		}

		var out []byte
		var err error
		switch {
		case detailedGetCmd:
			switch getResource.Type.Flavor {
			case "s3", "s3bucket":
				out, err = s3StorageDetails(getParams, getResource)
				if err != nil {
					return err
				}
			case "efs":
				log.Warn("efs is not supported yet")
				return nil
			default:
				return fmt.Errorf("unknown flavor: %s", getResource.Type.Flavor)
			}
		default:
			switch getResource.Type.Flavor {
			case "s3", "s3bucket":
				out, err = s3Storage(getParams, getResource)
				if err != nil {
					return err
				}
			case "efs":
				log.Warn("efs is not supported yet")
				return nil
			default:
				return fmt.Errorf("unknown flavor: %s", getResource.Type.Flavor)
			}
		}

		return formatOutput(out)
	},
}

func s3Storage(params map[string]string, resource *spinup.Resource) ([]byte, error) {
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

func s3StorageDetails(params map[string]string, resource *spinup.Resource) ([]byte, error) {
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
