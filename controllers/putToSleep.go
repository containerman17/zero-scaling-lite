package controllers

import (
	"context"
	"encoding/base64"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

func putToSleep(ingressName string, ingressNamespace string, r *ScalingBackInfoReconciler) {
	log := r.Log
	log.Info("putToSleep", "ingressName", ingressName, "ingressNamespace", ingressNamespace)

	ctx := context.Background()

	// get ingress

	namespacedIngressName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      ingressName,
	}
	ingress := &extensionsv1beta1.Ingress{}

	if err := r.Get(ctx, namespacedIngressName, ingress); err != nil {
		log.Error(err, "unable to get Ingress in putToSleep")
		return
	}

	// log.Info("debug ingress", "ingress", ingress.ObjectMeta.)

	//  create proxy service

	namespacedProxyServiceName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      "zero-scaling-proxy",
	}
	proxyService := &apiv1.Service{}

	if err := r.Get(ctx, namespacedProxyServiceName, proxyService); err != nil {
		if err.Error() == "Service \"zero-scaling-proxy\" not found" {
			// create service

			proxyService = &apiv1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "zero-scaling-proxy",
					Namespace: ingressNamespace,
				},
				Spec: apiv1.ServiceSpec{
					Type:         "ExternalName",
					ExternalName: "downscaler-kubebuider-controller-manager-metrics-service.downscaler-kubebuider-system.svc.cluster.local",
					Ports: []apiv1.ServicePort{
						{
							Port: 80,
						},
					},
				},
			}

			err = r.Create(ctx, proxyService)
		}

		if err != nil {
			log.Error(err, "unable to get service in putToSleep")
			return
		}

	}

	//  update ingress with proxy service and back up original service data
	portsBackup, err := ingress.Spec.Marshal()
	if err != nil {
		log.Error(err, "unable to marshal spec")
		return
	}

	for ruleIndex := range ingress.Spec.Rules {

		for pathIndex := range ingress.Spec.Rules[ruleIndex].HTTP.Paths {
			// //backup
			// portsBackup[strconv.Itoa(ruleIndex)+"_"+strconv.Itoa(pathIndex)] = ServicePort{
			// 	ServiceName: ingress.Spec.Rules[ruleIndex].HTTP.Paths[pathIndex].Backend.ServiceName,
			// 	ServicePort: ingress.Spec.Rules[ruleIndex].HTTP.Paths[pathIndex].Backend.ServicePort.IntValue(),
			// }
			//set proxy service
			ingress.Spec.Rules[ruleIndex].HTTP.Paths[pathIndex].Backend.ServiceName = "zero-scaling-proxy"
			ingress.Spec.Rules[ruleIndex].HTTP.Paths[pathIndex].Backend.ServicePort = intstr.FromInt(80)
		}
	}

	ingress.ObjectMeta.Annotations["zero-scaling/backup"] = base64.StdEncoding.EncodeToString(portsBackup)

	err = r.Update(ctx, ingress)

	if err != nil {
		log.Error(err, "unable to update ingress in putToSleep")
		return
	}

	log.Info("putToSleep complete", "ingressName", ingressName, "ingressNamespace", ingressNamespace)

	//  scale deployment to zero

	namespacedDeploymentName := client.ObjectKey{
		Namespace: ingressNamespace,
		Name:      ingress.ObjectMeta.Annotations["zero-scaling/deploymentName"],
	}
	deployment := &appsv1.Deployment{}

	if err := r.Get(ctx, namespacedDeploymentName, deployment); err != nil {
		log.Error(err, "unable to get Deployment "+namespacedDeploymentName.String()+" in putToSleep")
		return
	}

	zero := int32(0)
	deployment.Spec.Replicas = &zero

	err = r.Update(ctx, deployment)

	if err != nil {
		log.Error(err, "unable to update deployment in putToSleep")
		return
	}
}
