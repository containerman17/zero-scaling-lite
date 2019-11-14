package controllers

func checkDownLoop(r *ScalingBackInfoReconciler) {
	log := r.Log
	log.V(1).Info("List", "Ingresses on watch", len(ingressesCollection))
	//http://prometheus-server.ingress-nginx.svc.cluster.local:9090
}
