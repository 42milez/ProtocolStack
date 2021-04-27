#!/bin/bash

dnf -y update
dnf -y group install "Development Tools"

systemctl disable firewalld
