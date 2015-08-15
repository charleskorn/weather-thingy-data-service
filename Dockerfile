FROM busybox
COPY weather-thingy-data-service-amd64-linux /weather-thingy-data-service
CMD ["/weather-thingy-data-service"]
