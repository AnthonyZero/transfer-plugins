#源镜像
FROM alpine:latest
#作者
MAINTAINER pingjin "736252868@qq.com"

# 在容器内设置 /app 为当前工作目录
WORKDIR /app

# 把执行文件 配置文件复制到当前工作目录
ADD main /app
ADD configs /app/configs

#ENV PROFILE pro
#EXPOSE 8080
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN echo "Asia/Shanghai" > /etc/timezone
# 执行可执行文件
ENTRYPOINT ["./main", "-env=dev"]