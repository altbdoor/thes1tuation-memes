FROM ruby:3.3-alpine as base

ENV JEKYLL_VERSION="4.3.4"

# ==========

FROM base as builder

WORKDIR /app/

RUN apk add build-base
RUN bundle init \
    && bundle add jekyll -v "$JEKYLL_VERSION" \
    && bundle add minima jekyll-paginate-v2

# ==========

FROM base as app

COPY --from=builder /app /app
COPY --from=builder /usr/local/bundle /usr/local/bundle

WORKDIR /srv

EXPOSE 4000 35729
