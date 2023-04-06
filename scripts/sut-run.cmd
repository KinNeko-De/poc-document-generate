:: starts the system under test
docker network create restaurant-dev-net

call build-main.cmd

docker compose -f sut/docker-compose.yml up --build --remove-orphans --exit-code-from restaurant-document-generate-function

docker compose -f sut/docker-compose.yml down

docker image rm restaurant-document-generate-function

pause