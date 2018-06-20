FROM centos:7

ARG rpmfile
ADD $rpmfile /
RUN yum -y install /$rpmfile

CMD ["dkron", "agent"]
