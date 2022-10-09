FROM golang:1.19.0-alpine3.15 as build
WORKDIR /app
COPY . .
RUN cd cmd/myapp && go build -o app

FROM alpine:3.15 as prod
COPY --from=build /app/cmd/myapp/app ./app
EXPOSE 3001
ENTRYPOINT ["/app"]
