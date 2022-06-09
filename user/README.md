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

v0.1.0 - 2022-03-
- First official release. Working on the xmodules/base, Xamboo, XCore, XDominion and XMask of this date.
- Admin pages integrated in xmodule

v0.0.1 - 2020-05-08
- Compatible with xmodule context standard

v0.0.0 - 2020-03-05
- Support for multithread context implemented

v0.0.0 - 2020-01-23
- BuildCaches is now on parallel thread

v0.0.0 - 2020-01-03
- This document added
- Order added to synchronize tables in database (due to hierarchy of FK-PK)

v0.0.0 - 2019-12-18
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

-----------------------
3. Normalization

All the concept of the module must be called by se same ywa among messages, screens and data.
English:
- Group of right accesses
- Right access
- User Profile
- Administration User
Spanish:
- Grupo de permisos de accesos
- Permiso de Acceso
- Perfil de usuario
- Usuario Administrador
French:
- Groupe d'accès
- Droit d'accès
- Profil d'utilisateur
- Utilisateur administrateur
