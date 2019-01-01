FROM golang:1.8
RUN mkdir /app
WORKDIR /app
COPY pystuff pystuff
COPY main.go main.go
RUN go get -v -t -d
ENV PYTHONUNBUFFERED 1

CMD ["go", "run", "main.go"]
