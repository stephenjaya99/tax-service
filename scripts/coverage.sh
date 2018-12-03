#!/bin/bash
echo "mode: atomic" > coverage.out

PACKAGES=(controller database)

for pkg in ${PACKAGES[@]}; do
	echo $pkg
	go test -v -coverprofile=profile.out -covermode=atomic gitlab.com/stephenjaya99/tax-service/$pkg; exit_code=$?

	if [[ $exit_code -eq 1 ]]; then
		echo "Test Error"
		exit 1
	fi

	if [ -f profile.out ]; then
	  tail -n +2 profile.out >> coverage.out; rm profile.out
	  go tool cover -func=coverage.out
	fi
done

