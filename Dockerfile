FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine
COPY --from=build /bin/app /bin/app
ENTRYPOINT ["/bin/app"]