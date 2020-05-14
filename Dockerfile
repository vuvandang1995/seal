FROM golang:1.13 as build
WORKDIR /app
COPY ./ /app
RUN go build -o /dist/app ./main.go

FROM gcr.io/distroless/base:debug
COPY --from=build /dist/app /app
ENTRYPOINT ["/app"]
