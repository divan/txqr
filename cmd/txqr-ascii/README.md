# TXQR-Ascii

`txqr-ascii` is a command-line tool for showing TXQR-encoded animated QR frames in terminal

It allows you to specify chunk size (`-split`, bytes per QR), and delay between frames (`-delay`). Size is adapted to the terminal width.

### Installation

```
go get github.com/divan/txqr/cmd/txqr-ascii

```

### Usage
To encode file `file.jpg` with 5fps and 450 bytes per QR frame, run:

```
txqr-ascii -split 450 -delay 200ms file.jpg
```


### Licence

MIT
