#!/bin/sh
set -e

echo "Waiting for Postgres to be ready..."
until nc -z $PGSQL_HOST $PGSQL_PORT; do
  echo "Postgres is unavailable - sleeping"
  sleep 2
done

echo "Running migrations..."
./app migrate up

echo "Starting service: $1"
./app "$1"
