# GO基础

> Go是编译性语言，工具链将源代码及其依赖转换成计算机的机器指令。
>
> Go语言不需要在语句或者声明的末尾添加分号，除非一行上有多条语句。*编译器会主动把特定符号后的换行符转换为分号*。

## 循环语法

```go
for initialization; condition; post {
}
```

```go
for condition {
}
```

> 类似于 `while`

```go
for {
}
```

> 类似于 `while(true)`

# 程序结构

## 命名

> 1. Go语言的风格是尽量使用短小的名字。
> 2. 名字的开头字母的大小写决定了名字在包外的可见性。
> 3. 包本身的名字一般总是用小写字母。
> 4. 一个名字的作用域比较大，生命周期也比较长，那么用长的名字将会更有意义。
> 5. Go语言程序员推荐使用 **驼峰式** 命名，而像ASCII和HTML这样的缩略词则避免使用大小写混合的写法，它们可能被称为htmlEscape。

- **关键字**

  ```shell
  break      default       func     interface   select
  case       defer         go       map         struct
  chan       else          goto     package     switch
  const      fallthrough   if       range       type
  continue   for           import   return      var
  ```

  > 25个。不能用于**自定义**名字，只能在**特定语法结构**中使用。

- **预定义名称**

  - *内建常量*

    ```go
    true false iota nil
    ```

  - *内建类型*

    ```go
    int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64 uintptr float32 float64 complex128 complex64 bool byte rune string  error
    ```

  - *内建函数*

    ```go
    make len cap new append copy close delete complex real imag panic recover
    ```

  > 主要对应内建的常量、类型和函数。可以在定义中**重新使用**它们。

## 声明

> 1. 声明语句定义了程序的各种实体对象以及部分或全部的属性。
> 2. 每个源文件中以包的声明语句开始，说明该源文件是属于哪个包。
> 3. 包声明语句之后是import语句导入依赖的其它包。
> 4. 然后是包一级的类型「type」、变量「var」、常量「const」、函数「func」的声明语句。
> 5. 包一级的各种类型的声明语句的顺序无关紧要。
> 6. 一个声明语句将程序中的实体和一个名字关联，比如一个函数或一个变量。

## 变量

- **声明语法**

  ```go
  var i, j, k int                 // int, int, int
  var b, f, s = true, 2.3, "four" // bool, float64, string
  ```

  > var 变量名字 类型 = 表达式

  1. *其中 「类型」 或 「= 表达式」 两个部分可以省略其中的一个*。

     - **省略类型**：根据初始化表达式来推导变量的类型信息。
     
     - **省略初始化表达式**：用**零值**初始化该变量。*零值初始化机制可以确保每个声明的变量总是有一个良好定义的值*。
     
       >数值类型：0。
       >
       >布尔类型：false。
       >
       >字符串类型：空字符串。
       >
       >接口或引用类型「包括**slice**、**指针**、**map**、**chan**和函数」：nil。
       >
       >数组或结构体等聚合类型：每个元素或字段都是对应该类型的零值。

  2. 可以声明多个类型不同的变量。

  3. 初始化表达式可以是字面量或任意的表达式。
     - **包级别**声明的变量会在main入口函数**执行前**完成初始化。
     - **局部变量**将在声明语句被**执行到**的时候完成初始化。

- **简短变量声明**

  ```go
  i, j := 0, 1
  ```

  > 名字 := 表达式

  1. 局部变量的声明和初始化通常使用这种方式。

  2. var形式声明语句使用场景：

     >需要显式指定变量类型的地方。
     >
     >因变量稍后会被重新赋值而初始值无关紧要的地方。

  3. 简短变量声明左边的变量可能并非全部都是刚刚声明的。对于已经声明过的，就只有赋值行为了。

  4. 简短变量声明语句中必须至少要声明一个新的变量。

- **指针**

  > 指针之间可以进行**相等测试**：只有当它们*指向同*一个变量或*全部是nil*时才相等。
  >
  > 在Go语言中，返回函数中局部变量的地址也是安全的。
  >
  > 指针是实现标准库中flag包的关键技术。

  ```go
  func f() *int {
    v := 1
    return &v
  }
  ```

  > 每次调用f函数都将返回不同的结果
  >
  > ```go
  > fmt.Println(f() == f()) // "false"
  > ```

  - **new函数**

    >表达式 `new(T)` 将创建一个T类型的匿名变量。初始化为T类型的零值，然后返回变量地址。
    >
    >`new` 函数类似是一种语法糖，而不是一个新的基础概念。
    >
    >每次调用new函数都是返回一个新的变量的地址。

  - **生命周期**

    - *包变量*：等同于整个程序的运行周期。

    - *局部变量*：从创建的声明语句开始，直到该变量不再被引用为止。

      >**问**：垃圾收集器如何知道一个变量何时可以被回收？
      >
      >**答**：从每个包级的变量和每个当前运行函数的每一个局部变量开始，通过指针或引用的访问路径遍历，是否可以找到该变量。如果不存在这样的访问路径，那么说明该变量是不可达的，也就是说它是否存在并不会影响程序后续的计算结果。

    - **存储空间**：编译器自动选择在**栈**或**堆**上分配空间，并不由 `var` 或 `new` **声明变量的方式**决定的。

## 赋值

> 命名变量的赋值
>
> ```go
> x = 1
> ```
>
> 通过指针间接赋值
>
> ```go
> *p = true
> ```
>
> 结构体字段赋值
>
> ```go
> person.name = "bob"
> ```
>
> 数组、slice或map的元素赋值
>
> ```go
> count[x] = count[x] * scale
> ```
>
> **自增自减**
>
> ```go
> v := 1
> v++    // 等价方式 v = v + 1；v 变成 2
> v--    // 等价方式 v = v - 1；v 变成 1
> ```
>
> > 自增和自减是语句，而不是表达式，因此`x = i++`之类的表达式是错误的。

- **元组赋值**

  ```go
  x, y = y, x
  a[i], a[j] = a[j], a[i]
  ```

  > 交换两个变量的值

  1. 对于处理有些同时出现在元组赋值语句左右两边的变量很有帮助。

  2. 可以使一系列琐碎赋值更加紧凑。
  3. 如果表达式太复杂的话，应该尽量避免过度使用元组赋值。
  4. 可以用下划线空白标识符 `_` 来丢弃不需要的值。

- **可赋值性**

  > 只有右边的值对于左边的变量是可赋值的，赋值语句才是允许的。
  >
  > 对于两个值是否可以用 `==` 或 `!=` 进行相等比较的能力也和可赋值能力有关系。

## 类型

```go
type 类型名字 底层类型
```

- 变量或表达式的类型定义了对应存储值的属性特征。

- 新命名的类型提供了一个方法，用来分隔不同概念的类型，这样即使它们底层类型相同也是不兼容的。

  - 不可以 **被相互比较** 或 **混在一个表达式** 运算。

- 对于每一个类型T，都有一个对应的类型转换操作T(x)，用于将x转为T类型。

  - 如果T是指针类型，可能会需要用小括弧包装T，比如`(*int)(0)`

- 比较运算符 `==` 和 `<` 也可以用来比较 

  - **一个命名类型的变量** 和 **另一个有相同类型的变量** ，
  - 或有着 **相同底层类型的未命名类型** 的值之间做比较。

  ```go
  type Celsius float64    // 摄氏温度
  type Fahrenheit float64 // 华氏温度
  var c Celsius
  var f Fahrenheit
  fmt.Println(c == 0)          // "true"
  fmt.Println(f >= 0)          // "true"
  fmt.Println(c == f)          // compile error: type mismatch
  fmt.Println(c == Celsius(f)) // "true"!
  ```

- 如果两个值有着不同的类型，则不能直接进行比较。

## 包和文件

> 为了支持模块化、封装、单独编译和代码重用。
>
> 包中如果一个名字是大写字母开头的，那么该名字是导出的。

- **导入包**
  - 如果导入了一个包，但是又没有使用该包将被当作一个编译错误处理。
  - 导入语句将导入的包绑定到一个短小的名字，然后通过该短小的名字就可以引用包中导出的全部内容。
- **包初始化**
  - 包的初始化首先是解决包级变量的依赖顺序，然后按照包级变量声明出现的顺序依次初始化。
  - 每个文件都可以包含多个init初始化函数。
  - 每个包在解决依赖的前提下，以导入声明的顺序初始化，每个包只会被初始化一次。

## 作用域

> - 声明语句的作用域是指源代码中可以有效使用这个名字的范围。
>
> - **不要将作用域和生命周期混为一谈**。
>
>   - 声明语句的作用域对应的是一个源代码的文本区域；它是一个编译时的属性。
>
>   - 一个变量的生命周期是指程序运行时变量存在的有效时间段，在此时间区域内它可以被程序的其他部分引用；是一个运行时的概念。
>
> - 句法块内部声明的名字是无法被外部块访问的。*由花括弧所包含的一系列语句*。
>
> - 声明在代码中并未显式地使用花括号包裹起来，我们称之为词法块。
>
> - **控制流标号**：就是 `break`、`continue` 或 `goto` 语句后面跟着的那种标号，则是函数级的作用域。
> - 编译器遇到一个名字引用时，它会对其定义进行查找，查找过程从最内层的词法域向全局的作用域进行。如果查找失败，则报告“未声明的名字”这样的错误。如果该名字在内部和外部的块分别声明过，则内部块的声明首先被找到。

# 数据类型

> Go语言的数值类型包括几种不同大小的整数、浮点数和复数。每种数值类型都决定了对应的大小范围和是否支持正负符号。

