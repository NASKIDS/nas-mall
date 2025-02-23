#!/bin/bash

set -ex

source scripts/list_app.sh

get_app_list

for app_path in ${app_list[*]}; do
    pushd "${app_path}" || exit
    golangci-lint run -E gofumpt --path-prefix=. --fix --timeout=5m
    popd || exit
done
