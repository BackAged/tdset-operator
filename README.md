# TDSet
A simple kubernetes operator to set replica count based on the hour of the day, built for demonstration purpose for this article -
https://shahin-mahmud.medium.com/write-your-first-kubernetes-operator-in-go-177047337eae

## Description
Built with https://github.com/operator-framework/operator-sdk

This operator lets you dynamically scale up and down your deployment replicaset based on the hour of the day.

Example: <br>
	Everyday from 18:00 to 19:00  you want to have 5 replicas of the deployment, <br>
	from 19:00 to 21:00 you want to have 3 replicas of the deployment <br>
        and rest of the time you want to have 2 replicas of the deployment.
```
apiVersion: schedule.rs/v1
kind: TDSet
metadata:
  labels:
    app.kubernetes.io/name: tdset-sample
    app.kubernetes.io/instance: tdset-sample
    app.kubernetes.io/part-of: rs
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: rs
  name: tdset-sample
  #namespace: tdtest
spec:
  container:
    image: crccheck/hello-world
    port: 8000
  schedulingConfig:
    - startTime: 18
      endTime: 19
      replica: 5
    - startTime: 19
      endTime: 21
      replica: 3
  defaultReplica: 2
```

## Helm installation
```
helm repo add tdset https://backaged.github.io/tdset-operator/
helm install tdset tdset/tdset-controller -n tdset
```
