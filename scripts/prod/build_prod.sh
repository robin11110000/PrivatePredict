#!/bin/bash

# Make sure the script can only be run via PrivatePredict Script
[ -z "$CALLED_FROM_PrivatePredict" ] && { echo "Not called from PrivatePredict"; exit 42; }

# Pull images
echo "Pulling images ..."
$COMPOSE --env-file "$SCRIPT_DIR"/.env --file "$SCRIPT_DIR/scripts/docker-compose-prod.yaml" pull
echo

echo "Images pulled."
echo
echo "Your admin credentials are:"
echo "Username: admin"
echo "Password: $ADMIN_PASS"
