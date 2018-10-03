FROM alpine:3.7

RUN apk add --no-cache libc6-compat
RUN apk add --no-cache libstdc++
RUN mkdir /lf
ADD . /lf/
WORKDIR /lf/bin
RUN chmod +x ./Microsoft.LocalForwarder.ConsoleHost
EXPOSE 55678

# ENTRYPOINT ["ldd", "./Microsoft.LocalForwarder.ConsoleHost"]
ENTRYPOINT ["./Microsoft.LocalForwarder.ConsoleHost", "noninteractive"]
