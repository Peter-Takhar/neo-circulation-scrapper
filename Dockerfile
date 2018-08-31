FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y \
        curl           \
        dpkg           \
        wget           \
        unzip          \
        libleveldb-dev \
        sqlite3        \
        libsqlite3-dev \
        libunwind8-dev \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q https://packages.microsoft.com/config/ubuntu/18.04/packages-microsoft-prod.deb && \
        dpkg -i packages-microsoft-prod.deb

RUN apt-get install -y apt-transport-https && \
        apt-get update && \
            apt-get install -y aspnetcore-runtime-2.1

RUN wget -q https://github.com/neo-project/neo-cli/releases/download/v2.8.0/neo-cli-linux-x64.zip && \
        unzip neo-cli-linux-x64.zip -d /opt


ENTRYPOINT ["/usr/bin/dotnet", "/opt/neo-cli/neo-cli.dll", "/rpc"]
