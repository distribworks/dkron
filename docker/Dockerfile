FROM docker-registry.qiyi.virtual/adfe/adfe_python3:latest
RUN mkdir -p /data/apps/dkron/bin
RUN mkdir -p /data/apps/dkron/data
RUN curl --output /data/apps/dkron/bin/dkron "http://jfrog.cloud.qiyi.domain/artifactory/iqiyi-pypi-adssrv/dkron"
RUN chmod +x /data/apps/dkron/bin/dkron
COPY ./deploy/dkron.json /data/apps/dkron
RUN yum install -y libev-devel supervisor
COPY ./deploy/supervisor-cron-app.ini /etc/supervisord.d/