FROM ruby:3.3-alpine as builder

ENV JEKYLL_VERSION='4.3.3'

RUN apk add build-base

WORKDIR /app/

RUN bundle init \
    && bundle add jekyll -v "$JEKYLL_VERSION" \
    && bundle add minima \
    && bundle add jekyll-paginate-v2

# ==========

FROM ruby:3.3-alpine as app

COPY --from=builder /app /app
COPY --from=builder /usr/local/bundle /usr/local/bundle

WORKDIR /srv

EXPOSE 4000 35729
