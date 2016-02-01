# vCardGen
_A simple vCard generator in Go_

by Cathal Garvey, Copyright 2016, Licensed AGPLv3 or later.

### About

![Godoc Badge](https://godoc.org/github.com/cathalgarvey/vcardgen?status.svg)

This is just a vCard generator in Go, which isn't even entirely compliant.
It is based upon [vCards JS](https://github.com/enesser/vCards-JS) by
[Eric J Nesser](https://github.com/enesser).

It can be used to create a vCard string simply by creating and assigning to
various properties, then calling `card.GetFormattedString()`:

```go
cathal := vcardgen.New()
cathal.FirstName = "Cathal"
cathal.MiddleName = "Joseph"
cathal.LastName = "Garvey"
cathal.CellPhone = "+353863434567"
cathal.Email = "cathal@xyz.xyz"
cathalVcard := cathal.GetFormattedString()
```

That's it! This isn't a project to be excited or proud of, because vCard is a
pile of crap, as formats go. But, it's a common interchange format, so when
needs must..

