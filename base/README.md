[ ![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xmodules/base)](https://goreportcard.com/report/github.com/webability-go/xmodules/base)
[ ![GoDoc](https://godoc.org/github.com/webability-go/xmodules/base?status.png)](https://godoc.org/github.com/webability-go/xmodules/base)
[ ![GolangCI](https://golangci.com/badges/github.com/webability-go/xmodules/base.svg)](https://golangci.com)

xmodules/base for Xamboo - GO
================================

The base package is used to control the applications and compiled xmodules for the Xamboo CMS, so each site is build on top of its base.
The base contains links to databases, tables, logs, config params, supported languages and installed modules.
The base package controls also the installed package on each datasource.
The base package is compatible with database transactions and multithread.
The base package support english, spanish and french installations.

TO DO:
=========
- Finish translation of messages in spanish y French
- Add messages in container.go

Version Changes Control
=======================

v0.1.1 - 2021-01-17
- Implementation of languages and messages

v0.1.0 - 2020-05-25
- Renamed to base instead of context
- Now support new xamboo standart assets.Application for Applications and XModules

v0.0.1 - 2020-03-05
- Implemented support for multithread (Mutex) on each object of the Container and base.
- The data is now accessible by Get/Set/Add functions

v0.0.0 - 2019-12-18
- First release of module



Manual:
=======================

I. base
=======================

Intro

-----------------------
1. Overview

Example:

```
import "github.com/webability-go/xmodules/base"

```


-----------------------
2. Reference

To use the package:

import "github.com/webability-go/xmodules/xbase"




II. Modules
=======================

Intro

1. Overview
------------------------

Example:

```
import "github.com/webability-go/xmodules/base"

```

2. Reference
------------------------

To use the package:

import "github.com/webability-go/xmodules/xbase"


---
