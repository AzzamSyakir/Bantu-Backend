# docker command
start-docker:
	clear && docker compose --env-file ./.env -f ./docker/docker-compose.yml up -d 

stop-docker:
	clear && docker compose --env-file ./.env -f ./docker/docker-compose.yml down --remove-orphans

clean-docker:
	clear && docker system prune && docker volume prune && docker image prune -a --env-file ./.env -f && docker container prune && docker buildx prune

start-db:
	clear && docker compose --env-file ./.env -f ./docker/docker-compose.yml up user-db product-db order-db auth-db -d

start-test:
	clear && go test -v -count=1 ./src/gateway-service/test

generate-proto:
	clear && protoc --proto_path=grpc/proto grpc/proto/*.proto --go_out=grpc --go-grpc_out=grpc
