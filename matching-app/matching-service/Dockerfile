FROM golang:1.17-alpine AS build_base

RUN apk --no-cache add ca-certificates git

WORKDIR /src
COPY go.mod ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./mainapp/

FROM scratch
COPY --from=build_base /bin/app /bin/app

EXPOSE 8081
ENTRYPOINT ["/bin/app"]