# TXQR-Gif

`txqr-gif` is a command-line tool for generating TXQR-encoded animated GIFs.

It allows you to specify chunk size (`-split`, bytes per QR), output image size (`-size`) and delay between frames (`-delay`). Output is written into `out.gif` file (see `-out` flag).

### Installation

```
go get github.com/divan/txqr/cmd/txqr-gif

```

### Usage
To encode file `file.jpg` with 5fps, image size 600px and 450 bytes per QR frame, run:

```
txqr-gif -size 600 -split 450 -delay 200ms file.jpg
```


### Licence

MIT
