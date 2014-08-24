HTTP Package Filter
2014-9-13 V1.0
====================
1. 用途：
对HTTP GET/POST格式的数据包进行过滤
如果HTTP数据包含下以格式：
GET支持query格式：key1=name1&key2=name2&key3=name3...
POST支持JSON格式：{"key1":name1,"key2":name2,"key3":name3,...}
则可通过过滤器对数据进行过滤操作（允许/拒绝）。

2. 特性
*简单高效
*多种逻辑运算符
*条件分支
*正则匹配
*内置函数

3. 例子
package foo

import (
	"filter"
	"strings"
)

func foo() {
	// our filter rule script:
	// @ means 'is in list ?'
	//
	filter_rule_script := "gz @ ( '10', 'abc', '303' ) => 1; default => 0"

	// http query content:
	//
	querys := []string{
		"gz=10&id=123456",
		"gz=303&id=123456",
		"gz=100&id=123456",
		"gz=111&id=123456"}

	// create a parser
	//
	h, err := filter.NewParser(strings.NewReader(filter_rule_script))
	if err != nil {
		panic(err)
	}

	for _, query := range querys {
		// create symbol list
		//
		symlist, err := filter.QueryToSymlist(query)
		if err != nil {
			panic(err)
		}

		// get result value
		//
		result, err := h.Parse(symlist)
		if err != nil {
			panic(err)
		}

		if result == 1 {
			// deny this package
			// querys[0] and querys[1] are denied because of filter_rule_script
			//
		}
	}

	// after that, querys[2] and querys[3] passed our test
	// and can be delivered forward
	//
}

该例子包含四步操作：
首先，自定义过滤规则(filter_rule_script)，调用API得到解析器句柄(h)
其次，HTTP数据包(querys)通过调用API得到符号输入表(symlist)
再次，根据符号输入表，调用解析器句柄的解析API(h.Parse(symlist))，得到结果(ret)
最后，根据结果是否符合预期，对数据包做丢弃或进一步处理

4. 语法手册
4.1 类型
字符串、数值是仅有的两种基本类型
字符串由单引号引用，如'abc','hello,world'
数值类型由float64表示，整形亦被转化为float64类型，如123,3.14159265
bool类型最终会转化为数值0.0或1.0

4.2 变量
变量命名由以下正则表达式描述：
[_a-zA-Z][_a-zA-Z0-9]*

变量在HTTP GET/POST数据包中被定义和赋值
如"gz=10&id=123456"定义了两个变量gz、id
或{"gz":"10","id":"123456"}亦能达到同样目的

4.3 常量
两种类型的常理，字符型常量、数值型常量

4.4 函数
过滤器目前有5个内置函数：
函数名	 说明	                        举例
len()	 求变量值或常量字符串长度	    len(gz), len(‘abc’)
count()	 求变量个数	                    count()
atoi()	 字符串转换为数字	            atoi(gz), atoi('123')
itoa()	 变量值/常量字符串所代表类型	itoa(100)
md5()	 求32位md5值	                md5(gz, ‘somesalt’), md5(gz)

md5支持1个或多个参数，其值为所有字符串参数拼接后的md5串

4.5 表达式
过滤器支持以下操作，使用括号改变优先级
逻辑操作名   说明         求值后的对象      举例
&&           逻辑与       数值              x > 10 && y == 'abc'
||           逻辑或       数值              x > 10 || y == 'abc'
@            在列表       数值/字符串       x @ (202,303)
!@           不在列表     数值/字符串       x !@ ('abc','def')
>            大于         数值/字符串       x > len('abc')
<            小于         数值/字符串       x < 'def'
>=           大于等于     数值/字符串       x >= 10
<=           小于等于     数值/字符串       len('abc') <= 10
==           等于         数值/字符串       10 == x
!=           不等于       数值/字符串       count() != 10
#            正则匹配     字符串            x # '20.*'
!#           非正则匹配   字符串            itoa(20) !# '20'
()           括号运算     数值              x > 10 && ( y == 'abcd' || z == 9 )
=>           设置返回值   数值              x > 10 => 1000

比较操作(@,!@,>,<,>=,<=,==,!=)，支持字符串比较和数值比较，字符串比较与C标准库函数strcmp()返回结果约定一致
正则匹配操作(#,!#)右部只能为字符串（正则模式串）,支持POSIX-ERE正则匹配(regexp.CompilePOSIX())

4.6 注释
过滤器支持行注释，以 “//” 开头，直到行尾

4.7 空格
分号、空格、制表符、换行为分隔符，表达式会自动忽略这些符号

4.8 语句
*过滤器支持多条语句组合，使用空格进行分隔
*过滤器支仅支持条件分支语句，条件满足便返回

例1：测试参数个数：
count() == 10 =>1; count() == 9 =>2; default =>0;
解释：如果输入的参数为10个，则返回值为1，如果为9个，则返回值为2，否则返回0
符号'=>'用来设置返回值，不指定符号时返回0（条件不成立）或1（条件成立）
词法分析或语法分析发生错误时返回-1

例2：测试数据包中变量md5值
md5(x,'@163.com') == '4131bfb2bf25f5d9ef86ff9bf53e0055';
数据包中变量x的内容和'salt'组合后，生成的md5值是否等于右边字符串

例3：组合测试
count() == 3 && md5(key,'@163.com') == '4131bfb2bf25f5d9ef86ff9bf53e0055' && flag # '^[01]$' && len(value) == 10
匹配成功条件：
HTTP数据包中变量个数为3（key,flag,value)，其中变量key和常量字符串'salt'组合成后，内容md5值为4131bfb2bf25f5d9ef86ff9bf53e0055，
变量flag取值只能为0或1，变量value值长度为10
则成功的HTTP数据包格式可能为：
GET key=justhechuang&value=1234567890&flag=1
POST {"key":"justhechuang", "value":"1234567890", "flag":"1"}

4.9 符号输入
HTTP GET/POST数据包即为符号输入
过滤器支持两种类型的符号输入：
*Query格式，用'='和'&'分隔的字符串
*JSON格式，Object对象组成的数据

Query格式定义的变量类型全为字符串
而JSON格式定义的变量可以为字符串类型和数值类型
例如：
key=justhechuang&value=1234567890&flag=1
定义了三个字符串变量:
变量key,其值为'justhechuang'
变量value,其值为'123456789'
变量flag,其值为'1'
而
{"key":"justhechuang", "value":"1234567890", "flag":1.0}
中，变量key和value与前述Query格式一致，而flag变量则为float64类型

4.10 案例
sample目录下有测试用例，每行其格式为：
期望值%过滤规则%符号输入(expect_value%filter_rule%symbol_input)
根据过滤规则和符号输入求的值如果等于期望值，则测试成功，否则为失败

5. 安装
编译： make
测试： make test
清除： make clean
本程序采用nex加go tool yacc生成
编译nex二进制程序，进入nex目录:go build，然后将生成的nex文件拷贝到系统搜索路径,既能正常编译测试filter