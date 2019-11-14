FROM gitpod/workspace-full

USER root
COPY --from=lachlanevenson/k8s-kubectl:v1.16.3 /usr/local/bin/kubectl /usr/local/bin/kubectl
RUN free -h
# Install custom tools, runtime, etc. using apt-get
# For example, the command below would install "bastet" - a command line tetris clone:
#
# RUN apt-get update \
#    && apt-get install -y bastet \
#    && apt-get clean && rm -rf /var/cache/apt/* && rm -rf /var/lib/apt/lists/* && rm -rf /tmp/*
#
# More information: https://www.gitpod.io/docs/42_config_docker/
