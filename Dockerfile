FROM openshift/origin-base

ENV GOPATH /go
RUN mkdir $GOPATH

COPY . $GOPATH/src/github.com/brancz/kube-rbac-proxy

RUN yum install -y golang make && \
   cd $GOPATH/src/github.com/brancz/kube-rbac-proxy && \
<<<<<<< HEAD
   make build && cp $GOPATH/src/github.com/brancz/kube-rbac-proxy/_output/linux/$(go env ARCH)/kube-rbac-proxy /usr/bin/ && \
=======
   make build && cp $GOPATH/src/github.com/brancz/kube-rbac-proxy/_output/linux/amd64/kube-rbac-proxy /usr/bin/ && \
>>>>>>> ffaa0a72658fe4458e62cdef26f5f13816175862
   yum erase -y golang make && yum clean all

LABEL io.k8s.display-name="kube-rbac-proxy" \
      io.k8s.description="This is a proxy, that can perform Kubernetes RBAC authorization." \
      io.openshift.tags="kubernetes" \
      maintainer="Frederic Branczyk <fbranczy@redhat.com>"

# doesn't require a root user.
USER 1001

ENTRYPOINT ["/usr/bin/kube-rbac-proxy"]
EXPOSE 8080
