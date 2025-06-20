#!/bin/sh

CONFIG_FILE="/opt/stanza/config/admin.json"

# Replace the values in admin.json with environment variables
cat "$CONFIG_FILE" | sed "s|_ADMIN_LISTENER|$SERVICE_LISTENER|g" | sed "s|_ADMIN_PORT|$SERVICE_PORT|g" | sed "s|_ADMIN_HOST|$SERVICE_HOST|g" | tee "$CONFIG_FILE"

cat "$CONFIG_FILE" | sed "s|_DB_HOST|$DB_HOST|g" | sed "s|_DB_PORT|$DB_PORT|g" | sed "s|_DB_NAME|$DB_NAME|g" | sed "s|_DB_USERNAME|$DB_USER|g" | sed "s|_DB_PASSWORD|$DB_PASS|g" | tee "$CONFIG_FILE"

# Start the stanza-admin service
cd /opt/stanza && ./bin/stzadmin -config="$CONFIG_FILE"
