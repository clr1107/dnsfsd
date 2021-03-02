module github.com/clr1107/dnsfsd/daemon

go 1.16

replace github.com/clr1107/dnsfsd/pkg => ../pkg

require (
	github.com/clr1107/dnsfsd/pkg v0.0.0-00010101000000-000000000000
	github.com/miekg/dns v1.1.39
	github.com/spf13/viper v1.7.1
)
