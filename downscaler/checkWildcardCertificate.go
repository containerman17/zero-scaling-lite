package main

import (
	"bytes"
	"log"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var (
	originalSecret *apiv1.Secret
)

func getFreeDomainCertSecret(api clientCoreV1.CoreV1Interface) (*apiv1.Secret, error) {
	if originalSecret != nil {
		return originalSecret, nil
	}
	log.Print("Getting secret git2prod-system/free-domain-certificate")
	secret, err := api.Secrets("git2prod-system").Get("free-domain-certificate", metav1.GetOptions{})

	if err == nil {
		originalSecret = secret
	}

	return secret, err
}

func checkWildcardCertificate(api clientCoreV1.CoreV1Interface, ns *apiv1.Namespace) error {
	localSecret, err := api.Secrets(ns.Name).Get("free-domain-certificate", metav1.GetOptions{})

	originalSecret, originalSecretErr := getFreeDomainCertSecret(api)
	if originalSecretErr != nil {
		return originalSecretErr
	}
	//TODO проверить что контент у секретов одинаковый
	if err != nil {
		switch err.(type) {
		case *errors.StatusError:
			statusError, _ := err.(*errors.StatusError) //always ok since we already checked
			if statusError.ErrStatus.Code != 404 {
				return err
			}
			log.Printf("Secret %v/free-domain-certificate not found, creating", ns.Name)

			newSecret := originalSecret.DeepCopy()
			newSecret.SetNamespace(ns.Name)
			newSecret.SetResourceVersion("")

			_, err = api.Secrets(ns.Name).Create(newSecret)
			if err != nil {
				return err
			}

			log.Printf("Secret %v/free-domain-certificate created", ns.Name)
		default:
			return err
		}
	} else {
		// Update old certificate
		if !bytes.Equal(localSecret.Data["tls.crt"], originalSecret.Data["tls.crt"]) {
			log.Printf("Secret %v/free-domain-certificate is outdated!", ns.Name)
			localSecret.Data = originalSecret.Data
			_, err := api.Secrets(ns.Name).Update(localSecret)
			if err != nil {
				return err
			}
			log.Printf("Secret %v/free-domain-certificate is updated", ns.Name)
		} else {
			log.Printf("Secret %v/free-domain-certificate is up to date", ns.Name)
		}
	}
	return nil
}
