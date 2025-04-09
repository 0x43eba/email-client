bootstrap-db:
    bash bootstrap.sh

run: bootstrap-db
    docker build -t demoproject .
    docker run --rm -it demoproject
