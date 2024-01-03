FROM kasmweb/core-debian-bookworm:1.14.0
USER root
RUN apt-get update && \
  apt-get install -y wmctrl xdotool copyq
RUN wget https://dldir1.qq.com/qqfile/qq/QQNT/9681283b/linuxqq_3.2.3-20201_amd64.deb && \
  dpkg -i linuxqq_3.2.3-20201_amd64.deb; \
  apt install -y -f && \
  rm *.deb && \
  apt-get clean
COPY fuck-qq /opt/fuck-qq
USER 1000
ENV RUN_COPYQ=true
