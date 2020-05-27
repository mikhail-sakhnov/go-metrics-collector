#!/usr/bin/env bash

set -ex

DB_URI=$2

function init() {
  go run ./*.go --action=init --databaseUri=$DB_URI
  go run ./*.go --action=up --databaseUri=$DB_URI
}

function reset() {
  go run ./*.go --action=reset --databaseUri=$DB_URI
  init
}

$1

