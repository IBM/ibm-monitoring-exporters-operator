FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/ibm-monitoring-exporters-operator \
    USER_UID=1000600001 \
    USER_NAME=ibm-monitoring-exporters-operator

# install operator binary
COPY build/_output/bin/ibm-monitoring-exporters-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
