#!/usr/bin/env bash

project_id=$(gcloud config get-value project)
wd="/root/$(basename ${PWD})"

docker container run -it --rm \
  --name virtual-ec \
  -v $(pwd):${wd} \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd)/.config:/root/.config \
  -v $(pwd)/.cache/go/mod:/root/go/pkg/mod/cache \
  -w ${wd} \
  gcr.io/${project_id}/devenv bash
