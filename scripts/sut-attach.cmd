:: attach to the system under test, only works if you define a CMD instead of an ENTRYPOINT
docker network create restaurant-dev-net

call build-main.cmd

docker compose -f sut/docker-compose.yml build

docker run -v %cd%\sut\run/:/app/run/ -v %cd%\sut\run\generated/:/app/run/generated/ -t -i --name restaurant-document-generate-function restaurant-document-generate-function bash

docker rm restaurant-document-generate-function

docker image rm restaurant-document-generate-function

pause