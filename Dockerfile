FROM alpine
WORKDIR /home
ADD . .
EXPOSE 80
CMD ./app