package controllers

import (
	"encoding/base64"

	"github.com/prometheus/common/log"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

type CustomRequestData struct {
	IngressName string
	ServiceName string
	Namespace   string
}

func getIngressByDomain(domain string) *CustomRequestData {
	for _, ingress := range ingressesCollection {
		for ruleIndex := range ingress.Spec.Rules {
			if ingress.Spec.Rules[ruleIndex].Host != domain {
				continue
			}

			restoredIngress := restoreIngress(*ingress)
			log.Debug("restored ingress", "restoredIngress", restoredIngress)

			return &CustomRequestData{
				IngressName: ingress.Name,
				ServiceName: restoredIngress.Spec.Rules[ruleIndex].HTTP.Paths[0].Backend.ServiceName,
				Namespace:   ingress.Namespace,
			}
		}
	}

	return nil
}

func restoreIngress(original extensionsv1beta1.Ingress) *extensionsv1beta1.Ingress {
	ingress := original.DeepCopy()

	specBackup, err := base64.StdEncoding.DecodeString(original.ObjectMeta.Annotations["zero-scaling/backup"])
	if err != nil {
		log.Error(err, "unable to decode backup in getIngressByDomain")
		return nil
	}

	ingress.Spec.Rules = []extensionsv1beta1.IngressRule{}
	err = ingress.Spec.Unmarshal(specBackup)
	log.Info("Restored rules", "ingress.Spec.Rules", ingress.Spec.Rules, "original.Spec.Rules", original.Spec.Rules)
	if err != nil {
		log.Error(err, "unable to restore backup in getIngressByDomain")
		return nil
	}

	return ingress
}
