#!/bin/bash

set -Eeuo pipefail

traperr() {
  echo "ERROR: ${BASH_SOURCE[1]} at about line ${BASH_LINENO[0]} ${ERR}"
}

set -o errtrace
trap traperr ERR

initdb() {
    echo "Wait for servers to be up"
    sleep 10

    HOSTPARAMS="--host roach1 --insecure"
    SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

    $SQL -e "CREATE DATABASE titanic;"
    $SQL -d titanic -e "CREATE TABLE people (
    	id INT8 NOT NULL DEFAULT unique_rowid(),
    	created_at TIMESTAMPTZ NULL,
    	updated_at TIMESTAMPTZ NULL,
    	deleted_at TIMESTAMPTZ NULL,
    	uuid UUID NOT NULL,
    	survived BOOL NULL,
    	pclass INT8 NULL,
    	name STRING NULL,
    	sex STRING NULL,
    	age INT8 NULL,
    	siblings_spouses_abroad BOOL NULL,
    	parents_children_aboard BOOL NULL,
    	fare DECIMAL NULL,
    	CONSTRAINT \"primary\" PRIMARY KEY (id ASC, uuid ASC),
    	INDEX idx_peoples_deleted_at (deleted_at ASC),
    	FAMILY \"primary\" (id, created_at, updated_at, deleted_at, uuid, survived, pclass, name, sex, age, siblings_spouses_abroad, parents_children_aboard, fare)
    );"
}

initdb
