#!/bin/sh

CONFIG_FILE="/opt/stanza/config/http.json"

# Replace the values in admin.json with environment variables
cat "$CONFIG_FILE" | sed "s|_HTTP_LISTENER|$SERVICE_LISTENER|g" | sed "s|_HTTP_PORT|$SERVICE_PORT|g" | sed "s|_HTTP_HOST|$SERVICE_HOST|g" | sed "s|_HTTP_TOKEN|$SERVICE_TOKEN|g" | sed "s|_HTTP_CALLBACKS_PORT|$SERVICE_CALLBACKS_PORT|g" | tee "$CONFIG_FILE"

cat "$CONFIG_FILE" | sed "s|_DB_HOST|$DB_HOST|g" | sed "s|_DB_PORT|$DB_PORT|g" | sed "s|_DB_NAME|$DB_NAME|g" | sed "s|_DB_USERNAME|$DB_USER|g" | sed "s|_DB_PASSWORD|$DB_PASS|g" | tee "$CONFIG_FILE"

# Start the stanza-admin service
cd /opt/stanza && ./bin/stzhttp -config="$CONFIG_FILE"
