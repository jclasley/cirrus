FROM golang
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN make build
EXPOSE 8080
CMD ["./server"]
