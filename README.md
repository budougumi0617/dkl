# dkl
Pure Go implementation of [dtt][dtt].

## Description
Pure Go implementation of [dtt][dtt].  
dkl is the text-mode interface for docker and kubectl command.

## Demo & VS.

## Requirement

## Usage

```
$ dkl -h
Usage:
    dkl
    dkl -k | --kubectl
    dkl -h | --help
    dkl -c | --config
Options:
    -k --kubectl             kubectl mode
    -h --help                Show this screen and exit.
    -c --config              Show config
```

## Install
You can download binary from [release page](https://github.com/budougumi0617/dkl/releases) and place it in $PATH directory.

### MacOS
If you want to install on MacOS, you can use Homebrew.
```
brew tap budougumi0617/tap
brew install budougumi0617/dkl
```


## Contribution
1. Fork ([https://github.com/budougumi0617/dkl/fork](https://github.com/budougumi0617/dkl/fork)
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## License

[MIT](https://github.com/budougumi0617/dtt-go/blob/master/LICENSE)

[dtt]: https://github.com/ymizushi/dtt
