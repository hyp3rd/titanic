#!/bin/bash
set -Eeuo pipefail

traperr() {
  echo "ERROR: ${BASH_SOURCE[1]} at about line ${BASH_LINENO[0]} ${ERR}"
}

set -o errtrace
trap traperr ERR

report () {
	cat >&2 <<-'EOF'

A three nodes cluster is running:

DB Listening on port:           26257
Dashboard listening on port:    8080

Beaware, this is an insecure setup:

DO NOT USE IT IN ANY PRODUCTION ENVIRONMENT.

	EOF
}

cockroachdb_data_dir="${PWD}/cockroach-data"

init () {
    cockroachdb_image_version="v19.1.5"

    docker network create -d bridge roachnet || :

    docker run -d --name=roach1 --hostname=roach1 --net=roachnet \
        -p 26257:26257 -p 8080:8080  \
        -v "${cockroachdb_data_dir}/roach1:/cockroach/cockroach-data" \
        cockroachdb/cockroach:${cockroachdb_image_version} start --insecure

    for i in {2..3}; do
        echo "roach$i is joining the cluster:"
        docker run -d --name=roach$i --hostname=roach$i \
            --net=roachnet -v "${cockroachdb_data_dir}/roach$i:/cockroach/cockroach-data" \
            cockroachdb/cockroach:${cockroachdb_image_version} start --insecure --join=roach1
    done

    report
}

cleanup () {
    echo "Cleaning up Containers and Volumes:"
    docker stop $(docker ps -aq) ; docker rm -v $(docker ps -aq) ; rm -rf ${cockroachdb_data_dir}/*
    echo "All the Containers and Volumes are removed."
}

[[ "$(ls -A ${cockroachdb_data_dir})" ]] ; cleanup || init
