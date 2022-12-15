# TDSet
A simple kubernetes operator to set replica count based on the hour of the day, built for demonstration purpose.

## Description
Built with https://github.com/operator-framework/operator-sdk

This operator lets you dynamically scale up and down your deployment replicaset based on the hour of the day.<b>

Example:
	Everyday from 10:00 to 18:00  you want to have 5 replicas of the deployment, <br>
	from 18:00 to 21:00 you want to have 3 replicas of the deployment <br>
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
    - startTime: 1
      endTime: 2
      replica: 1
  defaultReplica: 3
```
