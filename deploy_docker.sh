#!/usr/bin/env bash

docker compose --profile init_app up init_mysql
docker compose --profile init_app up init_app_db
docker compose --profile init_app up init_app

docker compose up -d
