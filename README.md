# Install swag tool
go install github.com/swaggo/swag/cmd/swag@latest

# Build Swagger Documentation
swag init -g ./cmd/main.go -o ./docs

# Create MongoDb docker container
docker run -d \                                                                                                                                    ✔
--name mongodb \
-e MONGO_INITDB_ROOT_USERNAME=admin \
-e MONGO_INITDB_ROOT_PASSWORD=password \
-p 27017:27017 \
mongo:4.4.3

