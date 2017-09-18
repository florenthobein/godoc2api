**WORK IN PROGRESS**

godoc2api [![Build Status](https://travis-ci.org/florenthobein/godoc2api.svg?branch=master)](https://travis-ci.org/florenthobein/godoc2api)
===
Easily describe your API by writing JSdoc-like golang comments that renders into a [RAML 1.0](https://raml.org/) file

# About
## What is RAML?

> RESTful API Modeling Language (RAML) is a YAML-based language for describing RESTful APIs. It provides all the information necessary to describe RESTful or practically RESTful APIs.
[Wikipedia](https://en.wikipedia.org/wiki/RAML_(software))

## Documentation-oriented design of APIs

> todo

# Installation

```bash
go get github.com/florenthobein/godoc2api
```

# Usage

> todo

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
godoc2api.LogLevel = godoc2api.LOG_WARN
```

# Roadmap

Traits, security schemes & annotations

# Licence

> todo

# Credits
 
 * This library is inspired by the work on [github.com/Jumpscale/go-raml](github.com/Jumpscale)
 * And of course, the [RAML 1.0 specifications](https://github.com/raml-org/raml-spec/blob/master/versions/raml-10/raml-10.md)