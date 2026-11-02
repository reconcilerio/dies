// Command inittrace imports every die package, to measure their init cost
// with GODEBUG=inittrace=1.
package main

import _ "reconciler.io/dies/v2/apis"

func main() {}
