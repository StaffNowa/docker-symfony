#!/bin/bash
# https://www.stuartellis.name/articles/shell-scripting/#enabling-better-error-handling-with-set
set -Eeuo pipefail

# Based on mongo/docker-entrypoint.sh
# https://github.com/docker-library/mongo/blob/master/docker-entrypoint.sh#L303
if [ "${MONGODB_USERNAME}" ] && [ "${MONGODB_PASSWORD}" ]; then
    rootAuthDatabase='admin'

    "${mongo[@]}" -u "${MONGODB_ROOT_USERNAME}" -p "${MONGODB_ROOT_PASSWORD}" --authenticationDatabase "$rootAuthDatabase" "${MONGODB_DATABASE}" <<-EOJS
				db.createUser({
					user: $(_js_escape "$MONGODB_USERNAME"),
					pwd: $(_js_escape "$MONGODB_PASSWORD"),
					roles: [ { role: 'root', db: $(_js_escape "$rootAuthDatabase") } ]
				})
			EOJS

    printf " +%-55s+\n" "-----------------------------------------------------------"
		printf " | %-55s %-2s|\n" "MongoDB root username: ${MONGODB_ROOT_USERNAME}"
    printf " | %-55s %-2s|\n" "MongoDB root password: ${MONGODB_ROOT_PASSWORD}"
    printf " +%-55s+\n" "-----------------------------------------------------------"
    printf " | %-55s %-2s|\n" "MongoDB database name: ${MONGODB_DATABASE}"
    printf " | %-55s %-2s|\n" "MongoDB username: ${MONGODB_USERNAME}"
    printf " | %-55s %-2s|\n" "MongoDB password: ${MONGODB_PASSWORD}"
    printf " +%-55s+\n" "-----------------------------------------------------------"
fi