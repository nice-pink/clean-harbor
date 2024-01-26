#!/bin/bash

docker-build-push -i clean-harbor -p nice-pink -d Dockerfile -r docker -t runner
