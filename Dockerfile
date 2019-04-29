FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y \
        curl           \
        screen         \
        dpkg           \
        wget           \
        unzip          \
        libleveldb-dev \
        sqlite3        \
        libsqlite3-dev \
        libunwind8-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q https://packages.microsoft.com/config/ubuntu/18.04/packages-microsoft-prod.deb && \
        dpkg -i packages-microsoft-prod.deb

RUN apt-get install -y apt-transport-https && \
        apt-get update && \
            apt-get install -y aspnetcore-runtime-2.1

WORKDIR /opt
RUN wget -q https://github.com/neo-project/neo-cli/releases/download/v2.8.0/neo-cli-linux-x64.zip && \
        unzip neo-cli-linux-x64.zip && \
        rm neo-cli-linux-x64.zip

WORKDIR /opt/neo-cli
RUN chmod u+x ./neo-cli

EXPOSE 10332 10333

ENTRYPOINT ["screen"]
CMD ["-S", "neo","./neo-cli", "/rpc"]
