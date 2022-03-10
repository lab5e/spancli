module github.com/lab5e/spancli

go 1.15

require (
	github.com/jedib0t/go-pretty/v6 v6.2.4
	github.com/jessevdk/go-flags v1.5.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lab5e/go-spanapi/v4 v4.1.18
	github.com/lab5e/go-userapi v1.3.11
	github.com/stretchr/testify v1.7.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/lab5e/go-spanapi/v4 => ../go-spanapi