## 整型

> int8、int16、int32、int6 对应 8、16、32、64bit 大小的有符号整数。
>
> uint8、uint16、uint32、uint64 对应 8、16、32、64bit 大小的无符号整数。
>
> int和uint 对应 特定CPU平台机器字大小的有符号和无符号整数。

- **Unicode字符rune类型和int32等价**：通常用于表示一个Unicode码点。这两个名称可以互换使用。
- **byte和uint8类型等价**：一般用于强调数值是一个原始的数据而不是一个小的整数。
- **uintptr**：没有指定具体的bit大小但是足以容纳指针。只有在底层编程时才需要。
- **互转**：需要显式的类型转换操作。
- **溢出**：计算结果溢出，超出的高位的bit位部分将被丢弃。
- **比较**：布尔型、数字类型和字符串等基本类型都是可**比较**的。还可以根据比较结果**排序**。
- **无符号数**：只有在位运算或其它特殊的运算场景才会使用，就像bit集合、分析二进制文件格式或者是哈希和加密操作等。它们通常**并不用于仅仅是表达非负数**的场合。
- **改变数值或丢精度**：一个大尺寸的整数类型转为一个小尺寸的整数类型。
  - 浮点数到整数的转换将丢失任何小数部分，然后向数轴零方向截断。
- **进制**：
  - **八进制**：0；通常用于POSIX操作系统上的文件访问权限标志，
  - **十六进制**：0x；强调数字值的bit位模式。

## 浮点数

> float32
>
> - 最大数值：math.MaxFloat32，大约是 3.4e38。
> - 最小数值：1.4e-45。
>
> float64
>
> - 最大数值：math.MaxFloat64，大约是1.8e308。
> - 最小数值：4.9e-324。

- **NaN非数**：表示无效的除法操作结果`0/0或Sqrt(-1)`。

- **无穷大**：

  - *uvnan*：正无穷大。表示太大溢出的数字。
  - *uvneginf*：负无穷大。表示除零的结果。
  - Math.IsInf：判断是不是无穷大。

- **不唯一**：在浮点数中，NaN、正无穷大和负无穷大都不是唯一的。

  > 不可以用相等比较。

## 复数

> complex64和complex128，分别对应float32和float64两种浮点数精度。

```go
x := 1 + 2i
y := 3 + 4i
```

> 写法

- **real()**：获取实部。
- **imag()**：获取虚部 

## 布尔型

一个布尔类型的值只有两种：true和false。

**短路行为**：如果运算符左边值已经可以确定整个布尔表达式的值，那么运算符右边的值将不再被求值。

## 字符串

- **不可变**：一个字符串是一个不可改变的字节序列。修改字符串内部数据的操作也是被禁止的。

  > s[0] = 'L'。

- **panic**：访问超出字符串索引范围的字节。

- **字节&字符**：第i个字节并不一定是字符串的第i个字符，因为对于非ASCII字符的UTF8编码会要两个或多个字节。

### 字符串面值

> ```
> \a      响铃
> \b      退格
> \f      换页
> \n      换行
> \r      回车
> \t      制表符
> \v      垂直制表符
> \'      单引号（只用在 '\'' 形式的rune符号面值中）
> \"      双引号（只用在 "..." 形式的字符串面值中）
> \\      反斜杠
> ```
>
> Unicode码点。

- **原生**：一个原生的字符串面值形式是\`...`，使用反引号代替双引号。在原生的字符串面值中，没有转义操作。

  > 原生字符串面值用于编写正则表达式会很方便，因为正则表达式往往会包含很多反斜杠。
  >
  > 广泛应用于HTML模板、JSON面值、命令行提示信息、那些需要扩展到多行的场景。

- **Unicode**：收集了这个世界上所有的符号系统，包括重音符号和其它变音符号，制表符和回车符，还有很多神秘的符号，每个符号都分配一个唯一的Unicode码点，Unicode码点对应Go语言中的rune整数类型。

- **UTF-8**：UTF8是一个将Unicode码点编码为字节序列的变长编码。

  > ASCII部分字符只使用1个字节，常用字符部分使用2或3个字节表示。
  >
  > **unicode/utf8**包则提供了用于rune字符序列的UTF8编码和解码的功能。

  - *前缀、后缀、子串*：我们可以不用解码直接测试一个字符串是否是另一个字符串的前缀等。

  - *UTF8解码器*

    > ```go
    > s := "Hello, 世界"
    > 
    > utf8.RuneCountInString(s)
    > for i := 0; i < len(s); {
    >     r, size := utf8.DecodeRuneInString(s[i:])
    > }
    > ```

  - *循环统计字符串中字符的数目*

    > ```go
    > n := 0
    > for range s {
    >     n++
    > }
    > ```

  - *错误*：遇到一个错误的UTF8编码输入，将生成一个特别的Unicode字符`\uFFFD`，在印刷中这个符号通常是�。

  - *转换*

    > 将[]rune类型的Unicode字符slice或数组转为string，则对它们进行UTF8编码。
    >
    > 将整数转型为字符串意思是生成以只包含对应Unicode码点字符的UTF8字符串。

### 字符串和Byte切片

> - strings包提供了许多如字符串的查询、替换、比较、截断、拆分和合并等功能。
>
> - bytes包也提供了很多类似功能的函数，但是针对和字符串有着相同结构的[]byte类型。
>
> - strconv包提供了布尔型、整型数、浮点数和对应字符串的相互转换，还提供了双引号转义相关的转换。
>
> - unicode包提供了IsDigit、IsLetter、IsUpper和IsLower等类似功能，它们用于给字符分类。
>
> - path和path/filepath包提供了关于文件路径名更一般的函数操作。

```go
[]byte(s)
```

> 分配了一个新的字节数组用于保存字符串数据的拷贝，然后引用这个底层的字节数组。

## 常量

- **计算**：常量表达式的值在编译期计算，而不是在运行期。

  > 常量间的所有算术运算、逻辑运算和比较运算的结果也是常量。

- **缺省**：批量声明的常量。

  ```go
  const (
    a = 1
    b
    c = 2
    d
  )
  ```

  > 除了第一个外其它的常量右边的初始化表达式都可以省略。
  >
  > 省略表示使用前面常量的初始化表达式写法，对应的常量类型也一样的。

### iota 常量生成器

在一个const声明语句中，在第一个声明的常量所在的行，iota将会被置为0，然后在每一个有常量声明的行加一。

```go
const (
  Sunday Weekday = iota
  Monday
  Tuesday
  Wednesday
  Thursday
  Friday
  Saturday
)

const (
  FlagUp Flags = 1 << iota // is up
  FlagBroadcast            // supports broadcast access capability
  FlagLoopback             // is a loopback interface
  FlagPointToPoint         // belongs to a point-to-point link
  FlagMulticast            // supports multicast access capability
)
```

### 无类型常量

1. 编译器为这些没有明确基础类型的数字常量提供比基础类型更高精度的算术运算；你可以认为至少有256bit的运算精度。
2. 有六种未明确类型的常量类型，分别是无类型的布尔型、无类型的整数、无类型的字符、无类型的浮点数、无类型的复数、无类型的字符串。

## 数组

- **元素**：由固定长度的特定类型元素组成的序列，一个数组可以由零个或多个元素组成。

- **默认值**：数组的每个元素都被初始化为元素类型对应的零值。

- **长度缺省**：如果在数组的长度位置出现的是 `...` 省略号，则表示数组的长度是根据初始化值的个数来计算。

  ```go
  q := [...]int{1, 2, 3}
  fmt.Printf("%T\n", q) // "[3]int"
  ```

- **类型**：数组的长度是数组类型的一个组成部分，因此[3]int和[4]int是两种不同的数组类型。

- **长度**：长度必须是常量表达式，因为数组的长度需要在编译阶段确定。

- **其他写法**

  > ```go
  > const (
  >     USD Currency = iota // 美元
  >     EUR                 // 欧元
  >     GBP                 // 英镑
  >     RMB                 // 人民币
  > )
  > symbol := [...]string{USD: "$", EUR: "€", GBP: "￡", RMB: "￥"}
  > 
  > r := [...]int{99: -1}
  > ```

- **比较**：如果一个数组的元素类型是可以相互比较的，那么数组类型也是可以相互比较的。当两个数组的所有元素都是相等的时候数组才是相等的。

## Slice

> 我们一般使用slice来替代数组。

- **写法**：[]T，其中T代表slice中元素的类型。slice的语法和数组很像，只是没有固定长度而已。

- **组成**：指针、长度和容量。

- **共享**：多个slice之间可以共享底层的数据，并且引用的数组部分区间可能重叠。

- **切片操作**：超出cap(s)的上限将导致一个panic异常，但是超出len(s)则是意味着扩展了slice，因为新slice的长度会变大。

- **比较**：slice之间不能比较。可以用 `bytes.Equal` 函数比较字节型切片。

  > 第一个原因，一个slice的元素是间接引用的，一个slice甚至可以包含自身。
  >
  > 第二个原因，因为slice的元素是间接引用的，一个固定的slice值（译注：指slice本身的值，不是元素的值）在不同的时刻可能包含不同的元素，因为底层数组的元素可能会被修改。

- 一个零值的slice等于nil，空值的的slice不等于nil。

### append函数

> 用于向slice追加元素。

检测底层数组是否有足够的容量

- **有**：在原底层数组上扩展slice，新添加的元素复制到新扩展的空间，返回原来的指针。
- **没有**：会先分配一个足够大的slice用于保存新的结果，返回新的底层数组指针。

