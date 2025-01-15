FROM golang:1.23.4
WORKDIR /FlickMeter
RUN go install github.com/air-verse/air@latest
COPY go.* ./
RUN go mod download
COPY . .
CMD [ "air", "-c", ".air.toml" ]
