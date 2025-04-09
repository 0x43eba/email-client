FROM golang:1.23

RUN apt-get update && apt-get install -y build-essential && rm -rf /var/lib/apt/lists/*


WORKDIR /app

COPY go.mod go.sum ./

COPY . .

CMD [ "go", "run", "main.go" ]
