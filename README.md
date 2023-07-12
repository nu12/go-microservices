# go-microservices

This repo contains the source code from the excellent [Trevor Sawler](https://github.com/tsawler)'s [Working with Microservices in Go (Golang)](https://www.udemy.com/course/working-with-microservices-in-go/) course.


## Local deployment

```bash
cd project
docker compose up --build
```

Available endpoints:

* localhost:8080 - Frontend
* localhost:8081 - Backend broker service
* localhost:8025 - Mailhog
* localhost:9080 - pgAdmin (admin@example.com/password)
* localhost:8081 - Mongo Express


## Kubernetes deployment

```bash
cd project
helm dependency update charts/go-microservices
```

Create and edit a `my_values.yaml` with the specifications that suit your cluster. To review the values of each component of the deployment, run `helm show values charts/go-microservices/charts/<name-of-the-chart>`.

Example of a minimal `values.yaml` file configuration:

```yaml
front-end:
  ingress:
    enabled: true
    className: "nginx"
    hosts:
      - host: minikube.local
        paths:
          - path: /
            pathType: ImplementationSpecific
  env:
    - name: BACKEND_ADDRESS
      value: "http://backend.minikube.local"

broker-service:
  ingress:
    enabled: true
    className: "nginx"
    hosts:
      - host: backend.minikube.local
        paths:
          - path: /
            pathType: ImplementationSpecific
```

To deploy, run:

```bash
helm install myapp charts/go-microservices --values my_values.yaml -n namespace
```