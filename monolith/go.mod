module github.com/stackus/ftgogo/monolith

go 1.16

replace github.com/stackus/ftgogo/accounting => ./../accounting

replace github.com/stackus/ftgogo/consumer => ./../consumer

replace github.com/stackus/ftgogo/customer-web => ./../customer-web

replace github.com/stackus/ftgogo/delivery => ./../delivery

replace github.com/stackus/ftgogo/kitchen => ./../kitchen

replace github.com/stackus/ftgogo/order => ./../order

replace github.com/stackus/ftgogo/order-history => ./../order-history

replace github.com/stackus/ftgogo/restaurant => ./../restaurant

replace github.com/stackus/ftgogo/serviceapis => ./../serviceapis

replace shared-go => ../shared-go

require (
	github.com/stackus/ftgogo/accounting v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/consumer v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/customer-web v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/delivery v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/kitchen v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/order v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/order-history v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/restaurant v0.0.0-20210116185538-3dd9fbb69179
	github.com/stackus/ftgogo/serviceapis v0.0.0-20210116185538-3dd9fbb69179
	shared-go v0.0.0-00010101000000-000000000000
)
