package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func checkDownLoop(r *ScalingBackInfoReconciler) {
	log := r.Log
	log.V(1).Info("List", "Ingresses on watch", len(ingressesCollection))
	//	http://prometheus-server.ingress-nginx.svc.cluster.local:9090/api/v1/query?query=
	//TODO make time customizable for every deployment
	ingressMap, err := getIngressMap("1m", r)

	log.V(1).Info("Got ingress map", "ingressMap", ingressMap)

	for _, ingress := range ingressesCollection {
		namespacedName := ingress.Namespace + "/" + ingress.Name

		proxyWorkingOnIngress := (ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName == "zero-scaling-proxy")

		mapValue, keyExists := ingressMap[namespacedName]
		hasTraffic := mapValue && keyExists
		log.V(1).Info("Got ingress data", "hasTraffic", hasTraffic, "namespacedName", namespacedName)

		if proxyWorkingOnIngress && hasTraffic {
			wakeUp(ingress.Name, ingress.Namespace, r)
		}

		if !proxyWorkingOnIngress && !hasTraffic {

			//check it is not updated recently
			lastWakeup, err := time.Parse(time.RFC3339, ingress.ObjectMeta.Annotations["zero-scaling/last-wakeup"])
			if err != nil {
				log.Error(err, "Got parsing last wakeup error on"+namespacedName+" time = "+ingress.ObjectMeta.Annotations["zero-scaling/last-wakeup"])
			} else {
				secondsPassed := int(time.Since(lastWakeup).Seconds())
				if secondsPassed < 120 {
					log.Info("Skip - no enough time since last wakeup", "namespacedName", namespacedName)
				}
			}

			putToSleep(ingress.Name, ingress.Namespace, r)
		}
	}

	if err != nil {
		log.Error(err, "Got ingress map error")
		return
	}
}

func getIngressMap(timing string, r *ScalingBackInfoReconciler) (map[string]bool, error) {
	log := r.Log
	// log.V(1).Info("Get ingress data started")

	query := fmt.Sprintf("sum by (ingress, namespace) ( rate(nginx_ingress_controller_bytes_sent_sum[%s]) )", timing) //TODO unhardcode
	prometheusURL := "http://prometheus.test.vscodecloud.com/api/v1/query?query=" + url.QueryEscape(query)

	var response Response
	err := getJson(prometheusURL, &response)
	if err != nil {
		log.WithValues("prometheusURL", prometheusURL)
		log.Error(err, "unable to retrueve information from prometheus")
		return map[string]bool{}, err
	}

	// log.V(1).Debug("Got response", "response", response)

	result := make(map[string]bool)

	for _, resultLine := range response.Data.Result {
		ingressName, ok := resultLine.Metric["ingress"]
		if !ok {
			continue
		}

		ingressNamespace, ok := resultLine.Metric["namespace"]
		if !ok {
			continue
		}

		byteRate, err := strconv.ParseFloat(resultLine.Value[1].(string), 64)
		if err != nil {
			return map[string]bool{}, err
		}
		hasTraffic := byteRate > 0

		// log.V(1).Info("Got ingress byterate", "byteRate", byteRate, "ingressName", ingressName, "hasTraffic", hasTraffic)

		result[ingressNamespace+"/"+ingressName] = hasTraffic
	}

	return result, nil
}

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
	Value  []interface{}     `json:"value"`
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
