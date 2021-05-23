# Installation

To start using Open Integration we reccomend to use oictl to create and manage your pipelines.

## Download

### Homebrew
```
brew tap open-integration/oi
brew install oictl
```

### Linux
Download the binary from the [release](https://github.com/open-integration/oi/releases) page:
```
V=$(curl https://raw.githubusercontent.com/open-integration/oi/master/VERSION)

curl https://github.com/open-integration/oi/releases/download/v$V/oi_$V_Linux_x86_64.tar.gz
```