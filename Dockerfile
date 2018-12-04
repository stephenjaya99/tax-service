FROM ubuntu

RUN apt-get update
RUN apt-get install -y curl 
RUN curl -sL https://deb.nodesource.com/setup | sudo bash -
RUN apt-get update
RUN apt-get install -y build-essential golang