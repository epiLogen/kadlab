FROM larjim/kademlialab:latest



RUN apt-get update
RUN apt-get -y upgrade

RUN apt-get install -y \
net-tools inetutils-traceroute apt-utils \
iputils-ping xinetd telnetd golang-go

RUN mkdir /home/go/src/kadlab
RUN export GOROOT=/home/go/
ADD . /home/go/src/kadlab
WORKDIR /home/go/src/kadlab/
