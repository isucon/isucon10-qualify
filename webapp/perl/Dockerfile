FROM perl:5.32.0

EXPOSE 1323

WORKDIR /app

RUN apt-get update && apt-get install -y wget default-mysql-client

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY cpanfile .
RUN cpanm -nq App::cpm
RUN cpm install -g --show-build-log-on-failure
ENV PERL5LIB=/usr/local/lib/perl5
COPY . .

#CMD ["plackup", "app.psgi"]
