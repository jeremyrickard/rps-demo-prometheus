FROM microsoft/dotnet-nightly:2.1-runtime-alpine3.7

RUN apk add libc6-compat
RUN mkdir /lf
ADD . /lf/
WORKDIR /lf/bin

EXPOSE 55678

ENTRYPOINT ["dotnet", "./Microsoft.LocalForwarder.ConsoleHost.dll", "noninteractive"]
