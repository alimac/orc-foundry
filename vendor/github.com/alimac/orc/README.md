[![Build Status](https://travis-ci.org/alimac/orc.svg?branch=master)](https://travis-ci.org/alimac/orc)
[![Coverage Status](https://coveralls.io/repos/github/alimac/orc/badge.svg?branch=master)](https://coveralls.io/github/alimac/orc?branch=master)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/alimac/orc)

# Orc Library

Orc library generates Orc names, greetings, and weapons.

Inspired by [https://github.com/Pallinder/sillyname-go](sillyname-go) library.

Used in [https://github.com/alimac/orc-foundry](orc-foundry) project.

## Usage

``` go
package main

import (
	"fmt"

	"github.com/alimac/orc"
)

func main() {
	fmt.Printf("%s greets you with \"%s!\", holding a %s \n",
		orc.Forge("name"), orc.Forge("greeting"), orc.Forge("weapon"))
}
```

## Example output

```
Guldug greets you with "Time to die!", holding a DismalTrowel
```
