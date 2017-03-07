FROM nginx:latest
MAINTAINER Yusuf

RUN apt-get update && apt-get install -y supervisor

RUN apt-get clean && rm -rf /var/cache/apt/archives/*

ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf

COPY . /app/

RUN mkdir /etc/nginx/templating

ADD sha1check.sh /usr/local/bin/nginx-reload
ADD nginxy /usr/local/bin/nginxy
ADD default.conf /etc/nginx/conf.d/default.conf
ADD conf.tmpl /etc/nginx/templating/conf.tmpl


RUN chmod u+x /usr/local/bin/nginx-reload
RUN chmod u+x /usr/local/bin/nginxy

WORKDIR /app/

RUN chown nginx /var/log/nginx/*

ENV DOCKER_HOST unix:///tmp/docker.sock

VOLUME ["/etc/nginx/certs"]

ENTRYPOINT ["/usr/bin/supervisord"]
