FROM fedora:33

RUN \
    dnf -y update && \
    dnf -y group install "Development Tools" && \
    dnf clean all

ENV APPROOT /var/app

WORKDIR $APPROOT
