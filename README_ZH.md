压缩与打包扩展
=============

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/Unknwon/cae)

包 cae 实现了 PHP 风格的压缩与打包扩展。

但本包依据 Go 语言的风格进行了一些修改。

引用：[PHP:Compression and Archive Extensions](http://www.php.net/manual/en/refs.compression.php).

### 实现

- 包 `zip`：本包允许你轻易的读取或写入 ZIP 压缩档案和其内部文件。[Go Walker](http://gowalker.org/github.com/Unknwon/cae/zip).
	- 特性：
		- 允许将任意位置的文件或目录加入档案，没有一对一的操作限制。
		- 允许只解压部分文件，而非一次性解压全部。 

- 包 `tz`：本包允许你轻易的读取或写入 tar.gz 压缩档案和其内部文件。[Go Walker](http://gowalker.org/github.com/Unknwon/cae/tz).
	- 特性：
		- 允许将任意位置的文件或目录加入档案，没有一对一的操作限制。
		- 允许只解压部分文件，而非一次性解压全部。 

### 测试用例与覆盖率

所有子包均采用 [GoConvey](http://smartystreets.github.io/goconvey/) 来书写测试用例，覆盖率均超过百分之 85。