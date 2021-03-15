module github.com/clr1107/dnsfsd/dnsfsd-util

go 1.16

// this module should be compiled as `dnsfs` and ran as such
// `dnsfsd` is the daemon.

replace github.com/clr1107/dnsfsd/pkg => ../pkg

require (
	github.com/clr1107/dnsfsd/pkg v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
)
