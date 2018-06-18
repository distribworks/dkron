FROM centos:6

ARG rpmfile
ADD $rpmfile /
RUN yum -y install /$rpmfile

CMD ["dkron", "agent"]
