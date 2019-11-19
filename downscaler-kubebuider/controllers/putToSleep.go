package controllers

import (
	"context"

	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

func putToSleep(ingressName string, ingressNamespace string, r *ScalingBackInfoReconciler) {
	log := r.Log
	log.Info("putToSleep", "ingressName", ingressName, "ingressNamespace", ingressNamespace)

	ctx := context.Background()

	// get ingress
	// TODO check that iongress is updated not less than a minute ago

	namespacedIngressName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      ingressName,
	}
	ingress := &extensionsv1beta1.Ingress{}

	if err := r.Get(ctx, namespacedIngressName, ingress); err != nil {
		log.Error(err, "unable to get Ingress in putToSleep")
		return
	}

	// TODO create proxy service

	namespacedProxyServiceName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      "zero-scaling-proxy",
	}
	proxyService := &apiv1.Service{}

	if err := r.Get(ctx, namespacedProxyServiceName, proxyService); err != nil {
		if err.Error() == "Service \"zero-scaling-proxy\" not found" {
			//TODO create service

			proxyService = &apiv1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "zero-scaling-proxy",
					Namespace: ingressNamespace,
				},
				Spec: apiv1.ServiceSpec{
					Type:         "ExternalName",
					ExternalName: "google.com",
					Ports: []apiv1.ServicePort{
						{
							Port: 80,
						},
					},
				},
			}

			r.Create(ctx, proxyService)
		}

		if err != nil {
			log.Error(err, "unable to get service in putToSleep")
			return
		}

	}

	// TODO update ingress with proxy service and back up original service data

	for index1, _ := range ingress.Spec.Rules {
		for index2, _ := range ingress.Spec.Rules[index1].HTTP.Paths {
			ingress.Spec.Rules[index1].HTTP.Paths[index2].Backend.ServiceName = "zero-scaling-proxy"
			ingress.Spec.Rules[index1].HTTP.Paths[index2].Backend.ServicePort = apiv1.ServicePort{
				Port: 80,
			}

		}
	}

	// TODO scale deployment to zero
}
