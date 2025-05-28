FROM rockylinux:9

# Install dependencies
RUN dnf install -y golang && dnf clean all
# Install Git
RUN dnf install -y git && dnf clean all
# Install sudo
RUN dnf install -y sudo && dnf clean all

# Set PATH for Go
ENV PATH="/usr/local/go/bin:$PATH"


RUN useradd -ms /bin/bash vscode && \
    usermod -aG wheel vscode && \
    echo '%wheel ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
USER vscode
WORKDIR /workspace
CMD ["tail", "-f", "/dev/null"]
