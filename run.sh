#!/bin/sh

IMAGE_NAME="jekyll-dev:latest"

if [[ "$1" == "rebuild" || "$(docker images -q $IMAGE_NAME 2> /dev/null)" == "" ]]; then
    docker build -t $IMAGE_NAME .
fi

docker run --rm -it -p '4000:4000' -p '35729:35729' -v "$PWD:/srv" $IMAGE_NAME \
    jekyll serve --livereload --incremental --host 0.0.0.0 --force_polling
