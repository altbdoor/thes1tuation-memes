FROM docker.io/library/ruby:3.3-slim-bookworm AS base

ENV JEKYLL_VERSION="4.4.1" \
    RUBY_YJIT_ENABLE="true" \
    DEBIAN_FRONTEND=noninteractive \
    BUNDLE_GEMFILE=/app/Gemfile

# ==========

FROM base AS builder

WORKDIR /app/

RUN apt-get update -qq \
    && apt-get install -yqq --no-install-recommends build-essential

RUN bundle init \
    && echo "gem 'jekyll', '$JEKYLL_VERSION'" >> Gemfile \
    && echo "gem 'minima'" >> Gemfile \
    && echo "gem 'jekyll-paginate-v2'" >> Gemfile \
    && bundle install --jobs $(nproc)

# ==========

FROM base AS app

COPY --from=builder /app /app
COPY --from=builder /usr/local/bundle /usr/local/bundle

WORKDIR /srv

EXPOSE 4000 35729
