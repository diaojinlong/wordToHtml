#基于alpine版本镜像
FROM alpine:3.20

RUN apk add libreoffice libreoffice-lang-zh_cn libc6-compat && \
    apk cache clean && \
    mkdir /data && \
    mkdir /opt/libreoffice7.5 && \
    mkdir /opt/libreoffice7.5/program && \
    ln -s /usr/bin/soffice /opt/libreoffice7.5/program/soffice

#暴露端口
EXPOSE 8083

COPY wordToHtml /data

WORKDIR /data

CMD ["./wordToHtml"]