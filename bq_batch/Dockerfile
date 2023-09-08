FROM php:7.4-cli

RUN apt-get update \ 
&& apt install -y curl \ 
dnsutils \
iputils-ping \
netcat \
net-tools \
telnet \
libxml2-dev procps \
&& docker-php-ext-install pdo_mysql simplexml pcntl

WORKDIR /home/application/bq_batch
COPY ./ /home/application/bq_batch
