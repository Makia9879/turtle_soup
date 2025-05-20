#!/usr/bin/env bash

docker compose up mysql -d
docker compose up redis -d

docker compose run init_mysql
docker compose run init_app_db
docker compose run init_app

docker compose up -d
