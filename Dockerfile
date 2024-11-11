FROM ruby:3.3-alpine as base

ENV JEKYLL_VERSION="4.3.4" \
    RUBY_YJIT_ENABLE="true"

# ==========

FROM base as builder

WORKDIR /app/

RUN apk add build-base
RUN bundle init \
    && echo "gem 'jekyll', '$JEKYLL_VERSION'" >> Gemfile \
    && echo "gem 'minima'" >> Gemfile \
    && echo "gem 'jekyll-paginate-v2'" >> Gemfile \
    && bundle install --jobs 2

# ==========

FROM base as app

COPY --from=builder /app /app
COPY --from=builder /usr/local/bundle /usr/local/bundle

WORKDIR /srv

EXPOSE 4000 35729
