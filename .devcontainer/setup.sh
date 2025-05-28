#!/bin/bash
set -e

sudo dnf -y update && \
sudo dnf -y install git postgresql procps-ng  && \
sudo dnf module enable -y nodejs:18 && \
sudo dnf install -y nodejs && \
sudo dnf clean all

# create ~/bin
mkdir -p ~/bin

# download and extract just to ~/bin/just
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ~/bin

# add `~/bin` to the paths that your shell searches for executables
# this line should be added to your shells initialization file,
# e.g. `~/.bashrc` or `~/.zshrc`
export PATH="$PATH:$HOME/bin"

# just should now be executable
just --help

# install golang-migrate CLI
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | sudo tar xvz -C /usr/local/bin migrate && \
sudo chmod +x /usr/local/bin/migrate