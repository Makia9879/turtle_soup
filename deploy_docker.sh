#!/usr/bin/env bash

docker compose up --remove-orphans mysql -d
docker compose up --remove-orphans redis -d

docker compose run --remove-orphans init_mysql
docker compose run --remove-orphans init_app_db
docker compose run --remove-orphans init_app

docker compose up --remove-orphans -d
