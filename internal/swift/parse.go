package swift

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v3"
	"github.com/NateMartes/go-swift-tui/pkg/util"
	"github.com/NateMartes/go-swift-tui/pkg/errors"
)

// Loads a clouds.yaml file into a struct for keystone login
func LoadCloudsYAML(path string) (*CloudsYAML, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read clouds.yaml: %w", err)
    }
    var clouds CloudsYAML
    if err := yaml.Unmarshal(data, &clouds); err != nil {
        return nil, fmt.Errorf("failed to parse clouds.yaml: %w", err)
    }
    return &clouds, nil
}

// Using an OpenStackClient clouds.yaml file, get a Cloud struct to login to OpenStack Swift with
func GetCloudFromCloudsFile(filepath string) Cloud {
	
	util.LogDebug(fmt.Sprintf("Marshaling file %s into clouds.yaml structure", filepath))
	cloudsYAML, err := LoadCloudsYAML(filepath)
	if err != nil {
		util.LogFatal(err.Error(), errors.PARSE_ERROR)
	}

	cloudName, err := util.CloudNameVal()
	if util.CloudNameSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if cloudName == "" {
		util.LogFatal(
			"No name specified for cloud, use --cloud-name to specify",
			errors.ARGUMENT_ERROR,
		)
	}
	
	if len(cloudsYAML.Clouds) == 0 {
		util.LogFatal(fmt.Sprintf("No clouds found in file %s", filepath), errors.PARSE_ERROR)
	}
	
	util.LogDebug(fmt.Sprintf("Using cloud '%s' from %s", cloudName, filepath))
	output, ok := cloudsYAML.Clouds[cloudName]
	if !ok {
		util.LogFatal(
			fmt.Sprintf(
				"Cloud '%s' not found in '%s'",
				cloudName,
				filepath,
			),
			errors.ARGUMENT_ERROR,
		)
	}
	return output
}