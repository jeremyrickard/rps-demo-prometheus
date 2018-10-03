SETLOCAL

SET APPINSIGHTS_INSTRUMENTATIONKEY=879547bf-f2f4-48f5-9d9b-5e7e48fb1cc8
docker-compose up
docker rm local-forwarder
docker image list
docker rmi local-forwarder-img
