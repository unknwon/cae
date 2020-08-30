Compression and Archive Extensions
==================================

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/unknwon/cae/Go?logo=github&style=for-the-badge)](https://github.com/unknwon/cae/actions?query=workflow%3AGo)
[![codecov](https://img.shields.io/codecov/c/github/unknwon/cae/master?logo=codecov&style=for-the-badge)](https://codecov.io/gh/unknwon/cae)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/unknwon/cae?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/unknwon/cae)

[中文文档](README_ZH.md)

Package cae implements PHP-like Compression and Archive Extensions.

But this package has some modifications depends on Go-style.

Reference: [PHP:Compression and Archive Extensions](http://www.php.net/manual/en/refs.compression.php).

Code Convention: based on [Go Code Convention](https://github.com/unknwon/go-code-convention).

### Implementations

Package `zip`([Go Walker](http://gowalker.org/github.com/unknwon/cae/zip)) and `tz`([Go Walker](http://gowalker.org/github.com/unknwon/cae/tz)) both enable you to transparently read or write ZIP/TAR.GZ compressed archives and the files inside them.

- Features:
	- Add file or directory from everywhere to archive, no one-to-one limitation.
	- Extract part of entries, not all at once. 
	- Stream data directly into `io.Writer` without any file system storage.

### Test cases and Coverage

All subpackages use [GoConvey](http://goconvey.co/) to write test cases, and coverage is more than 80 percent.

### Use cases

- [Gogs](https://github.com/gogits/gogs): self hosted Git service in the Go Programming Language.
- [GoBlog](https://github.com/fuxiaohei/GoBlog): personal blogging application.
- [GoBuild](https://github.com/shxsun/gobuild/): online Go cross-platform compilation and download service.

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
