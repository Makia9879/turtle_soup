#!/bin/bash

docker compose down redis
docker compose up redis -d
