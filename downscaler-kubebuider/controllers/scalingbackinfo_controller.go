/*

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScalingBackInfoReconciler reconciles a ScaleBackInfo object
type ScalingBackInfoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	ingressesCollection [0]string
)

// +kubebuilder:rbac:groups=zero-scaling.controllers.dockerize.io,resources=scalebackinfoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zero-scaling.controllers.dockerize.io,resources=scalebackinfoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extensions;networking.k8s.io,resources=ingresses,verbs=get;list;watch

func (r *ScalingBackInfoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("Ingress", req.NamespacedName)

	// get ingress

	ingress := &extensionsv1beta1.Ingress{}

	if err := r.Get(ctx, req.NamespacedName, ingress); err != nil {
		log.Error(err, "unable to get Ingress ")
		return ctrl.Result{}, err
	}

	log.V(1).Info("Got ingress")

	return ctrl.Result{}, nil
}

func (r *ScalingBackInfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionsv1beta1.Ingress{}).
		Complete(r)
}
