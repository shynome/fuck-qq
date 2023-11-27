FROM kasmweb/core-debian-bookworm:1.14.0
USER root
RUN apt-get update && \
  apt-get install -y xclip wmctrl xdotool
RUN wget https://dldir1.qq.com/qqfile/qq/QQNT/fd2e886e/linuxqq_3.2.2-18394_amd64.deb && \
  dpkg -i linuxqq_3.2.2-18394_amd64.deb; \
  apt install -y -f && \
  rm linuxqq_3.2.2-18394_amd64.deb && \
  apt-get clean
USER 1000
