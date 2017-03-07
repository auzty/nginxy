# Nginxy

Nginxy is a nginx configuration generator that using **docker swarm** , nginxy watch the docker swarm event for service created and then generating the nginx configuration automatically

## How to Compile

```bash
CGO_ENABLED=0 go build -o nginxy nginxy.go nginx-conf.go
```

## How to Run

By default **nginxy** will set the endpoint of docker socket to **/tmp/docker.sock** but can be adjustable using the parameter

```bash
Usage of nginxy:
  -endpoint string
        Docker Endpoint (socket only) (default "/tmp/docker.sock")
```

## How to Use

### 1. Non-Docker Environment

Just execute the **nginxy** using the -endpoint parameter 

```bash
nginxy -endpoint /var/run/docker.sock
```

### 2. Docker Environment

- Build this **nginxy** image

```bash
docker build -t nginxy .
```

- Make sure the docker is swarm enabled

- Create the swarm Overlay Network first

```bash
docker network create --driver overlay --subnet 1.1.1.0/24 \
nginxy-backend --attachable
```

- Make sure the **nginxy** container run on **manager** node 

```bash
docker service create --network nginxy-backend --name nginxy \
--constraint 'node.role == manager' --publish 80:80 --publish 443:443 \
--mount type=bind,src=/var/run/docker.sock,dst=/tmp/docker.sock,ro=true \
--mount type=bind,src=/opt/ssl,dst=/etc/nginx/certs nginxy
```


## Other Service Creating

If there are new service(s) created on swarm, **nginxy** will watch it and update the nginx configuration and reload it based on the service labels:

- example running hello-world app (with SSL)
```bash
docker service create --label nginxy.domain=hello-world.io \
--label nginxy.port=3000 --label nginxy.ssl=true \
--label nginxy.ssl.cert=hello.cert --label nginxy.ssl.key=hello.key \
--name hello registry.dwp.io/go-hello:latest
```

- example running hello-world app (non SSL)
```bash
docker service create --label nginxy.domain=hello-world.io \
--label nginxy.port=3000 --name hello registry.dwp.io/go-hello:latest
```

**Labels Description :**
- **nginxy.domain** ``set the domain name``
- **nginxy.port** ``set the exposed container port``
- **nginxy.ssl** ``set to true if using ssl``
- **nginxy.ssl.cert** ``name of certificate``
- **nginxy.ssl.key** ``name of the certificate key``


