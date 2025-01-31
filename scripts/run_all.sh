#!/bin/bash

. scripts/list_app.sh

get_app_list

readonly root_path=`pwd`
for app_path in ${app_list[*]}; do
    # 切换到应用程序目录
    pushd "${root_path}/${app_path}" > /dev/null || exit

    # 使用 nohup 启动应用程序
    nohup go run ./*.go &

    # 返回原来的目录
    popd > /dev/null || exit
done