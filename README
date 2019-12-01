# demo readiness gate for kubernetes
Using controller runtime to run a mutation admission controller and reconcile pod
update pod readiness gate condition to true.

still under developing, just a demo for pod readiness gate.
test kubernetes version: 1.16.0

# usage
run local
- create server crt (tls.crt) and key (tls.key)
put them in dir /tmp/k8s-webhook-server/serving-certs, using cfssl is recommanded.
- run as follow
```
sudo ./demo-readiness-gate --kubeconfig ~/.kube/config
```
- try to run pods
```
# kubectl run demo --image=nginx --replicas=3
```
- check readiness status
```
# kubectl get pod -o wide
NAME                   READY   STATUS    RESTARTS   AGE   IP            NODE                     NOMINATED NODE   READINESS GATES
demo-9c94c674c-4nnzv   1/1     Running   0          11m   10.244.0.59   cluster1-control-plane   <none>           1/1
demo-9c94c674c-fds5z   1/1     Running   0          11m   10.244.0.58   cluster1-control-plane   <none>           1/1
demo-9c94c674c-svff9   1/1     Running   0          11m   10.244.0.60   cluster1-control-plane   <none>           1/1
```



