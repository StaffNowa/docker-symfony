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
mkdir -p ${COMPOSER_PATH}/cache
mkdir -p ${SSH_KEY_PATH}

if [ ! -f "${SSH_KEY_PATH}/id_rsa" ]
then
    ssh-keygen -b 4096 -t rsa -f ${SSH_KEY_PATH}/id_rsa -q -P ""
fi

if [ ! -f "${SSH_KEY_PATH}/known_hosts" ]
then
    touch ${SSH_KEY_PATH}/known_hosts
fi

mkdir -p ${NGINX_LOG_PATH}
mkdir -p ${MYSQL_DATA_PATH}
mkdir -p ${USER_CONFIG_PATH}
touch ${USER_CONFIG_PATH}/.bash_history

docker-compose build
docker-compose up -d