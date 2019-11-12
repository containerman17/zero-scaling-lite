package main

import (
	"encoding/base64"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateRegistryPasswordSecret(nsName string, login string, password string) (v1.Secret, error) {

	dockerJSONStr, err := formDockerConfigJSON(getenv("REGISTRY", "registry-gate.registry.svc.cluster.local"), login, password)
	if err != nil {
		return v1.Secret{}, err
	}

	return v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "registry-password",
			Namespace: nsName,
		},
		Data: map[string][]byte{
			"login":             []byte(login),
			"password":          []byte(password),
			".dockerconfigjson": []byte(dockerJSONStr),
		},
		Type: "kubernetes.io/dockerconfigjson",
	}, nil
}

func formDockerConfigJSON(server string, login string, password string) (string, error) {
	type Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Auth     string `json:"auth"`
	}
	type DockerAuth struct {
		Auths map[string]Auth `json:"auths"`
	}

	logLassTogether := login + ":" + password

	auth := &Auth{
		Username: login,
		Password: password,
		Auth:     base64.StdEncoding.EncodeToString([]byte(logLassTogether)),
	}

	dockerAuth := &DockerAuth{
		Auths: map[string]Auth{
			server: *auth,
		},
	}

	res1B, err := json.Marshal(dockerAuth)
	return string(res1B), err
}
