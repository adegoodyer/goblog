![GitHub last commit](https://img.shields.io/github/last-commit/adegoodyer/goblog)
![GitHub issues](https://img.shields.io/github/issues/adegoodyer/goblog)
![GitHub pull requests](https://img.shields.io/github/issues-pr/adegoodyer/goblog)
![Docker Pulls](https://img.shields.io/docker/pulls/adegoodyer/goblog)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/adegoodyer/goblog/dev)
![License](https://img.shields.io/github/license/adegoodyer/goblog)

## Overview
- fork of [mariomac/goblog](https://github.com/mariomac/goblog)
- tried following along to [Grafana Beyla Docker Setup](https://grafana.com/docs/beyla/latest/setup/docker/) but was getting container error [issue #6](https://github.com/mariomac/goblog/issues/6)
- PR raised [pull-request #7](https://github.com/mariomac/goblog/pull/7)
- enhanced Dockerfile to use multi stage build [pull-request #8](https://github.com/mariomac/goblog/pull/8)
- added security scan and SBOM

## Commands
```bash
# compile go
go build -o bin/goblog ./src

# execute binary
./goblog

# build container image
d build -t adegoodyer/goblog:dev .

# build container image (debug mode)
d build -t adegoodyer/goblog:dev . --progress=plain

# run container locally (insecure)
docker run -p 8080:8080 --name goblog adegoodyer/goblog:dev

# run container locally (tls)
# need to change src/assets/install/config.go -> HTTPSRedirect to true
# docker run -p 8443:8443 --name goblog adegoodyer/goblog:dev

# sec scan
grype adegoodyer/goblog:dev

# generate SBOM
syft adegoodyer/goblog:dev

# publish image
d logout && d login --username=adegoodyer
d push adegoodyer/goblog --all-tags
```

## Sec Scan
```bash
NAME        INSTALLED  FIXED-IN  TYPE  VULNERABILITY  SEVERITY
libcrypto3  3.3.2-r0   3.3.2-r1  apk   CVE-2024-9143  Medium
libssl3     3.3.2-r0   3.3.2-r1  apk   CVE-2024-9143  Medium
```

## SBOM
```
AME                                      VERSION                             TYPE
alpine-baselayout                         3.6.5-r0                            apk
alpine-baselayout-data                    3.6.5-r0                            apk
alpine-keys                               2.4-r1                              apk
apk-tools                                 2.14.4-r0                           apk
busybox                                   1.36.1-r29                          apk
busybox-binsh                             1.36.1-r29                          apk
ca-certificates-bundle                    20240705-r0                         apk
github.com/alecthomas/chroma/v2           v2.2.0                              go-module
github.com/caarlos0/env/v6                v6.10.1                             go-module
github.com/dlclark/regexp2                v1.7.0                              go-module
github.com/fsnotify/fsnotify              v1.7.0                              go-module
github.com/mariomac/goblog                (devel)                             go-module
github.com/mariomac/guara                 v0.0.0-20221222112709-f95b15506aee  go-module
github.com/yuin/goldmark                  v1.7.4                              go-module
github.com/yuin/goldmark-highlighting/v2  v2.0.0-20230729083705-37449abec8cc  go-module
golang.org/x/net                          v0.29.0                             go-module
golang.org/x/sys                          v0.25.0                             go-module
golang.org/x/tools                        v0.25.0                             go-module
gopkg.in/yaml.v2                          v2.4.0                              go-module
libcrypto3                                3.3.2-r0                            apk
libssl3                                   3.3.2-r0                            apk
musl                                      1.2.5-r0                            apk
musl-utils                                1.2.5-r0                            apk
scanelf                                   1.3.7-r2                            apk
ssl_client                                1.36.1-r29                          apk
stdlib                                    go1.23.3                            go-module
zlib                                      1.3.1-r1                            apk
```

## Acknowledgments

- [mariomac/goblog](https://github.com/mariomac/goblog)
- [Syft](https://github.com/anchore/syft)
- [Grype](https://github.com/anchore/grype)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
