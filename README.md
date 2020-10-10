# What

Detect appropriate package manager by parsing os-release.

```console
❯ go build . && docker run -it --rm -v `pwd`:`pwd` -w `pwd` centos:7 ./go-os-parse
Detected package manager: "yum"

❯ go build . && docker run -it --rm -v `pwd`:`pwd` -w `pwd` fedora ./go-os-parse
Detected package manager: "dnf"

❯ go build . && docker run -it --rm -v `pwd`:`pwd` -w `pwd` debian ./go-os-parse
Detected package manager: "dpkg"
```

# Why

Avoid any dependencies on external commands including shell built-ins.

# References

- https://www.freedesktop.org/software/systemd/man/os-release.html
- https://github.com/tresf/whut
- https://github.com/sonatype-nexus-community/ahab/tree/master/docker
