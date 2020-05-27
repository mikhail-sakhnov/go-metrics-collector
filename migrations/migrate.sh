#!/usr/bin/env bash

set -ex

go run ./*.go --action=init --databaseUri=$1
go run ./*.go --action=up --databaseUri=$1