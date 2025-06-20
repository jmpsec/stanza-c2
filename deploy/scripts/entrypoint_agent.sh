#!/bin/sh

# Give some time for the server to start
echo "Waiting for the server to start..."
echo "Using CALLBACKS_URL: $CALLBACKS_URL"
sleep 5

# Start the agent
cd /opt/stanza && STZ_CALLBACKS_CODE=$STANZA_TOKEN STZ_UUID=$STANZA_UUID STZ_CALLBACKS=$CALLBACKS_URL STZ_MIN=10 STZ_MAX=30 STZ_DEBUG=1 ./bin/stzagent
