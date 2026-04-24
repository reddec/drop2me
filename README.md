# drop2me

The dumbest possible software for the dumbest possible use case: you have a file on device A and need it on device B, and they're on the same network. That's it.

A 100% vibecoded (yet reviewed) piece of junk (yet useful) software.

No accounts. No cloud. No "sign up with Google." No 47MB Node.js runtime. Just a single binary that spawns a web page where you dump a file and it lands on disk. Revolutionary technology, truly.

## How it works

```
$ drop2me
Open one of these URLs to upload files:
  http://192.168.1.42:8080/

[QR code appears here, scan with your phone or whatever]

Listening on :8080
```

You open the URL, pick a file, hit upload. The file appears in the working directory. That's the entire feature set.

## Install

```
go install github.com/reddec/drop2me@latest
```

Or grab a binary from releases. Or build it yourself. It's one file, you'll survive.

## Usage

```
drop2me [flags]
```

| Flag        | Env                | Default | Description                         |
|-------------|--------------------|---------|-------------------------------------|
| `-bind`     | `DROP2ME_BIND`     | `:8080` | Listen address                      |
| `-dir`      | `DROP2ME_DIR`      | `.`     | Upload directory                    |
| `-max-size` | `DROP2ME_MAX_SIZE` | `0`     | Max upload size in bytes (0 = unlimited) |

```bash
# defaults
drop2me

# custom port and directory
drop2me -bind :3000 -dir ~/Downloads

# same but with env vars because you're that kind of person
DROP2ME_BIND=:3000 DROP2ME_DIR=~/Downloads drop2me
```

## Features

- Streams files to disk. Doesn't eat your RAM.
- Works in airgapped networks. Zero external dependencies in the UI.
- Mobile-friendly upload form. No JavaScript.
- Prints a QR code in the terminal so you don't have to type IPs on your phone.
- The binary is like 5MB. Your browser tab with ChatGPT open uses more memory.

## Docker

```bash
docker run --rm -p 8080:8080 -v $(pwd)/uploads:/data ghcr.io/reddec/drop2me

# run as non-root user
docker run --rm -p 8080:8080 -v $(pwd)/uploads:/data -u 1000:1000 ghcr.io/reddec/drop2me
```

## Why

Because sometimes `scp` is too much typing, `python -m http.server` doesn't upload, and setting up Nextcloud to move one PDF is a cry for help.
