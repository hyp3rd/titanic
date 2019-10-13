#!/bin/bash
set -Eeuo pipefail

traperr() {
  echo "ERROR: ${BASH_SOURCE[1]} at about line ${BASH_LINENO[0]} ${ERR}"
}

set -o errtrace
trap traperr ERR

report () {
	cat >&2 <<-'EOF'

The Docker images are now up to date;
remember to change the sha256 checksum in the k8s deployment file.

	EOF
}

# build it's a simple golang cross-compiler that generates alpine linux compatible binaries
build () {
    docker run --rm -it -v "$PWD":/usr/src/app -w /usr/src/app golang:latest bash -c '
    for GOOS in darwin linux; do
        for GOARCH in 386 amd64; do
         export GOOS GOARCH
         CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix "static" -o releases/titanic-$GOOS-$GOARCH cmd/titanic/main.go
        done
    done
    '
}

push_to_scm() {
  git status
  git add . ; git commit -m "builds update" ; git push
}

update_docker_images () {
    cd $(pwd)/releases

    if [[ -z "$PROJECT_ID" ]] || [[ -z "$REGION" ]]; then
      echo "To run this deployment you need to export PROJECT_ID and REGION as follows:
      export REGION=<region e.g. europe-west1>
      export PROJECT_ID=<project name e.g. hyperd-titanic-api>";
      exit 1
    fi

    # build the cs-api image with our modifications (see Dockerfile) and tag for private GCR
    docker build --no-cache --file ../docker/Dockerfile -t gcr.io/$PROJECT_ID/cs-api .

    # configure pushing to private GCR, and push our image
    gcloud auth configure-docker -q
    docker push gcr.io/$PROJECT_ID/cs-api

    report
}

build

push_to_scm

update_docker_images
