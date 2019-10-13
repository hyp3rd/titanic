#!/bin/bash
set -Eeuo pipefail

traperr() {
  echo "ERROR: ${BASH_SOURCE[1]} at about line ${BASH_LINENO[0]} ${ERR}"
}

set -o errtrace
trap traperr ERR

report () {
	cat >&2 <<-'EOF'

I've successfully initialized all the modules.

	EOF
}

ini_modules () {
    modules=('transport' 'transport/http' 'inmemory' 'implementation' 'cmd/titanic')

    git add . ; git commit -m "v1.0.0" ; git push

    for i in "${modules[@]}"; do
      cd $i ; go mod init ; go get ; go mod tidy ; go build ; cd -
    done
}

ini_modules
