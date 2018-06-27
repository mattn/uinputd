# uinputd

uinput daemon

## Usage

Run uinputd with specifying configuration file like below.

```
$ uinputd
```

A format of configuration file is:

```
- device: XX:XX:XX:XX:XX:XX
  release: echo PLAY NEXT

- device: YY:YY:YY:YY:YY:YY
  release: echo PUSHED BUTTON
```

Configuration file is located at `~/.config/uinput`.

## Installation

```
$ go get github.com/mattn/uinputd
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
