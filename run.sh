#!/bin/sh

set -e

IMAGE_NAME="jekyll-dev:latest"

if [[ "$1" == "rebuild" || "$(docker images -q $IMAGE_NAME 2> /dev/null)" == "" ]]; then
    docker build --progress plain -t $IMAGE_NAME .
fi

if [[ "$1" == "sh" ]]; then
    docker run --rm -it -v "$PWD:/srv" $IMAGE_NAME sh
elif [[ "$1" == "build" ]]; then
    docker run --rm -it -v "$PWD:/srv" $IMAGE_NAME bundle exec jekyll build
else
    docker run --rm -it -p '4000:4000' -p '35729:35729' -v "$PWD:/srv" $IMAGE_NAME \
        bundle exec jekyll serve --livereload --host 0.0.0.0 --force_polling
fi
