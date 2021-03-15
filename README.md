# Go-Combinate-Images

> Merge Images With Every Combination

## Introduction

It is tool that merge images from different combination.

## Build

```sh
  # build windows
  GOOS=windows GOARCH=amd64 go build -tags release .
  # build mac
  GOOS=darwin GOARCH=amd64 go build -tags release .

  # others refs https://golang.org/doc/install/source#environment
```

## Run

Output files to ./assets/output .

## Setting Example

```json
  {
    "folders": [
      "background",
      "cats"
    ]
  }
```

It is folders from assets that you want to merge images in here. It will merge images for the sort.