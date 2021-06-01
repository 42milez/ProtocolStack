#!/bin/bash

dnf -y update
dnf -y group install "Development Tools"
dnf -y install tcpdump

systemctl disable firewalld

readonly BRIDGE_NAME=br0
readonly ETH_NAME=eth0
readonly TAP_NAME=tap0

readonly BRIDGE_IP=192.168.42.42/24
readonly TAP_IP=192.0.2.1/24

nmcli connection add type tun ifname "${TAP_NAME}" con-name "${TAP_NAME}" mode tap ip4 "${TAP_IP}"

nmcli con add type bridge ifname "${BRIDGE_NAME}"
nmcli con mod "bridge-${BRIDGE_NAME}" bridge.stp no
nmcli con mod "bridge-${BRIDGE_NAME}" ipv4.method manual ipv4.address "${BRIDGE_IP}"

nmcli con add type bridge-slave ifname ${ETH_NAME} master "bridge-${BRIDGE_NAME}"
nmcli con add type bridge-slave ifname ${TAP_NAME} master "bridge-${BRIDGE_NAME}"
