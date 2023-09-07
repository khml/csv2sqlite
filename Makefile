.PHONY: build-linux build-darwin build-windows

build-linux:
	cd ./src && GOOS=windows GOARCH=amd64 go build -o ../csv2sqlite.exe -v && cd -

build-darwin:
	cd ./src && GOOS=darwin GOARCH=arm64 go build -o ../csv2sqlite_darwind_arm64 -v && cd -

build-windows:
	cd ./src && GOOS=windows GOARCH=amd64 go build -o ../csv2sqlite_linux_amd64 -v && cd -

