# uinputd

uinput daemon

## Usage

Run uinputd with specifying configuration file like below.

```
$ uinputd -c ~/.config/uinputd
```

A format of configuration file is:

```
device: XX:XX:XX:XX:XX:XX
press: echo PRESS
release: echo RELEASE
longpress: echo LONGPRESS
```

## Installation

```
$ go get github.com/mattn/uinputd
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
