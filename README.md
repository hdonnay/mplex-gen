mplex-gen generates code for a channel multiplexer.

The function generated accepts a slice of channels, and spawns n+1
goroutines to service the supplied channel. The returned channel will be
closed if all supplied channels are closed. It is the caller's
responsibility to service the returned channel. Failure to do so will
put backpressue on the supplier channels.

This tool is meant to be called via `go generate`:

	//go:generate mplex-gen -o mplex.go main *bytes.Buffer


