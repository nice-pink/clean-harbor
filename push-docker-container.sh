#!/bin/bash

docker-build-push -i clean-harbor -p nicepink -d Dockerfile -r docker -t runner
