# create-memo

## summary

Creates a file with the current timestamp while simultaneously archiving files older than one month and deleting empty files.

## Install

`go install github.com/hxrxchang/create-memo/cmd/cm@latest`

## usage

### Create a file in ~/memo with the default extension (md)

```sh
create-memo
```

### Create a file with the specified extension (txt)

```sh
create-memo -ext txt
```

### Specify the directory to save the file

```sh
create-memo -path ~/Documents/memos
```

### Create a file in the specified directory with the txt extension

```sh
create-memo -path ~/Documents/memos -ext txt
```
