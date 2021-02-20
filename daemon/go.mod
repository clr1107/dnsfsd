module github.com/clr1107/dnsfsd/daemon

go 1.14

replace github.com/clr1107/dnsfsd/pkg => ../pkg

require (
	github.com/clr1107/dnsfsd/pkg v0.0.0-00010101000000-000000000000
	github.com/miekg/dns v1.1.38
	github.com/spf13/viper v1.7.1
	golang.org/x/net v0.0.0-20210220033124-5f55cee0dc0d // indirect
)
