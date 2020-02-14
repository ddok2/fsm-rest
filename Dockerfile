FROM ubuntu:18.04
MAINTAINER Wonkyu Lee <kurzweil@nuritelecom.com>
WORKDIR /usr/local/automation_tester
COPY main /usr/local/automation_tester
COPY configs /usr/local/automation_tester/configs
CMD ./main