# Container Solution API-exercise

Container Solution API-exercise, to assess technical proficiency with Software Engineering, DevOps, and Infrastructure tasks.

## The Titanic API

[![GoDoc](https://godoc.org/gitlab.com/hyperd/titanic?status.svg)](https://godoc.org/gitlab.com/hyperd/titanic)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/hyperd/titanic)](https://goreportcard.com/report/gitlab.com/hyperd/titanic)
[![wercker status](https://app.wercker.com/status/4bf9af6e8c03d46082ba0d40c192e5b6/s/master "wercker status")](https://app.wercker.com/project/byKey/4bf9af6e8c03d46082ba0d40c192e5b6)

The Titanic API is written in **golang**. It leverages **go-kit** to grant better modularity and micro-services support out-of-the-box.

The data layer is designed around [**CockroachDB**](https://www.cockroachlabs.com), deployed to the K8S cluster leveraging a [StatefulSet configuration](./deploy/k8s/cockroachdb/cockroachdb-statefulset-secure.yaml).

### Build the API

There are two ways here available to build the API code; a targetted method and a [cross-platform builder script](./build.bash). Both allow to create portable executables, compatible with [Alpine Linux](https://www.alpinelinux.org/), compiled statically linking C bindings `-installsuffix 'static'`, and omitting the symbol and debug info `-ldflags "-s -w"`.

#### Targetted build, based on your system/architecture

```bash
# change according to your system/architecture
CGO_ENABLED=0 GOARCH=[amd64|386] GOOS=[linux|darwin] go build -ldflags="-w -s" -a -installsuffix 'static' -o titanic cmd/titanic/main.go
```

#### Cross-platform build, leveraging the [build.bash](./build.bash) script

```bash
chmod +x build.bash && ./build.bash
```

**The [build.bash](./build.bash) script will also re-build and push the docker images to our private registry [GCR](https://cloud.google.com/container-registry/).**

Currently, the builds in the [releases](./releases/) folder are available for the following platforms and architectures:

- darwin / amd64;
- darwin / 386;
- linux / amd64;
- linux / 386.

### Build the Docker images locally

The Docker images are in the [docker](./docker/) folder.
To build the *production* API Docker image, run these commands in your terminal:

```bash
cd docker
docker build --no-cache --file Dockerfile -t gcr.io/$PROJECT_ID/titanic-api:latest .
```

To build the *development* Docker image, you must extend the build context to the top-level folder of this repo, and include the files [go.mod](./go.mod) and [go.sum](./go.sum), along with the source code and run the API:

```bash
# the Dockerfile the build context must be extended one level up to allow the proper files inclusions
docker build --no-cache --file ./docker/devel.Dockerfile -t gcr.io/$PROJECT_ID/titanic-api:devel .
```

### Run the API locally in Docker

Running the API locally is quite simple; it does not require any particular language or library installed on your system, other than **Docker** and **docker-compose**.
The golang server is listening both on port `8443/TCP` over **TLS** and `3000/TCP` over **http**; to properly run the API locally, before spawning the docker stack, you need to generate the certs as follow, from the top dir of this repo:

```bash
mkdir -p $(pwd)/tls; \
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
-subj "/C=NL/ST=Amsterdam/L=Amsterdam/O=hyperd/CN=titanic-api.hyperd.sh" \
-keyout $(pwd)/tls/tls.key -out $(pwd)/tls/tls.crt; \
chmod 444 tls/*
```

Once the _cert_ and the _key_ are in place, it's enough to lift the docker-compose file with the command described here below:

```bash
docker-compose up -d
```

### API Walkthrough

The Titanic API exposes the following methods:

`POST /people/` adds another passenger to the people table:

```bash
payload='
{
  "survived": true,
  "pclass": 1,
  "name": "Francesco",
  "sex": "M",
  "age": 30,
  "siblings_spouses_abroad": 1,
  "parents_children_aboard": 1,
  "fare": 7.34
}
'
curl -d "$payload" -H "Content-Type: application/json" -X POST http://localhost:3000/people/ | jq
{
  "id": "bcf1d1e9-056d-46cf-9baa-aed0e6ffd219"
}
```

`GET /people/:uuid` retrieves the given passenger by uuid from the people collection:

```bash
curl http://localhost:3000/people/35d4ab59-fa9d-478d-a57e-61b526ee0a33 | jq
{
  "people": {
    "uuid": "35d4ab59-fa9d-478d-a57e-61b526ee0a33",
    "survived": true,
    "pclass": 1,
    "name": "Francesco",
    "sex": "M",
    "age": 30,
    "siblings_spouses_abroad": 1,
    "parents_children_aboard": 1,
    "fare": 7.34
  }
}
```

`DELETE /people/:uuid` removes the given passenger:

```bash
curl -X "DELETE" http://localhost:3000/people/35d4ab59-fa9d-478d-a57e-61b526ee0a33
{
  "id": "35d4ab59-fa9d-478d-a57e-61b526ee0a33"
}
```

`PATCH /people/:uuid` partial update of the passenger information:

```bash
payload='
{
  "siblings_spouses_abroad": 1,
  "parents_children_aboard": 2
}
'
curl -d "$payload" -H "Content-Type: application/json" -X PATCH -k http://localhost/people/35d4ab59-fa9d-478d-a57e-61b526ee0a33
{}
```

`PUT /people/:uuid` posts updated information about a given passenger:

```bash
payload='
{
  "survived": true,
  "pclass": 1,
  "name": "Francesco",
  "sex": "M",
  "age": 30,
  "siblings_spouses_abroad": 1,
  "parents_children_aboard": 1,
  "fare": 9.81
}
'
curl -d "$payload" -H "Content-Type: application/json" -X PUT http://localhost:3000/people/35d4ab59-fa9d-478d-a57e-61b526ee0a33
{}
```

`GET /people/` retrieves all the passengers of the Titanic:

```bash
curl http://localhost:3000/people/ | jq
{
  "people": [
    {
      "uuid": "30615024-ada8-4af6-8611-882c006d17f4",
      "survived": true,
      "pclass": 1,
      "name": "Francesco",
      "sex": "M",
      "age": 30,
      "siblings_spouses_abroad": 1,
      "parents_children_aboard": 2,
      "fare": 7.34
    },
    {
      "uuid": "363f558a-eeb1-4bf6-b570-33e61e60b867",
      "survived": true,
      "pclass": 1,
      "name": "Anne McLeod",
      "sex": "F",
      "age": 49,
      "siblings_spouses_abroad": 1,
      "parents_children_aboard": 3,
      "fare": 9.34
    },
    ...
  ]
}
```

## Deploy the API to GCP

To deploy the stack to **GKE** on [GCP](https://cloud.google.com) follow this [documentation](./deploy/README.md).
