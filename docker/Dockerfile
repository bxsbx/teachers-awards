FROM golang:1.20 AS gobuild
ARG APP_ROOT
ARG APP_NAME
ARG TAG_NAME
ARG CONFIG_SERVER
ARG CONFIG_ENV
ENV GOPROXY https://goproxy.cn,direct
WORKDIR ${APP_ROOT}
COPY ./ ./
COPY docker/run.sh ./
RUN go mod tidy
RUN printf "\n  ${CONFIG_ENV}: ${TAG_NAME}" >> ./.config.yml
RUN printf "\nenv: ${CONFIG_ENV}" >> ./.config.yml
RUN curl "${CONFIG_SERVER}" -fd "`cat ./.config.yml`" > ./conf/dev/configs
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP_NAME} main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate migrate.go

FROM registry-vpc.cn-shenzhen.aliyuncs.com/ebag/dream-ebag-new-alpine:v3
ARG APP_ROOT
ARG EXPOSE_PORT
ARG DREAM_ENV
ARG APP_NAME
WORKDIR /app
ENV DREAMENV ${DREAM_ENV}
ENV APPNAMEENV ${APP_NAME}
EXPOSE ${EXPOSE_PORT}
USER root
RUN mkdir -p ./migration
RUN mkdir -p ./conf/dev
COPY --from=gobuild ${APP_ROOT}/${APP_NAME} ./${APP_NAME}
COPY --from=gobuild ${APP_ROOT}/migrate ./migrate
COPY --from=gobuild ${APP_ROOT}/run.sh ./run.sh
COPY --from=gobuild ${APP_ROOT}/conf/ ./conf/
COPY --from=gobuild ${APP_ROOT}/migration/ ./migration/
RUN linux-decrypter64 ./conf/dev/app.conf "./conf/dev/configs"
RUN chmod +x ./run.sh
RUN chmod +x ./${APP_NAME}
RUN chmod +x ./migrate
ENTRYPOINT ["/bin/sh","-c","/app/./run.sh"]