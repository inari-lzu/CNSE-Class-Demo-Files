#!/bin/bash
docker buildx create --use 
docker buildx build --platform linux/amd64,linux/arm64 -f ./voter-api/Dockerfile ./voter-api -t xf2000/cst680:multiplatform --push