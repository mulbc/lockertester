FROM alpine
ARG main
ENV MAIN ${main}
COPY ${main} /app/
WORKDIR /app
ENTRYPOINT ./${MAIN}