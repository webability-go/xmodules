[ ![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xmodules/userlink)](https://goreportcard.com/report/github.com/webability-go/xmodules/userlink)
[ ![GoDoc](https://godoc.org/github.com/webability-go/xmodules/userlink?status.png)](https://godoc.org/github.com/webability-go/xmodules/userlink)
[ ![GolangCI](https://golangci.com/badges/github.com/webability-go/xmodules/userlink.svg)](https://golangci.com)

xmodules/userlink for Xamboo - GO
================================

The userlink package is used to build a replicated set of administration users in a distant database, mainly to link with FK/PK to this table, when you do not need to scan all the accesses packages, but just need the user ID


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

I. Userlink
=======================

Intro

-----------------------
1. Overview

Example:

```
import "github.com/webability-go/xmodules/userlink"

```


-----------------------
2. Reference

To use the package:

import "github.com/webability-go/xmodules/userlink"
