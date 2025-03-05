#!/bin/bash
set -xe
. scripts/list_app.sh

get_app_list

readonly root_path=`pwd`
for app_path in ${app_list[*]}; do
    # 切换到应用程序目录
    pushd "${root_path}/${app_path}" > /dev/null || exit

    sh build.sh

    # 返回原来的目录
    popd > /dev/null || exit
done
