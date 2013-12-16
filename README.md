Compression and Archive Extensions
==================================

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/Unknwon/cae)

[中文文档](README_ZH.md)

Package cae implements PHP-like Compression and Archive Extensions.

But this package has some modifications depends on Go-style.

Reference: [PHP:Compression and Archive Extensions](http://www.php.net/manual/en/refs.compression.php).

### Implementations

- Package `zip`: this package enables you to transparently read or write ZIP compressed archives and the files inside them. [Go Walker](http://gowalker.org/github.com/Unknwon/cae/zip).
	- Features:
		- Add file or directory from everywhere to archive, no one-to-one limitation.
		- Able to extract part of entries, not all at once. 

### Test cases and Coverage

All subpackages use [GoConvey](http://smartystreets.github.io/goconvey/) to write test cases, and coverage is more than 85 percent.