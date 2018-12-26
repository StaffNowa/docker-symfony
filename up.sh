#!/bin/bash

if [ ! -f .env ]; then
    cp .env.dist .env

    TMP_USER_ID=`id -u`
    TMP_GROUP_ID=`id -g`
    sed -i 's#USER_ID=#'"USER_ID=${TMP_USER_ID}"'#g' .env
    sed -i 's#GROUP_ID=#'"GROUP_ID=${TMP_GROUP_ID}"'#g' .env
fi

source ./.env

mkdir -p ${PROJECT_PATH}
mkdir -p ${PROJECT_PATH}/${SYMFONY_LOG_PATH}
mkdir -p ${SYMFONY_LOG_PATH}
mkdir -p ${COMPOSER_PATH}
mkdir -p ${SSH_KEY_PATH}
mkdir -p ${NGINX_LOG_PATH}
mkdir -p ${MYSQL_DATA_PATH}

docker-compose build
docker-compose up -d