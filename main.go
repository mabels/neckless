package main

import "github.com/mabels/neckless/cmd/neckless"


// GitCommit is injected during compile time
var GitCommit string

// Version is injected during compile time
var Version string

func main() { 
   neckless.Neckless(GitCommit, Version)
}
