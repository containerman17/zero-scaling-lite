package main

import (
	"log"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func checkRegistryPassword(api clientCoreV1.CoreV1Interface, ns *apiv1.Namespace) error {
	//start func checkSecrets
	secrets, err := api.Secrets(ns.Name).Get("registry-password", metav1.GetOptions{})
	if err != nil {
		switch err.(type) {
		case *errors.StatusError:
			statusError, _ := err.(*errors.StatusError) //always ok since we already checked
			if statusError.ErrStatus.Code != 404 {
				return err
			}
			log.Println("registry-password not found, creating")

			login := ns.Name
			password := genPassword()

			newSecret, err := generateRegistryPasswordSecret(ns.Name, login, password)
			if err != nil {
				return err
			}

			_, err = api.Secrets(ns.Name).Create(&newSecret)
			if err != nil {
				return err
			}
			secretCreatedCounter.Inc() //Prometheus
			log.Println("registry-password created", starPassword(string(secrets.Data["password"]), 3))
		default:
			return err
		}
	} else {
		log.Println("registry-password exists", starPassword(string(secrets.Data["password"]), 3))
	}
	return nil
}
