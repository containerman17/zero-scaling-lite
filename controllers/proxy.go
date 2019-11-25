package controllers

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

var wakeUpLocks = make(map[string]string)
var wakeUpLocksMutex = sync.Mutex{}

func startProxy(r *ScalingBackInfoReconciler) {
	log := r.Log

	director := func(req *http.Request) {
		customRequestData := getIngressByDomain(strings.Split(req.Host, ":")[0])
		log.Info("Got request data", "CustomRequestData", customRequestData)
		// CustomRequestData{
		// 	IngressName: "echo-ingress-1",
		// 	ServiceName: "echo1",
		// 	Namespace:   "default",
		// }

		req.Header.Add("X-Forwarded-Host", req.Host)
		// req.Header.Add("X-Origin-Host", origin.Host)
		log.Info("Request", "host", req.Host)

		log.Info("Wait for lock started", "host", req.Host)
		for {
			wakeUpLocksMutex.Lock()
			if wakeUpLocks[req.Host] != "locked" {
				wakeUpLocksMutex.Unlock()
				break
			}
			wakeUpLocksMutex.Unlock()
			log.Info("Locked...", "host", req.Host)
			time.Sleep(1 * time.Second)
		}
		log.Info("Wait for lock finished", "host", req.Host)

		//lock to prevent multiple requests
		wakeUpLocks[req.Host] = "locked"

		wakeUp(customRequestData.IngressName, customRequestData.Namespace, r)
		waitForWakeUp(customRequestData.ServiceName, customRequestData.Namespace)

		//unlock
		wakeUpLocksMutex.Lock()
		delete(wakeUpLocks, req.Host)
		wakeUpLocksMutex.Unlock()

		req.URL.Scheme = "http"
		req.URL.Host = customRequestData.ServiceName + "." + customRequestData.Namespace + ".svc.cluster.local"
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
		// Transport: &http.Transport{
		// 	Dial: (&net.Dialer{
		// 		Timeout: 10 * time.Second,
		// 	}).Dial,
		// },
		// ModifyResponse: func(r *http.Response) error {
		// 	// return nil
		// 	//
		// 	// purposefully return an error so ErrorHandler gets called
		// 	return errors.New("uh-oh")
		// },
		// ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
		// 	fmt.Printf("error was: %+v", err)
		// 	rw.WriteHeader(http.StatusInternalServerError)
		// 	rw.Write([]byte(err.Error()))
		// },
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":8081", nil)
	log.Error(err, "Impossible to listen on 8081 port")
}

func waitForWakeUp(serviceName string, serviceNamespace string) error {
	//TODO push all waiting connections into one queue
	timeout := time.After(30 * time.Second)
	tick := time.Tick(2000 * time.Millisecond)

	err := errors.New("timed out")

	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return err
		// Got a tick, we should check on doSomething()
		case <-tick:
			err = makeOneRequest(serviceName, serviceNamespace)
			if err == nil {
				return nil
			}
		}
	}
}

func makeOneRequest(serviceName string, serviceNamespace string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	_, err := client.Get("http://" + serviceName + "." + serviceNamespace + ".svc.cluster.local/healthz")
	if err != nil {
		return err
	}
	return nil
}
