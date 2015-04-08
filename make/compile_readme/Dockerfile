FROM ubuntu:12.04
MAINTAINER oliveagle@gmail.com

COPY ./sources.list /etc/apt/sources.list
MKDIR /root/.pip
COPY ./pip.conf /root/.pip/

RUN apt-get update && apt-get install -y python-pip && pip install grip && apt-get autoclean

# let grip do caching
RUN echo "test" > /tmp/test.md && grip --export /tmp/test.md /tmp/test.html && rm /tmp/test.*
