FROM ubuntu:16.04

ARG debfile
ADD $debfile /
RUN dpkg -i /$debfile

CMD ["dkron", "agent"]
