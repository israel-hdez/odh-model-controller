/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"fmt"

	servingv1beta1 "github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/opendatahub-io/odh-model-controller/internal/webhook/serving"
)

// nolint:unused
// log is for logging in this package.
var inferenceservicelog = logf.Log.WithName("inferenceservice-resource")

// SetupInferenceServiceWebhookWithManager registers the webhook for InferenceService in the manager.
func SetupInferenceServiceWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&servingv1beta1.InferenceService{}).
		WithDefaulter(&InferenceServiceCustomDefaulter{client: mgr.GetClient()}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-serving-kserve-io-v1beta1-inferenceservice,mutating=true,failurePolicy=fail,sideEffects=None,groups=serving.kserve.io,resources=inferenceservices,verbs=create;update,versions=v1beta1,name=minferenceservice-v1beta1.kb.io,admissionReviewVersions=v1

// InferenceServiceCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind InferenceService when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type InferenceServiceCustomDefaulter struct {
	client client.Client
}

var _ webhook.CustomDefaulter = &InferenceServiceCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind InferenceService.
func (d *InferenceServiceCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	inferenceservice, ok := obj.(*servingv1beta1.InferenceService)

	if !ok {
		return fmt.Errorf("expected an InferenceService object but got %T", obj)
	}
	logger := inferenceservicelog.WithValues("name", inferenceservice.GetName())
	logger.Info("Defaulting for InferenceService", "name", inferenceservice.GetName())

	err := serving.ApplyDefaultServerlessAnnotations(ctx, d.client, inferenceservice.GetName(), inferenceservice.GetAnnotations(), logger)
	if err != nil {
		return err
	}

	return nil
}
