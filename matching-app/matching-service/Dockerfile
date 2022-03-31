FROM golang:1.12-alpine AS build_base

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./mainapp/

FROM scratch
COPY --from=build_base /bin/app /bin/app

EXPOSE 8081
ENTRYPOINT ["/bin/app"]