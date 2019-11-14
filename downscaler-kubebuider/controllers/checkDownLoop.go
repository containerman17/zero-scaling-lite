package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

func checkDownLoop(r *ScalingBackInfoReconciler) {
	log := r.Log
	log.V(1).Info("List", "Ingresses on watch", len(ingressesCollection))
	//	http://prometheus-server.ingress-nginx.svc.cluster.local:9090/api/v1/query?query=
	prometheusURL := "http://prometheus.test.vscodecloud.com/api/v1/query?query=sum%20by%20(host)%20(%0A%20%20rate(nginx_ingress_controller_bytes_sent_sum%5B3m%5D)%0A)"

	var response Response
	err := getJson(prometheusURL, &response)
	if err != nil {
		log.WithValues("prometheusURL", prometheusURL)
		log.Error(err, "unable to retrueve information from prometheus")
		return
	}

	log.V(1).Info("Got response", "response", response)

}

func get

var myClient = &http.Client{Timeout: 10 * time.Second}

type UntypedJson map[string][]interface{}

type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric map[string]string `json:"metric"`
	Value  interface{}       `json:"value"`
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
