# tson
`tson` is JSON viewer written in Go.
This tool displays JSON as a tree and you can search and edit key or values.

![](https://i.imgur.com/tBGLEsT.gif)

## Support OS
- Mac
- Linux

## Installtion
```sh
$ git clone https://github.com/skanehira/tson
$ cd && go install
```

## Usage
```sh
$ tson < test.json
```

## keybinding
### json tree

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
| Enter  | edit node           |
| /      | search nodes        |

# Author
skanehira
