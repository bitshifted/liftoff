# Copyright 2024 Bitshift D.O.O
# SPDX-License-Identifier: MPL-2.0

FROM ubuntu:24.04

ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

ENV ANSIBLE_PACKAGE_VERSION=11.9.0-1ppa~noble 
ENV TF_VERSION=1.10.3
ENV TF_URL=https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip

# Labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="bitshifted/liftoff"
LABEL org.label-schema.description="Tool for repid infrastructure deployment to cloud"
LABEL org.label-schema.url="https://github.com/bitshifted/liftoff"
LABEL org.label-schema.vcs-url="https://github.com/bitshifted/liftoff"
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.vendor="WSO2"
LABEL org.label-schema.version=$BUILD_VERSION

# Install required components for Ansible
RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt update && \
    apt install --no-install-recommends -y unzip wget ca-certificates software-properties-common && \
    add-apt-repository --yes --update ppa:ansible/ansible && \
    apt install --no-install-recommends -y ansible=${ANSIBLE_PACKAGE_VERSION} && \
    apt clean all
RUN mkdir /workspace

# Install Terraform
RUN wget ${TF_URL} && unzip terraform_${TF_VERSION}_linux_amd64.zip && \
    mv terraform /usr/bin && \
    rm terraform_${TF_VERSION}_linux_amd64.zip


ARG binary_location=target/linux-amd64/liftoff

COPY ${binary_location} /usr/bin
RUN chmod 755 /usr/bin/liftoff

WORKDIR /workspace
ENTRYPOINT [ "/usr/bin/liftoff" ]
