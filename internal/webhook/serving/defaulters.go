package serving

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/opendatahub-io/odh-model-controller/internal/controller/utils"
)

func ApplyDefaultServerlessAnnotations(ctx context.Context, client client.Client, resourceName string, resourceAnnotations map[string]string, logger logr.Logger) error {
	deploymentMode, err := utils.GetDeploymentModeForKServeResource(ctx, client, resourceAnnotations)
	if err != nil {
		return fmt.Errorf("error resolving deployment mode for resource %s: %w", resourceName, err)
	}

	if deploymentMode == utils.Serverless {
		logAnnotationsAdded := make([]string, 0, 3)

		if _, exists := resourceAnnotations["serving.knative.openshift.io/enablePassthrough"]; !exists {
			resourceAnnotations["serving.knative.openshift.io/enablePassthrough"] = "true"
			logAnnotationsAdded = append(logAnnotationsAdded, "serving.knative.openshift.io/enablePassthrough")
		}

		if _, exists := resourceAnnotations["sidecar.istio.io/inject"]; !exists {
			resourceAnnotations["sidecar.istio.io/inject"] = "true"
			logAnnotationsAdded = append(logAnnotationsAdded, "sidecar.istio.io/inject")
		}

		if _, exists := resourceAnnotations["sidecar.istio.io/rewriteAppHTTPProbers"]; !exists {
			resourceAnnotations["sidecar.istio.io/rewriteAppHTTPProbers"] = "true"
			logAnnotationsAdded = append(logAnnotationsAdded, "sidecar.istio.io/rewriteAppHTTPProbers")
		}

		if len(logAnnotationsAdded) > 0 {
			logger.V(1).Info("Annotations added", "annotations", strings.Join(logAnnotationsAdded, ","))
		}
	}
	return nil
}
