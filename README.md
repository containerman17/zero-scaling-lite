# Zero scaler for Kubernetes
## Simple, lightweight
### Works on Ingress

Checks ingress with `zero-scaling\sleep_enabled` metrics every 60 seconds. If ingress has traffic after `zero-scaling\sleep_after` seconds, scales down deployment `zero-scaling\deployment_name` to zero and redirect all traffic to proxy.

On request received, holds it, rescales the deployment and sends traffic to `zero-scaling\service_name`.

Depends on Prometheus. 

## Proxy component


## Scaler service

### HTTP /wakeup/[domainName]
### On Ingress update
### Every minute
