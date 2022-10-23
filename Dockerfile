# syntax=docker/dockerfile:1

## Build
FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY pkg ./pkg
COPY *.go ./

ARG VERSION="0.0.0.0"
ARG GIT_COMMIT="unknown"

RUN go build -v -o /poiskzooCrawler -ldflags "-X 'github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version.GitCommit=$GIT_COMMIT' -X 'github.com/LostPetInitiative/poiskzoo-ru-crawler/pkg/version.AppVersion=$VERSION'"

## Deploy
FROM ubuntu as final
# for some reason the https certs are outdated 🤷‍♂️
RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates
#RUN apt-get update && apt-get install --no-install-recommends -y libxss1
#RUN apt-get update && apt-get upgrade --no-install-recommends -y
WORKDIR /
COPY --from=build /poiskzooCrawler /poiskzooCrawler

# ENV CARDS_DIR=xxxx
# ENV PIPELINE_NOTIFICATION_URL=xxx

CMD ["/poiskzooCrawler"]
