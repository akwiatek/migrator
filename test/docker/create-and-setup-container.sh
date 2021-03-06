#!/usr/bin/env bash

if [[ ! -z "$DEBUG" ]]; then
  set -x
fi

cd `dirname $0`

CONTAINER_TYPE=$1

for SCRIPT in scripts/*.sh; do
  source "$SCRIPT"
done

case $CONTAINER_TYPE in
  postgres )
    postgresql_start $CONTAINER_TYPE
    ;;
  mysql|mariadb|percona )
    mysql_start $CONTAINER_TYPE
    ;;
  mssql )
    mssql_start $CONTAINER_TYPE
    ;;
  * )
    >&2 echo "Unknown container type $CONTAINER_TYPE"
    exit 1

esac
