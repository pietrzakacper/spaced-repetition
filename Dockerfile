FROM golang:1.21.7

COPY ./controller /app/controller
COPY ./flashcard /app/flashcard
COPY ./persistance /app/persistance
COPY ./user /app/user
COPY ./web /app/web
COPY go.work /app
COPY go.work.sum /app

WORKDIR /app

RUN ls -la
RUN go mod download

RUN go build -o ./bin/spaced-repetition web
RUN cp -r ./web/templates ./bin/templates
RUN cp -r ./web/static ./bin/static

WORKDIR /app/bin
ENV PORT 8080
EXPOSE 8080
CMD ["./spaced-repetition"]