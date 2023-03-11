FROM golang:1.19

COPY ./controller /app/controller
COPY ./csv /app/csv
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
EXPOSE 3000
CMD ["./spaced-repetition"]