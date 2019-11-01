# tson
`tson` is JSON viewer and editor written in Go.
This tool displays JSON as a tree and you can search and edit key or values.

![](https://i.imgur.com/tBGLEsT.gif)

## Support OS
- Mac
- Linux

## Installation
```sh
$ git clone https://github.com/skanehira/tson
$ cd tson && go install
```

## Usage
```sh
# fromstdin
$ tson < test.json

# from url(only run http get)
$ tson -url http://gorilla/likes/json
```

## Keybinding
### JSON tree

| key    | description         |
|--------|---------------------|
| j      | move down           |
| k      | move up             |
| g      | move to the top     |
| G      | move to the bottom  |
| ctrl-f | page up             |
| ctrl-b | page down           |
| h      | hide current node   |
| H      | collaspe all nodes  |
| l      | expand current node |
| L      | expand all nodes    |
| r      | read from file      |
| s      | save to file        |
| Enter  | edit node           |
| /      | search nodes        |

# Author
skanehira
