package helm

import (
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider"
	"github.com/stretchr/testify/require"
)

func TestDownloadNifi(t *testing.T) {
	src := HelmChartSrc{
		HelmChartOpts: provider.HelmChartOpts{
			Chart: "oci://ghcr.io/konpyutaika/helm-charts/nifikop",
		},
	}
	err := src.Download()
	require.NoError(t, err)
}

func TestDownloadCertManager(t *testing.T) {
	src := HelmChartSrc{
		HelmChartOpts: provider.HelmChartOpts{
			Chart: "cert-manager",
			HelmFetchOpts: provider.HelmFetchOpts{
				Repo: "https://charts.jetstack.io",
			},
		},
	}
	err := src.Download()
	require.NoError(t, err)
}
