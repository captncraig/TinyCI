CID=$(docker build docker/build)
docker cp $CID:/go/bin/tinyci docker/tinyci