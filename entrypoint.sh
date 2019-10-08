#!/bin/sh

export DATA_SOURCE_NAME="${DB_USER}:${DB_PASSWORD}@(${DB_DNS}:${DB_PORT:3306})/" && mysqld_exporter --collect.all
