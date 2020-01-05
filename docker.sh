#!/bin/bash

DOCKER_SCRIPT_VER=1.0.0

WORK_DIR=`pwd`

cd ${WORK_DIR}

#Added for security
chmod 700 ${WORK_DIR}/docker.sh >/dev/null 2>&1
chmod 700 ${WORK_DIR} >/dev/null 2>&1

##################################################
# Variables for bolded text
boldon="`tput -Txterm bold`"
boldoff="`tput -Txterm sgr0`"
##################################################

# Help sections
HELP_SECTIONS="DOCKER MYSQL PASSWD"
DOCKER_SET="start stop exec logs"
DOCKER_DESC="Docker for Symfony (PHP-FPM - NGINX - MySQL)"
DOCKER_REQADD="start stop exec logs"
START_ADDIT=""
STOP_ADDIT=""
EXEC_ADDIT=""
LOGS_ADDIT=""

MYSQL_SET="mysql_dump mysql_restore"
MYSQL_DESC="Backup and Restore a MySQL Database"
MYSQL_REQADD="mysql_dump mysql_restore"
DUMP_ADDIT=""
RESTORE_ADDIT=""

PASSWD_SET="passwd"
PASSWD_DESC="Show password sensitive information"
PASSWD_REQADD="passwd"
PASSWD_ADDIT="show"

##################################################

# Emulate ${!variable}
eval_var() {
    var=${1}
        if [ -z ${var} ]; then
            echo ""
        else
            eval newval="\$${var}"
            echo $newval
        fi
}

showVersion() {
    echo "${DOCKER_SCRIPT_VER} (rev: 0001)"
}

generateHelp() {
    for section in ${HELP_SECTIONS}; do
        DESC=${section}_DESC
        echo " +-----------------------------------------------------------+"
        printf " | %-55s %-2s|\n" "$(eval_var ${DESC}):"

        BUILDSET="${section}_SET"
        BUILD_ADD_TO_OUT="${section}_REQADD"
        for setting in $(eval_var ${BUILDSET}); do
            ADDIT=""
            for i in $(eval_var ${BUILD_ADD_TO_OUT}); do
                ADDIT_VAR=`echo "${i}_ADDIT" | tr "[a-z]" "[A-Z]"`
                if [ "$i" = "${setting}" ]; then
                    ADDIT="$(eval_var ${ADDIT_VAR})"
                fi
                done
                printf " | %-55s %-2s|\n" "   $0 ${setting} ${ADDIT}"
            done
    done
}

doChecks() {
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

    docker-compose build

    # Clears the screen.
    clear
}

start() {
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
    echo "MailHog: http://${PROJECT_DOMAIN}:8025"
    echo ""
    echo "Thank you for using Docker for Symfony. Should you have any questions, don't hesitate to contact us at support@prado.lt"
}

stop() {
    docker-compose down
}

dockerComposeExec() {
    # $1 command
    docker-compose exec $1 bash
}

doMysqlDump() {
    # Load .env file into the current shell script
    source ${WORK_DIR}/.env

    docker-compose exec mysql bash /tmp/db/mysql_dump.sh
}

doPasswd() {
    # $1 command
    COMMAND=$1

    if [ "${COMMAND}" = "show" ]; then
        showPasswords
    fi
}

mainHeader() {
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "Docker for Symfony (PHP-FPM - NGINX - MySQL)"
    printf " | %-55s %-2s|\n" "Written by Vasilij Dusko"
    printf " | %-55s %-2s|\n" "Version: $(showVersion)"
    printf " +%-55s+\n" "-----------------------------------------------------------"
}

showPasswords() {
    source ${WORK_DIR}/.env

    mainHeader
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "The following information has been set:"
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "Server IP: 127.0.0.1"
    printf " | %-55s %-2s|\n" "Server Hostname: ${PROJECT_DOMAIN}"
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "MySQL root username: root"
    printf " | %-55s %-2s|\n" "MySQL root password: ${MYSQL_ROOT_PASSWORD}"
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "MySQL database name: ${MYSQL_DATABASE}"
    printf " | %-55s %-2s|\n" "MySQL username: ${MYSQL_USER}"
    printf " | %-55s %-2s|\n" "MySQL password: ${MYSQL_PASSWORD}"
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "To login now, follow this link:"
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "Project URL: http://${PROJECT_DOMAIN}"
    printf " | %-55s %-2s|\n" "phpMyAdmin: http://${PROJECT_DOMAIN}:8080"
    printf " | %-55s %-2s|\n" "MailHog: http://${PROJECT_DOMAIN}:8025"
    printf " +%-55s+\n" "-----------------------------------------------------------"
}

showHelp() {
    mainHeader
    generateHelp
    printf " +%-55s+\n" "-----------------------------------------------------------"
}

case "$1" in
    start) doChecks; start
        ;;
    stop) stop
        ;;
    exec) dockerComposeExec $2
        ;;
    mysql_dump) doMysqlDump
        ;;
    passwd) doPasswd $2
        ;;
    * ) showHelp
        exit 0
        ;;
esac
