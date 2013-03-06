go-compresser
===============

Sources Statement:
- go-tar.gz.go / go-zip.go
	- Can do basic pack and unpack operations
	- Can handle with single file instead of directory, but unpack as directory always
	- Total length of path cannot be too long, otherwise error "header too long" will occur 