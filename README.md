## newest

Basic CLI tool to print the most recently modified file in a directory:

```
usage: newest [OPTIONS] [DIR]

Prints the newest file in a directory

If a positional argument is not provided, prints the newest file in this directory
If the current directory could not be determined, this fails

Optional arguments:
  -ignore-hidden
    	Ignore hidden files
  -include-dirs
    	Include directories in addition to files
```

### Install

Using `go get` to put it on your `$GOBIN`:

`go get github.com/seanbreckenridge/newest`

Manually:

```bash
git clone https://github.com/seanbreckenridge/newest
cd ./newest
go build .
# copy binary somewhere on your $PATH
sudo cp ./newest /usr/local/bin
```

### Examples

I use this in a lot of small scripts/aliases bound to key bindings:

To move the last file from my Downloads folder to my current directory:

`mv -v "$(newest "${HOME}/Downloads")" './'`

To preview my latest screenshot:

`sxiv -a "$(newest "${HOME}/Pictures/Screenshots")"`

To grab my most recent screenshot, and upload it to imgur:

`imgur-uploader "$(newest "${HOME}/Pictures/Screenshots")"`
