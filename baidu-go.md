# 基础

## 通用规范

1. **命名规范**

   - **包名**
     - *advice*：出_test包外，包名不能包含下划线、应该全部大写或全部小写。


   - **变量名**
     - *advice*：类型为time.Duration 时，变量名不应该以时间单位作为后缀。

     - *warning*：名称中包含专用词【IP】时，应全部大写或小写。

     - *warning*：recevier 的名字不能是下划线。

     - *warning*：recevier name 应该反映出功能、角色，不应该使用 this 或 self。

     - *warning*：recevier name 命名应该保持一致。

     - **error**：error类型的变量名应该以 err 或者 Err 开头。


   - **结构体命名**
     - *warning*：结构体中不能包含仅大小写不同的名称定义。


   - **函数名**
     - *advice*：不能出现含仅大小写不同的名称定义。

2. **注释规范**

   - **包注释**
     - *advice*：package声明都应该有对应的包说明性注释，一个包内有一个文件有即可。

     - *advice*：应该以 Package开头，后面跟注释。

     - *warning*：// 后不要加空格或tab。

     - **error**：包注释和包声明之间不能有空行。


   - **类型注释**
     - *advice*：可导出类型的注释应该以 type name 开头。

     - *warning*：可导出类型应该包含注释。


   - **变量注释**
     - *warning*：可导出变量应该包含注释。

     - *warning*：可导出变量应该单独定义。


   - **结构体注释**
     - *warning*：可导出结构体的注释应该以变量名开头。


   - **函数注释**
     - *advice*：可导出函数应该包含注释。
     - *advice*：可导出函数的注释应该以函数名开头。

3. **可维护性规范**

   - **代码复杂度**
     - *advice*：函数**圈复杂度**不应该大于15。

     - *warning*：代码锁紧深度 <= 5。


   - **参数与返回值**
     - *warning*：返回值小于等于3个，大于三个时必须通过 struct 传递。

     - *warning*：参数不建议超过3个，大于三个时建议通过 struct 传递。


   - **程序规模**
     - *advice*：每行代码不超过100字符。
     - *advice*：每行注释不超过100字符。
     - *advice*：函数不超过100行。
     - *advice*：文件不超过2000行。

4. **包引用规范**

   - **error**：空引用只能存在于`main`或`test`包中，除非有对应的comments说明理由。该规则为了明确所导入包的用途。

