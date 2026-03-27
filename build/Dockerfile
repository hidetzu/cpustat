FROM golang:1.23 AS build

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /out/cpustat .

FROM build AS bin
COPY --from=build /out/cpustat /cpustat
