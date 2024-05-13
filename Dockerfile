FROM golang:1.21 as BUILDER
ARG MERLIN_GOPROXY=https://goproxy.cn,direct
ARG GH_USER
ARG GH_TOKEN
RUN go env -w GOPROXY=${MERLIN_GOPROXY} && go env -w GOPRIVATE=github.com/openmerlin
RUN echo "machine github.com login ${GH_USER} password ${GH_TOKEN}" > $HOME/.netrc

# build binary
COPY . /go/src/github.com/openmerlin/cronjob
RUN cd /go/src/github.com/openmerlin/cronjob && GO111MODULE=on CGO_ENABLED=0 go build -buildmode=pie --ldflags "-s -linkmode 'external' -extldflags '-Wl,-z,now'"

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update --repo OS --repo update && \
    dnf in -y shadow --repo OS --repo update && \
    dnf remove -y gdb-gdbserver && \
    groupadd -g 1000 cronjob && \
    useradd -u 1000 -g cronjob -s /sbin/nologin -m cronjob && \
    echo > /etc/issue && echo > /etc/issue.net && echo > /etc/motd && \
    echo "umask 027" >> /root/.bashrc &&\
    echo 'set +o history' >> /root/.bashrc

USER cronjob
WORKDIR /home/cronjob

COPY  --chown=cronjob --from=BUILDER /go/src/github.com/openmerlin/cronjob/cronjob /home/cronjob

ARG MODE=release

ENTRYPOINT ["/home/cronjob/cronjob"]