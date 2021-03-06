#!/bin/bash

set -Eeuo pipefail

traperr() {
  echo "ERROR: ${BASH_SOURCE[1]} at about line ${BASH_LINENO[0]} ${ERR}"
}

set -o errtrace
trap traperr ERR

initdb() {
	# TODO: implement a more elegant way to wait for the cluster
	# to be up and running.
    echo "Wait for servers to be up"
    sleep 10

    HOSTPARAMS="--host roach1 --insecure"
    SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

    $SQL -e "CREATE DATABASE titanic; CREATE USER IF NOT EXISTS d4gh0s7; GRANT ALL ON DATABASE titanic TO d4gh0s7;"
    $SQL -d titanic -e "CREATE TABLE people (
    	created_at TIMESTAMPTZ NULL,
    	updated_at TIMESTAMPTZ NULL,
    	deleted_at TIMESTAMPTZ NULL,
    	id UUID NOT NULL,
    	survived BOOL NULL,
    	pclass INT8 NULL,
    	name STRING NULL,
    	sex STRING NULL,
    	age INT8 NULL,
    	siblings_spouses_abroad INT8 NULL,
    	parents_children_aboard INT8 NULL,
    	fare DECIMAL NULL,
		CONSTRAINT \"primary\" PRIMARY KEY (id ASC),
    	INDEX idx_peoples_deleted_at (deleted_at ASC),
    	FAMILY \"primary\" (created_at, updated_at, deleted_at, id, survived, pclass, name, sex, age, siblings_spouses_abroad, parents_children_aboard, fare)
    );"
}

initdb
