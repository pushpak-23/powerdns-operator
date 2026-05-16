#!/bin/sh
set -eu

: "${PDNS_GPGSQL_HOST:=postgresql}"
: "${PDNS_GPGSQL_DBNAME:=powerdns}"
: "${PDNS_GPGSQL_USER:=powerdns}"
: "${PDNS_GPGSQL_PASSWORD:=powerdns}"
: "${PDNS_API_KEY:=powerdns}"
: "${PDNS_WEBSERVER_ADDRESS:=0.0.0.0}"
: "${PDNS_WEBSERVER_PORT:=8081}"

mkdir -p /etc/powerdns
cat >/etc/powerdns/pdns.conf <<EOF
launch=gpgsql
gpgsql-host=${PDNS_GPGSQL_HOST}
gpgsql-dbname=${PDNS_GPGSQL_DBNAME}
gpgsql-user=${PDNS_GPGSQL_USER}
gpgsql-password=${PDNS_GPGSQL_PASSWORD}
api=yes
api-key=${PDNS_API_KEY}
webserver=yes
webserver-address=${PDNS_WEBSERVER_ADDRESS}
webserver-port=${PDNS_WEBSERVER_PORT}
local-address=0.0.0.0
local-port=53
EOF

exec pdns_server --daemon=no --guardian=no --disable-syslog