FROM golang:1.16-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/app .
RUN chmod +x wait-for-it.sh

FROM alpine AS bin
COPY --from=build /out/app /src/wait-for-it.sh /bin/
RUN apk update
RUN apk upgrade
RUN apk add bash
CMD app