```go
runes = append(runes, r)
```

> 我们也不能确认新的slice和原始的slice是否引用的是相同的底层数组空间。同样，我们不能确认在原先的slice上的操作是否会影响到新的slice。因此，通常是将append返回的结果直接赋值给输入的slice变量。

```go
x = append(x, 4, 5, 6)
x = append(x, x...)
```

> append函数则可以追加多个元素，甚至追加一个slice。

## Map

- **类型**：所有的key都有相同的类型，所有的value也有相同的类型。但是key和value之间可以是不同的数据类型。

- **key**：必须是支持==比较运算符的数据类型。

- **value**：map中的元素并不是一个变量，不能对map的元素进行取址操作。

  ```go
  _ = &ages["bob"] // compile error: cannot take address of map element
  ```

  > 原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而可能导致之前的地址无效。

- **迭代**：Map的迭代顺序是不确定的，并且不同的哈希函数实现可能导致不同的遍历顺序。

- **排序**：可以使用sort包的Strings函数对字符串slice进行排序。

  - 使用“_”空白标识符来忽略第一个循环变量。

- **零值**：map类型的零值是nil，向一个nil值的map存入元素将导致一个panic异常。

- **是否存在**

  ```go
  age, ok := ages["bob"]
  ```

- **比较**：map之间也不能进行相等比较；唯一的例外是和nil进行比较。

- **set**：Go语言中并没有提供一个set类型，但是map中的key也是不相同的，可以用map实现类似set的功能。

- **ReadRune**：执行UTF-8解码并返回三个值：解码的rune字符的值，字符UTF-8编码后的长度，和一个错误值。

## 结构体

> 结构体是一种聚合的数据类型，是由零个或多个任意类型的值聚合成的实体。每个值称为结构体的成员。

- **访问成员**：结构体变量的成员可以通过 `.` 操作符访问。

- **定义**：通常一行对应一个成员，**名字在前类型在后**，如果相邻的成员类型相同可以被合并到一行。

  ```go
  type Employee struct {
      ID            int
      Name, Address string
      DoB           time.Time
  }
  ```

- **顺序**：结构体成员的输入顺序也有重要的意义。

- **导出**：结构体成员名字以大写字母开头。一个结构体可能**同时包含**导出和未导出的成员。

- **非递归**：一个聚合的值不能包含它自身。

- **零值**：结构体类型的零值是每个成员都是零值。通常会将零值作为最合理的默认值。

- **效率**：较大的结构体通常会用指针的方式传入和返回。

- **比较**：如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的。

  > 可比较的结构体可用于map的key。

- **匿名成员**：只声明一个成员对应的数据类型而不指名成员的名字。匿名成员的数据类型必须是命名的类型或指向一个命名的类型的指针。

  > 我们可以直接访问叶子属性而不需要给出完整的路径。
  >
  > 结构体字面值必须遵循形状类型声明时的结构。
  >
  > 外层的结构体不仅仅是获得了匿名成员类型的所有成员，而且也获得了该类型导出的全部的方法。

## json

- **json.Marshal**：返回字符串，并且没有空白缩进。
- **json.MarshalIndent**：产生整齐缩进的输出。有两个额外的字符串参数用于表示每一行输出的前缀和每一个层级的缩进。
- **json.Unmarshal**：将JSON格式的字符串解码为字节slice。
- **json.Decoder**：从一个输入流解码JSON数据。

## 文本和HTML模板

- **.**：当前值“.”最初被初始化为调用模板时的参数。
- **.field**：结构体中 field 成员以默认的方式打印的值。
- **range end**：对应一个循环action。
- **|**：表示将前一个表达式的结果作为后一个函数的输入，类似于UNIX中管道的概念。

# 函数

## 函数声明

```go
func name(parameter-list) (result-list) {
  body
}
```

> 函数声明包括**函数名**、**形式参数列表**、**返回值列表**（可省略）以及**函数体**。
>
> - 形式参数列表描述了函数的参数名以及参数类型。
> - 返回值列表描述了函数返回值的变量名以及类型。
>
> - 如果一组形参或返回值有相同的类型，我们不必为每个形参都写出参数类型。

- **函数类型**

  > 函数的类型被称为 **函数的签名**。

  ```go
  func f(i, j, k int, s, t string)                 { /* ... */ }
  func f(i int, j int, k int,  s string, t string) { /* ... */ }
  ```

  > 如果两个函数形式参数列表和返回值列表中的变量类型一一对应，那么这两个函数被认为有相同的类型或签名。

- **无参数默认值**

  > 在函数调用时，Go语言**没有默认参数值**，也**不能通过参数名指定形参**，因此形参和返回值的变量名对于函数调用者而言没有意义。

## 多返回值

1. 当调用接受多参数的函数时，可以将一个返回多参数的函数调用作为该函数的参数。

2. 准确的变量名可以传达函数返回值的含义。尤其在返回值的类型都相同时。

   ```go
   func Size(rect image.Rectangle) (width, height int)
   func Split(path string) (dir, file string)
   func HourMinSec(t time.Time) (hour, minute, second int)
   ```

3. 按照惯例，函数的**最后一个bool类型**的返回值表示函数**是否运行成功**，**error类型**的返回值代表**函数的错误信息**，对于这些类似的惯例，我们不必思考合适的命名，它们都无需解释。

4. 如果一个函数所有的返回值都有显式的变量名，那么该函数的return语句可以**省略操作数**。这称之为**bare return**。

5. 当一个函数有多处return语句以及许多返回值时，bare return 可以减少代码的重复，但是使得代码难以被理解。**不宜过度使用bare return**。

## 错误

将运行失败看作是预期结果的函数，它们会返回一个额外的返回值，通常是最后一个，来传递错误信息。

> 1. 如果导致失败的原因只有一个，额外的返回值可以是一个布尔值，通常被命名为ok。
> 2. 通常，导致失败的原因不止一种，因此，额外的返回值不再是简单的布尔类型，而是error类型。

- **错误处理策略**：5种方式

  1. *传播错误*。使用 `fmt.Errorf` 函数添加额外的上下文信息到原始错误信息。

     ```go
     doc, err := html.Parse(resp.Body)
     resp.Body.Close()
     if err != nil {
         return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
     }
     ```

  2. *重试*。偶然性错误或不可预知问题。重新尝试失败的操作。

     ```go
     func WaitForServer(url string) error {
       const timeout = 1 * time.Minute
       deadline := time.Now().Add(timeout)
       for tries := 0; time.Now().Before(deadline); tries++ {
         _, err := http.Head(url)
         if err == nil {
           return nil // success
         }
         log.Printf("server not responding (%s);retrying…", err)
         time.Sleep(time.Second << uint(tries)) // exponential back-off
       }
       return fmt.Errorf("server %s failed to respond after %s", url, timeout)
     }
     ```

  3. *输出错误信息并结束程序*：错误发生后，程序无法继续运行。

     > 注意，这种策略只应在main中执行。
     >
     > 调用log.Fatalf可以更简洁的代码达到与上文相同的效果。log中的所有函数，都默认会在错误信息之前输出时间信息。

     ```go
     if err := WaitForServer(url); err != nil {
       log.Fatalf("Site is down: %v\n", err)
       os.Exit(1)
     }
     ```

  4. *只输出错误信息*：通过log包提供的函数。

     ```go
     if err := Ping(); err != nil {
       log.Printf("ping failed: %v; networking disabled",err)
     }
     ```

  5. *忽略错误*：当你决定忽略某个错误时，你应该清晰地写下你的意图。

     ```go
     os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodically
     ```

- **文件结尾错误**

  >io包保证任何由文件结束引起的读取失败都返回同一个错误——io.EOF，该错误在io包中定义。

  ```go
  var EOF = errors.New("EOF")
  ```

## 函数值

1. **函数被看作第一类值**【first-class values】：函数像其他值一样，拥有类型，可以被赋值给其他变量，传递给函数，从函数返回。

2. **函数类型的零值是nil**：调用值为nil的函数值会引起panic错误，函数值可以与nil比较：

   ```go
   var f func(int) int
   f(3) // 此处f的值为nil, 会引起panic错误
   ```

3. **函数值之间不可比较**：也不能用函数值作为map的key。

4. 函数值使得我们不仅仅可以**通过数据来参数化函数**，亦可**通过行为**。

## 匿名函数

1. **有函数名的函数**：只能在包级语法块中被声明，通过函数字面量【function literal】。
2. **函数值字面量**：是一种表达式，它的值被称为匿名函数【anonymous function】，区别在于func关键字后没有函数名。
3. **匿名函数访问**：可以访问完整的词法环境，在函数中定义的内部函数可以引用该函数的变量。
4. **函数值**：Go使用闭包【closures】技术实现函数值，Go程序员也把函数值叫做闭包。

### 警告：捕获迭代变量

> *循环变量的作用域*：在循环中生成的所有函数值都共享相同的循环变量。
>
> 函数值中记录的是循环变量的内存地址，而不是循环变量某一时刻的值。
>
> 通常，为了解决这个问题，我们会**引入一个与循环变量同名的局部变量**，作为循环变量的副本。

## 可变参数

1. 声明可变参数函数时，需要在参数列表的最后一个参数类型之前加上省略符号“...”，这表示该函数会接收任意数量的该类型参数。
2. 函数名的后缀f是一种通用的命名规范，代表该可变参数函数可以接收Printf风格的格式化字符串。

## Deferred函数

1. **语法**：在调用普通函数或方法前加上关键字defer。

