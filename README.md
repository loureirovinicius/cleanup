# Cleanup (Cloud Provider Sanitization tool)

Cleanup is a tool designed to accomplish **effective costs on Cloud Providers** (AWS, GCP, etc...) without wasting money on unused resources - an empty Load Balancer, for example. Such tool was thought to be one of the greatest allies in a FinOps culture for its simplicity, efficiency and security.\
\
Everything can be packaged into a single binary that can run on serverless structures, pipelines and even manually. Contracts were established through interfaces so it could also be expanded to other providers and resources, so you're highly encouraged to contribute with it since the codebase is pretty simple.

## Supported resources:

| Provider | Resource | Validation method |
| -------- | -------- | ----------------- |
| AWS | eni | Checks if ENI status is "available"
| AWS | eip | Checks if EIP has an association ID
| AWS | ebs | Checks if EBS disk has its state as "Available" and if tag "cleanup-ignore" equals to "true"
| AWS | targetGroup | Checks if TargetGroup has no LoadBalancer attached
| AWS | loadBalancer | Checks if LoadBalancer has no Listener attached


## Usage:

1. Populate the configuration file as required:

config.yaml
```yaml
aws:
  region: # AWS Region
  authentication:
    profile:
      name: # AWS Profile name (optional)
      path: # AWS Config file path (optional)
    credentials:
      access_key: # AWS Access Key (optional)
      secret_key: # AWS Secret Key (optional)
```

2. Compile or run it using Docker or Go:
    - Docker
      - Building the binary using a Docker container:
        ```bash
          docker run --rm --mount type=bind,source=$(pwd),target=/bin -w /app golang:alpine go build -o cleanup cmd/main.go
        ```
      - Building and running it inside a Docker container:
        ```bash
          docker build -t cleanup:latest
          docker run -i --rm --mount type=bind,source="$(pwd)/config.yaml",target=/var/cleanup/config.yaml cleanup:latest
          /app/cleanup help
        ```
      - Pulling image from Dockerhub:
        ```bash
          docker run -i --rm --mount type=bind,source="$(pwd)/config.yaml",target=/var/cleanup/config.yaml loureirovinicius/cleanup:latest
        ```

    - Go
        ```bash
          go build -o cleanup ./cmd/main.go
        ```

3. Use it:
```bash
cleanup list --service targetGroup
cleanup validate --service targetGroup
cleanup delete --service targetGroup
```


