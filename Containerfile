FROM docker.io/golang
ADD . /sources
RUN cd /sources && go build && go install
CMD code-gen