package mocks

import (
	"encoding/json"

	pkgerrors "github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sdataplatform/pulumi-helper/mocks/provider"
)

type Mocks int

func (Mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (Mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	if args.Token == "kubernetes:helm:template" {
		k8sProvider, err := provider.MakeKubeProvider(nil, "test", "v1.28", []byte{}, []byte{})
		if err != nil {
			return nil, err
		}

		var jsonOpts string
		if jsonOptsArgs := args.Args["jsonOpts"]; jsonOptsArgs.HasValue() && jsonOptsArgs.IsString() {
			jsonOpts = jsonOptsArgs.StringValue()
		} else {
			return nil, pkgerrors.New("missing required field 'jsonOpts' of type string")
		}

		var opts provider.HelmChartOpts
		err = json.Unmarshal([]byte(jsonOpts), &opts)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to unmarshal 'jsonOpts'")
		}

		text, err := k8sProvider.HelmTemplate(opts)
		if err != nil {
			// TODO: in unit tests this error is not seen / caught
			return nil, pkgerrors.Wrap(err, "failed to generate YAML for specified Helm chart")
		}

		// Decode the generated YAML here to avoid an extra invoke in the client.
		result, err := k8sProvider.DecodeYaml(text, opts.Namespace)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to decode YAML for specified Helm chart")
		}

		return resource.NewPropertyMapFromMap(map[string]interface{}{"result": result}), nil
	}

	return args.Args, nil
}

func WithMocks(project, stack string, mocks pulumi.MockResourceMonitor) pulumi.RunOption {
	return func(info *pulumi.RunInfo) {
		info.Project, info.Stack, info.Mocks = project, stack, mocks
	}
}

func WithMocksAndConfig(project, stack string, config map[string]string, mocks pulumi.MockResourceMonitor) pulumi.RunOption {
	return func(info *pulumi.RunInfo) {
		info.Project, info.Stack, info.Mocks, info.Config = project, stack, mocks, config
	}
}
