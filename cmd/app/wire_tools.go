//go:build tools
// +build tools

// solve the following problem:

// missing go.sum entry for module providing package github.com/google/subcommands (imported by github.com/google/wire/cmd/wire); to add:
//        go get github.com/google/wire/cmd/wire@v0.5.0

package main

import _ "github.com/google/wire/cmd/wire"
