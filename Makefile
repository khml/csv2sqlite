.PHONY: build-linux build-darwin build-windows build-release-artifacts

build-release-artifacts:
	make build-linux && \
	make build-darwin && \
	make build-windows && \
	zip csv2sqlite.linux_amd64.zip csv2sqlite_linux_amd64 && \
	zip csv2sqlite.darwind_arm64.zip csv2sqlite_darwind_arm64 && \
	zip csv2sqlite.windows_amd64.zip csv2sqlite.exe

build-linux:
	cd ./src && GOOS=linux GOARCH=amd64 go build -o ../csv2sqlite_linux_amd64 -v && cd -

build-darwin:
	cd ./src && GOOS=darwin GOARCH=arm64 go build -o ../csv2sqlite_darwind_arm64 -v && cd -

build-windows:
	cd ./src && GOOS=windows GOARCH=amd64 go build -o ../csv2sqlite.exe -v && cd -

