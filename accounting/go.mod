module github.com/stackus/ftgogo/account

go 1.16

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/stackus/edat v0.0.3
	github.com/stackus/errors v0.0.2
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	shared-go v0.0.0-00010101000000-000000000000
)
