FROM gcr.io/distroless/static-debian11

ENTRYPOINT ["/ec2id"]
COPY ec2id /ec2id
