FROM microsoft/dotnet:2.1-sdk

RUN mkdir /lf
ADD . /lf/

RUN chmod +x /lf/LF-ConsoleHost/Microsoft.LocalForwarder.ConsoleHost

EXPOSE 55678

CMD bash -c 'cd /lf/LF-ConsoleHost && ./Microsoft.LocalForwarder.ConsoleHost noninteractive'

