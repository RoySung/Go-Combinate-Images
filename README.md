# Go-Combinate-Images

> Merge Images With Every Combination

## Introduction

It is tool that merge images from different combination.

## How To Use Executable File

1. Download file from release

2. Create `assets/` and `assets/output` folders

3. Create `settings.json` file

4. Using command line to execute it

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

settings.json
```json
  {
    "folders": [
      "background",
      "cats"
    ]
  }
```

It is folders from assets that you want to merge images in here. It will merge images for the sort.
