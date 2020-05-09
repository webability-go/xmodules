[ ![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xmodules/clientlink)](https://goreportcard.com/report/github.com/webability-go/xmodules/clientlink)
[ ![GoDoc](https://godoc.org/github.com/webability-go/xmodules/clientlink?status.png)](https://godoc.org/github.com/webability-go/xmodules/clientlink)
[ ![GolangCI](https://golangci.com/badges/github.com/webability-go/xmodules/clientlink.svg)](https://golangci.com)

xmodules/clientlink for Xamboo - GO
================================

The clientlink package is used to synchronize a local table of clients with a master table in another datasource. This is used to build primary key id clients to foreign keys tables in the local database


Version Changes Control
=======================

v0.0.1 - 2020-05-08
- Compatible with xmodule context standard

v0.0.0 - 2020-03-05
- Support for multithread context implemented

v0.0.0 - 2020-01-23
- Function to synchronize with origin database added

v0.0.0 - 2020-01-03
- This document added

v0.0.0 - 2019-12-18
- First release of module


Manual:
=======================

I. clientlink
=======================

Intro

-----------------------
1. Overview

Example:

```
import "github.com/webability-go/xmodules/clientlink"

```


-----------------------
2. Reference

To use the package:

import "github.com/webability-go/xmodules/clientlink"
