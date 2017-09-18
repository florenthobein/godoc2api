**WORK IN PROGRESS**

godoc2api [![Build Status](https://travis-ci.org/florenthobein/godoc2api.svg?branch=master)](https://travis-ci.org/florenthobein/godoc2api)
===
> Easily describe your API by writing JSdoc-like golang comments that render into a [RAML 1.0](https://raml.org/) file

# About
## What is RAML?

> RESTful API Modeling Language (RAML) is a YAML-based language for describing RESTful APIs. It provides all the information necessary to describe RESTful or practically RESTful APIs.
[Wikipedia](https://en.wikipedia.org/wiki/RAML_(software))

## Documentation-oriented design of APIs

An API is a living thing -- wether it's from a designer, a builder or a consumer point of view, maintaining a coherent and up-to-date documentation of its endpoints is mandatory.

`godoc2api` simplifies this process by peeping into your go code to extract and generate a structured, comprehensible API documentation.

* Overly simple to implement on existing APIs
* Binds code documentation to API behaviour for feature integrity
* Enforces good code documentation practices
* Produces a standardised and consommable output
* Code documentation + auto-generated public doc + testable API = 👍

# Installation

```bash
go get github.com/florenthobein/godoc2api
```

# Limitations

For now only RAML 1.0 specification is supported. It also mainly focuses on a full `application/json` API.

# Usage

## Examples

Detailed examples are written on the [godoc page](https://godoc.org/github.com/florenthobein/godoc2api/examples)

# Comment parsing

> todo

# Configuration

## Defining types

> todo

## Defining traits

> todo

## Defining security schemes

> todo

## Defining annotations

> todo

# Debugging

Change the log level of the library to display warnings:
```golang
// Possible values:
// - LOG_DEBUG
// - LOG_WARN
// - LOG_ERR
// - LOG_NOTHING (default)
godoc2api.LogLevel = godoc2api.LOG_WARN
```

# Roadmap

- [ ] Implementation of traits
- [ ] Implementation of security schemes
- [ ] Implementation of annotations
- [ ] Exportation in multiple files & includes
- [ ] RAML structure validation
- [ ] Support for other standards

# Credits
 
 * This library is inspired by the work on [github.com/Jumpscale/go-raml](github.com/Jumpscale)
 * And of course, the [RAML 1.0 specifications](https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md)

# License

Copyright (c) 2017 Florent Hobein. Licensed under the MIT license.