2. **执行时间**：直到包含该defer语句的函数执行完毕时，defer后的函数才会被执行。【**defer 必须要能执行到**】

3. **执行情况**：通过return正常结束，或由于panic导致的异常结束，都会执行defer。

4. **执行顺序**：可以在一个函数中执行多条defer语句，执行顺序与声明顺序相反。

5. **使用场景**：经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。

6. **修改返回值**：对匿名函数采用defer机制，可以使其观察函数的返回值。延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值。

   ```go
   func double(x int) (result int) {
       defer func() { fmt.Printf("double(%d) = %d\n", x,result) }()
       return x + x
   }
   _ = double(4)
   // Output:
   // "double(4) = 8"
   
   func triple(x int) (result int) {
       defer func() { result += x }()
       return double(x)
   }
   fmt.Println(triple(4)) // "12"
   ```

7. **循环中defer**：解决方法是将循环体中的defer语句移至另外一个函数。在每次循环时，调用这个函数。

## Panic异常

>  Go的类型系统会在编译时捕获很多错误，但有些错误只能在运行时检查，如数组访问越界、空指针引用等。这些运行时错误会引起painc异常。

1. **panic异常发生**：程序会中断运行，并立即执行在该goroutine中被延迟的函数。

2. **使用场景**：panic一般用于严重错误，如程序内部的逻辑不一致。所以对于大部分漏洞，我们应该使用Go提供的错误机制，而不是panic，尽量避免程序的崩溃。

3. **输出堆栈**：runtime包允许程序员输出堆栈信息。

   ```go
   func main() {
       defer printStack()
       f(3)
   }
   func printStack() {
       var buf [4096]byte
       n := runtime.Stack(buf[:], false)
       os.Stdout.Write(buf[:n])
   }
   ```

   > 在Go的panic机制中，延迟函数的调用在释放堆栈信息之前。

## Recover捕获异常

1. 在deferred函数内部，panic value被附加到错误信息中；并用err变量接收错误信息，返回给调用者。
3. 为了标识某个panic是否应该被恢复，我们可以将panic value设置成特殊类型。在recover时对panic value进行检查，如果发现panic value是特殊类型，就将这个panic作为error处理，如果不是，则按照正常的panic进行处理。

# 方法

## 方法声明

1. **接收器**

   ```go
   func (p Point) Distance(q Point) float64 {
       return math.Hypot(q.X-p.X, q.Y-p.Y)
   }
   ```

   > p 为方法的接收器。

2. **选择器**：`p.Distance`。

3. **定义局限**：能够给任意类型定义方法，只要这个命名类型的底层类型不是指针或者interface。

4. **方法名**：单类型内部的方法名都必须唯一，但不同的类型可以有同样的方法名。

5. **优势**：方法名可以简短。当我们在包外调用的时候这种好处就会被放大，因为我们可以使用这个短名字，而可以省略掉包的名字。

## 指针接收器

1. **隐式类型转化**：不管method的receiver是指针类型还是非指针类型，都是可以通过指针/非指针类型进行调用的，编译器会帮你做类型转换。
2. **避免拷贝**：如果一个方法使用指针作为接收器，你需要避免对其进行拷贝。
3. **nil**：nil也是一个合法的接收器类型。

## 嵌入结构体扩展类型

1. **内嵌**：内嵌可以使我们定义字段特别多的复杂类型，我们可以将字段先按小类型分组，然后定义小类型的方法，之后再把它们组合起来。
2. **命名类型的指针**：在类型中内嵌的匿名字段也可能是一个命名类型的指针，这种情况下字段和方法会被间接地引入到当前的类型中。
3. **同名方法**：选择器有二义性的话编译器会报错，比如你在同一级里有两个同名的方法。

## bit数组

> 在数据流分析领域，集合元素通常是一个非负整数，集合会包含很多元素，并且集合会经常进行并集、交集操作，这种情况下，bit数组会比map表现更加理想。

## 封装

> 一个对象的变量或者方法如果对调用方是不可见的话，一般就被定义为“封装”。封装有时候也被叫做信息隐藏，同时也是面向对象编程最关键的一个方面。

1. **控制可见性**：大写首字母的标识符会从定义它们的包中被导出，小写字母的则不会。
2. **优点**
   1. 首先，因为调用方不能直接修改对象的变量值，其只需要关注少量的语句并且只要弄懂少量变量的可能的值即可。
   2. 第二，隐藏实现的细节，可以防止调用方依赖那些可能变化的具体实现，这样使设计包的程序员在不破坏对外的api情况下能得到更大的自由。
   3. 阻止了外部调用方对对象内部的值任意地进行修改。
3. **注意**：Go的编码风格不禁止直接导出字段。当然，一旦进行了导出，就没有办法在保证API兼容的情况下去除对其的导出，所以在一开始的选择一定要经过深思熟虑并且要考虑到包内部的一些不变量的保证，未来可能的变化，以及调用方的代码质量是否会因为包的一点修改而变差。

# 接口

> 接口类型只能定义方法，不会暴露出它所代表的对象的内部值的结构和这个对象支持的基础操作的集合。

## 接口类型

1. 只能定义方法，不会暴露出它所代表的对象的内部值的结构和这个对象支持的基础操作的集合。

2. 具体描述了一系列方法的集合，一个实现了这些方法的具体类型是这个接口类型的实例。

3. Go语言中的命名习惯是单方法接口。

4. *接口内嵌*

   ```go
   type ReadWriter interface {
     Reader
     Writer
   }
   type ReadWriter interface {
     Read(p []byte) (n int, err error)
     Write(p []byte) (n int, err error)
   }
   type ReadWriter interface {
     Read(p []byte) (n int, err error)
     Writer
   }
   ```

   > 3种定义方式都是一样的效果。

5. 接口类型封装和隐藏具体类型和它的值。

### 类型断言

> 类型断言是一个使用在接口值上的操作。
>
> 对一个接口类型的类型断言改变了类型的表述方式，改变了可以获取的方法集合（通常更大），但是它保留了接口值内部的动态类型和值的部分。
>
> ```go
> x.(T)
> ```

- **类型T是一个具体类型**

  > 类型断言检查x的动态类型是否和T相同。
  >
  > - 检查成功：结果是x的动态。
  >
  > - 检查失败：
  >
  >   - 单赋值抛出panic。
  >
  >   - ```go
  >     f, ok := w.(*os.File) 
  >     ```
  >
  >     > 不会抛出panic。

- **断言的类型T是一个接口类型**

  > 类型断言检查是否x的动态类型满足T。
  >
  > - 检查成功：结果是一个有相同动态类型和值部分的接口值，但是结果为类型T。
  > - 检查失败：同上。

- **x 是一个 nil 接口值**

  > 那么不论被断言的类型是什么这个类型断言都会失败。

- **if 中变量名重用**

  ```go
  if w, ok := w.(*os.File); ok {
      // ...use w...
  }
  ```

  > 其实是声明了一个同名的新的本地变量，外层原来的w不会被改变。

### 询问行为

```go
package fmt

func formatOneValue(x interface{}) string {
  if err, ok := x.(error); ok {
    return err.Error()
  }
  if str, ok := x.(Stringer); ok {
    return str.String()
  }
  // ...all other types...
}
```

### 类型分支

```go
switch x.(type) {
case nil:       // ...
case int, uint: // ...
case bool:      // ...
case string:    // ...
default:        // ...
}
```

## 接口值

```go
var w io.Writer
```

> <img src="./imgs/接口初始化.png" alt="img" style="zoom:75%;" />

```go
w = os.Stdout
```

> <img src="./imgs/File接口.png" alt="img" style="zoom:75%;" />

```go
w = new(bytes.Buffer)
w = nil
```

> ![img](./imgs/buffer接口.png)

```go
w = nil
```

> <img src="./imgs/接口初始化.png" alt="img" style="zoom:75%;" />

- **组成**：由两个部分组成，一个具体的类型和那个类型的值。它们被称为接口的动态类型和动态值。

- **零值**：对于一个接口的零值就是它的类型和值的部分都是nil。

- **panic**：调用一个空接口值上的任意方法都会产生panic。

- **比较**：

  - 如果两个接口值的动态类型相同，但是这个动态类型是不可比较的（比如切片），将它们进行比较就会失败并且panic。
  - 两个接口值相等仅当它们都是nil值，或者它们的动态类型相同并且动态值也根据这个动态类型的==操作相等。
  - 因为接口值是可比较的，所以它们可以用在map的键或者作为switch语句的操作数。

- **警告**：一个包含nil指针的接口不是nil接口。

  ```go
  var buf *bytes.Buffer
  f(buf)
  func f(out io.Writer) {
    // ...do something...
    if out != nil {
      out.Write([]byte("done!\n"))
    }
  }
  ```

  > 会发生panic，因为接口不为nil，但接口值为nil。
  >
  > ![img](./imgs/包含nil指针的接口.png)

## 接口实现

> 一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。

1. 接口指定的规则非常简单：表达一个类型属于某个接口只要这个类型实现这个接口。
2. 因为空接口类型对实现它的类型没有要求，所以我们可以将任意一个值赋给空接口类型。
3. 因为接口与实现只依赖于判断两个类型的方法，所以没有必要定义一个具体类型和它实现的接口之间的关系。
4. 非空的接口类型比如io.Writer经常被指针类型实现，尤其当一个或多个接口方法像Write方法那样隐式的给接收者带来变化的时候。一个结构体的指针是非常常见的承载方法的类型。
5. 在Go语言中我们可以在需要的时候定义一个新的抽象或者特定特点的组，而不需要修改具体类型的定义。

