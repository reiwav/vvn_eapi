FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y wget git gcc

RUN wget -P /tmp https://dl.google.com/go/go1.17.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf /tmp/go1.17.linux-amd64.tar.gz
RUN rm /tmp/go1.17.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install tzdata libc6-dev gcc libc-dev g++ libffi-dev libxml2 unixodbc-dev unixodbc-dev
COPY ./binary/libtbodbc.so /usr/lib/odbc/
COPY ./binary/odbc.ini ./binary/odbcinst.ini /etc/

WORKDIR /go/src/app/eapi

COPY . /go/src/app/eapi

RUN go build -o ./eapirun .
#RUN 
RUN chmod +x ./eapirun
EXPOSE 3000



#FROM node:14.17.6

WORKDIR /go/src/app/eapi/wooribank_cms/

COPY ./wooribank_cms/package.json ./wooribank_cms/yarn.lock ./
RUN apt-get update
RUN apt-get install -y nodejs npm
RUN npm install
COPY ./wooribank_cms ./
# COPY . .
RUN npm run build --prod

RUN mv ./dist/ /go/src/app/eapi/admin

# # FROM nginx:1.21.3

# # COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf
#COPY --from=build /go/src/app/eapi/wooribank_cms/dist /go/src/app/eapi/admin


# # EXPOSE 80
WORKDIR /go/src/app/eapi
ENTRYPOINT  ["./eapirun"]


