package main

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "go.uber.org/automaxprocs"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/rest"
)

func main() {
	go startMetricsServer()

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := rest.InClusterConfig()

	if _, fileErr := os.Stat(kubeconfig); fileErr == nil {
		//using ~/.kube/config
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		log.Println("Using kubeconfig file", kubeconfig)
	} else if os.IsNotExist(fileErr) {
		//no ~/.kube/config, using service account
		log.Println("Using service account (default)")
	} else {
		log.Fatal("Accesing file error", fileErr)
	}

	//config
	if err != nil {
		log.Fatal("creating config error: ", err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	api := clientset.ExtensionsV1beta1()

	log.Println("--- watch updates ---")

	watcher, err := api.Ingresses("").Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal("api.Ingresses().Watch error: ", err)
	}
	ch := watcher.ResultChan()
	for event := range ch {
		ingress, ok := event.Object.(*extensionsv1beta1.Ingress)
		if !ok {
			errorsCounter.Inc() //Prometheus
			log.Println("Conversion error, skip...", event)
			continue
		}

		if event.Type == watch.Added {
			log.Printf("Added %s/%s", ingress.Namespace, ingress.Name)
		} else if event.Type == watch.Modified {
			log.Printf("Modified %s/%s", ingress.Namespace, ingress.Name)
		} else if event.Type == watch.Deleted {
			log.Printf("Deleted %s/%s", ingress.Namespace, ingress.Name)
		} else if event.Type == watch.Error {
			log.Println("Error: event.type == watch.Error", event)
		} else {
			log.Println("Error: unexpected type", event)
		}
	}
}
