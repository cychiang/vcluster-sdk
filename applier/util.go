package applier

import (
	"context"
	"fmt"
	"os"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func ApplyManifestFile(inClusterConfig *rest.Config, filename string) error {
	manifest, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("function ApplyManifestFile failed, unable to read %s file: %v", filename, err)
	}

	return ApplyManifest(inClusterConfig, &manifest)
}

func ApplyManifest(inClusterConfig *rest.Config, manifest *[]byte) error {
	httpClient, err := rest.HTTPClientFor(inClusterConfig)
	if err != nil {
		return fmt.Errorf("unable to initialize httpClient: %w", err)
	}
	restMapper, err := apiutil.NewDynamicRESTMapper(inClusterConfig, httpClient)
	if err != nil {
		return fmt.Errorf("unable to initialize NewDynamicRESTMapper")
	}

	a := DirectApplier{}
	opts := ApplierOptions{
		Manifest:   string(*manifest),
		RESTConfig: inClusterConfig,
		RESTMapper: restMapper,
	}
	return a.Apply(context.Background(), opts)
}
