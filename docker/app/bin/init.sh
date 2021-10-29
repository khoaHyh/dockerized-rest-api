#!/usr/bin/env bash

echo 'Running migrations...'
/dockerized-rest-api/migrate up > /dev/null 2>&1 &

echo 'Deleting mysql-client...'
apk del mysql-client

echo 'Start application...'
/dockerized-rest-api/app
