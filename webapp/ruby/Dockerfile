FROM ruby:2.7.1-slim

RUN apt-get update && \
    apt-get install -y build-essential default-mysql-client default-libmysqlclient-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir /app
COPY Gemfile Gemfile.lock /tmp/
RUN cd /tmp && \
  bundle config set deployment true && \
  bundle config set path /gems && \
  bundle config set without 'development test' && \
  bundle install -j4

WORKDIR /app
COPY . /app

CMD ["bundle", "exec", "unicorn", "-c", "unicorn.conf.rb"]
