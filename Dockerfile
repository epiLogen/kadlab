FROM larjim/kademlialab:latest

RUN apt-get update
RUN apt-get -y upgrade

RUN apt-get install -y \
net-tools inetutils-traceroute apt-utils \
iputils-ping xinetd telnetd