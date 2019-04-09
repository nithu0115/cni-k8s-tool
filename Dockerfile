FROM scratch
COPY cni-k8s-tool .
CMD ["/cni-k8s-tool"]
