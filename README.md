# kubernetes-ingress-based-zero-scaling

Checks ingres with `zero-scaling\sleep_enabled` metrics every `zero-scaling\check_period` seconds. If ingress has traffic after `zero-scaling\sleep_after` seconds, scales down deployment `zero-scaling\deployment_name` to zero and redirect all traffic to proxy. 

On request received, holds it, rescales the deployment and sends traffic to `zero-scaling\service_name`.

Depends on Prometheus. 
