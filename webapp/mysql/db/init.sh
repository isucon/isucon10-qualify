#!/bin/bash
set -xe
set -o pipefail

CURRENT_DIR=$(cd $(dirname $0);pwd)
export MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
export MYSQL_PORT=${MYSQL_PORT:-3306}
export MYSQL_USER=${MYSQL_USER:-isucon}
export MYSQL_DBNAME=${MYSQL_DBNAME:-isuumo}
export MYSQL_PWD=${MYSQL_PASS:-isucon}
export LANG="C.UTF-8"
cd $CURRENT_DIR

cat 0_Schema.sql 1_DummyEstateData.sql 2_DummyChairData.sql 3_AddIndex.sql | mysql --defaults-file=/dev/null -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER $MYSQL_DBNAME
echo "set global slow_query_log_file = '/var/lib/mysql/slow-query.log';set global long_query_time=0;set global slow_query_log = ON;set global log_output = 'FILE';" | mysql --defaults-file=/dev/null -h $MYSQL_HOST -P $MYSQL_PORT -u root $MYSQL_DBNAME