## 举例

### flag.Value

flag.Duration函数创建一个time.Duration类型的标记变量并且允许用户通过多种用户友好的方式来设置这个变量的大小，这种方式还包括和String方法相同的符号排版形式。

```go
var period = flag.Duration("period", 1*time.Second, "sleep period")
```

> ```shell
> # Sleeping for 50ms...
> ./sleep -period 50ms
> # Sleeping for 2m30s...
> ./sleep -period 2m30s
> # Sleeping for 1h30m0s...
> ./sleep -period 1.5h
> # invalid value "1 day" for flag -period: time: invalid duration 1 day
> ./sleep -period "1 day"
> ```

我们为我们自己的数据类型定义新的标记符号是简单容易的。我们只需要定义一个实现flag.Value接口的类型。

### sort.Interface

```go
package sort

type Interface interface {
  Len() int
  Less(i, j int) bool // i, j are indices of sequence elements
  Swap(i, j int)
}
```

- sort包为[]int、[]string和[]float64的正常排序提供了特定版本的函数和类型。
- sort.Sort：排序。【快排】
- sort.Stable：稳定排序【插入排序】。
- sort.Reverse：sort包定义了一个不公开的struct类型reverse，它嵌入了一个sort.Interface。reverse的Less方法调用了内嵌的sort.Interface值的Less方法，但是通过交换索引的方式使排序结果变成逆序。

### http.Handler

- **ListenAndServe**

  > 需要一个例如“localhost:8000”的服务器地址，和一个所有请求都可以分派的Handler接口实例。
  >
  > ```go
  > func ListenAndServe(address string, h Handler) error
  > ```
  >
  > 更真实的服务器会定义多个不同的URL，每一个都会触发一个不同的行为。
  >
  > ```go
  > func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  >   switch req.URL.Path {
  >     case "/list":
  >     fmt.Println(111)
  >     case "/price":
  >     fmt.Println(222)
  >     default:
  >     w.WriteHeader(http.StatusNotFound) // 404
  >     fmt.Fprintf(w, "no such page: %s\n", req.URL)
  >   }
  > }
  > ```

- **ServeMux**

  > 以通过组合来处理更加错综复杂的路由需求。
  >
  > ```go
  > mux := http.NewServeMux()
  > mux.Handle("/list", http.HandlerFunc(db.list))
  > mux.Handle("/price", http.HandlerFunc(db.price))
  > log.Fatal(http.ListenAndServe("localhost:8000", mux))
  > ```

- **Handler** 和 **HandlerFunc**

  > Handler是个接口，必须实现ServeHTTP方法。
  >
  > HandlerFunc是一个让函数值满足一个接口的适配器，这里函数和这个接口仅有的方法有相同的函数签名。

- **DefaultServerMux**

  > net/http包提供了一个全局的ServeMux实例DefaultServerMux和包级别的http.Handle和http.HandleFunc函数。
  >
  > 服务器的主函数可以如下简化。
  >
  > ```go
  > db := database{"shoes": 50, "socks": 5}
  > http.HandleFunc("/list", db.list)
  > http.HandleFunc("/price", db.price)
  > log.Fatal(http.ListenAndServe("localhost:8000", nil))
  > ```

### error

```go
type error interface {
  Error() string
}
```

- **errors.New**

  ```go
  package errors
  func New(text string) error { return &errorString{text} }
  type errorString struct { text string }
  func (e *errorString) Error() string { return e.text }
  ```

- **fmt.Errorf**：更加方便，它还会处理字符串格式化。

  ```go
  func Errorf(format string, args ...interface{}) error {
    return errors.New(Sprintf(format, args...))
  }
  ```

- **syscall.Errno**：syscall包提供了Go语言底层系统调用API。在多个平台上，它定义一个实现error接口的数字类型Errno，并且在Unix平台上，Errno的Error方法会从一个字符串表中查找错误消息。

  ```go
  package syscall
  
  type Errno uintptr // operating system error code
  
  var errors = [...]string{
    1:   "operation not permitted",   // EPERM
    2:   "no such file or directory", // ENOENT
    3:   "no such process",           // ESRCH
    // ...
  }
  
  func (e Errno) Error() string {
    if 0 <= int(e) && int(e) < len(errors) {
      return errors[e]
    }
    return fmt.Sprintf("errno %d", e)
  }
  ```

- **os.PathError**

  ```go
  package os
  
  // PathError records an error and the operation and file path that caused it.
  type PathError struct {
    Op   string
    Path string
    Err  error
  }
  
  func (e *PathError) Error() string {
    return e.Op + " " + e.Path + ": " + e.Err.Error()
  }
  
  import (
    "errors"
    "syscall"
  )
  
  var ErrNotExist = errors.New("file does not exist")
  
  // IsNotExist returns a boolean indicating whether the error is known to
  // report that a file or directory does not exist. It is satisfied by
  // ErrNotExist as well as some syscall errors.
  func IsNotExist(err error) bool {
    if pe, ok := err.(*PathError); ok {
      err = pe.Err
    }
    return err == syscall.ENOENT || err == ErrNotExist
  }
  ```

## 建议

1. 接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要。
2. 当一个接口只被一个单一的具体类型实现时有一个例外，就是由于它的依赖，这个具体类型不能和这个接口存在在一个相同的包中。这种情况下，一个接口是解耦这两个包的一个好方式。

# Goroutines

> 每一个并发的执行单元叫作一个goroutine。
>
> `go` 语句会使其语句中的函数在一个新创建的goroutine中运行。

## 动态栈

> - 每一个OS线程都有**2MB的内存块**来做栈，这个栈用来存储正在被调用或挂起（指在调用其它函数时）的函数的内部变量。
> - 一个goroutine会以一个很小的栈开始其生命周期，一般只需要2KB。保存其活跃或挂起的函数调用的本地变量。
> - 一个goroutine的栈大小并不是固定的；栈的大小会根据需要动态地伸缩。最大值有1GB。

## Goroutine调度

- **上下文切换**：不需要进入内核的上下文，所以重新调度一个goroutine比调度一个线程**代价低得多**。
- **m:n线程**：在n个操作系统线程上多工（调度）m个goroutine。
- **休眠**：当Goroutine执行时间较长时，GO调度器会将其置为休眠。
- **优势**：榨干cpu的油水，在一个线程中的goroutine阻塞时，不进行线程切换，直接执行另一个goroutine。

## GOMAXPROCS

- **默认值**：运行机器上的**CPU的核心数**。
- **修改**
  - **环境变量**：`GOMAXPROCS`
  - **程序函数**：`runtime.GOMAXPROCS`

## Goroutine没有ID号

> goroutine没有可以被程序员获取到的身份（id）的概念。这一点是设计上故意而为之，由于thread-local storage总是会被滥用。

# Channels

> - channels是它们goroutine的**通信机制**。
>
> - **类型**：每个channel都有一个特殊的类型，也就是channels可发送数据的类型。channel是**引用类型**。
>
> - **发送**：channel对象 <- 值
>   
> - **接收**：值 <- channel对象。
>   
> - **关闭**：close(channel)
>   
>   - **不需要关闭每一个channel**。只有当需要告诉接收者goroutine，所有的数据已经全部发送时才需要关闭channel。
>   - **重复关闭**一个channel将导致**panic异常**；关闭一个nil值的channel也将导致**panic异常**。
>   - **发送**将导致panic异常。
>   - 可以**接收**到之前已经成功发送的数据，空channel会接受到零值。
>   
> - **range**
>
>   > 当channel被关闭并且没有值可接收时跳出循环。
>
> - **广播**
>
>   > 不要向channel发送值，而是用关闭一个channel来进行广播。

- **不带缓存的Channels**

  - **发送**：在值没被接收时导致发送者goroutine阻塞。

  - **接收**：在值还没发送时导致接收者goroutine阻塞。

  - **同步Channels**：因为会导致两个goroutine做一次同步操作。

  - **goroutines泄漏**：慢的goroutines因为没有人接收而被永远卡住。泄漏的goroutines并不会被自动回收。

- **串联的Channels（Pipeline）**

  > Channels也可以用于将多个goroutine连接在一起，一个Channel的输出作为下一个Channel的输入。


- **单方向的Channel**

  - **只能发送**

    ```go
    chan<- [type]
    ```

  - **只能接收**：

    ```go
    <-chan [type]
    ```

    **无法close**【编译错误】。

- **带缓存的Channels**

  - **队列**：内部持有元素队列。**最大容量**在**make**时创建。

    > *发送*：尾部插入元素
    >
    > *接收*：头部删除元素

  - **cap函数**：获取channel内部缓存的容量。

  - **len函数**：获取内部缓存队列中有效元素的个数。

  - **性能**：缓存也可能影响程序的性能。
  
  - 可用于做**计数信号量**。


## sync.WaitGroup

- **Add**：必须在worker goroutine开始之前调用。
- **Done**：和Add(-1)是等价的。
- 一般用来等待所有的`Goroutines`返回。

## Select

```go
select {
  case <-ch1:
  // ...
  case x := <-ch2:
  // ...use x...
  case ch3 <- y:
  // ...
  default:
  // ...
}
```

- **多路复用**：每个case代表一个通信操作。

- **接收表达式**：只包含接收表达式自身，或者包含在一个简短的变量声明中。

- **case就绪**：通信并执行case之后的语句；这时候其它通信是不会执行的。

- **多个case同时就绪**：随机地选择一个执行，每一个channel都有平等的机会。

