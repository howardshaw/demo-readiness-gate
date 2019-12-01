FROM alpine

ADD bin/demo-readiness-gate /webhook
ENTRYPOINT ["/demo-readiness-gate"]
