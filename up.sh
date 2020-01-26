#!/bin/bash

# Make us independent from working directory
pushd `dirname $0` > /dev/null
WORK_DIR=`pwd`
popd > /dev/null

# Prerequisites
pwgen > /dev/null 2>&1 || { echo >&2 "pwgen not found."; exit 1; }
docker --version > /dev/null 2>&1 || { echo >&2 "Docker not found. Please install it via https://docs.docker.com/install/"; exit 1; }
docker-machine --version > /dev/null 2>&1 || { echo >&2 "Docker machine not found. https://docs.docker.com/machine/install-machine/"; exit 1; }
docker-compose --version > /dev/null 2>&1 || { echo >&2 "Docker compose not found. Please install it via https://docs.docker.com/compose/install/"; exit 1; }

if [ ! -f "${WORK_DIR}/.env" ]; then
    if [ -f "${WORK_DIR}/.env.dist" ]; then
        cp "${WORK_DIR}/.env.dist" "${WORK_DIR}/.env"
    else
        echo >&2 "The .env file does not exist. Project setup will not work"
        exit 1
    fi
fi

# Assign user id and group id into variables
TMP_USER_ID=`id -u`
TMP_GROUP_ID=`id -g`

# Always validate user id and group id before start using .env file
sed -i 's#USER_ID=.*#'"USER_ID=${TMP_USER_ID}"'#g' ${WORK_DIR}/.env
sed -i 's#GROUP_ID=.*#'"GROUP_ID=${TMP_USER_ID}"'#g' ${WORK_DIR}/.env
sed -i 's#MYSQL_ROOT_PASSWORD=root#'"MYSQL_ROOT_PASSWORD=`pwgen -s 20 1`"'#g' ${WORK_DIR}/.env
sed -i 's#MYSQL_PASSWORD=db_password#'"MYSQL_PASSWORD=`pwgen -s 20 1`"'#g' ${WORK_DIR}/.env

# Load .env file into the current shell script
source ${WORK_DIR}/.env

# Ensure all folders exists
mkdir -p ${PROJECT_PATH}
mkdir -p ${PROJECT_PATH}/${SYMFONY_LOG_PATH}
mkdir -p ${SYMFONY_LOG_PATH}
mkdir -p ${COMPOSER_PATH}
mkdir -p ${COMPOSER_PATH}/cache
mkdir -p ${SSH_KEY_PATH}
mkdir -p ${MYSQL_DUMP_PATH}

# Create an SSH private and public keys if we do not have it
if [ ! -f "${SSH_KEY_PATH}/id_rsa" ]; then
    ssh-keygen -b 4096 -t rsa -f ${SSH_KEY_PATH}/id_rsa -q -P ""
fi

# Create a file if it does not exist
if [ ! -f "${SSH_KEY_PATH}/known_hosts" ]
then
    touch ${SSH_KEY_PATH}/known_hosts
fi

# Ensure all folders exists
mkdir -p ${NGINX_LOG_PATH}
mkdir -p ${MYSQL_DATA_PATH}
mkdir -p ${USER_CONFIG_PATH}

# Create a file if it does not exist
touch ${USER_CONFIG_PATH}/.bash_history

if [ ! -f "${USER_CONFIG_PATH}/.my.cnf" ]; then
    printf "[client]\nuser=${MYSQL_USER}\npassword=${MYSQL_PASSWORD}\n" >> ${USER_CONFIG_PATH}/.my.cnf
fi

if [ ! -f "config/php/conf.d/xdebug.ini" ]; then
    cp config/php/conf.d/xdebug.d4d config/php/conf.d/xdebug.ini
fi

docker-compose build

# Clears the screen.
clear

# Start server
echo "Starting docker containers..."
docker-compose up -d

# Documentation for end user
echo ""
echo "The following information has been set:"
echo ""
echo "Server IP: 127.0.0.1"
echo "Server Hostname: ${PROJECT_DOMAIN}"
echo ""
echo "To login now, follow this link:"
echo ""
echo "Project URL: http://${PROJECT_DOMAIN}"
echo "phpMyAdmin: http://${PROJECT_DOMAIN}:8080"
echo ""
echo "Thank you for using Docker Symfony. Should you have any questions, don't hesitate to contact us at support@prado.lt"