- **channel零值**：对一个nil的channel发送和接收操作会**永远阻塞**。

  > 可以用nil来**激活**或者**禁用**case。来达成处理其它输入或输出事件时超时和取消的逻辑。

# 共享变量

## 竞争条件

> 指的是程序在多个goroutine交叉执行操作时，没有给出正确的结果。

- **非常恶劣**：非常难以复现而且难以分析诊断。

- **数据竞争**：无论任何时候，只要有两个goroutine并发访问同一变量，且至少其中的一个是写操作的时候就会发生数据竞争。

  > 根本就没有良性数据竞争。我们一定要避免数据竞争。

  - *避免方式*

    > - 不要去写变量：尽量能初始化解决。
    > - 避免从多个goroutine访问变量。**不要使用共享数据来通信；使用通信来共享数据**
    > - 允许很多goroutine去访问变量，但是在**同一个时刻最多只有一个**goroutine在访问。

### 检测

```shell
go build -race
go run -race
go test -race
```

> 创建一个应用修改版，附带了能够记录所有运行期对共享变量访问工具的test。并且会记录下每一个读或者写共享变量的goroutine的身份信息。
>
> 工具会打印一份报告，内容包含变量身份，读取和写入的goroutine中活跃的函数的调用栈。这些信息在定位问题时通常很有用。
>
> 加了竞争检测的程序跑起来会**慢**一些，且需要**更大的内存**。

## sync.Mutex

- 如果可能的话**尽量使用defer**来将临界区扩展到函数的结束。

- **不能重入**：没法对一个已经锁上的 mutex 来再次上锁——这会导致**程序死锁**。

- **通用的解决方案**：将一个函数分离为多个函数

  > 一个不导出的函数，这个函数假设锁总是会被保持并去做实际的操作。
  >
  > 一个导出的函数，调用不导出的函数，但在调用前会先去获取锁。
  >
  > ```go
  > func Withdraw(amount int) bool {
  >      mu.Lock()
  >      defer mu.Unlock()
  >      deposit(-amount)
  >      if balance < 0 {
  >        deposit(amount)
  >        return false // insufficient funds
  >      }
  >      return true
  > }
  > 
  > func Deposit(amount int) {
  >      mu.Lock()
  >      defer mu.Unlock()
  >      deposit(amount)
  > }
  > 
  > func Balance() int {
  >      mu.Lock()
  >      defer mu.Unlock()
  >      return balance
  > }
  > 
  > // This function requires that the lock be held.
  > func deposit(amount int) { balance += amount }
  > ```

## sync.RWMutex

- RLock只能在临界区共享变量**没有任何写入操作**时可用。
- RWMutex只有当获得锁的大部分goroutine都是读操作，而锁在竞争条件下，也就是说，goroutine们必须等待才能获取到锁的时候，RWMutex才是最能带来好处的。
- RWMutex需要更复杂的内部记录，**会比mutex慢一些**。

## 内存同步

```go
var x, y int
go func() {
  x = 1 // A1
  fmt.Print("y:", y, " ") // A2
}()
go func() {
  y = 1                   // B1
  fmt.Print("x:", x, " ") // B2
}()
```

> 可能存在如下输出：
>
> ```go
> x:0 y:0
> y:0 x:0
> ```
>
> 原因可能是：
>
> 1. 处理器本地缓存没来的及同步到主存，两个处理器看到的不一样。
> 2. 指令重排序。
>
> 解决方案：
>
> 1. 将变量限制在goroutine内部。
> 2. 使用互斥条件来访问。

- 现代计算机中有一堆处理器，每一个都有其**本地缓存**「local cache」。为了效率，对内存的写入一般会在每一个处理器中缓冲，并在必要时一起flush到主存。
- **channel通信**或者**互斥量操作** *原语* 会使处理器聚集写入flush并commit。

## sync.Once

```go
func loadIcons() {
  icons = map[string]image.Image{
    "spades.png":   loadIcon("spades.png"),
    "hearts.png":   loadIcon("hearts.png"),
    "diamonds.png": loadIcon("diamonds.png"),
    "clubs.png":    loadIcon("clubs.png"),
  }
}

// NOTE: not concurrency-safe!
func Icon(name string) image.Image {
  if icons == nil {
    loadIcons() // one-time initialization
  }
  return icons[name]
}
```

> 因为缺少显式的同步，编译器和CPU是可以随意地去更改访问内存的指令顺序，以任意方式，只要保证每一个goroutine自己的执行顺序一致。
>
> ```go
> func loadIcons() {
>      icons = make(map[string]image.Image)
>      icons["spades.png"] = loadIcon("spades.png")
>      icons["hearts.png"] = loadIcon("hearts.png")
>      icons["diamonds.png"] = loadIcon("diamonds.png")
>      icons["clubs.png"] = loadIcon("clubs.png")
> }
> ```
>
> 这就导致一个初始化后的 goroutine 可能仍然读不到数据。

- 使用互斥访问icons的代价就是没有办法对该变量进行并发访问。
- **sync.Once**：解决这种一次性初始化的问题。
- **RWMutex**：太复杂且容易出错。

# 包

> 每个包都定义一个不同的名字空间用于它内部的每个标识符的访问。
>
> 每个包还通过控制包内名字的可见性和是否导出来实现封装特性。
>
> 当我们修改了一个源文件，我们必须重新编译该源文件对应的包和所有依赖该包的其他包。

- **编译速度快**

  > 1. 所有导入的包必须在每个文件的开头**显式声明**，编译器就有必要读取和分析整个源文件来判断包的依赖关系。
  > 2. 禁止包的**环状依赖**，包的依赖关系是一个有向无环图，每个包可以被独立编译，而且很可能是被**并发编译**。
  > 3. 编译后包的目标文件不仅仅记录包本身的导出信息，目标文件同时还记录了包的依赖关系。因此，在编译一个包的时候，编译器只需要读取每个直接导入包的目标文件，而不需要遍历所有依赖的的文件。

- **导入路径**

  > 如果你计划分享或发布包，那么导入路径最好是全球唯一的。为了避免冲突，所有非标准库包的导入路径建议以所在组织的互联网域名为前缀；而且这样也有利于包的检索。

- **包声明**

  > 1. 在每个Go语言源文件的开头都必须有包声明语句。包声明语句的主要目的是确定当前包被其它包导入时默认的标识符（也称为包名）。
  > 2. 默认的包名就是包导入路径名的最后一段。
  > 3. 有三种**例外**
  >    1. 包对应一个可执行程序，也就是main包，这时候main包本身的导入路径是无关紧要的。
  >    2. 包所在的目录中可能有一些文件名是以`_test.go`为后缀的Go源文件（译注：前面必须有其它的字符，因为以`_`或`.`开头的源文件会被构建工具忽略）。
  >    3. 一些依赖版本号的管理工具会在导入路径后追加版本号信息。

- **导入声明**

  > 1. **导入包的重命名**：同时导入两个名字相同的包，例如math/rand包和crypto/rand包，那么导入声明必须至少为一个同名包指定一个新的包名以避免冲突。
  >    1. 导入包的重命名只影响当前的源文件。
  >    2. 简短名称会更方便。选择用简短名称重命名导入包时候**最好统一**。
  >    3. 帮助避免和本地普通变量名产生冲突。
  > 2. 每个导入声明语句都明确指定了当前包和被导入包之间的依赖关系。

- **匿名导入**

  > 1. 导入一个包而并不使用导入的包将会导致一个编译错误。可以用下划线`_`来重命名导入的包。
  > 2. 只计算包级变量的初始化表达式和执行导入包的init初始化函数。

- **包和命名**

  > 1. 一般要用短小的包名，但也不能太短导致难以理解。
  > 2. 尽可能让命名有描述性且无歧义。
  > 3. 包名一般采用单数的形式。
  > 4. 要避免包名有其它的含义。
  > 5. 当设计一个包的时候，需要考虑包名和成员名两个部分如何很好地配合。
  > 6. 包中最重要的成员名字要简单明了。

## 工具

```shell
build            编译包和依赖项
clean            删除对象文件
doc              显示包或符号的文档
env              打印go环境信息
fmt              在包源上运行gofmt
get              下载并安装软件包和依赖项
install          编译并安装软件包和依赖项
list             列出包
run              编译并运行Go程序
test             测试包
version          打印go版本信息
vet              在软件包上运行go tool vet
```

- **工作区结构**

  - *GOPATH*：指定当前工作目录。

    > - src子目录用于存储源代码。
    >
    > - pkg子目录用于保存编译后的包的目标文件。
    > - bin子目录用于保存编译后的可执行程序。

  - *GOROOT*：指定Go的安装目录，还有它自带的标准库包的位置。

- **下载包**

  - `go get`可以下载一个单一的包或者用`...`下载整个子目录里面的每个包。

  - 一旦`go get`命令下载了包，然后就是安装包或包对应的可执行的程序。

  - `go get`命令支持当前流行的托管网站GitHub、Bitbucket和Launchpad，可以直接向它们的版本控制系统请求代码。

  - `go get`命令获取的代码是真实的本地存储仓库，而不仅仅只是复制源文件，因此你依然可以使用版本管理工具比较本地代码的变更或者切换到其它的版本。

  - 如果指定`-u`命令行标志参数，`go get`命令将确保所有的包和依赖的包的版本都是最新的，然后重新编译和安装它们。如果不包含该标志参数的话，而且如果包已经在本地存在，那么代码将不会被自动更新。

  - 导入路径含有的网站域名和本地Git仓库对应远程服务地址并不相同。

    ```html
    <meta name="go-import" content="golang.org/x/net git https://go.googlesource.com/net">
    ```

