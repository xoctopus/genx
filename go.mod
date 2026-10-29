module github.com/xoctopus/genx

go 1.27.3

tool (
	github.com/xoctopus/genx/internal/cmd/example
	github.com/xoctopus/genx/internal/cmd/skill-install
)

require (
	github.com/xoctopus/pkgx v0.4.4
	github.com/xoctopus/typx v0.4.7
	// +skill:testx
	github.com/xoctopus/x v0.6.0
	golang.org/x/mod v0.42.0
	golang.org/x/tools v0.52.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
)
