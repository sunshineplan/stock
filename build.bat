@echo off
go build
go build -ldflags "-s -w -H windowsgui" -o mystocks.exe
