# corkscrew-go
Corkscrew with go

## usage
go install

ssh -o ProxyCommand="corkscrew-go 127.0.0.1:8080 %h:%p" hostname
