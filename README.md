# sf6vid

Censor player information in a Street Fighter 6 game play video.

- [Install](#install)
  - [Prerequisites / Dependencies](#prerequisites--dependencies)
- [Uninstall](#uninstall)

![](screenshot.png)


## Install

Download the release for your operating system from the [releases page](https://github.com/techygrrrl/sf6vid/releases). You can put this somewhere on your path, e.g. in macOS and Linux, you can put it in `/usr/local/bin` or `/usr/bin`.

If you use Go and would like to install it with `go install` you can do the following:

    go install github.com/techygrrrl/sf6vid@latest


### Prerequisites / Dependencies

This tool uses [ffmpeg](https://ffmpeg.org/) on your system so the `ffmpeg` and `ffprobe` commands must be available.


## Uninstall

Just remove the binary from wherever you installed it.

If you installed it with go, you can do `go uninstall sf6vid`.
