# mediascan

A simple and fast Go (golang) command-line utility to recursively scan a directory for media files, extract metadata (including ID3v2 tags from both MP3 and M4A files), and save the output in a simple yaml format e.g. [files.yaml](files.yaml). 

- Reads configuration from yaml file e.g. [conf.yaml](conf.yaml)
- Has only two required command-line arguments: `{config yaml filepath}` and `{output yaml filepath}`
- Created specifically to run fast on a Raspberry Pi single-board computer as part of my other (Python) project: [timebox](https://github.com/bretttolbert/timebox)

## Dependencies (Go Modules)
- [tag](https://github.com/dhowden/tag) (Used for reading ID3 tags)
- [yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) (Used for generating files.yaml)
- [mp3](github.com/tcolgate/mp3) (Used for calculating MP3 duration)

### Installing Go modules

```bash
cd $GOPATH
mkdir -p src/github.com/tcolgate
cd src/github.com/tcolgate
git clone git@github.com:tcolgate/mp3.git
cd mp3
go install
```

## Performance Demo 1 - desktop PC with decade old i7 CPU
```bash
$ time go run mediascan.go conf.yaml files.yaml
...
2022/07/02 17:19:53 Successfully loaded 8376 media files
2022/07/02 17:19:53 Skipped 179 excluded media files

real	0m2.852s
user	0m1.832s
sys	0m1.232s
```

## Performance Demo 1.5 - desktop PC with decade old i7 CPU
Much slower on first run. Why?

First run:
```
2022/12/14 06:55:32 Successfully loaded 9124 media files
2022/12/14 06:55:32 Skipped 179 excluded media files

real	3m6.427s
user	0m4.687s
sys	0m5.567s

Second run:
```
2022/12/14 07:09:53 Successfully loaded 9124 media files
2022/12/14 07:09:53 Skipped 179 excluded media files

real	0m3.637s
user	0m2.410s
sys	0m1.299s

```


## Performance Demo 2 - Raspberry Pi 4 model B
```bash
$ time go run mediascan.go conf.yaml files.yaml
2022/07/04 15:06:39 Successfully loaded 8376 media files

real	0m18.652s
user	0m8.251s
sys	0m11.406s
```

## Performance Demo 3 - Raspberry Pi Zero 2W
```bash
$ time go run mediascan.go conf.yaml files.yaml
2022/07/04 15:21:14 Successfully loaded 8376 media files

real	1m18.217s
user	0m18.136s
sys	0m14.305s
```

## Configuration YAML Reference

| Property | Description |
| -------- | ----------- |
| `mediadir` | the path to the directory to be scanned |
| `mediaexts` | the file extensions to be included in the scan |
| `excludepath` | path substrings to exclude from scan (case-insensitive) |
| `excludetitle` | id3 title substrings to exclude from scan results (case-insensitive) |
| `excludetitle` | id3 artist substrings to exclude from scan results (case-insensitive) |
| `excludealbum` | id3 album substrings to exclude from scan results (case-insensitive) |
| `excludegenre` | id3 genre substrings to exclude from scan results (case-insensitive) |
| `sortby` | (`year`, `artist`, `none`) media file sort options |
| `groupby` | (`none`, `year` ) group media files into playlists using the specified tag |

## Installing Dependencies

Install Go programming language compiler, linker, compiled stdlib and supplementary Go tools
```bash
sudo apt update
sudo apt upgrade
sudo apt install golang-go golang-src golang-doc golang-golang-x-tools
```

Edit `~/.bashrc` by adding the following lines:

```bash
export GOPATH=~/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
```
Then reload it: `source ~/.bashrc`

Note: `GOROOT` directory may vary. To confirm, check the `go` symlink, e.g.:
```bash
$ which go
/usr/bin/go
$ ls -l /usr/bin/go
lrwxrwxrwx 1 root root 21 Sep 16  2020 /usr/bin/go -> ../lib/go-1.15/bin/go
```
So based on the above output on the rpi we need: `export GOROOT=/usr/lib/go-1.15`



To test it (Raspberry Pi):
```bash
$ $GOROOT/bin/go version
go version go1.15.15 linux/arm64
```

To test it (Linux desktop):
```bash
$ $GOROOT/bin/go version
go version go1.18.1 linux/amd64
```

Install Go dependency 1 of 2: `tag`
```bash
go get github.com/dhowden/tag/cmd/tag
```

Install Go dependency 2 of 2: `yaml.v3`
```bash
go get gopkg.in/yaml.v3
```

Now we can see that these packages were installed into `$GOPATH/src`:
```bash
$ ls $GOPATH/src
github.com  gopkg.in
```

Now install mediascan:
```bash
go get github.com/bretttolbert/mediascan
cd ~/go/src/github.com/bretttolbert/mediascan/ && go install
```

Edit [conf.yaml](conf.yaml) and set the `mediadir` to the desired directory path.

Now you should be all set to run mediascan.

You can run it from the mediascan directory like this:
```bash
go run mediascan.go conf.yaml files.yaml
```

Or you can run it from any directory like this:
```bash
go run github.com/bretttolbert/mediascan conf.yaml files.yaml
```

Happy scanning!
