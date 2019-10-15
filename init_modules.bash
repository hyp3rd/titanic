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

deps () {
	go list ./...
}

unique_repos () {
	cut -d '/' -f-3 | sort | uniq
}

no_titanic () {
	grep -v '*titanic*'
}

go_get_update () {
	while read d
	do
		echo $d
		go get -u $d/... || echo "failed, trying again with master" && cd $GOPATH/src/$d && git checkout master && go get -u -x $d
	done
}

ini_modules () {
    modules=('.' 'transport' 'transport/http' 'inmemory' 'implementation' 'cmd/titanic')

    for i in "${modules[@]}"; do
        cd $i ; rm -rf go.* ; go mod init ; cd - #go mod init ; go get -u -x ; go mod tidy ; go build ; cd -
    done

	git add -A . ; git commit -m "modules update" || : ; git push || :

	report
}

# deps | unique_repos | go_get_update

ini_modules
