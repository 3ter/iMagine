# iMagine

For running the demo application execute the following in a terminal (e.g. directly in vscode):

`cd cmd/ && go run iMagine.go`

Build Windows executable from Linux:
```
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build
```
(You need to have `gcc-mingw-w64-x86-64` installed for that)
