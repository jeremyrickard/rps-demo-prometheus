FROM microsoft/dotnet-nightly:2.1-runtime-alpine3.7

RUN apk add libc6-compat
RUN mkdir /lf
ADD . /lf/
WORKDIR /lf/bin
RUN chmod +x ./Microsoft.LocalForwarder.ConsoleHost
EXPOSE 55678

# ENTRYPOINT ["./Microsoft.LocalForwarder.ConsoleHost", "noninteractive"]
ENTRYPOINT ["dotnet", "./Microsoft.LocalForwarder.ConsoleHost.dll", "noninteractive"]
