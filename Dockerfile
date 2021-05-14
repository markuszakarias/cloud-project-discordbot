# Start from golang base image
FROM golang:latest

ENV WEATHER_KEY=f6a8e67b1a5f1d5be2bffe4d461cc155
ENV NEWS_KEY=03b8fc7d5add4ac98eb2330004fbb45c
ENV MEALS_KEY=eeb5e8160efb4bedb1ccc4aa441b0102
ENV DB_SERVER=vmdata.database.windows.net
ENV DB_PORT=1433
ENV DB_USER=eriksen
ENV DB_PASSWORD=Tanzania1994!
ENV DB=VM_Data
ENV DC_TOKEN=ODM2OTgzNjUyMjUxMzM2Nzc1.YIl7xQ.cuxQXG5lW9Sqmylm6rx4INNiLpc

# Setting current working directory
WORKDIR /build

# Caching all dependencies by downloading them so we dont have to download them every time we build image
COPY go.mod ./
COPY go.sum ./

# Downloading all dependencies
RUN go mod download

# Copying the source code
COPY . .

# Build the application
RUN go build -o ./main .

# Run the executable
CMD ["./main"]