- **构建包**

  - `go build`命令编译命令行参数指定的每个包。
    - **包是库**：则**忽略输出结果**；这可以用于**检测**包是可以**正确编译**的。
    - **包名字是main**：`go build`调用链接器在当前目录创建一个可执行程序；以导入路径的最后一段作为可执行程序的名字。
    - `go build`：构建指定的包和它依赖的包，然后丢弃除了最后的可执行文件之外所有的中间编译结果。
    - `go install`：和`go build`命令很相似，但是它会保存每个包的编译成果，而不是将它们都丢弃。被编译的包会被保存到$GOPATH/pkg目录下，目录路径和 src目录路径对应，可执行程序被保存到$GOPATH/bin目录。
    - `go build -i`：将安装每个目标所依赖的包。
  - *指定包*
    - 绝对路径。
    - 相对路径，必须以`.`或`..`开头。
    - 不指定，当前目录对应的包。

- **包文档**

  > Go语言中的文档注释一般是完整的句子。
  >
  > 1. 第一行通常是摘要说明，以被注释者的名字开头。
  > 2. 注释中**函数参数**或**其它标识符**并不需要额外的引号或其它标记注明。
  >
  > 如果包的注释内容比较长，一般会放到一个独立的源文件中。

  - *go doc*：打印其后所指定的实体的声明与文档注释。

    > 并不需要输入完整的包导入路径或正确的大小写。

    - ```shell
      go doc time
      ```

      >包
      >
      >```shell
      >package time // import "time"
      >
      >Package time provides functionality for measuring and displaying time.
      >
      >const Nanosecond Duration = 1 ...
      >func After(d Duration) <-chan Time
      >func Sleep(d Duration)
      >func Since(t Time) Duration
      >func Now() Time
      >type Duration int64
      >type Time struct { ... }
      >...many more...
      >```

    - ```shell
      go doc time.Since
      ```

      > 包成员
      >
      > ```shell
      > func Since(t Time) Duration
      > 
      >    Since returns the time elapsed since t.
      >    It is shorthand for time.Now().Sub(t).
      > ```

    - ```shell
      go doc time.Duration.Seconds
      ```

      > 方法
      >
      > ```shell
      > func (d Duration) Seconds() float64
      > 
      >    Seconds returns the duration as a floating-point number of seconds.
      > ```

  - *godoc*

    > 运行：`godoc -http :8000`。
    >
    > 查看：`http://localhost:8000/pkg`

- **内部包**

  > **internal包**：包含**internal**名字的路径段的包，Go语言的构建工具其做了特殊处理。
  >
  > **导入**：只能被和internal目录有**同一个父目录**的包所导入。

  - *eg.* 

    > net/http/internal/chunked内部包只能被net/http/httputil或net/http包导入，但是不能被net/url包导入。不过net/url包却可以导入net/http/httputil包。

- **查询包**

  > go list

  - *测试包是否在工作区并打印它的导入路径*

    ```shell
    go list github.com/go-sql-driver/mysql
    ```

    > ```shell
    > github.com/go-sql-driver/mysql
    > ```

  - *特定子目录下的所有包*

    ```shell
    go list gopl.io/ch3/...
    ```

    > ```shell
    > gopl.io/ch3/basename1
    > gopl.io/ch3/basename2
    > gopl.io/ch3/comma
    > gopl.io/ch3/mandelbrot
    > gopl.io/ch3/netflag
    > gopl.io/ch3/printints
    > gopl.io/ch3/surface
    > ```

  - *和某个主题相关的所有包*

    ```shell
    go list ...xml...
    ```

    > ```shell
    > encoding/xml
    > gopl.io/ch7/xmlselect
    > ```

  - *获取每个包完整的元信息，不仅仅只是导入路径*

    - **JSON格式打印每个包的元信息**

      ```shell
      go list -json hash
      ```

      > ```javascript
      > {
      >     "Dir": "/home/gopher/go/src/hash",
      >     "ImportPath": "hash",
      >     "Name": "hash",
      >     "Doc": "Package hash provides interfaces for hash functions.",
      >     "Target": "/home/gopher/go/pkg/darwin_amd64/hash.a",
      >     "Goroot": true,
      >     "Standard": true,
      >     "Root": "/home/gopher/go",
      >     "GoFiles": [
      >             "hash.go"
      >     ],
      >     "Imports": [
      >         "io"
      >     ],
      >     "Deps": [
      >         "errors",
      >         "io",
      >         "runtime",
      >         "sync",
      >         "sync/atomic",
      >         "unsafe"
      >     ]
      > }
      > ```

    - **用户使用text/template包的模板语言定义输出文本的格式**

      ```shell
      go list -f '{{join .Deps " "}}' strconv
      ```

      > ```shell
      > go list -f "{{join .Deps \" \"}}" strconv
      > ```

# 测试

## go test

> -v：打印每个测试函数的名字和运行时间。
>
> -run：对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会执行。

- **运行**

  > 1. 遍历所有的`*_test.go`文件中符合测试函数命名规则的函数。
  > 2. 生成临时的main包用于调用相应的测试函数。
  > 3. 接着构建并运行、报告测试结果，最后清理测试中生成的临时文件。

- **main包**

  > 在包目录内，所有以`_test.go`为后缀名的源文件在执行go build时不会被构建成包的一部分，它们是go test测试的一部分。

- **默认**：采用当前目录对应的包。

## 测试函数

```go
func TestName(t *testing.T) {}
```

> t参数用于报告测试失败和附加的日志信息。

- **t.Error**：报告失败信息。
- **t.Fatal或t.Fatalf**：停止当前测试函数

### 随机测试

> 对于一个随机的输入，我们如何能知道希望的输出结果呢？
>
> 1. 编写另一个对照函数，使用简单和清晰的算法，虽然效率较低但是行为和要测试的函数是一致的，然后针对相同的随机输入检查两者的输出结果。
> 2. 生成的随机输入的数据遵循特定的模式，这样我们就可以知道期望的输出的模式。

### 测试一个命令

> 虽然是main包，也有对应的main入口函数，但是在测试的时候main包只是TestEcho测试函数导入的一个普通包，里面main函数并没有被导出，而是被忽略的。
>
> 测试代码中并没有调用log.Fatal或os.Exit，因为调用这类函数会导致程序提前退出；调用这些函数的特权应该放在main函数中。

### 白盒测试

- 黑盒测试只需要测试包公开的文档和API行为，内部实现对测试代码是透明的。
- 白盒测试有访问包内部函数和数据结构的权限，因此可以做到一些普通客户端无法实现的测试。

```go
saved := notifyUser
defer func() { notifyUser = saved }()
```

> 修改全局变量后要主意恢复。

### 外部测试包

- **循环依赖**：可以通过外部测试包的方式解决循环依赖的问题，也就是在net/url包所在的目录声明一个独立的url_test测试包。其中包名的`_test`后缀告诉go test工具它应该建立一个额外的包来运行测试。
- 在设计层面，外部测试包是在所有它依赖的包的上层。

### 编写有效的测试

不仅报告调用的具体函数、它的输入和结果的意义；并且打印的真实返回的值和期望返回的值；并且即使断言失败依然会继续尝试运行更多的测试。

### 避免脆弱的测试

- 只检测你真正关心的属性。
- 保持测试代码的简洁和内部结构的稳定。
- 对断言部分要有所选择。
- 不要对字符串进行全字匹配，而是针对那些在项目的发展中是比较稳定不变的子串。
- 很多时候值得花力气来编写一个从复杂输出中提取用于断言的必要信息的函数。

## 测试覆盖率

> 语句的覆盖率是指在测试中至少被运行一次的代码占总代码数的比例。

- 测试覆盖率工具用法

  ```go
  go tool cover
  ```

- 生成覆盖率文件

  ```go
  go test -coverprofile=c.out
  ```

- 显示html报告

  ```go
  go tool cover -html=c.out
  ```

## 基准测试

> 基准测试是测量一个程序在固定工作负载下的性能。

```go
func BenchmarkIsPalindrome(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsPalindrome("A man, a plan, a canal: Panama")
    }
}
```

> 1. 以Benchmark为前缀名，带有一个`*testing.B`类型的参数。
> 2. `*testing.B`参数除了提供和`*testing.T`类似的方法，还有额外一些和性能测量相关的方法。它还提供了一个整数N，用于指定操作执行的循环次数。

- **默认**：默认情况下不运行任何基准测试。我们需要通过`-bench`命令行标志参数手工指定要运行的基准测试函数。

- ```sh
  PASS
  BenchmarkIsPalindrome-8 1000000                1035 ns/op
  ok      gopl.io/ch11/word2      2.179s
  ```

  > - 8：表示运行时对应的GOMAXPROCS的值。
  > - 1.035：每次调用函数花费1.035微秒。是执行1,000,000次的平均时间。
  > - 1000000：执行次数。

- **比较型基准测试**

  ```go
  func benchmark(b *testing.B, size int) { /* ... */ }
  func Benchmark10(b *testing.B)         { benchmark(b, 10) }
  func Benchmark100(b *testing.B)        { benchmark(b, 100) }
  func Benchmark1000(b *testing.B)       { benchmark(b, 1000) }
  ```

  > 通常是单参数的函数，由几个不同数量级的基准测试函数调用。

## 剖析

- **CPU剖析**：标识了最耗CPU时间的函数。

  ```sh
  go test -cpuprofile=cpu.out
  ```

