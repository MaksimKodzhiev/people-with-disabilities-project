FROM amd64/golang:1.22

WORKDIR /project

EXPOSE 80

RUN adduser --disabled-password --gecos "" backend && \
    echo "backend:m8L4fFJr0szPMyYD5RC2iTA7HvNxt1Od" | chpasswd && \
    chmod u-s $(which passwd) && \
    chown backend:backend /project && \
    chmod u+rwx /project

COPY --chown=backend:backend --chmod=600 . .

RUN go mod tidy && \
    go mod download && \
    go mod verify && \
    go build -v -o /server /project/cmd/backend/main.go && \
    chown backend:backend /server && \
    chmod u+rwx /server

USER backend

CMD /server