5. **声明规范**

   - **变量声明**
     - *advice*：变量在定义时如果被赋零值，应省略赋值操作。

     - *advice*：定义时指定了某类型初值，定义时无需再显示标注变量类型。

     - *advice*：error字符串首字母【除专有名词】不应该大写，或者以标点符号或换行符结束。

     - *warning*：声明slice时，建议使用var方式声明，不建议使用大括号的方式。


   - **结构体声明**
     - *warning*：结构体的标签应该包含键和值，且之间要用英文冒号分隔。
     - *warning*：标签最外层是`，里面是"。
     - *warning*：标签多个键值对是，要用空格分隔。

6. **数据流规范**

   - *advice*：不应该使用基础类型为key调用context.WithValue。

   - *advice*：context.Context应该是函数的第一个参数。

   - *advice*：不要在go代码块中处理在迭代过程中被赋予值的迭代变量。不要在可能被延时处理的代码块中直接使用迭代变量。

   - *advice*：对 sync/atomic 使用不能破坏原子性的使用方式。

   - *warning*：不要修改形参的值。

   - *warning*：可导出的函数或方法应该返回可导出的类型。

   - *warning*：不能为值传递 `receiver` 的 `field` 赋值。

7. **控制流规范**

   - *advice*：if 块如果是 ruturn 结束的话，不该有else。

   - *advice*：在range循环中，第二个值不该用_，可直接忽略。

   - *warning*：函数和方法中不应该有无效代码。

   - **error**：函数中返回的 error 应该是最后一个返回值。

8. **其他规范**

   - *advice*：自增自减用`++`、`--`。不要用 `+=1`。

   - *warning*：error应该用 fmt.Errorf()，不要组合使用。

   - **error**：当 expr 是 boor类型时，不应该使用 ==。


## 细节规范

1. **语言规范**

   - **true/false求值**
     - **必须遵守**：禁止使用==或!=与true/false比较，应该使用expr或!expr。
     - **必须遵守**：判断某个整数表达式expr是否为零时，禁止使用!expr，应该使用expr == 0。

   - **receiver**
     - *receiver type*
       - **必须遵守**：receiver是map、函数或者chan类型，类型不可以是指针
       - **必须遵守**：receiver是slice，并且方法不会进行reslice或者重新分配slice，类型不可以是指针。
       - **必须遵守**：receiver是struct，且包含sync.Mutex类型字段，则必须使用指针避免拷贝。
       - *建议遵守*：receiver是比较大的struct/array，建议使用指针，这样会更有效率。
       - *建议遵守*：receiver是struct、array或slice，其中指针元素所指的内容可能在方法内被修改，建议使用指针类型。
       - *建议遵守*：receiver是比较小的struct/array，建议使用value类型。
     - *receiver 命名*
       - **必须遵守**：尽量简短并有意义。
       - *建议遵守*：禁止使用`this`、`self`等面向对象语言中特定的叫法。
       - *建议遵守*：receiver的命名要保持一致性。

   - **申明空Slices**
     - *建议遵守*：申明slice时，建议使用var方式申明，不建议使用大括号的方式。【var方式申明在slice不被append的情况下避免了内存分配】。

   - **Error Handler**
     - **必须遵守**：对于返回值中的error，一定要进行判断和处理，不可以使用 ”_“ 变量忽略error。

   - **`{` 的使用**
     - **必须遵守**：struct、函数、条件判断中的 `{` ，不可以作为独立的一行。

   - **`embedding`【嵌入】的使用**
     - *建议遵守*：embedding只用于「is a」的语义下，而不用于「has a」的语义下。
     - *建议遵守*：一个定义内，多于一个的embedding尽量少用。

2. **风格规范**

   - **Go文件Layout**
     - **必须遵守**：对于的布局的各个部分，采用单个空行分割。
       - 多个类型定义采用单个空行分割。
       - 多个函数采用单个空行分割。
     - *建议遵守*：布局方式
       - General Documentation： 对整个模块和功能的完整描述注释，写在文件头部。
       - package：当前package定义。
       - imports：包含的头文件。
       - Constants：常量。
       - Typedefs： 类型定义。
       - Globals：全局变量定义。
       - functions：函数实现。
     - *建议遵守*：函数内不同的业务逻辑处理建议采用单个空行分割。
     - *建议遵守*：常量或者变量如果较多，建议按照业务进行分组，组间用单个空行分割。

   - **General Documentation Layout**
     - *建议遵守*：每个文件开头部分包括文件版权说明（Copyright）
     - *建议遵守*：每个文件开头部分包括文件标题。【Title包含文件的名称和文件的简单说明，Title应该在一行内完成】。
     - *建议遵守*：每个文件开头部分包括修改记录（Modification History）。
       - 当书写新的函数模块时，只需要使用形如"Add func1()"这样的说明。
       - 如果后面有对函数中的算法进行了修改，需要指出修改的具体位置和修改方法。
       - 具体格式为：<修改时间>, <修改人>, <修改动作 >。
     - *建议遵守*：每个文件开头部分包括文件描述（Description）。详细描述文件的功能和作用。

   - **import规范**
     - **必须遵守**：按照如下顺序进行头文件import，并且每个import部分内的package需按照字母升序排列。
       - 系统package。
       - 第三方的package。
       - 程序自己的package。
     - **必须遵守**：每部分import间用单个空行进行分隔。


   - **Go函数Layout**
     - *函数注释*
       - *建议遵守*：函数的注释，建议包括以下内容：
         - Description：对函数的完整描述，主要包括函数功能和使用方法。
         - Params：对参数的说明。
         - Returns：对返回值的说明。
     - *函数参数和返回值*
       - **必须遵守**：函数返回值小于等于3个，大于3个时必须通过struct进行包装。
       - *建议遵守*：函数参数不建议超过3个，大于3个时建议通过struct进行包装。
       - *建议遵守*：对于「逻辑判断型」的函数，返回值的意义代表「真」或「假」，返回值类型定义为`bool`。
       - *建议遵守*：对于「操作型」的函数，返回值的意义代表「成功」或「失败」，返回值类型定义为`error`。
       - *建议遵守*：对于「获取数据型」的函数，返回值的意义代表「有数据」或「无数据/获取数据失败」，返回值类型定义为`(data, error)`。


   - **程序规模**
     - **必须遵守**：每行代码不超过100个字符。
     - **必须遵守**：每行注释不超过100个字符。
     - *建议遵守*：函数不超过100行。
     - *建议遵守*：文件不超过2000行。


   - **命名规范**
     - *文件名*
       - **必须遵守**：文件名都使用小写字母，如果需要，可以使用下划线分割。
       - **必须遵守**：文件名的后缀使用小写字母。
     - *函数名/变量名*
       - **必须遵守**：采用驼峰方式命名，禁止使用下划线命名。首字母是否大写，根据是否需要外部访问来决定。
     - *常量*
       - *建议遵守*：不要在程序中直接写数字，特殊字符串，全部用常量替代 。
     - *缩写词*
       - **必须遵守**：保持命名的一致性。
         - 同一变量字母大小写的一致性。
         - 不同变量间的一致性。

   - **缩进**
     - **必须遵守**：使用tab进行缩进。
     - **必须遵守**：跨行的缩进使用gofmt的缩进方式。
     - **必须遵守**：设置tabstop=4。

   - **空格**
     - **必须遵守**：圆括号、方括号、花括号内侧都不加空格。
     - **必须遵守**：逗号、冒号（slice中冒号除外）前不加空格，后边加一个空格。
     - **必须遵守**：所有二元运算符前后各加一个空格（**作为函数参数时除外**）`func(1+2)`。

   - **括号**
     - *建议遵守*：除非用于明确算术表达式优先级，否则尽量避免冗余的括号。

   - **注释**
     - *建议遵守*：单行注释，采取`//`或者`/*...*/`的注释方式。
     - *建议遵守*：多行注释，采取每行开头`//`或者用`/* ... */`包括起来的注释（`/*`和`*/`作为独立的行）。
     - *建议遵守*：紧跟在代码之后的注释，使用`//`。
     - 大多数情况下，使用”//"更方便。

3. **编程实践**

   - **error string**
     - *建议遵守*：error string尽量使用小写字母，并且结尾不带标点符号。【因为可能error string会用于其它上下文中】。

   - **Don't panic**
     - *建议遵守*：除非出现不可恢复的程序错误，不要使用panic，用多返回值和error。

   - **关于lock的保护**
     - *建议遵守*：如果临界区内的逻辑较复杂、无法完全避免panic的发生，则要求适用defer来调用Unlock，即使在临界区过程中发生了panic，也会在函数退出时调用Unlock释放锁。
     - *建议遵守*：上述操作如果造成临界区扩大后，需要建立单独的一个函数访问临界区。

   - **日志的处理**
     - *建议遵守*：使用公司golang-lib中的log库。

   - **unsafe package**
     - *建议遵守*：除非特殊原因，不建议使用unsafe package。


# 单元测试

> 用来对模块、函数和类进行正确性校验。
>
> 开发人员在提测前完成。

- **定义**

  - 和轻量级测试框架 testing一起使用。

  - 以Test为函数名前缀，以*testing.T为单一参数的函数。

  - 以Benchmark为前缀，以*testing.B为单一参数的函数。

  - 测试函数名 TestXxx。

  - 测试源文件名称： xxx_test.go，在执行go build时不会被构建成包的一部分。


- **参数**

  ```she
  go test [-c] [-i] [build flags] [packages] [flags for test binary]
  ```

  > - -c：编译成可执行的二进制文件，但不运行测试。
  > - -i：安装测试包依赖的package，但不运行测试。
  > - -v：输出全部单元测试用例，不加只输出失败的。
  > - -run regexp：指定运行正则匹配的测试函数。
  > - -bench regexp：指定运行正则匹配的性能测试函数。
  > - -test.bench="regexp"：性能测试。
  > - -benchtime t：性能测试运行的时间，默认1s。
  > - -cover：单测覆盖率统计。
  > - -timeout t：如果测试用例运行时间超过 t，则抛出 panic。
  > - -short：支持跳过测试标识。
  > - -coverprofile  生成html文件。

  ```shell
  go test add_test.go add.go
  ```

  > 指定单元测试文件，需要加上对应的源代码

- **报告**

  - t.Error、t.Errorf、t.FailNow、t.Fatal、t.FatalIf 表示测试不通过。

    > t.Error、t.Errorf：报告出错继续。
    >
    > t.Fatal、t.FatalIf：报告出错终止。
    >
    > t.FatalIf：失败终止。

  - t.Log 用来记录测试信息。


## 优雅的单元测试

- **test main**

  ```go
  func TestMain(m *testing.M) {
    // begin
    m.Run()
    // end
  }
  ```

  >begin部分在所有单测里最先执行，end部分最后执行。
  >
  >测试之前初始化操作「打开连接」，测试之后清理工作「关闭连接」。

- **Parallel**

  ```go
  func TestParallel(t *testing.T) {
    names := []string{"aa", "bb", "cc"}
    for _, name := range names {
      tName := name
      t.Run(tName, func(t *testing.T) {
      	t.Parallel()
      	fmt.Println(tName) // 返回names的乱序
    	})
    }
  }
  ```

  > 使用parallel并行运行，提高运行速度。

- **GoConvey**

  > 可以管理和运行测试用例，同时提供了丰富的断言函数，并支持很多 Web 界面特性。

  ```go
  func TestStringSliceEqual(t *testing.T) {
    Convey("a,b相等", t, func() {
      a := []string{"hello", "goconvey"}
      b := []string{"hello", "goconvey"}
      So(StringSliceEqual(a, b), ShouldBeTrue)
    })
    Convey("a,b都为空", t, func() {
      So(StringSliceEqual(nil, nil), ShouldBeTrue)
    })
    Convey("a,b不相等", t, func() {
      a := []string(nil)
      b := []string{}
      So(StringSliceEqual(a, b), ShouldBeFalse)
    })
  }
  ```

  > 执行./goconvey 命令，在浏览器可看到ui界面。
  >
  > ```shell
  > go get github.com/smartystreets/goconvey && go install
  > ```

- **gotests**

  ```go
  func TestStringSliceEqual(t *testing.T) {
    type args struct {
      a []string
      b []string
    }
    tests := []struct {
      name string
      args args
      want bool
    }{
      // TODO: Add test cases.
    }
    for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
        if got := StringSliceEqual(tt.args.a, tt.args.b); got != tt.want {
          t.Errorf("StringSliceEqual() = %v, want %v", got, tt.want)
        }
      })
    }
  }
  ```

  > 规划一个数组、循环遍历每一个 case，通过gotests可以自动生成单测。
  >
  > ```shell
  > gotests –all  [filename]
  > gotests –w –only  [funcname] [filename]
  > gotests –w –all origin.go, origin_test.go
  > ```
  >
  > >ide集成了gotests工具
  > >
  > >- goland快捷键 command+shirt+T
  > >- vscode安装Go插件后，右键Go: Generate Unit Test For Function 即可生成单测代码。

- **GoMock**

  >生成接口测试代码。

  ```shell
  go get github.com/golang/mock/gomock && go install
  go github.com/golang/mock/mockgen && go install 
  ```

  > 基础库依赖

  ```shell
  mockgen -destination=mocks/mock_doer.go -package=mocks gomock Doer
  ```

  > 用mockgen为要模拟的接口生成模拟。
  >
  > - destination=mocks/mock_doer.go ： 将生成的模拟接口放入指定文件中
  > - package=mocks ： 将生成的模拟接口放入包mocks中
  > - gomock 接口定义在gomock目录下
  > - Doer 接口名

  ```go
  type Doer interface {
    DoSomething(int, string) error
  }
  
  // 一个包含Doer接口的结构体
  type User struct {
    d Doer
  }
  
  func (u *User) Use(i int, s string) error {
    return u.d.DoSomething(i, s)
  }
  ```

  > Demo

## 打桩

- **GoStub**

  > 轻量级的测试框架，接口友好，可以对全局变量、函数或过程打桩。

  - **变量打桩**

    ```go
    stubs := Stub(&commNum, 100)
    ```

  - **函数打桩**

    ```go
    var (
      commFunc = func(cmd string) string {
        return cmd
      }
    )
    
    stubs := Stub(&commFunc, func(cmd string) string {
      return "bbb"
    })
    ```

  - **第三方函数打桩**

    ```go
    var Marshal = json.Marshal
    stubs := StubFunc(&Marshal, []byte(`{"name":"aaa"}`), nil)
    ```

    > 需要自定义方法函数变量。

- **Monkey**

  > 一个补丁框架，在运行时通过汇编语句重写可执行文件，将待打桩函数或方法的实现跳转到桩实现。

  - **对象方法打桩**

    ```go
    var client = &Client{}
    monkey.PatchInstanceMethod(reflect.TypeOf(client), "Test", func(c *Client, param string) bool {
      return false
    })
    ```

  - **方法打桩**

    ```go
    monkey.Patch(fmt.Println, func(a ...interface{}) (n int, err error) {
    	s := make([]interface{}, len(a))
    	for i, v := range a {
    		s[i] = strings.Replace(fmt.Sprint(v), "apple", "banana", -1)
    	}
    	return fmt.Fprintln(os.Stdout, s...)
    })
    ```

  - **第三方函数打桩**

    ```go
    var guard *monkey.PatchGuard
    guard = monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
      guard.Unpatch()
      refer guard.Restore()
      if !strings.HasPrefix(url, "https://") {
        return nil, fmt.Errorf("only https requests allowed")
      }
      return c.Get(url)
    })
    ```

    >Monkey对内联函数不生效，执行时可以通过命令行参数 `-gcflags=-l` 禁止inline。
    >
    >Monkey不是线程安全的，不要将Monkey用于并发的测试中。

# GoDoc

> 在声明之前写一个常规注释，中间没有空行。

- **Package文档**

  - 对整个包，一个目录的文档说明，一般整体描述一个包的功能。

  - 只需要在一个文件上写包注释。若同一个包下，多个文件都写了注释，将按照文件名顺序合并注释。

  - 一般情况下写一个doc.go，不写代码，只写注释。


- **变量、常量、函数文档**

  >以变量名为开头，后面写注释。

- **函数、变量过期标记**

  - `Deprecate`用于标记已过期，不推荐使用。会给调用发提示。

  - 调用方会显示删除线。


- **标记BUG**

  > `BUG(who):`标签，可以用来记录一些bug标签和对bug的说明。

- **插入代码块**

  - 使用一个空行注释，可以换行。

  - 注释前加一个tab，将后面的内容变成代码块。


- **添加example**

  > go test examples
  >
  > `test`文件中的`ExampleXxx`函数。

## 查看go文档

1. 安装`godoc`工具：`go get -u golang.org/x/tools/cmd/godoc`。

2. 启动`http`服务：`godoc -http=:6060 `。
3. 浏览器访问模块位置即可查看模块的api手册。`http://127.0.0.1/pkg/pkg下的模块路径`。

## 检查文档风格

- 安装`golint`工具：`go get -u golang.org/x/lint/golint`。
- 检查目录下所有代码：`golint -min_confidence=0.3 ./..`。