#1. butuh base image golang
FROM golang:1.19 as baseGolang

#2. copas semua kodingan project ke dalam (wadah) dockerfile
WORKDIR /app
COPY . .

#3. ngebuild jadi binary
RUN go build -tags netgo -o main.app .
#CMD ["/app/main.app"]

#4. butuh wadah ke 2 lagi (base image) yg lebih kecil
FROM alpine:latest

#5. copas hasil build an ke wadah ke 2
COPY --from=baseGolang /app/main.app .

#6. run
CMD ["/main.app"]