FROM golang:1-alpine AS builder

# https://github.com/hashicorp/terraform/releases
ARG TF_VERSION=1.1.7

RUN apk --update --no-cache add git

# build tfconsole
WORKDIR /go/src/github.com/lingrino/tfconsole
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/tfconsole

# install terraform
RUN wget -q https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip && \
    unzip terraform_${TF_VERSION}_linux_amd64.zip && rm terraform_${TF_VERSION}_linux_amd64.zip && \
    mv terraform /bin/terraform

FROM scratch

COPY --from=builder /bin/terraform /bin/
COPY --from=builder /bin/tfconsole /bin/

COPY --from=builder /go/src/github.com/lingrino/tfconsole/static    /static
COPY --from=builder /go/src/github.com/lingrino/tfconsole/templates /templates

ENTRYPOINT ["tfconsole"]
