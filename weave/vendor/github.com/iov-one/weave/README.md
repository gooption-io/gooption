# Iov Weave
[![Build Status TravisCI](https://api.travis-ci.com/iov-one/weave.svg?branch=master)](https://travis-ci.com/iov-one/weave)
[![codecov](https://codecov.io/gh/iov-one/weave/branch/master/graph/badge.svg)](https://codecov.io/gh/iov-one/weave/branch/master)
[![LoC](https://tokei.rs/b1/github/iov-one/weave)](https://github.com/iov-one/weave)
[![Go Report Card](https://goreportcard.com/badge/github.com/iov-one/weave)](https://goreportcard.com/report/github.com/iov-one/weave)
[![API Reference](https://godoc.org/github.com/iov-one/weave?status.svg
)](https://godoc.org/github.com/iov-one/weave)
[![ReadTheDocs](https://readthedocs.org/projects/weave/badge/?version=latest)](http://weave.readthedocs.io/en/latest/)
[![license](https://img.shields.io/github/license/iov-one/weave.svg)](https://github.com/iov-one/weave/blob/master/LICENSE)

IOV Weave is a framework for quickly building your custom
[ABCI application](https://github.com/tendermint/abci)
to run a blockchain on top of the best-of-class
BFT Proof-of-stake [Tendermint consensus engine](https://tendermint.com).
It provides much commonly used functionality that can
quickly be imported in your custom chain, as well as a
simple framework for adding the custom functionality unique
to your project.

**Note: Requires Go 1.9+**

It is inspired by the routing and middleware model of many web
application frameworks, and informed by years of wrestling with
blockchain state machines. More directly, it is based on the
official cosmos-sdk, both the
[0.8 release](https://github.com/cosmos/cosmos-sdk/tree/v0.8.0) as well as the
future [0.9 rewrite](https://github.com/cosmos/cosmos-sdk/tree/develop). Naturally, as I was the main author of 0.8.

While both of those are extremely powerful and flexible
and contain advanced features, they have a steep learning
curve for novice users. Thus, this library aims to favor
simplicity over power when there is a choice. If you hit
limitations in the design of this library (such as
maintaining multiple merkle stores in one app), I highly
advise you to use
[the official cosmos sdk](https://github.com/cosmos/cosmos-sdk).

On the other hand, if you want to try out tendermint, or have a
design that doesn't require an advanced setup, you should try
this library and give feedback, especially on ease-of-use.
The end goal is to make blockchain development almost as
productive as web development (in golang), by providing
defaults and best practices for many choices, while allowing
extreme flexibility in business logic and data modelling.

For more details on the design goals, see the
[Design Document](./docs/design.rst)

## Prerequisites

* [golang 1.9+](https://golang.org/doc/install)


## Instructions

First, make sure you have
[set up the requirements](https://weave.readthedocs.io/en/latest/mycoind/setup.html).
If you have a solid go and node developer setup, you may skip this,
but good to go through it to be sure.

Once you are set up, you should be able to run something
like the following to compile both `mycoind` (sample app)
and `tendermint` (a [BFT consensus engine](https://tendermint.com)):

```
go get github.com/iov-one/weave
cd $GOPATH/src/github.com/iov-one/weave
make deps
make install
# test it built properly
tendermint version
# 0.21.0-46369a1a
mycoind version
# v0.6.0-7-g3aa16c9
```

Note that this app relies on a separate tendermint process
to drive it. It is helpful to first read a primer on
[tendermint](https://tendermint.readthedocs.io/en/master/introduction.html)
as well as the documentation on the
[tendermint cli commands](https://tendermint.readthedocs.io/en/master/using-tendermint.html).

Once it compiles, I highly suggest going through the
[tutorials on readthedocs](https://weave.readthedocs.io/en/latest/index.html#mycoin-tutorial)

## History

The original version, until `v0.6.0` was released under
`confio/weave`. The original author, Ethan Frey, had
previously worked on the
[Cosmos SDK](https://github.com/cosmos/cosmos-sdk)
and wanted to make a simpler framework he could use to
start building demo apps, while the main sdk matured.
Thus, `confio/weave` was born the first few months of 2018.
This framework was designed to be open source and shared,
but the only real usage and development was by
[IOV](https://github.com/iov-one), so it was donated to
that organization in August 2018 to be developed further
for their BNS blockchain, as well as a companion to
[iov-core](https://github.com/iov-one/iov-core)
client libraries that deprecated `confio/weave-js`