- **堆剖析**：标识了最耗内存的语句。剖析库会记录调用内部内存分配的操作，平均每512KB的内存申请会触发一个剖析数据。

  ```sh
  go test -memprofile=mem.out
  ```

- **阻塞剖析**：记录阻塞goroutine最久的操作，例如系统调用、管道发送和接收，还有获取锁等。

  ```sh
  go test -blockprofile=block.out
  ```

- **分析工具**

  ```go
  go tool pprof
  ```

## 示例函数

> 以Example为函数名开头。示例函数没有函数参数和返回值。

1. **作为文档**：一个包的例子可以更简洁直观的方式来演示函数的用法，比文字描述更直接易懂，特别是作为一个提醒或快速参考时。

2. **示例函数测试**：在`go test`执行测试的时候也会运行示例函数测试。如果示例函数内含有类似上面例子中的`// Output:`格式的注释，那么测试工具会执行这个示例函数，然后检查示例函数的标准输出与注释是否匹配。

   ```go
   func ExampleIsPalindrome() {
       fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
       fmt.Println(IsPalindrome("palindrome"))
       // Output:
       // true
       // false
   }
   ```

3. **提供演练场**。

   > <img src="./imgs/演练场.png" alt="img" style="zoom:75%;" />

# 反射

> 需要反射的原因：没有办法来检查未知类型的表示方式。

## reflect.Type

> reflect.TypeOf 返回的是一个动态类型的接口值，它总是返回具体的类型。
>
> fmt.Printf 提供了一个缩写 %T 参数，内部使用 reflect.TypeOf 来输出。

```go
t := reflect.TypeOf(3)  // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int"
```

```go
var w io.Writer = os.Stdout
fmt.Println(reflect.TypeOf(w)) // "*os.File"
```

- **reflect.Type.Field**

  > 返回一个reflect.StructField，里面含有每个成员的名字、类型和可选的成员标签等信息。其中成员标签信息对应reflect.StructTag类型的字符串。

- **reflect.StructTag.Get**

  > 用于解析和根据特定key提取的子串。

- **reflect.Type.Method(i)**

  > 返回一个reflect.Method的实例，对应一个用于描述一个方法的名称和类型的结构体。

## reflect.Value

> reflect.ValueOf 返回的结果也是具体的类型。
>
> 除非 Value 持有的是字符串，否则 String 方法只返回其类型。
>
> **使用 fmt 包的 %v 标志参数会对 reflect.Values 特殊处理**。

- **reflect.Value.Interface**

  > reflect.ValueOf 的逆操作，返回一个 interface{} 类型。

- **reflect.Value.Type**

  > 返回具体类型所对应的 reflect.Type。

- **reflect.Value.Kind**

  > 类型是有限的。Kind 只关心底层表示，可以解放 *switch*。
  >
  > - Bool、String 和 所有数字类型的基础类型。
  > - Array 和 Struct 对应的聚合类型。
  > - Chan、Func、Ptr、Slice 和 Map 对应的引用类型。
  > - interface 类型。
  > - 还有表示空值的 Invalid 类型。（空的 reflect.Value 的 kind 即为 Invalid。）

- **reflect.Value.CanAddr**

  > 判断其是否可以被取地址。每当我们通过指针间接地获取的reflect.Value都是可取地址的。

- **reflect.Value.UnsafeAddr**

  > 获取不安全的变量地址，类型是uintptr。

- **reflect.Value.Addr**

  > 获取指向结构体字段的指针。

- **reflect.Value.CanSet**

  > 用于检查对应的reflect.Value是否是可取地址并可被修改的。

- **reflect.Value.Set**：更新对应的值

  > 要确保改类型的变量可以接受对应的值，否则panic。
  >
  > 对一个不可取地址的reflect.Value调用Set方法也会导致panic异常。

- **reflect.Value.Method(i)**

  > 返回一个reflect.Value以表示对应的值。

- **reflect.Value.Call**

  > 调用一个Func类型的Value。

### 更新值

1. ```go
   x := 2
   d := reflect.ValueOf(&x).Elem()   // d refers to the variable x
   px := d.Addr().Interface().(*int) // px := &x
   *px = 3                           // x = 3
   ```

   > 1. 调用Addr()方法，它返回一个Value，里面保存了指向变量的指针。
   > 2. 调用Interface()方法，也就是返回一个interface{}，里面包含指向变量的指针。
   > 3. 使用类型的断言机制将得到的interface{}类型的接口强制转为普通的类型指针。

2. ```go
   d.Set(reflect.ValueOf(4))
   d.SetInt(3)
   ```

   > SetInt、SetUint、SetString和SetFloat。
   >
   > 以SetInt为例，只要变量是某种类型的有符号整数就可以工作，即使是一些命名的类型、甚至只要底层数据类型是有符号整数就可以，而且如果对于变量类型值太大的话会被自动截断。
   >
   > 对于一个引用interface{}类型的reflect.Value调用SetInt会导致panic异常，即使那个interface{}变量对于整数类型也不行。

3. 反射可以**读取**结构体的**未导出**成员，不能**修改未导出**成员。

4. 对不可取地址的reflect.Value调用Set方法也会导致**panic异常**。

   ```go
   x := 2
   b := reflect.ValueOf(x)
   b.Set(reflect.ValueOf(3)) // panic: Set using unaddressable value
   ```

   > **重点**：：：：：需要ValueOf(&x).Elem()。

### 类型处理

- **Slice和数组**

  - **Len**：返回slice或数组值中的元素个数。

  - **Index(i)**：获得索引i对应的元素，返回的也是一个reflect.Value；索引i超出范围的话将导致panic异常。

- **结构体**

  - **NumField**：返回结构体中成员的数量。

  - **Field(i)**：以reflect.Value类型返回第i个成员的值。

- **Maps**

  - **MapKeys**：返回一个reflect.Value类型的slice，每一个元素对应map的一个key。

  - **MapIndex(key)**：返回map中key对应的value。

- **指针**

  - **Elem**：返回指针指向的变量，依然是reflect.Value类型。nil也安全。

  - **IsNil**：测试空指针。

- **接口**

  - **Elem()**：来获取接口对应的动态值，并且打印对应的类型和值。

  - **IsNil**：测试接口是否是nil

## reflect.DeepEqual

> 深度判断相等

## 警告

- 基于反射的代码是比较脆弱的；在编译时检查不出错误。降低了程序的安全性，还影响了自动化重构和分析工具的准确性，因为它们无法识别运行时才能确认的类型信息。
- 反射的操作不能做静态类型检查，而且大量反射的代码通常难以理解。
- 基于反射的代码通常比正常的代码运行速度慢一到两个数量级。对于性能关键路径的函数，最好避免使用反射。

# 底层编程

## unsafe

- **Sizeof**：返回操作数在内存中的字节大小，参数可以是任意类型的表达式，但是它并不会对表达式进行求值。

  | 类型                            | 大小                              |
  | ------------------------------- | --------------------------------- |
  | `bool`                          | 1个字节                           |
  | `intN, uintN, floatN, complexN` | N/8个字节（例如float64是8个字节） |
  | `int, uint, uintptr`            | 1个机器字                         |
  | `*T`                            | 1个机器字                         |
  | `string`                        | 2个机器字（data、len）            |
  | `[]T`                           | 3个机器字（data、len、cap）       |
  | `map`                           | 1个机器字                         |
  | `func`                          | 1个机器字                         |
  | `chan`                          | 1个机器字                         |
  | `interface`                     | 2个机器字（type、value）          |

  > 机器字：操作系统位数，对齐位。
  >
  > 有效的包装可以使数据结构更加紧凑，内存使用率和性能都可能会受益。

- **Alignof**

  > 返回对应参数的类型需要对齐的倍数。返回一个常量表达式，对应一个常量。
  >
  > 通常情况下布尔和数字类型需要对齐到它们本身的大小（最多8个字节）；
  >
  > 其它的类型对齐到机器字大小。

- **Offsetof**

  > 参数必须是一个字段 `x.f`，然后返回 `f` 字段相对于 `x` 起始地址的偏移量，包括可能的空洞。

- **Pointer**

  > 可以包含任意类型变量的地址。

  - **解地址**：不可以直接通过`*p`来获取unsafe.Pointer指针指向的真实变量的值，因为我们并不知道变量的具体类型。

  - **比较**：可以比较，并且支持和nil常量比较判断是否为空指针。
  - **深度比较**：通过比较地址值，避免递归。
  - **uintptr**：
    - unsafe.Pointer指针也可以被转化为uintptr类型，然后保存到指针型数值变量中。
    - 但是将uintptr转为unsafe.Pointer指针可能会破坏类型系统，因为并不是所有的数字都是有效的内存地址。
    - 将unsafe.Pointer指针转为原生数字，然后再转回为unsafe.Pointer类型指针的操作也是不安全的。

  ```go
  pT := uintptr(unsafe.Pointer(new(T))) // 提示: 错误!
  ```

  > 并没有指针引用`new`新创建的变量，因此该语句执行完成之后，垃圾收集器有权马上回收其内存空间，所以返回的pT将是无效的地址。

  ```go
  // NOTE: subtly incorrect!
  tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
  pb := (*int16)(unsafe.Pointer(tmp))
  *pb = 42
  ```

  > 产生错误的原因很微妙。有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为**移动GC**【虽然目前的Go语言实现还没有使用移动GC】。当一个变量被移动，所有的保存该变量旧地址的指针必须同时被更新为变量移动后的新地址。
