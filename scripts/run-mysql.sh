#!/bin/bash

shellpath=`dirname "$0"`
cd ${shellpath}

CURRENT_DIR=$(pwd)

docker run -itd --restart=always \
--name paopao-db \
--network host \
-v ${CURRENT_DIR}/paopao-mysql.sql:/docker-entrypoint-initdb.d/paopao.sql \
-v ${CURRENT_DIR}/data:/var/lib/mysql \
-e MYSQL_DATABASE=paopao \
-e MYSQL_USER=paopao \
-e MYSQL_PASSWORD=paopao \
-e MYSQL_RANDOM_ROOT_PASSWORD=yes \
mysql:8.0