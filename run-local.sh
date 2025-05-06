docker rm -f $(docker ps -a -q)
docker volume rm $(docker volume ls -q)

# cd src/
# go mod vendor
# tar -czvf vendor.tar.gz vendor/
# cd ..

docker compose up --build

# docker-compose -f docker-compose.yml up -d
