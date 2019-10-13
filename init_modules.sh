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
    modules=('.' 'transport' 'transport/http' 'inmemory' 'implementation' 'cmd/titanic')

    git add . ; git commit -m "modules update" ; git push || :

    for i in "${modules[@]}"; do
        cd $i ; rm -rf go.*; go mod init ; go get ; go mod tidy ; go build ; cd -
    done

    report
}

ini_modules
