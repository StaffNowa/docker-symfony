#!/bin/sh

#
# This script do backup all mysql databases in serveral files in exists directory
# file name: backup+db_name+date.sql.gz
#

mysql -Ne "show databases" | grep -v "schema\|mysql\|information_schema" | while read db; do mysqldump $db | gzip > "/tmp/db/backup_"$db"_$(date +%Y%m%d).sql.gz"; done