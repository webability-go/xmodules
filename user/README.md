[ ![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xmodules/user)](https://goreportcard.com/report/github.com/webability-go/xmodules/user)
[ ![GoDoc](https://godoc.org/github.com/webability-go/xmodules/user?status.png)](https://godoc.org/github.com/webability-go/xmodules/user)
[ ![GolangCI](https://golangci.com/badges/github.com/webability-go/xmodules/user.svg)](https://golangci.com)

xmodules/user for Xamboo - GO
================================

The user package is used to build a set of administration users with all the security controls needed:
- accesses (simple named accesses)
- extended accesses (based on the records of a table)
- sessions and history
- profiles
- per-user accesses adjustements


Version Changes Control
=======================

v1.0.1 - 2020-01-03
- This document added
- Order added to synchronize tables in database (due to hierarchy of FK-PK)

v1.0.0 - 2019-12-18
- First release of module



Manual:
=======================

I. Users
=======================

Intro

-----------------------
1. Overview

Example:

```
import "github.com/webability-go/xmodules/user"

```


-----------------------
2. Reference

To use the package:

import "github.com/webability-go/xmodules/user"
