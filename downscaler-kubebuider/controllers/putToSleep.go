package controllers

import (
	"context"
	"encoding/base64"

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

			err = r.Create(ctx, proxyService)
		}

		if err != nil {
			log.Error(err, "unable to get service in putToSleep")
			return
		}

	}

	// TODO update ingress with proxy service and back up original service data
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

	//TODO make sure zero-scaling/is-sleeping are not called

	err = r.Update(ctx, ingress)

	if err != nil {
		log.Error(err, "unable to update ingress in putToSleep")
		return
	}

	log.Info("putToSleep complete", "ingressName", ingressName, "ingressNamespace", ingressNamespace)

	// TODO scale deployment to zero
}

type ServicePort struct {
	ServiceName string `json:"serviceName"`
	ServicePort int    `json:"servicePort"`
}
