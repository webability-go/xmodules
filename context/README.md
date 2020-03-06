[ ![Go Report Card](https://goreportcard.com/badge/github.com/webability-go/xmodules/context)](https://goreportcard.com/report/github.com/webability-go/xmodules/context)
[ ![GoDoc](https://godoc.org/github.com/webability-go/xmodules/context?status.png)](https://godoc.org/github.com/webability-go/xmodules/context)
[ ![GolangCI](https://golangci.com/badges/github.com/webability-go/xmodules/context.svg)](https://golangci.com)

xmodules/context for Xamboo - GO
================================

The context package is used to build a set of contexts for the Xamboo CMS, so each site is build on top of its context.
The context contains links to databases, tables, logs, config params, supported languages and installed modules.
The context package controls also the installed package on each context.


Version Changes Control
=======================

v2.0.0 - 2020-03-05
- Implemented support for multithread (Mutex) on each object of the Container and Context.
- The data is now accessible by Get/Set/Add functions

v1.0.0 - 2019-12-18
- First release of module



Manual:
=======================

I. Context
=======================

Intro

-----------------------
1. Overview

Example:

```
import "github.com/webability-go/xmodules/context"

```


-----------------------
2. Reference

To use the package:

import "github.com/webability-go/xmodules/xcontext"




II. Modules
=======================

Intro

1. Overview
------------------------

Example:

```
import "github.com/webability-go/xmodules/context"

```

2. Reference
------------------------

To use the package:

import "github.com/webability-go/xmodules/xcontext"


---
