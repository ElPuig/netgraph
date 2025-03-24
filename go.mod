module github.com/ElPuig/netgraph

go 1.24.1

require (
	github.com/ElPuig/netgraph/pkg/xml_loader v0.0.0-00010101000000-000000000000
	github.com/alexflint/go-arg v1.5.1
)

require (
	github.com/PuerkitoBio/goquery v1.10.2 // indirect
	github.com/alexflint/go-scalar v1.2.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	golang.org/x/net v0.35.0 // indirect
)

replace github.com/ElPuig/netgraph/pkg/xml_loader => ./pkg/xml_loader
