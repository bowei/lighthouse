FROM debian
RUN apt-get update && apt-get install -y tcpdump curl
COPY bin/lightc /lightc
