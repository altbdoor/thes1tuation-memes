FROM docker.io/library/ruby:3.3-slim-bookworm AS base

ENV JEKYLL_VERSION="4.4.1" \
    RUBY_YJIT_ENABLE="true"

# ==========

FROM base AS builder

WORKDIR /app/

RUN apt-get update -qq \
    && apt-get install -yqq build-essential

RUN bundle init \
    && echo "gem 'jekyll', '$JEKYLL_VERSION'" >> Gemfile \
    && echo "gem 'minima'" >> Gemfile \
    && echo "gem 'jekyll-paginate-v2'" >> Gemfile \
    # https://github.com/protocolbuffers/protobuf/issues/19932
    && echo "gem 'google-protobuf', '< 4.29.3'" >> Gemfile \
    && bundle install --jobs 2

# ==========

FROM base AS app

COPY --from=builder /app /app
COPY --from=builder /usr/local/bundle /usr/local/bundle

WORKDIR /srv

EXPOSE 4000 35729
