# Docs

## [faker](#)

Struct Data Fake Generator

Faker will generate you a fake data based on your Struct.

[![Build Status](https://travis-ci.org/togglhire/faker.svg?branch=master)](https://travis-ci.org/togglhire/faker)
[![codecov](https://codecov.io/gh/togglhire/faker/branch/master/graph/badge.svg)](https://codecov.io/gh/togglhire/faker)
[![Go Report Card](https://goreportcard.com/badge/github.com/togglhire/faker)](https://goreportcard.com/report/github.com/togglhire/faker)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/togglhire/faker/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/togglhire/faker?status.svg)](https://godoc.org/github.com/togglhire/faker)

## Index

* [Support](#support)
* [Getting Started](#getting-started)
* [Example](#example)
* [Limitation](#limitation)
* [Contribution](#contribution)


## Support

You can file an [Issue](https://github.com/togglhire/faker/issues/new).
See documentation in [Godoc](https://godoc.org/github.com/togglhire/faker)


## Getting Started

#### Download

```shell
go get -u github.com/togglhire/faker/v3
```
# Example

---
 
 - Using Struct's tag: 
   - [basic tags: example_with_tags_test.go](/example_with_tags_test.go)
   - [length and bounds: example_with_tags_lenbounds_test.go](/example_with_tags_lenbounds_test.go)
   - [unique: example_with_tags_unique_test.go](example_with_tags_unique_test.go)
 - Custom Struct's tag (define your own faker data): [example_custom_faker_test.go](/example_custom_faker_test.go)
 - Without struct's tag: [example_without_tag_test.go](/example_without_tag_test.go)
 - Single Fake Data Function: [example_single_fake_data_test.go](/example_single_fake_data_test.go)
 
## DEMO

---

![Example to use Faker](https://cdn-images-1.medium.com/max/800/1*AkMbxngg7zfvtWiuvFb4Mg.gif)

## Benchmark

---

Bench To Generate Fake Data
#### Without Tag
```bash
BenchmarkFakerDataNOTTagged-4             500000              3049 ns/op             488 B/op         20 allocs/op
```

#### Using Tag
```bash
 BenchmarkFakerDataTagged-4                100000             17470 ns/op             380 B/op         26 allocs/op
```

### MUST KNOW

---

The Struct Field must be PUBLIC.<br>
Support Only For :

* `int`, `int8`, `int16`, `int32` & `int64`
* `[]int`, `[]int8`, `[]int16`, `[]int32` & `[]int64`
* `bool` & `[]bool`
* `string` & `[]string`
* `float32`, `float64`, `[]float32` &`[]float64`
* `time.Time` & `[]time.Time`
* Nested Struct Field

## Limitation

---

Unfortunately this library has some limitation
* It does not support private fields. Make sure your structs fields you intend to generate fake data for are public, it would otherwise trigger a panic. You can however omit fields using a tag skip `faker:"-"` on your private fields.
* It does not support the `interface{}` data type. How could we generate anything without knowing its data type?
* It does not support the `map[interface{}]interface{}`, `map[any_type]interface{}` & `map[interface{}]any_type` data types. Once again, we cannot generate values for an unknown data type.
* Custom types are not fully supported. However some custom types are already supported: we are still investigating how to do this the correct way. For now, if you use `faker`, it's safer not to use any custom types in order to avoid panics.

## Contribution

---

To contrib to this project, you can open a PR or an issue.
