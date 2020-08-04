#!/usr/bin/env bash

cd $(dirname $0)

set -e

project_id=$(gcloud config get-value project)

gcloud builds submit -t gcr.io/${project_id}/devenv
