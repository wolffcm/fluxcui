module github.com/wolffcm/fluxcui

go 1.13

require (
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59
	github.com/google/go-cmp v0.3.1
	github.com/influxdata/flux v0.57.0
	github.com/influxdata/influxdb v1.5.1-0.20191213220711-88468822e23d
	github.com/jroimartin/gocui v0.4.0
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/nsf/termbox-go v0.0.0-20190817171036-93860e161317 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	github.com/wolffcm/drawille-go v0.0.0-20191221022539-4afb83ac080b
)

replace github.com/wolffcm/drawille-go => ../drawille-go
