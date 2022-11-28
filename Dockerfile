FROM golang as build_image

WORKDIR /build

COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o eblog main.go \
    && cp conf/app.example.conf conf/app.conf \
    && mkdir build_finsh \
    && cp -r conf build_finsh/conf \
    && cp eblog build_finsh/eblog \
    && cp -r profile build_finsh/profile \
    && cp -r static build_finsh/static \
    && cp -r views build_finsh/views 



FROM ubuntu

# VOLUME [ "/app/conf",'/app/static/upload' ]

EXPOSE 8080

WORKDIR /app

COPY --from=build_image /build/build_finsh /app

CMD ./eblog