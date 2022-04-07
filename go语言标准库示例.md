# 输入输出

## io — 基本的 IO 接口

> io 包为 I/O 原语提供了基本的接口。它主要包装了这些原语的已有实现。

### 接口

#### Reader 和 Writer 接口

- **Reader 接口**

  ```go
  type Reader interface {
    Read(p []byte) (n int, err error)
  }
  ```

  > Read 将 len(p) 个字节读取到 p 中。它返回读取的字节数 n（0 <= n <= len(p)） 以及任何遇到的错误。即使 Read 返回的 n < len(p)，它也会在调用过程中占用 len(p) 个字节作为暂存空间。若可读取的数据不到 len(p) 个字节，Read 会返回可用数据，而不是等待更多数据。
  >
  > 当 Read 在成功读取 n > 0 个字节后遇到一个错误或 EOF (end-of-file)，它会返回读取的字节数。它可能会同时在本次的调用中返回一个non-nil错误,或在下一次的调用中返回这个错误（且 n 为 0）。 一般情况下, Reader会返回一个非0字节数n, 若 n = len(p) 个字节从输入源的结尾处由 Read 返回，Read可能返回 err == EOF 或者 err == nil。
  >
  > **并且之后的 Read() 都应该返回 (n:0, err:EOF)**。
  >
  > 调用者在考虑错误之前应当首先处理返回的数据。这样做可以正确地处理在读取一些字节后产生的 I/O 错误，同时允许EOF的出现。

  - *示例*

    ```go
    func ReadFrom(reader io.Reader, num int) ([]byte, error) {
        p := make([]byte, num)
        n, err := reader.Read(p)
        if n > 0 {
            return p[:n], nil
        }
        return p, err
    }
    // 从标准输入读取
    data, err = ReadFrom(os.Stdin, 11)
    // 从普通文件读取，其中 file 是 os.File 的实例
    data, err = ReadFrom(file, 9)
    // 从字符串读取
    data, err = ReadFrom(strings.NewReader("from string"), 12)
    ```

  - *io.EOF*

    ```go
    var EOF = errors.New("EOF")
    ```

    > 根据 reader 接口的说明，在 n > 0 且数据被读完了的情况下，返回的 error 有可能是 EOF 也有可能是 nil。

- **Writer 接口**

  ```go
  type Writer interface {
    Write(p []byte) (n int, err error)
  }
  ```

  > Write 将 len(p) 个字节从 p 中写入到基本数据流中。它返回从 p 中被写入的字节数 n（0 <= n <= len(p)）以及任何遇到的引起写入提前停止的错误。若 Write 返回的 n < len(p)，它就必须返回一个 非nil 的错误。

-  **实现场景**

  - *常见类型*：`os.File`、`strings.Reader`、`bufio.Reader/Writer`、`bytes.Buffer`、`bytes.Reader`。

  - *所有*

    > - os.File 同时实现了 io.Reader 和 io.Writer
    > - strings.Reader 实现了 io.Reader
    > - bufio.Reader/Writer 分别实现了 io.Reader 和 io.Writer
    > - bytes.Buffer 同时实现了 io.Reader 和 io.Writer
    > - bytes.Reader 实现了 io.Reader
    > - compress/gzip.Reader/Writer 分别实现了 io.Reader 和 io.Writer
    > - crypto/cipher.StreamReader/StreamWriter 分别实现了 io.Reader 和 io.Writer
    > - crypto/tls.Conn 同时实现了 io.Reader 和 io.Writer
    > - encoding/csv.Reader/Writer 分别实现了 io.Reader 和 io.Writer
    > - mime/multipart.Part 实现了 io.Reader
    > - net/conn 分别实现了 io.Reader 和 io.Writer(Conn接口定义了Read/Write)

  - *os.File*

    ```go
    var (
      Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
      Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
      Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
    )
    ```

#### ReaderAt 和 WriterAt 接口

- **ReaderAt 接口**

  ```go
  type ReaderAt interface {
    ReadAt(p []byte, off int64) (n int, err error)
  }
  ```

  > ReadAt 从基本输入源的偏移量 off 处开始，将 len(p) 个字节读取到 p 中。它返回读取的字节数 n（0 <= n <= len(p)）以及任何遇到的错误。
  >
  > 当 ReadAt 返回的 n < len(p) 时，它就会返回一个 非nil 的错误来解释 为什么没有返回更多的字节。在这一点上，ReadAt 比 Read 更严格。
  >
  > 即使 ReadAt 返回的 n < len(p)，它也会在调用过程中使用 p 的全部作为暂存空间。若可读取的数据不到 len(p) 字节，ReadAt 就会**阻塞**,直到所有数据都可用或一个错误发生。 在这一点上 ReadAt 不同于 Read。
  >
  > 若 n = len(p) 个字节从输入源的结尾处由 ReadAt 返回，Read可能返回 err == EOF 或者 err == nil
  >
  > 若 ReadAt 携带一个偏移量从输入源读取，ReadAt 应当**既不影响偏移量也不被它所影响**。
  >
  > 可对相同的输入源并行执行 ReadAt 调用。

- **WriterAt 接口**

  ```go
  type WriterAt interface {
    WriteAt(p []byte, off int64) (n int, err error)
  }
  ```

  > WriteAt 从 p 中将 len(p) 个字节写入到偏移量 off 处的基本数据流中。它返回从 p 中被写入的字节数 n（0 <= n <= len(p)）以及任何遇到的引起写入提前停止的错误。若 WriteAt 返回的 n < len(p)，它就必须返回一个 非nil 的错误。
  >
  > 若 WriteAt 携带一个偏移量写入到目标中，WriteAt 应当**既不影响偏移量也不被它所影响**。
  >
  > 若被写区域没有重叠，**可**对相同的目标**并行**执行 WriteAt 调用。

#### ReaderFrom 和 WriterTo 接口

- **ReaderFrom接口**

  ```go
  type ReaderFrom interface {
    ReadFrom(r Reader) (n int64, err error)
  }
  ```

  > ReadFrom 从 r 中读取数据，直到 EOF 或发生错误。其返回值 n 为读取的字节数。除 io.EOF 之外，在读取过程中遇到的任何错误也将被返回。
  >
  > 如果 ReaderFrom 可用，**Copy 函数就会使用它**。

  - **如果不通过 ReadFrom 接口读取数据**，使用 io.Reader 接口，**有两种思路**：
    1. 先获取文件的大小（File 的 Stat 方法），之后定义一个该大小的 []byte，通过 Read 一次性读取
    2. 定义一个小的 []byte，不断的调用 Read 方法直到遇到 EOF，将所有读取到的 []byte 连接到一起

- **WriterTo接口**

  ```go
  type WriterTo interface {
    WriteTo(w Writer) (n int64, err error)
  }
  ```

  > WriteTo 将数据写入 w 中，直到没有数据可写或发生错误。其返回值 n 为写入的字节数。 在写入过程中遇到的任何错误也将被返回。
  >
  > 如果 WriterTo 可用，**Copy 函数就会使用它**。

#### Seeker 接口

```go
type Seeker interface {
  Seek(offset int64, whence int) (ret int64, err error)
}
```

> Seek 设置下一次 Read 或 Write 的偏移量为 offset，它的解释取决于 whence： 0 表示相对于文件的起始处，1 表示相对于当前的偏移，而 2 表示相对于其结尾处。 Seek 返回新的偏移量和一个错误，如果有的话。

- **示例**：获取倒数第二个字符。

  ```go
  reader := strings.NewReader("Go语言中文网")
  reader.Seek(-6, io.SeekEnd)
  r, _, _ := reader.ReadRune()
  ```

- **whence**

  ```go
  SeekStart   = 0 // seek relative to the origin of the file
  SeekCurrent = 1 // seek relative to the current offset
  SeekEnd     = 2 // seek relative to the end
  ```

#### Closer接口

```go
type Closer interface {
  Close() error
}
```

> 文件 (os.File)、归档（压缩包）、数据库连接、Socket 等需要手动关闭的资源都实现了 Closer 接口。
>
> 在close之前校验错误。

#### ByteReader 和 ByteWriter 接口

- **ByteReader 接口**

  ```go
  type ByteReader interface {
    ReadByte() (c byte, err error)
  }
  ```

  > 读一个字节

- **ByteWriter 接口**

  ```go
  type ByteWriter interface {
    WriteByte(c byte) error
  }
  ```

  > 写一个字节

- **实现场景**

  > - bufio.Reader/Writer 分别实现了io.ByteReader 和 io.ByteWriter
  > - bytes.Buffer 同时实现了 io.ByteReader 和 io.ByteWriter
  > - bytes.Reader 实现了 io.ByteReader
  > - strings.Reader 实现了 io.ByteReader

#### ByteScanner、RuneReader 和 RuneScanner 接口

- **ByteScanner 接口**

  ```go
  type ByteScanner interface {
    ByteReader
    UnreadByte() error
  }
  ```

  > UnreadByte 是重置上一次的 ReadByte。
  >
  > UnreadByte 调用**之前必须调用**了 ReadByte，且**不能连续调用** UnreadByte。

- **RuneScanner 接口和 ByteScanner 类似**

#### ReadCloser、ReadSeeker、ReadWriteCloser、ReadWriteSeeker、ReadWriter、WriteCloser 和 WriteSeeker 接口

> 这些接口是上面介绍的接口的两个或三个组合而成的新接口

- **ReadWriter 接口**

  ```go
  type ReadWriter interface {
    Reader
    Writer
  }
  ```

### 类型

#### SectionReader 类型

```go
type SectionReader struct {
  r     ReaderAt    // 该类型最终的 Read/ReadAt 最终都是通过 r 的 ReadAt 实现
  base  int64        // NewSectionReader 会将 base 设置为 off
  off   int64        // 从 r 中的 off 偏移处开始读取数据
  limit int64        // limit - off = SectionReader 流的长度
}

func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader
```

> NewSectionReader 返回一个 SectionReader，它从 r 中的偏移量 off 处读取 n 个字节后以 EOF 停止。

#### LimitedReader 类型

```go
type LimitedReader struct {
  R Reader // underlying reader，最终的读取操作通过 R.Read 完成
  N int64  // max bytes remaining
}

func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }
```

> 从 R 读取但将返回的数据量限制为 N 字节。每调用一次 Read 都将更新 N 来反应新的剩余数量。

#### PipeReader 和 PipeWriter 类型

- **PipeReader 类型**

  ```go
  type PipeReader struct {
    p *pipe
  }
  ```

  > 是管道的读取端。实现了 io.Reader 和 io.Closer 接口

  - **PipeReader.Read**

    > 从管道中读取数据。该方法会堵塞，直到管道写入端开始写入数据或写入端被关闭。如果写入端关闭时带有 error（即调用 CloseWithError 关闭），该Read返回的 err 就是写入端传递的error；否则 err 为 EOF。

- **PipeWriter 类型**

  ```go
  type PipeWriter struct {
    p *pipe
  }
  ```

  > 是管道的写入端。它实现了 io.Writer 和 io.Closer 接口。

  - **PipeWriter.Write**

    > 写数据到管道中。该方法会堵塞，直到管道读取端读完所有数据或读取端被关闭。如果读取端关闭时带有 error（即调用 CloseWithError 关闭），该Write返回的 err 就是读取端传递的error；否则 err 为 ErrClosedPipe。

- **创建同步的内存管道**

  ```go
  func Pipe() (*PipeReader, *PipeWriter)
  ```

### 函数

#### Copy 和 CopyN 函数

- **Copy 函数**

  ```go
  func Copy(dst Writer, src Reader) (written int64, err error)
  ```

  > Copy 将 src 复制到 dst，直到在 src 上到达 EOF 或发生错误。它返回复制的字节数，如果有错误的话，还会返回在复制时遇到的第一个错误。
  >
  > 成功的 Copy 返回 err == nil，而非 err == EOF。由于 Copy 被定义为从 src 读取直到 EOF 为止，因此它不会将来自 Read 的 EOF 当做错误来报告。
  >
  > 若 dst 实现了 ReaderFrom 接口，其复制操作可通过调用 dst.ReadFrom(src) 实现。此外，若 src 实现了 WriterTo 接口，其复制操作可通过调用 src.WriteTo(dst) 实现。

- **CopyN 函数**

  ```go
  func CopyN(dst Writer, src Reader, n int64) (written int64, err error)
  ```

  > CopyN 将 n 个字节(或到一个error)从 src 复制到 dst。 它返回复制的字节数以及在复制时遇到的最早的错误。当且仅当err == nil时,written == n 。
  >
  > 若 dst 实现了 ReaderFrom 接口，复制操作也就会使用它来实现。

#### ReadAtLeast 和 ReadFull 函数

> **区别**：ReadFull 将 buf 读满；而 ReadAtLeast 是最少读取 min 个字节。

- **ReadAtLeast 函数**

  ```go
  func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error)
  ```

  > ReadAtLeast 将 r 读取到 buf 中，直到读了最少 min 个字节为止。它返回复制的字节数，
  >
  > - 如果读取的字节较少，还会返回一个错误。
  > - 若没有读取到字节，错误就只是 EOF。
  > - 如果一个 EOF 发生在读取了少于 min 个字节之后，ReadAtLeast 就会返回 ErrUnexpectedEOF。
  > - 若 min 大于 buf 的长度，ReadAtLeast 就会返回 ErrShortBuffer。
  > - 对于返回值，当且仅当 err == nil 时，才有 n >= min。

- **ReadFull 函数**

  ```go
  func ReadFull(r Reader, buf []byte) (n int, err error)
  ```

  > ReadFull 精确地从 r 中将 len(buf) 个字节读取到 buf 中。它返回复制的字节数，
  >
  > - 如果读取的字节较少，还会返回一个错误。
  >
  > - 若没有读取到字节，错误就只是 EOF。
  > - 如果一个 EOF 发生在读取了一些但不是所有的字节后，
  > - ReadFull 就会返回 ErrUnexpectedEOF。
  > - 对于返回值，当且仅当 err == nil 时，才有 n == len(buf)。

#### WriteString 函数

> 方便写入 string 类型提供的函数。

```go
func WriteString(w Writer, s string) (n int, err error)
```

> WriteString 将s的内容写入w中，当 w 实现了 WriteString 方法时，会直接调用该方法，否则执行 w.Write([]byte(s))。

#### MultiReader 和 MultiWriter 函数

- **MultiReader 函数**

  ```go
  func MultiReader(readers ...Reader) Reader
  ```

  > 接收多个 Reader，返回一个 Reader。

- **MultiWriter 函数**

  ```go
  func MultiWriter(writers ...Writer) Writer
  ```

  > 接收多个 Writer，返回一个 Writer。

#### TeeReader函数

```go
func TeeReader(r Reader, w Writer) Reader
```

> TeeReader 返回一个 Reader，它将从 r 中读到的数据写入 w 中。所有经由它处理的从 r 的读取都匹配于对应的对 w 的写入。它没有内部缓存，即写入必须在读取完成前完成。任何在写入时遇到的错误都将作为读取错误返回。
>
> 这种功能的实现其实挺简单，无非是在 Read 完后执行 Write。

- **使用**

  ```go
  reader := io.TeeReader(strings.NewReader("Go语言中文网"), os.Stdout)
  reader.Read(make([]byte, 20))
  ```

  > 输出结果：Go语言中文网

## ioutil — 方便的IO操作函数集

> 标准库中提供了一些常用、方便的IO操作函数。

### 函数

#### NopCloser 函数

```go
func NopCloser(r io.Reader) io.ReadCloser
```

> 有时候我们需要传递一个 io.ReadCloser 的实例，而我们现在有一个 io.Reader 的实例。
>
> 它包装一个io.Reader，返回一个 io.ReadCloser ，而相应的 Close 方法啥也不做，只是返回 nil。

#### ReadAll 函数

```go
func ReadAll(r io.Reader) ([]byte, error)
```

> 很多时候，我们需要一次性读取 io.Reader 中的数据。
>
> 通过 bytes.Buffer 中的 [ReadFrom](http://docscn.studygolang.com/src/bytes/buffer.go?s=5385:5444#L144) 来实现读取所有数据的。

#### ReadDir 函数

```go
func ReadDir(dirname string) ([]fs.FileInfo, error)
```

> 读取目录并返回排好序的文件和子目录名。

#### ReadFile 和 WriteFile 函数

- **ReadFile 函数**

  ```go
  func ReadFile(filename string) ([]byte, error)
  ```

  > ReadFile 从 filename 指定的文件中读取数据并返回文件的内容。
  >
  > 成功的调用返回的err 为 nil 而非 EOF。
  >
  > 因为本函数定义为读取整个文件，它不会将读取返回的 EOF 视为应报告的错误。(同 ReadAll )

- **WriteFile 函数**

  ```go
  func WriteFile(filename string, data []byte, perm os.FileMode) error
  ```

  > WriteFile 将data写入filename文件中，
  >
  > - 文件不存在时会根据perm指定的权限进行创建一个。
  > - 文件存在时会先清空文件内容。
  >
  > 对于 perm 参数，我们一般可以指定为：0666，具体含义 os 包中讲解。

#### TempDir 和 TempFile 函数

- **TempDir 函数**

  ```go
  func TempDir(dir, pattern string) (name string, err error)
  ```

  > 第一个参数如果为空，表明在系统默认的临时目录（ os.TempDir ）中创建临时目录；第二个参数指定临时目录名的前缀，该函数返回临时目录的路径。

- **TempFile 函数**

  ```go
  func TempFile(dir, pattern string) (f *os.File, err error)
  ```

### 变量

#### Discard 变量

> Discard 对应的类型（`type devNull int`）实现了 io.Writer 接口，同时，为了优化 io.Copy 到 Discard，避免不必要的工作，实现了 io.ReaderFrom 接口。
>
> **用于清空数据**。

## fmt — 格式化IO

### 接口

#### Stringer 接口

```go
type Stringer interface {
  String() string
}
```

> 如果格式化输出某种类型的值，只要它实现了 String() 方法，那么会调用 String() 方法进行处理。

#### Formatter 接口

```go
type Formatter interface {
  Format(f State, c rune)
}
```

> Formatter 接口由带有定制的格式化器的值所实现。 Format 的实现可调用 Sprintf 或 Fprintf(f) 等函数来生成其输出。

- **fmt.State 是一个接口**：由于 Format 方法是被 fmt 包调用的，它内部会实例化好一个 fmt.State 接口的实例，我们不需要关心该接口
- **自定义占位符**：同时 fmt 包中和类型相对应的预定义占位符会无效。
- **Stringer 接口失效**：实现了 Formatter 接口，相应的 Stringer 接口不起作用。
- **第二个参数**：是占位符中%后的字母（有精度和宽度会被忽略，只保留字母）。

#### GoStringer 接口

```go
type GoStringer interface {
  GoString() string
}
```

> 定义了类型的Go语法格式。用于打印(Printf)格式化占位符为 %#v 的值。
>
> 一般的，我们不需要实现该接口。

#### Scanner 和 ScanState 接口

- **Scanner 接口**

  > 任何实现了 Scan 方法的对象都实现了 Scanner 接口，Scan 方法会从输入读取数据并将处理结果存入接收端，接收端必须是有效的指针。

- **ScanState 接口**

  > 是一个交给用户定制的 Scanner 接口的参数的接口。

### 函数

#### Printing

> - `S/F/Printf`函数通过**指定的格式**输出或格式化内容；
>
> - `S/F/Print`函数只是使用**默认的格式**输出或格式化内容；
>
> - `S/F/Println`函数使用**默认的格式**输出或格式化内容，同时会**在最后加上**"换行符"。

##### Fprint、Fprintf、Fprintln

> 第一个参数接收一个io.Writer类型，会将内容输出到 io.Writer 中去。

##### Sprint、Sprintf、Sprintln

> 格式化内容为 string 类型，而并不输出到某处，需要格式化字符串并返回时，可以用这组函数。

##### Print、Printf、Println

> 调用相应的F开头一类函数。将内容输出到标准输出中。

##### 占位符

###### 普通占位符

| 占位符 | 说明                                                         | 举例                                              | 输出                                |
| ------ | ------------------------------------------------------------ | ------------------------------------------------- | ----------------------------------- |
| %v     | 相应值的默认格式。<br>**在打印结构体时，“加号”标记（%+v）会添加字段名** | **Printf("%v", site)**<br>**Printf("%+v", site)** | {studygolang}<br>{Name:studygolang} |
| %#v    | 相应值的Go语法表示                                           | Printf("%#v", site)                               | main.Website{Name:"studygolang"}    |
| %T     | 相应值的类型的Go语法表示                                     | Printf("%T", site)                                | main.Website                        |
| %%     | 字面上的百分号，并非值的占位符                               | Printf("%%")                                      | %                                   |

###### 布尔占位符

| 占位符 | 说明                 | 举例               | 输出 |
| ------ | -------------------- | ------------------ | ---- |
| %t     | 单词 true 或 false。 | Printf("%t", true) | true |

###### 整数占位符

| 占位符 | 说明                                       | 举例                 | 输出   |
| ------ | ------------------------------------------ | -------------------- | ------ |
| %b     | 二进制表示                                 | Printf("%b", 5)      | 101    |
| %c     | 相应Unicode码点所表示的字符                | Printf("%c", 0x4E2D) | 中     |
| %d     | 十进制表示                                 | Printf("%d", 0x12)   | 18     |
| %o     | 八进制表示                                 | Printf("%o", 10)     | 12     |
| %q     | 单引号围绕的字符字面值，由Go语法安全地转义 | Printf("%q", 0x4E2D) | '中'   |
| %x     | 十六进制表示，字母形式为小写 a-f           | Printf("%x", 13)     | d      |
| %X     | 十六进制表示，字母形式为大写 A-F           | Printf("%x", 13)     | D      |
| %U     | Unicode格式：U+1234，等同于 "U+%04X"       | Printf("%U", 0x4E2D) | U+4E2D |

###### 浮点数和复数的组成部分（实部和虚部）

| 占位符 | 说明                                                         | 举例                   | 输出         |
| ------ | ------------------------------------------------------------ | ---------------------- | ------------ |
| %b     | 无小数部分的，指数为二的幂的科学计数法，与 strconv.FormatFloat 的 'b' 转换格式一致。例如 -123456p-78 |                        |              |
| %e     | 科学计数法，例如 -1234.456e+78                               | Printf("%e", 10.2)     | 1.020000e+01 |
| %E     | 科学计数法，例如 -1234.456E+78                               | Printf("%E", 10.2)     | 1.020000E+01 |
| %f     | 有小数点而无指数，例如 123.456                               | Printf("%f", 10.2)     | 10.200000    |
| %g     | 根据情况选择 %e 或 %f 以产生更紧凑的（无末尾的0）输出        | Printf("%g", 10.20)    | 10.2         |
| %G     | 根据情况选择 %E 或 %f 以产生更紧凑的（无末尾的0）输出        | Printf("%G", 10.20+2i) | (10.2+2i)    |

###### 字符串与字节切片

| 占位符 | 说明                                   | 举例                                 | 输出           |
| ------ | -------------------------------------- | ------------------------------------ | -------------- |
| %s     | 输出字符串表示（string类型或[]byte)    | Printf("%s", []byte("Go语言中文网")) | Go语言中文网   |
| %q     | 双引号围绕的字符串，由Go语法安全地转义 | Printf("%q", "Go语言中文网")         | "Go语言中文网" |
| %x     | 十六进制，小写字母，每字节两个字符     | Printf("%x", "golang")               | 676f6c616e67   |
| %X     | 十六进制，大写字母，每字节两个字符     | Printf("%X", "golang")               | 676F6C616E67   |

###### 指针

| 占位符 | 说明                  | 举例                | 输出     |
| ------ | --------------------- | ------------------- | -------- |
| %p     | 十六进制表示，前缀 0x | Printf("%p", &site) | 0x4f57f0 |

###### 其它标记

| 占位符 | 说明                                                         | 举例                  | 输出           |
| ------ | ------------------------------------------------------------ | --------------------- | -------------- |
| +      | 总打印数值的正负号<br>对于%q（%+q）保证只输出ASCII编码的字符。 | Printf("%+q", "中文") | "\u4e2d\u6587" |
| -      | 在右侧而非左侧填充空格（左对齐该区域）                       |                       |                |
| #      | 备用格式：<br>八进制：添加前导 0（%#o）<br>十六进制：添加前导 0x（%#x）或0X（%#X）<br>%p：（%#p）去掉前导 0x<br>%q：（%#q）打印原始（即反引号围绕的）字符串<br>如果是可打印字符，%U：（%#U）会写出该字符的Unicode 编码形式（如字符 x 会被打印成 U+0078 'x'）。 | Printf("%#U", '中')   | U+4E2D '中'    |
| ' '    | （空格）为数值中省略的正负号留出空白（% d）<br>以十六进制（% x, % X）打印字符串或切片时，在字节之间用空格隔开 |                       |                |
| 0      | 填充前导的0而非空格；对于数字，这会将填充移到正负号之后      |                       |                |

###### 宽度与精度

> 宽度与精度的控制格式以 Unicode 码点为单位。
>
> 宽度为输出的最小字符数，如果必要的话会为已格式化的形式填充空格

- **数值**：宽度为该数值占用区域的最小宽度；精度为小数点之后的位数。
  - %g/%G：精度为所有数字的总数。**默认精度**为确定该值所必须的最小位数。
  - %e 和 %f 的**默认精度**为6。
- **字符串**：精度为输出的最大字符数，如果必要的话会直接截断。

###### 格式化错误

> 如果给占位符提供了无效的实参（例如将一个字符串提供给 %d），所生成的字符串会包含该问题的描述。

| 错误类型              | 格式                        | 举例                                                 | 输出                            |
| --------------------- | --------------------------- | ---------------------------------------------------- | ------------------------------- |
| 类型错误或占位符未知  | %!verb(**type**=value)      | Printf("%d", hi)                                     | %!d(string=hi)                  |
| 实参太多              | %!(EXTRA **type**=value)    | Printf("hi", "guys")                                 | hi%!(EXTRA string=guys)         |
| 实参太少              | %!verb(MISSING)             | Printf("hi%d")                                       | hi %!d(MISSING)                 |
| 宽度或精度不是int类型 | %!(BADWIDTH) 或 %!(BADPREC) | Printf("%*s", 4.5, "hi")<br>Printf("%.s", 4.5, "hi") | %!(BADWIDTH)hi<br>%!(BADPREC)hi |

#### Scanning

> 1. **Scan、Scanf 、Scanln**：从 os.Stdin 中读取。
> 2. **Fscan、Fscanf、Fscanln**：从指定的 io.Reader 中读取。
> 3. **Sscan、Sscanf、Sscanln**：从实参字符串中读取。

##### Scanln、Fscanln 和 Sscanln

> 在换行符处停止扫描。且需要条目紧随换行符之后。

##### Scanf、Fscanf 和 Sscanf

> 需要输入换行符来匹配格式中的换行符

- 根据格式字符串解析实参，类似于 Printf，但有例外。

  - **%p、%T**： 没有实现。
  - **%e、%E、%f、%F、%g、%G**：完全等价，且可扫描任何浮点数或复数数值。
  - **%s、%v**：在扫描字符串时会将其中的空格作为分隔符。
    - **%v**：扫描整数时，可接受友好的进制前缀0（八进制）和0x（十六进制）。
  - **#、+**：没有实现

- **宽度**：被解释为输入的文本

  > %5s 意为最多从输入中读取5个 rune 来扫描成字符串。

- **精度**：没有精度的语法

- **连续空白字符**：（除换行符外）都等价于单个空格。

##### Scan、Fscan、Sscan

> 将换行符视为空格。
>
> 这组函数将连续由空格分隔的值存储为连续的实参。

## bufio — 缓存IO

### Reader 类型和方法

```go
type Reader struct {
  buf          []byte        // 缓存
  rd           io.Reader    // 底层的io.Reader
  // r:从buf中读走的字节（偏移）；w:buf中填充内容的偏移；
  // w - r 是buf中可被读的长度（缓存数据的大小），也是Buffered()方法的返回值
  r, w         int
  err          error        // 读过程中遇到的错误
  lastByte     int        // 最后一次读到的字节（ReadByte/UnreadByte)
  lastRuneSize int        // 最后一次读到的Rune的大小 (ReadRune/UnreadRune)
}

func NewReader(rd io.Reader) *Reader {
  // 默认缓存大小：defaultBufSize=4096
  return NewReaderSize(rd, defaultBufSize)
}
```

#### ReadSlice 方法

```go
func (b *Reader) ReadSlice(delim byte) (line []byte, err error)
```

> 从输入中读取，直到遇到第一个界定符（delim）为止【包含界定符】，返回一个指向缓存中字节的 slice。
>
> 在下次调用读操作（read）时，这些**字节会无效**：ReadSlice 返回的 []byte 是指向 Reader 中的 buffer ，而不是 copy 一份返回。

- 如果 ReadSlice 在找到界定符之前遇到了 error ，它就会返回缓存中所有的数据和错误本身（经常是 io.EOF）。
- 如果在找到界定符之前缓存已经满了，ReadSlice 会返回 bufio.ErrBufferFull 错误。
- 当且仅当返回的结果（line）没有以界定符结束的时候，ReadSlice 返回err != nil，也就是说，如果ReadSlice 返回的结果 line 不是以界定符 delim 结尾，那么返回的 er r也一定不等于 nil（可能是bufio.ErrBufferFull或io.EOF）。

#### ReadBytes 方法

```go
func (b *Reader) ReadBytes(delim byte) (line []byte, err error)
```

> ReadBytes 返回的 []byte 不会是指向 Reader 中的 buffer。**是一份拷贝**。

#### ReadString 方法

```go
func (b *Reader) ReadString(delim byte) (line string, err error) {
  bytes, err := b.ReadBytes(delim)
  return string(bytes), err
}
```

> 调用了 ReadBytes 方法，并将结果的 []byte 转为 string 类型。

#### ReadLine 方法

```go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
```

> ReadLine 是一个底层的原始行读取命令。
>
> ReadLine 尝试返回单独的行，不包括行尾的换行符。
>
> 如果一行大于缓存，isPrefix 会被设置为 true，同时返回该行的开始部分（等于缓存大小的部分）。
>
> 该行剩余的部分就会在下次调用的时候返回。当下次调用返回该行剩余部分时，isPrefix 将会是 false 。
>
> 跟 ReadSlice 一样，返回的 line 只是 buffer 的引用，在下次执行IO操作时，line 会无效。
>
> 返回值中，要么 line 不是 nil，要么 err 非 nil，两者不会同时非 nil

#### Peek 方法

```go
func (b *Reader) Peek(n int) ([]byte, error)
```

> 该方法只是“窥探”一下 Reader 中没有读取的 n 个字节。好比栈数据结构中的取栈顶元素，但不出栈。
>
> 返回的 []byte 只是 buffer 中的引用，在下次IO操作后会无效。

#### 其他方法

```go
func (b *Reader) Read(p []byte) (n int, err error)
func (b *Reader) ReadByte() (c byte, err error)
func (b *Reader) ReadRune() (r rune, size int, err error)
func (b *Reader) UnreadByte() error
func (b *Reader) UnreadRune() error
func (b *Reader) WriteTo(w io.Writer) (n int64, err error)
```

### Scanner 类型和方法

> 在 bufio 包中有多种方式获取文本输入，ReadBytes、ReadString 和独特的 ReadLine，对于简单的目的这些都有些过于复杂了。

```go
type Scanner struct {
  r            io.Reader // The reader provided by the client.
  split        SplitFunc // The function to split the tokens.
  maxTokenSize int       // Maximum size of a token; modified by tests.
  token        []byte    // Last token returned by split.
  buf          []byte    // Buffer used as argument to split.
  start        int       // First non-processed byte in buf.
  end          int       // End of data in buf.
  err          error     // Sticky error.
}

func NewScanner(r io.Reader) *Scanner {
  return &Scanner{
    r:            r,
    split:        ScanLines,
    maxTokenSize: MaxScanTokenSize,
    buf:          make([]byte, 4096), // Plausible starting size; needn't be large.
  }
}
```

- **SplitFunc**

  ```go
  type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
  ```

  > 用于对输入进行分词的 split 函数的签名。
  >
  > - 参数
  >   - data：是还未处理的数据
  >   - atEOF：标识 Reader 是否还有更多数据（是否到了EOF）。
  > - 返回值
  >   - advance：表示从输入中读取的字节数。
  >   - token：表示下一个结果数据。
  >   - err：则代表可能的错误。

  - **ScanBytes**：返回单个字节作为一个 token。
  - **ScanRunes**：返回单个 UTF-8 编码的 rune 作为一个 token。
  - **ScanWords**：返回通过“空格”分词的单词。即包括：`'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)`。
  - **ScanLines**：返回一行文本，不包括行尾的换行符。这里的换行包括了Windows下的"\r\n"和Unix下的"\n"。

- **split 字段**：代表了当前 Scanner 使用的分词策略，可以使用上面介绍的预定义 SplitFunc 实例赋值，也可以自定义 SplitFunc 实例。

  > 调用 Scanner 的 Split 方法赋值。

- **maxTokenSize 字段**：通过 split 分词后的一个 token 允许的最大长度。默认【64K】。

- **token 字段**：下一个结果数据。

#### Split 方法

```go
func (s *Scanner) Split(split SplitFunc)
```

> 可以通过 Split 方法为 Scanner 实例设置分词行为。默认 split 总是 ScanLines。

#### Scan 方法

```go
func (s *Scanner) Scan() bool
```

> 好比 iterator 中的 Next 方法，它用于将 Scanner 获取下一个 token，以便 Bytes 和 Text 方法可用。
>
> 当扫描停止时，它返回false，这时候，要么是到了输入的末尾要么是遇到了一个错误。

#### Bytes 和 Text 方法

```go
func (s *Scanner) Bytes() []byte
func (s *Scanner) Text() string
```

> 都是返回最近的 token。
>
> 该方法应该在 Scan 调用后调用，而且，下次调用 Scan 会覆盖这次的 token。

#### Err 方法

```go
func (s *Scanner) Err() error
```

> 通过 Err 方法可以获取第一个遇到的错误（**但如果错误是 io.EOF，Err 方法会返回 nil**）。

### Writer 类型和方法

```go
type Writer struct {
  err error        // 写过程中遇到的错误
  buf []byte       // 缓存
  n   int          // 当前缓存中的字节数
  wr  io.Writer    // 底层的 io.Writer 对象
}

func NewWriter(wr io.Writer) *Writer {
  // 默认缓存大小：defaultBufSize=4096
  return NewWriterSize(wr, defaultBufSize)
}
```

#### Available 和 Buffered 方法

> **Available 方法**：获取缓存中还未使用的字节数（缓存大小 - 字段 n 的值）。
>
> **Buffered 方法**：获取写入当前缓存中的字节数（字段 n 的值）

#### Flush 方法

> 该方法将缓存中的所有数据写入底层的 io.Writer 对象中。
>
> 使用 bufio.Writer 时，在**所有的 Write 操作完成之后**，应该调用 Flush 方法使得缓存都写入 io.Writer 对象中。

#### 其他方法

```go
// 实现了 io.ReaderFrom 接口
func (b *Writer) ReadFrom(r io.Reader) (n int64, err error)

// 实现了 io.Writer 接口
func (b *Writer) Write(p []byte) (nn int, err error)

// 实现了 io.ByteWriter 接口
func (b *Writer) WriteByte(c byte) error

// io 中没有该方法的接口，它用于写入单个 Unicode 码点，返回写入的字节数（码点占用的字节），内部实现会根据当前 rune 的范围调用 WriteByte 或 WriteString
func (b *Writer) WriteRune(r rune) (size int, err error)

// 写入字符串，如果返回写入的字节数比 len(s) 小，返回的error会解释原因
func (b *Writer) WriteString(s string) (int, error)
```

### ReadWriter 类型

```go
type ReadWriter struct {
  *Reader
  *Writer
}

func NewReadWriter(r *Reader, w *Writer) *ReadWriter {
  return &ReadWriter{r, w}
}
```

# 文本

> Go 标准库中有几个包专门用于处理文本。
>
> - *strings* 包提供了很多操作字符串的简单函数，通常一般的字符串操作需求都可以在这个包中找到。
> - *strconv* 包提供了基本数据类型和字符串之间的转换。
> - *regexp* 包提供了正则表达式功能，它的语法基于 [RE2](http://code.google.com/p/re2/wiki/Syntax) ，*regexp/syntax* 子包进行正则表达式解析。
> - 在标准库 *unicode* 包及其子包 utf8、utf16 中，提供了对 Unicode 相关编码、解码的支持，同时提供了测试 Unicode 码点（Unicode code points）属性的功能。

## strings — 字符串操作

### 类型

#### Replacer 类型

```go
type Replacer struct {
	once   sync.Once // guards buildOnce method
	r      replacer
	oldnew []string
}

func NewReplacer(oldnew ...string) *Replacer
```

- **WriteString**

  ```go
  func (r *Replacer) WriteString(w io.Writer, s string) (n int, err error)
  ```

  > 替换之后将结果写入 io.Writer 中。

#### Reader 类型

```go
type Reader struct {
  s        string    // Reader 读取的数据来源
  i        int // current reading index（当前读的索引位置）
  prevRune int // index of previous rune; or < 0（前一个读取的 rune 索引位置）
}

func NewReader(s string) *Reader
```

> 如果只是为了读取，NewReader 会更高效。

#### Builder 类型

```go
type Builder struct {
    addr *Builder // of receiver, to detect copies by value
    buf  []byte
}
```

> 该类型实现了 io 包下的 Writer, ByteWriter, StringWriter 等接口，可以向该对象内写入数据，
>
> Builder 没有实现 Reader 等接口，所以该类型不可读，但提供了 String 方法可以获取对象内的数据。

##### 方法

```go
// 该方法向 b 写入一个字节
func (b *Builder) WriteByte(c byte) error
// WriteRune 方法向 b 写入一个字符
func (b *Builder) WriteRune(r rune) (int, error)
// WriteRune 方法向 b 写入字节数组 p
func (b *Builder) Write(p []byte) (int, error)
// WriteRune 方法向 b 写入字符串 s
func (b *Builder) WriteString(s string) (int, error)
// Len 方法返回 b 的数据长度。
func (b *Builder) Len() int
// Cap 方法返回 b 的 cap。
func (b *Builder) Cap() int
// Grow 方法将 b 的 cap 至少增加 n (可能会更多)。如果 n 为负数，会导致 panic。
func (b *Builder) Grow(n int)
// Reset 方法将 b 清空 b 的所有内容。
func (b *Builder) Reset()
// String 方法将 b 的数据以 string 类型返回。
func (b *Builder) String() string
```

### 函数

#### 字符串比较

- **Compare**

  ```go
  func Compare(a, b string) int 
  ```

  > 用于比较两个字符串的大小，如果两个字符串相等，返回为 0。
  >
  > - 如果 a 小于 b ，返回 -1 ，反之返回 1 。
  > - **不推荐使用这个函数**，直接使用 == != > < >= <= 等一系列运算符更加直观。

- **EqualFold**

  ```go
  func EqualFold(s, t string) bool
  ```

  > 计算 s 与 t 忽略字母大小写后是否相等。

#### 是否存在某个字符或子串

- **Contains**

  ```go
  func Contains(s, substr string) bool
  ```

  > 子串 substr 在 s 中，返回 true。

- **ContainsAny**

  ```go
  func ContainsAny(s, chars string) bool
  ```

  > chars 中任何一个 Unicode 代码点在 s 中，返回 true。

- **ContainsRune**

  ```go
  func ContainsRune(s string, r rune) bool
  ```

  > Unicode 代码点 r 在 s 中，返回 true。

#### 子串出现次数 ( 字符串匹配 )

```go
func Count(s, sep string) int
```

> 查找子串出现次数即字符串模式匹配，实现的是 Rabin-Karp 算法。
>
> **当 sep 为空时，Count 的返回值是：utf8.RuneCountInString(s) + 1**

#### 字符串分割为[]string

- **Fields 和 FieldsFunc**

  ```go
  func Fields(s string) []string
  func FieldsFunc(s string, f func(rune) bool) []string
  ```

  > - **Fields** 用一个或多个连续的空格分隔字符串 s，返回子字符串的数组（slice）。如果字符串 s 只包含空格，则返回空列表 ([]string 的长度为 0）。
  >
  > - **FieldsFunc** 用这样的 Unicode 代码点 c 进行分隔：满足 f(c) 返回 true。该函数返回[]string。如果字符串 s 中所有的代码点 (unicode code points) 都满足 f(c) 或者 s 是空，则 FieldsFunc 返回空 slice。

- **Split 和 SplitAfter**、**SplitN 和 SplitAfterN**

  ```go
  func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }
  func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
  func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) }
  func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }
  ```

  > - Split 会将 s 中的 sep 去掉，而 SplitAfter 会保留 sep。
  > - 带 N 的方法可以通过最后一个参数 n 控制返回的结果中的 slice 中的元素个数
  >   - 当 n < 0 时，返回所有的子字符串；
  >   - 当 n == 0 时，返回的结果是 nil；
  >   - 当 n > 0 时，表示返回的 slice 中最多只有 n 个元素，其中，最后一个元素不会分割。

#### 字符串是否有某个前缀或后缀

- **HasPrefix**

  ```go
  func HasPrefix(s, prefix string) bool {
    return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
  }
  ```

  > s 中是否以 prefix 开始。

- **HasSuffix**

  ```go
  func HasSuffix(s, suffix string) bool {
    return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
  }
  ```

  > s 中是否以 suffix 结尾。

#### 字符或子串在字符串中出现的位置

```go
// 在 s 中查找 sep 的第一次出现，返回第一次出现的索引
func Index(s, sep string) int
// 在 s 中查找字节 c 的第一次出现，返回第一次出现的索引
func IndexByte(s string, c byte) int
// chars 中任何一个 Unicode 代码点在 s 中首次出现的位置
func IndexAny(s, chars string) int
// 查找字符 c 在 s 中第一次出现的位置，其中 c 满足 f(c) 返回 true
func IndexFunc(s string, f func(rune) bool) int
// Unicode 代码点 r 在 s 中第一次出现的位置
func IndexRune(s string, r rune) int

// 有三个对应的查找最后一次出现的位置
func LastIndex(s, sep string) int
func LastIndexByte(s string, c byte) int
func LastIndexAny(s, chars string) int
func LastIndexFunc(s string, f func(rune) bool) int
```

#### 字符串 JOIN 操作

```go
func Join(a []string, sep string) string
```

> 将字符串数组（或 slice）连接起来。

#### 字符串重复几次

```go
func Repeat(s string, count int) string
```

> 将 s 重复 count 次，如果 count 为负数或返回值长度 len(s)*count 超出 string 上限会导致 panic。

#### 字符替换

```go
func Map(mapping func(rune) rune, s string) string
```

> Map 函数，将 s 的每一个字符按照 mapping 的规则做映射替换，如果 mapping 返回值 <0 ，则舍弃该字符。

#### 字符串子串替换

- **Replace**

  ```go
  func Replace(s, old, new string, n int) string
  ```

  > 用 new 替换 s 中的 old，一共替换 n 个。
  >
  > 如果 n < 0，则不限制替换次数，即全部替换。

- **ReplaceAll**

  ```go
  func ReplaceAll(s, old, new string) string { return Replace(s, old, new , -1) }
  ```

#### 大小写转换

```go
func ToLower(s string) string
func ToLowerSpecial(c unicode.SpecialCase, s string) string
func ToUpper(s string) string
func ToUpperSpecial(c unicode.SpecialCase, s string) string
```

> ToLowerSpecial,ToUpperSpecial 可以转换特殊字符的大小写。
>
> ToLower,ToUpper 用于大小写转换。

#### 标题处理

- **Title**

  ```go
  func Title(s string) string
  ```

  > Title 会将 s 每个单词的首字母大写，不处理该单词的后续字符。

- **ToTitle**

  ```go
  func ToTitle(s string) string
  ```

  > 将 s 的每个字母大写。

- **ToTitleSpecial**

  ```go
  func ToTitleSpecial(c unicode.SpecialCase, s string) string
  ```

  > 将 s 的每个字母大写，并且会将一些特殊字母转换为其对应的特殊大写字母。

#### 修剪

```go
// 将 s 左侧和右侧中匹配 cutset 中的任一字符的字符去掉
func Trim(s string, cutset string) string
// 将 s 左侧的匹配 cutset 中的任一字符的字符去掉
func TrimLeft(s string, cutset string) string
// 将 s 右侧的匹配 cutset 中的任一字符的字符去掉
func TrimRight(s string, cutset string) string
// 如果 s 的前缀为 prefix 则返回去掉前缀后的 string , 否则 s 没有变化。
func TrimPrefix(s, prefix string) string
// 如果 s 的后缀为 suffix 则返回去掉后缀后的 string , 否则 s 没有变化。
func TrimSuffix(s, suffix string) string
// 将 s 左侧和右侧的间隔符去掉。常见间隔符包括：'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL)
func TrimSpace(s string) string
// 将 s 左侧和右侧的匹配 f 的字符去掉
func TrimFunc(s string, f func(rune) bool) string
// 将 s 左侧的匹配 f 的字符去掉
func TrimLeftFunc(s string, f func(rune) bool) string
// 将 s 右侧的匹配 f 的字符去掉
func TrimRightFunc(s string, f func(rune) bool) string
```

## bytes — byte slice 便利操作

### 类型

#### Reader 类型

```go
type Reader struct {
  s        []byte
  i        int64 // 当前读取下标
  prevRune int   // 前一个字符的下标，也可能 < 0
}

func NewReader(b []byte) *Reader
```

##### 方法

```go
// 读取数据至 b 
func (r *Reader) Read(b []byte) (n int, err error) 
// 读取一个字节
func (r *Reader) ReadByte() (byte, error)
// 读取一个字符
func (r *Reader) ReadRune() (ch rune, size int, err error)
// 读取数据至 w
func (r *Reader) WriteTo(w io.Writer) (n int64, err error)
// 进度下标指向前一个字节，如果 r.i <= 0 返回错误。
func (r *Reader) UnreadByte() 
// 进度下标指向前一个字符，如果 r.i <= 0 返回错误，且只能在每次 ReadRune 方法后使用一次，否则返回错误。
func (r *Reader) UnreadRune() 
// 读取 r.s[off:] 的数据至b，该方法忽略进度下标 i，不使用也不修改。
func (r *Reader) ReadAt(b []byte, off int64) (n int, err error) 
// 根据 whence 的值，修改并返回进度下标 i ，当 whence == 0 ，进度下标修改为 off，当 whence == 1 ，进度下标修改为 i+off，当 whence == 2 ，进度下标修改为 len[s]+off.
// off 可以为负数，whence 的只能为 0，1，2，当 whence 为其他值或计算后的进度下标越界，则返回错误。
func (r *Reader) Seek(offset int64, whence int) (int64, error)
```

#### Buffer 类型

```go
type Buffer struct {
  buf      []byte
  off      int   
  lastRead readOp 
}

func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }
func NewBufferString(s string) *Buffer {return &Buffer{buf: []byte(s)}}
```

> Buffer 包含了 21 个读写相关的方法，大部分同名方法的用法与前面的类似。

##### 方法

```go
// 读取到字节 delim 后，以字节数组的形式返回该字节及前面读取到的字节。如果遍历 b.buf 也找不到匹配的字节，则返回错误(一般是 EOF)
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error)
// 读取到字节 delim 后，以字符串的形式返回该字节及前面读取到的字节。如果遍历 b.buf 也找不到匹配的字节，则返回错误(一般是 EOF)
func (b *Buffer) ReadString(delim byte) (line string, err error)
// 截断 b.buf , 舍弃 b.off+n 之后的数据。n == 0 时，调用 Reset 方法重置该对象，当 n 越界时（n < 0 || n > b.Len() ）方法会触发 panic.
func (b *Buffer) Truncate(n int)
```

### 函数

#### 是否存在某个子 slice

```go
func Contains(b, subslice []byte) bool
```

> 是否存在某个子 slice。

#### []byte 出现次数

```go
func Count(s, sep []byte) int
```

> slice sep 在 s 中出现的次数（无重叠）。

#### Runes 类型转换

```go
func Runes(s []byte) []rune
```

> 将 []byte 转换为 []rune。

## strconv — 字符串和基本数据类型之间转换

### strconv 包转换错误处理

- *ErrRange* ：值超过了类型能表示的最大范围，比如将 "128" 转为 int8 就会返回这个错误。

- *ErrSyntax*：语法错误，比如将 "" 转为 int 类型会返回这个错误。

### 字符串和整型

#### 字符串 -> 整型

- **ParseInt**

  ```go
  func ParseInt(s string, base int, bitSize int) (i int64, err error)
  ```

  > *base*：字符串按照给定的进制进行解释。可以取 2~36（0-9，a-z），如果 base 的值为 0，则会根据字符串的前缀来确定 base 的值："0x" 表示 16 进制； "0" 表示 8 进制；否则就是 10 进制。
  >
  > *bitSize*：的是整数取值范围，或者说整数的具体类型。取值 0、8、16、32 和 64 分别代表 int、int8、int16、int32 和 int64。

- **ParseUint**

  ```go
  func ParseUint(s string, base int, bitSize int) (n uint64, err error)
  ```

  > 转为无符号整型。

- **Atoi**

  ```go
  func Atoi(s string) (i int, err error)
  ```

  > ParseInt 的便捷版，内部通过调用 *ParseInt(s, 10, 0)* 来实现的。

#### 整型 -> 字符串

- **FormatInt**

  ```go
  func FormatInt(i int64, base int) string
  ```

  > 有符号整型转字符串

- **FormatUint**

  ```go
  func FormatInt(i int64, base int) string
  ```

  > 有符号整型转字符串

- **Itoa**

  ```go
  func Itoa(i int) string
  ```

  > *Itoa* 内部直接调用 *FormatInt(i, 10)* 实现的。

### 字符串和布尔值

- **ParseBool**

  ```go
  func ParseBool(str string) (value bool, err error)
  ```

  > 1. 接受 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False 等字符串。
  > 2. 其他形式的字符串会返回错误。

- **FormatBool**

  ```go
  func FormatBool(b bool) string
  ```

  > 直接返回 "true" 或 "false"

- **AppendBool**

  ```go
  func AppendBool(dst []byte, b bool)
  ```

  > - 将 "true" 或 "false" append 到 dst 中。
  >
  > - 这里用了一个 append 函数对于字符串的特殊形式：append(dst, "true"...)

### 字符串和浮点数

- **ParseFloat**

  ```go
  func ParseFloat(s string, bitSize int) (f float64, err error)
  ```

  > 字符串转为浮点数

- **FormatFloat**

  ```go
  func FormatFloat(f float64, fmt byte, prec, bitSize int) string
  ```

  > *prec*：表示有效数字。
  >
  > - 对 *fmt='b'* 无效。
  >
  > - 对于 'e', 'E' 和 'f'，有效数字用于小数点之后的位数。
  > - 对于 'g' 和 'G'，则是所有的有效数字。
  >
  > 由于浮点数有精度的问题，精度不一样，ParseFloat 和 FormatFloat 可能达不到互逆的效果。
  >
  > **基于性能的考虑**，应该使用 *FormatFloat* 而不是 *fmt.Sprintf*。

- **AppendFloat**

### Quote

```go
func Quote(s string) string
```

> 为字符串增加双引号，即字面值字符串。

## regexp — 正则表达式

> 正则表达式实现保证运行时间随着输入大小线性增长的（即复杂度为 O(n)，其中 n 为输入的长度）。

Regexp 类型提供了多达 16 个方法，用于匹配正则表达式并获取匹配的结果。

## unicode — Unicode 码点、UTF-8/16 编码

> go 语言的所有代码都是 UTF8 的，所以如果我们在程序中的字符串都是 utf8 编码的，但是我们的单个字符（单引号扩起来的）却是 unicode 的。

### unicode 包

```go
// 是否控制字符
func IsControl(r rune) bool
// 是否阿拉伯数字字符，即 0-9
func IsDigit(r rune) bool
// 是否图形字符
func IsGraphic(r rune) bool
// 是否字母
func IsLetter(r rune) bool
// 是否小写字符
func IsLower(r rune) bool
// 是否符号字符
func IsMark(r rune) bool
// 是否数字字符，比如罗马数字Ⅷ也是数字字符
func IsNumber(r rune) bool
// 是否是 RangeTable 中的一个
func IsOneOf(ranges []*RangeTable, r rune) bool
// 是否可打印字符
func IsPrint(r rune) bool
// 是否标点符号
func IsPunct(r rune) bool
// 是否空格
func IsSpace(r rune) bool
// 是否符号字符
func IsSymbol(r rune) bool
// 是否 title case
func IsTitle(r rune) bool
// 是否大写字符
func IsUpper(r rune) bool
// r 是否为 rangeTab 类型的字符
func Is(rangeTab *RangeTable, r rune) bool
// r 是否为 ranges 中任意一个类型的字符
func In(r rune, ranges ...*RangeTable) bool
```

### utf8 包

- **判断是否符合 utf8 编码**

  ```go
  func Valid(p []byte) bool
  func ValidRune(r rune) bool
  func ValidString(s string) bool
  ```

- **判断 rune 所占字节数**

  ```go
  func RuneLen(r rune) int
  ```

- **判断字节数组或者字符串的 rune 数**

  ```go
  func RuneCount(p []byte) int
  func RuneCountInString(s string) (n int)
  ```

- **编码和解码到 rune**

  ```go
  func EncodeRune(p []byte, r rune) int
  func DecodeRune(p []byte) (r rune, size int)
  func DecodeRuneInString(s string) (r rune, size int)
  func DecodeLastRune(p []byte) (r rune, size int)
  func DecodeLastRuneInString(s string) (r rune, size int)
  ```

- **是否为完整 rune**

  ```go
  func FullRune(p []byte) bool
  func FullRuneInString(s string) bool
  ```

- **是否为 rune 第一个字节**

  ```go
  func RuneStart(b byte) bool
  ```

### utf16 包

```go
func Encode(s []rune) []uint16
func EncodeRune(r rune) (r1, r2 rune)
func Decode(s []uint16) []rune
func DecodeRune(r1, r2 rune) rune
func IsSurrogate(r rune) bool // 是否为有效代理对
```

# 数据结构与算法

> - *sort* 包：包含基本的排序方法，支持切片数据排序以及用户自定义数据集合排序。
> - *index/suffixary* 包：实现了后缀数组相关算法以支持许多常见的字符串操作。
> - *container* 包：提供了对 heap、list 和 ring 这 3 种数据结构的底层支持。

## sort — 排序算法

> 该包实现了四种基本排序算法：插入排序、归并排序、堆排序和快速排序。
>
> 但是这四种排序方法是不公开的，它们只被用于 sort 包内部使用。
>
> 为了方便对常用数据类型的操作，sort 包提供了对[]int 切片、[]float64 切片和[]string 切片完整支持。

### 数据集合排序

```go
type Interface interface {
  // 获取数据集合元素个数
  Len() int
  // 如果 i 索引的数据小于 j 索引的数据，返回 true，且不会调用下面的 Swap()，即数据升序排序。
  Less(i, j int) bool
  // 交换 i 和 j 索引的两个元素的位置
  Swap(i, j int)
}
```

- **IsSorted**

  ```go
  func IsSorted(data Interface) bool
  ```

  > 判断数据集合是否已经排好顺序。

- **Reverse**

  ```go
  func Reverse(data Interface) Interface
  ```

  > 将数据按 Less() 定义的排序方式逆序排序。

- **Search**

  ```go
  func Search(n int, f func(int) bool) int
  ```

  > 该方法会使用“二分查找”算法来找出能使 f(x)(0<=x<n) 返回 ture 的最小值 i。
  >
  > 前提条件 : 
  >
  > - f(x)(0<=x<i) 均返回 false。
  > - f(x)(i<=x<n) 均返回 ture。
  > - 如果不存在 i 可以使 f(i) 返回 ture, 则返回 n。

### sort包已经支持的内部数据类型排序

#### IntSlice 类型及[]int 排序

```go
type IntSlice []int
```

- **Ints**

  ```go
  func Ints(a []int) { Sort(IntSlice(a)) }
  ```

  > 对[]int 切片排序更常使用 sort.Ints()。

- **SearchInts**

  ```go
  func SearchInts(a []int, x int) int
  ```

  > 查找整数 x 在切片 a 中的位置。
  >
  > 使用条件为：**切片 a 已经升序排序** 。

#### Float64Slice 类型及[]float64 排序

```go
type Float64Slice []float64
```

- **Float64s**

  ```go
  func Float64s(a []float64)  
  ```

  > 排序

- **Float64sAreSorted**

  ```go
  func Float64sAreSorted(a []float64) bool
  ```

- **SearchFloat64s**

  ```go
  func SearchFloat64s(a []float64, x float64) int
  ```

#### StringSlice 类型及[]string 排序

> 两个 string 对象之间的大小比较是基于 「字典序」 的。

```go
type StringSlice []string
```

- **Strings**

  ```go
  func Strings(a []string)
  ```

  > 排序。

- **StringsAreSorted**

  ```go
  func StringsAreSorted(a []string) bool
  ```

- **SearchStrings**

  ```go
  func SearchStrings(a []string, x string) int
  ```

### []interface 排序与查找

- **sort.Slice**

  ```go
  func Slice(slice interface{}, less func(i, j int) bool) 
  ```

  > ```go
  > people := []struct {
  >     Name string
  >     Age  int
  > }{
  >     {"Gopher", 7},
  >     {"Alice", 55},
  >     {"Vera", 24},
  >     {"Bob", 75},
  > }
  > sort.Slice(people, func(i, j int) bool { return people[i].Age < people[j].Age }) // 按年龄升序排序
  > ```

- **sort.SliceStable**

  ```go
  func SliceStable(slice interface{}, less func(i, j int) bool) 
  ```

  > 该函数完成 []interface 的稳定排序。

- **sort.SliceIsSorted**

  ```go
  func SliceIsSorted(slice interface{}, less func(i, j int) bool) bool 
  ```

- **sort.Search**

  ```go
  func Search(n int, f func(int) bool) int
  ```

## index/suffixarray — 后缀数组实现子字符串查询

## container — 容器数据类型

### heap — 堆

> 堆使用的数据结构是最小二叉树，即根节点比左边子树和右边子树的所有值都小。

```go
type Interface interface {
  sort.Interface
  Push(x interface{}) // add x as element Len()
  Pop() interface{}   // remove and return element Len() - 1.
}
```

> ```go
> type IntHeap []int
> 
> func (h IntHeap) Len() int           { return len(h) }
> func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
> func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
> 
> func (h *IntHeap) Push(x interface{}) {
>   *h = append(*h, x.(int))
> }
> 
> func (h *IntHeap) Pop() interface{} {
>   old := *h
>   n := len(old)
>   x := old[n-1]
>   *h = old[0 : n-1]
>   return x
> }
> ```

### list — 链表

```go
type Element struct {
  next, prev *Element  // 上一个元素和下一个元素
  list *List  // 元素所在链表
  Value interface{}  // 元素
}

type List struct {
  root Element  // 链表的根元素
  len  int      // 链表的长度
}
```

#### Element

```go
func (e *Element) Next() *Element
func (e *Element) Prev() *Element
```

#### List

```go
func New() *List
func (l *List) Back() *Element   // 最后一个元素
func (l *List) Front() *Element  // 第一个元素
func (l *List) Init() *List  // 链表初始化
func (l *List) InsertAfter(v interface{}, mark *Element) *Element // 在某个元素后插入
func (l *List) InsertBefore(v interface{}, mark *Element) *Element  // 在某个元素前插入
func (l *List) Len() int // 在链表长度
func (l *List) MoveAfter(e, mark *Element)  // 把 e 元素移动到 mark 之后
func (l *List) MoveBefore(e, mark *Element)  // 把 e 元素移动到 mark 之前
func (l *List) MoveToBack(e *Element) // 把 e 元素移动到队列最后
func (l *List) MoveToFront(e *Element) // 把 e 元素移动到队列最头部
func (l *List) PushBack(v interface{}) *Element  // 在队列最后插入元素
func (l *List) PushBackList(other *List)  // 在队列最后插入接上新队列
func (l *List) PushFront(v interface{}) *Element  // 在队列头部插入元素
func (l *List) PushFrontList(other *List) // 在队列头部插入接上新队列
func (l *List) Remove(e *Element) interface{} // 删除某个元素
```

### ring — 环

```go
type Ring struct {
  next, prev *Ring
  Value      interface{}
}

func New(n int) *Ring  // 初始化环
```

```go
func (r *Ring) Do(f func(interface{}))  // 循环环进行操作
func (r *Ring) Len() int // 环长度
func (r *Ring) Link(s *Ring) *Ring // 连接两个环
func (r *Ring) Move(n int) *Ring // 指针从当前元素开始向后移动或者向前（n 可以为负数）
func (r *Ring) Next() *Ring // 当前元素的下个元素
func (r *Ring) Prev() *Ring // 当前元素的上个元素
func (r *Ring) Unlink(n int) *Ring // 从当前元素开始，删除 n 个元素
```

# 日期与时间 - time

## Location - 时区

> 代表一个地区，并表示该地区所在的时区（可能多个）。`Location` 通常代表地理位置的偏移，比如 CEST 和 CET 表示中欧。下一节将详细讲解 Location。
>
> Unix 系统以标准格式存于文件中，这些文件位于 /usr/share/zoneinfo，而本地时区可以通过 /etc/localtime 获取，这是一个符号链接，指向 /usr/share/zoneinfo 中某一个时区。

```go
type Location struct {
  name string
  zone []zone
  tx   []zoneTrans
  extend string
  cacheStart int64
  cacheEnd   int64
  cacheZone  *zone
}
```

- **LoadLocation**

  ```go
  func LoadLocation(name string) (*Location, error)
  ```

  > 获得特定时区的实例。

## Time

> 代表一个纳秒精度的时间点

```go
type Time struct {
  // sec gives the number of seconds elapsed since
  // January 1, year 1 00:00:00 UTC.
  sec int64

  // nsec specifies a non-negative nanosecond
  // offset within the second named by Seconds.
  // It must be in the range [0, 999999999].
  nsec int32

  // loc specifies the Location that should be used to
  // determine the minute, hour, month, day, and year
  // that correspond to this Time.
  // Only the zero Time has a nil Location.
  // In that case it is interpreted to mean UTC.
  loc *Location
}

func Now() Time {
  sec, nsec := now()
  return Time{sec + unixToInternal, nsec, Local}
}
```

> sec：并非 Unix 时间戳，加上的 `unixToInternal` 是 1-1-1 到 1970-1-1 经历的秒数。也就是 `Time` 中的 sec 是从 1-1-1 算起的秒数，而不是 Unix 时间戳。
>
> nsec：Unix 时间戳。

### 零值的判断

```go
func (t Time) IsZero() bool
```

> 零值是 sec 和 nsec 都是 0，表示 1 年 1 月 1 日。
>
> 判断 Time 表示的时间是否是 0 值。

### 与 Unix 时间戳的转换

```go
func Unix(sec int64, nsec int64) Time
```

> 通过 Unix 时间戳生成 `time.Time` 实例。

```go
func (t Time) Unix() int64
```

> 得到 Unix 时间戳。

```go
func (t Time) UnixNano() int64
```

> 得到 Unix 时间戳的纳秒表示。

### 格式化和解析

> **2006-01-02 15:04:05**，layout固定写法

- **解析**

  ```go
  func Parse(layout, value string) (Time, error)
  ```

  > 解析出来的时区却是 time.UTC。

  ```go
  func ParseInLocation(layout, value string, loc *Location) (Time, error)
  ```

  > 一般的，我们应该总是使用 `time.ParseInLocation` 来解析时间，并给第三个参数传递 `time.Local`。

- **格式化**

  ```go
  func (t Time) Format(layout string) string
  ```

  > 解析时间

  ```go
  func (t Time) String() string
  ```

  > 使用了 `2006-01-02 15:04:05.999999999 -0700 MST` 的layout。

### 序列化 / 反序列化 相关接口

> `Time` 实现了 `encoding` 包中的 `BinaryMarshaler`、`BinaryUnmarshaler`、`TextMarshaler` 和 `TextUnmarshaler` 接口；`encoding/json` 包中的 `Marshaler` 和 `Unmarshaler` 接口。

### Round 和 Truncate 方法

- **Round**

  ```go
  func (t Time) Round(d Duration) Time
  ```

  > 取最接近的时间整数。e.g.
  >
  > ```go
  > fmt.Println(t.Round(1 * time.Hour))
  > ```

- **Truncate**

  ```go
  func (t Time) Truncate(d Duration) Time
  ```

  > 向下取整。

## Duration — 时间常量

> 代表两个时间点之间经过的时间，以纳秒为单位。

```go
const (
  Nanosecond  Duration = 1
  Microsecond          = 1000 * Nanosecond
  Millisecond          = 1000 * Microsecond
  Second               = 1000 * Millisecond
  Minute               = 60 * Second
  Hour                 = 60 * Minute
)
```

- **输出**：直接输出时，因为实现了 `fmt.Stringer` 接口，会输出人类友好的可读形式，如：72h3m0.5s。

## Timer 和 Ticker — 定时器

### Timer

```go
type Timer struct {
  C <-chan Time     // The channel on which the time is delivered.
  r runtimeTimer
}

type timer struct {
  i int // heap index

  // Timer wakes up at when, and then at when+period, ... (period > 0 only)
  // each time calling f(now, arg) in the timer goroutine, so f must be
  // a well-behaved function and not block.
  when   int64
  period int64
  f      func(interface{}, uintptr)
  arg    interface{}
  seq    uintptr
}

func NewTimer(d Duration) *Timer {
  c := make(chan Time, 1)
  t := &Timer{
    C: c,
    r: runtimeTimer{
      when: when(d),
      f:    sendTime,
      arg:  c,
    },
  }
  startTimer(&t.r)
  return t
}
```

> - **when**：表示的时间到时，会往 Timer.C 中发送当前时间。`when` 表示的时间是纳秒时间，正常通过 `runtimeNano() + int64(d)` 赋值。
> - **f**：值是 `sendTime`，定时器时间到时，会调用 f，并将 `arg` 和 `seq` 传给 `f`。
> - **period**：因为 `Timer` 是一次性的，所以 `period` 保留默认值 0。

- **After**

  ```go
  func After(d Duration) <-chan Time
  ```

  > 模拟超时

- **Stop**

  ```go
  func (t *Timer) Stop() bool
  ```

  > 停止定时器

- **Reset**

  ```go
  func (t *Timer) Reset(d Duration) bool
  ```

  > 重置定时器

### Ticker

> `Ticker` 和 `Timer` 类似，区别是：`Ticker` 中的 `runtimeTimer` 字段的 `period` 字段会赋值为 `NewTicker(d Duration)` 中的 `d`，表示每间隔 `d` 纳秒，定时器就会触发一次。

- **停止**：除非程序终止前定时器一直需要触发，否则，不需要时应该调用 `Ticker.Stop` 来释放相关资源。
- **始终触发**：如果程序终止前需要定时器一直触发，可以使用更简单方便的 `time.Tick` 函数，因为 `Ticker` 实例隐藏起来了，因此，该函数启动的定时器无法停止。

## Weekday 和 Month — 常量类型

> 语义更明确，同时，实现 `fmt.Stringer` 接口，方便输出【英文名称】。

- **Weekday**

  ```go
  type Weekday int
  
  const (
    Sunday Weekday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
  )
  ```

- **Month**

  ```go
  type Month int
  
  const (
  	January Month = 1 + iota
  	February
  	March
  	April
  	May
  	June
  	July
  	August
  	September
  	October
  	November
  	December
  )
  ```

# 数学计算

## math — 基本数学函数

### 三角函数

```go
// 正弦函数
func Sin(x float64) float64
// 反正弦函数
func Asin(x float64) float64
// 双曲正弦
func Sinh(x float64) float64
// 反双曲正弦
func Asinh(x float64) float64

// 一次性返回 sin,cos
func Sincos(x float64) (sin, cos float64)

// 余弦函数
func Cos(x float64) float64
// 反余弦函数
func Acos(x float64) float64
// 双曲余弦
func Cosh(x float64) float64
// 反双曲余弦
func Acosh(x float64) float64

// 正切函数
func Tan(x float64) float64
// 反正切函数
func Atan(x float64) float64
func Atan2(y, x float64) float64
// 双曲正切
func Tanh(x float64) float64
// 反双曲正切
func Atanh(x float64) float64
```

### 幂次函数

```go
// 立方根函数
func Cbrt(x float64) float64
// x 的幂函数
func Pow(x, y float64) float64
// 10 根的幂函数
func Pow10(e int) float64
// 平方根
func Sqrt(x float64) float64
// 对数函数
func Log(x float64) float64
// 10 为底的对数函数
func Log10(x float64) float64
// 2 为底的对数函数
func Log2(x float64) float64
// log(1 + x)
func Log1p(x float64) float64
// 相当于 log2(x) 的绝对值
func Logb(x float64) float64
// 相当于 log2(x) 的绝对值的整数部分
func Ilogb(x float64) int
// 指数函数
func Exp(x float64) float64
// 2 为底的指数函数
func Exp2(x float64) float64
// Exp(x) - 1
func Expm1(x float64) float64
```

### 特殊函数

```go
// 正无穷
func Inf(sign int) float64
// 是否正无穷
func IsInf(f float64, sign int) bool
// 无穷值
func NaN() float64
// 是否是无穷值
func IsNaN(f float64) (is bool)
// 计算直角三角形的斜边长
func Hypot(p, q float64) float64
```

### 类型转化函数

```go
// float32 和 unit32 的转换
func Float32bits(f float32) uint32
// uint32 和 float32 的转换
func Float32frombits(b uint32) float32
// float64 和 uint64 的转换
func Float64bits(f float64) uint64
// uint64 和 float64 的转换
func Float64frombits(b uint64) float64
```

### 其他函数

```go
// 绝对值函数
func Abs(x float64) float64
// 向上取整
func Ceil(x float64) float64
// 向下取整
func Floor(x float64) float64
// 取模
func Mod(x, y float64) float64
// 分解 f，以得到 f 的整数和小数部分
func Modf(f float64) (int float64, frac float64)
// 分解 f，得到 f 的位数和指数
func Frexp(f float64) (frac float64, exp int)
// 取大值
func Max(x, y float64) float64
// 取小值
func Min(x, y float64) float64
// 复数的维数
func Dim(x, y float64) float64
// 0 阶贝塞尔函数
func J0(x float64) float64
// 1 阶贝塞尔函数
func J1(x float64) float64
// n 阶贝塞尔函数
func Jn(n int, x float64) float64
// 第二类贝塞尔函数 0 阶
func Y0(x float64) float64
// 第二类贝塞尔函数 1 阶
func Y1(x float64) float64
// 第二类贝塞尔函数 n 阶
func Yn(n int, x float64) float64
// 误差函数
func Erf(x float64) float64
// 余补误差函数
func Erfc(x float64) float64
// 以 y 的符号返回 x 值
func Copysign(x, y float64) float64
// 获取 x 的符号
func Signbit(x float64) bool
// 伽玛函数
func Gamma(x float64) float64
// 伽玛函数的自然对数
func Lgamma(x float64) (lgamma float64, sign int)
// value 乘以 2 的 exp 次幂
func Ldexp(frac float64, exp int) float64
// 返回参数 x 在参数 y 方向上可以表示的最接近的数值，若 x 等于 y，则返回 x
func Nextafter(x, y float64) (r float64)
// 返回参数 x 在参数 y 方向上可以表示的最接近的数值，若 x 等于 y，则返回 x
func Nextafter32(x, y float32) (r float32)
// 取余运算
func Remainder(x, y float64) float64
// 截取函数
func Trunc(x float64) float64
```

## math/big — 大数实现

## math/cmply — 复数基本函数

## math/rand — 伪随机数生成器

### source

```go
type Source interface {
  Int63() int64
  Seed(seed int64)
}

func NewSource(seed int64) Source
```

### rand

```go
type Rand struct {
  src Source
  s64 Source64 // non-nil if src is source64

  readVal int64
  readPos int8
}

func New(src Source) *Rand
```

# 文件系统

## os — 平台无关的操作系统功能实现

### 文件 I/O

#### OpenFile

```go
func OpenFile(name string, flag int, perm FileMode) (*File, error)
```

> 既能打开一个已经存在的文件，也能创建并打开一个新文件。
>
> - **flag**：于指定文件的访问模式，可用的值在 `os` 中定义为常量。
>
>   ```go
>   const (
>     O_RDONLY int = syscall.O_RDONLY // 只读模式打开文件
>     O_WRONLY int = syscall.O_WRONLY // 只写模式打开文件
>     O_RDWR   int = syscall.O_RDWR   // 读写模式打开文件
>     O_APPEND int = syscall.O_APPEND // 写操作时将数据附加到文件尾部
>     O_CREATE int = syscall.O_CREAT  // 如果不存在将创建一个新文件
>     O_EXCL   int = syscall.O_EXCL   // 和 O_CREATE 配合使用，文件必须不存在
>     O_SYNC   int = syscall.O_SYNC   // 打开文件用于同步 I/O
>     O_TRUNC  int = syscall.O_TRUNC  // 如果可能，打开时清空文件
>   )
>   ```
>
>   其中，`O_RDONLY`、`O_WRONLY`、`O_RDWR` 应该只指定一个，剩下的通过 `|` 操作符来指定。该函数内部会给 `flags` 加上 `syscall.O_CLOEXEC`，在 fork 子进程时会关闭通过 `OpenFile` 打开的文件，即子进程不会重用该文件描述符。
>
> - **perm**：指定了文件的模式和权限位，类型是 `os.FileMode`，文件模式位常量定义在 `os` 中。
>
>   ```go
>   const (
>     // 单字符是被 String 方法用于格式化的属性缩写。
>     ModeDir        FileMode = 1 << (32 - 1 - iota) // d: 目录
>     ModeAppend                                     // a: 只能写入，且只能写入到末尾
>     ModeExclusive                                  // l: 用于执行
>     ModeTemporary                                  // T: 临时文件（非备份文件）
>     ModeSymlink                                    // L: 符号链接（不是快捷方式文件）
>     ModeDevice                                     // D: 设备
>     ModeNamedPipe                                  // p: 命名管道（FIFO）
>     ModeSocket                                     // S: Unix 域 socket
>     ModeSetuid                                     // u: 表示文件具有其创建者用户 id 权限
>     ModeSetgid                                     // g: 表示文件具有其创建者组 id 的权限
>     ModeCharDevice                                 // c: 字符设备，需已设置 ModeDevice
>     ModeSticky                                     // t: 只有 root/ 创建者能删除 / 移动文件
>   
>     // 覆盖所有类型位（用于通过 & 获取类型位），对普通文件，所有这些位都不应被设置
>     ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
>     ModePerm FileMode = 0777 // 覆盖所有 Unix 权限位（用于通过 & 获取类型位）
>   )
>   ```

- **FileMode**

  ```go
  // 判断是否为dir
  func (m FileMode) IsDir() bool
  // 判断是否为普通文件
  func (m FileMode) IsRegular() bool
  // 获取权限位
  func (m FileMode) Perm() FileMode
  // 获取文件类型，通过类型位获取
  func (m FileMode) Type() FileMode
  ```

- **相关函数**

  ```go
  func Open(name string) (*File, error) {return OpenFile(name, O_RDONLY, 0)}
  func Create(name string) (*File, error) {return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)}
  ```

#### File.Read

```go
func (f *File) Read(b []byte) (n int, err error)
```

> **读取文件内容**
>
> `Read` 方法从 `f` 中读取最多 `len(b)` 字节数据并写入 `b`。它返回读取的字节数和可能遇到的任何错误。文件终止标志是读取 0 个字节且返回值 err 为 `io.EOF`。

#### File.Write

```go
func (f *File) Write(b []byte) (n int, err error)
```

> **数据写入文件**
>
> `Write` 向文件中写入 `len(b)` 字节数据。它返回写入的字节数和可能遇到的任何错误。如果返回值 `n!=len(b)`，本方法会返回一个非 nil 的错误。
>
> 注意：`Write` 调用成功并不能保证数据已经写入磁盘，因为内核会缓存磁盘的 I/O 操作。如果希望立刻将数据写入磁盘（一般场景不建议这么做，因为会影响性能），有两种办法：
>
> - 打开文件时指定 `os.O_SYNC`；
> - 调用 `File.Sync()` 方法。

#### File.Close

```go
func (f *File) Close() error
```

> 以下两种情况会导致 `Close` 返回错误：
>
> - 关闭一个未打开的文件。
> - 两次关闭同一个文件。
>
> 通常，我们不会去检查 `Close` 的错误。

#### File.Seek

```go
func (f *File) Seek(offset int64, whence int) (ret int64, err error)
```

> **改变文件偏移量**
>
> `Seek` 设置下一次读 / 写的位置。offset 为相对偏移量，而 whence 决定相对位置：0 为相对文件开头，1 为相对当前位置，2 为相对文件结尾。它返回新的偏移量（相对开头）和可能的错误。使用中，whence 应该使用 `os` 包中的常量：`SEEK_SET`、`SEEK_CUR` 和 `SEEK_END`。
>
> 注意：`Seek` 只是调整内核中与文件描述符相关的文件偏移量记录，并没有引起对任何物理设备的访问。

```go
file.Seek(0, os.SEEK_SET)    // 文件开始处
file.Seek(0, SEEK_END)       // 文件结尾处的下一个字节
file.Seek(-1, SEEK_END)      // 文件最后一个字节
file.Seek(-10, SEEK_CUR)     // 当前位置前 10 个字节
file.Seek(1000, SEEK_END)    // 文件结尾处的下 1001 个字节
```

### 截断文件

```go
func Truncate(name string, size int64) error
func (f *File) Truncate(size int64) error
```

> 如果文件当前长度大于参数 `size`，调用将丢弃超出部分，若小于参数 `size`，调用将在文件尾部添加一系列空字节或是一个文件空洞。

### 文件属性

```go
type FileInfo interface {
  Name() string       // 文件的名字（不含扩展名）
  Size() int64        // 普通文件返回值表示其大小；其他文件的返回值含义各系统不同
  Mode() FileMode     // 文件的模式位
  ModTime() time.Time // 文件的修改时间
  IsDir() bool        // 等价于 Mode().IsDir()
  Sys() interface{}   // 底层数据来源（可以返回 nil）
}
```

#### 改变文件时间戳

```go
func Chtimes(name string, atime time.Time, mtime time.Time) error
```

> 修改 name 指定的文件对象的访问时间和修改时间，类似 Unix 的 utime() 或 utimes() 函数。
>
> 底层的文件系统可能会截断 / 舍入时间单位到更低的精确度。

#### 文件属主

```go
func Chown(name string, uid, gid int) error
func Lchown(name string, uid, gid int) error
func (f *File) Chown(uid, gid int) error
```

> 符号链接下，lchown更改符号链接本身的所有者，而不是该符号链接所指向的文件的所有者。

#### 文件权限

```go
func IsPermission(err error) bool
```

> 在文件相关操作报错时，检查是否是权限的问题。

### 目录与链接

#### 创建和移除（硬）链接

- **创建**

  ```go
  func Link(oldname, newname string) error
  ```

  > `Link` 创建一个名为 newname 指向 oldname 的硬链接。如果出错，会返回 `*LinkError` 类型的错误。

- **移除**

  ```go
  func Remove(name string) error
  ```

  > `Remove` 删除 name 指定的文件或目录。如果出错，会返回 `*PathError` 类型的错误。
  >
  > 如果目录不为空，`Remove` 会返回失败。

#### 更改文件名

```go
func Rename(oldpath, newpath string) error
```

> `Rename` 修改一个文件的名字或移动一个文件。如果 `newpath` 已经存在，则替换它。注意，可能会有一些个操作系统特定的限制。

#### 使用符号链接

```go
func Symlink(oldname, newname string) error
```

> `Symlink` 创建一个名为 `newname` 指向 `oldname` 的符号链接。如果出错，会返回 `*LinkError` 类型的错误。
>
> 由 `oldname` 所命名的文件或目录**在调用时无需存在**。因为即便当时存在，也无法阻止后来将其删除。

```go
func Readlink(name string) (string, error)
```

> `Readlink` 获取 `name` 指定的符号链接指向的文件的路径。如果出错，会返回 `*PathError` 类型的错误。

#### 创建和移除目录

- **创建**

  ```go
  func Mkdir(name string, perm FileMode) error
  ```

  > `Mkdir` 使用指定的权限和名称创建一个目录。如果出错，会返回 `*PathError` 类型的错误。

- **删除**

  ```go
  func RemoveAll(path string) error
  ```

  > `RemoveAll` 删除 `path` 指定的文件，或目录及它包含的任何下级对象。它会尝试删除所有东西，除非遇到错误并返回。如果 `path` 指定的对象不存在，`RemoveAll` 会返回 nil 而不返回错误。

#### 读目录

```go
func (f *File) Readdirnames(n int) (names []string, err error)
```

> `Readdirnames` 读取目录 `f` 的内容，返回一个最多有 `n` 个成员的[]string，切片成员为目录中文件对象的名字，采用目录顺序。对本函数的下一次调用会返回上一次调用未读取的内容的信息。
>
> - 如果 n>0，`Readdirnames` 函数会返回一个最多 n 个成员的切片。这时，如果 `Readdirnames` 返回一个空切片，它会返回一个非 nil 的错误说明原因。如果到达了目录 `f` 的结尾，返回值 err 会是 `io.EOF`。
>
> - 如果 n<=0，`Readdirnames` 函数返回目录中剩余所有文件对象的名字构成的切片。此时，如果 `Readdirnames` 调用成功（读取所有内容直到结尾），它会返回该切片和 nil 的错误值。如果在到达结尾前遇到错误，会返回之前成功读取的名字构成的切片和该错误。

```go
func (f *File) Readdir(n int) (fi []FileInfo, err error)
```

> `Readdir` 内部会调用 `Readdirnames`，将得到的 `names` 构造路径，通过 `Lstat` 构造出 `[]FileInfo`。

## path/filepath — 操作路径

> 路径操作函数并不会校验路径是否真实存在。

### 解析路径名字符串

- **Dir**

  ```go
  func Dir(path string) string
  ```

  > 返回路径中除去最后一个路径元素的部分，即该路径最后一个元素所在的目录。
  >
  > 在使用 `Split` 去掉最后一个元素后，会简化路径并去掉末尾的斜杠。
  >
  > - 如果路径是空字符串，会返回 "."；
  >
  > - 如果路径由 1 到多个斜杠后跟 0 到多个非斜杠字符组成，会返回 "/"；
  >
  > - 其他任何情况下都不会返回以斜杠结尾的路径。

- **Base**

  ```go
  func Base(path string) string
  ```

  > 返回路径的最后一个元素。
  >
  > 在提取元素前会去掉末尾的斜杠。
  >
  > - 如果路径是 ""，会返回 "."；
  >
  > - 如果路径是只有一个斜杆构成的，会返回 "/"。

### 相对路径和绝对路径

- **IsAbs**

  ```go
  func IsAbs(path string) bool
  ```

  > 返回路径是否是一个绝对路径。

- **Abs**

  ```go
  func Abs(path string) (string, error)
  ```

  > 返回 `path` 代表的绝对路径，如果 `path` 不是绝对路径，会加入当前工作目录以使之成为绝对路径。
  >
  > 因为硬链接的存在，不能保证返回的绝对路径是唯一指向该地址的绝对路径。
  >
  > **出错**：在 `os.Getwd` 出错时，`Abs` 会返回该错误，一般不会出错，如果路径名长度超过系统限制，则会报错。

- **Rel**

  ```go
  func Rel(basepath, targpath string) (string, error)
  ```

  > `Rel` 函数返回一个相对路径，将 `basepath` 和该路径用路径分隔符连起来的新路径在词法上等价于 `targpath`。
  >
  > 也就是说，`Join(basepath, Rel(basepath, targpath))` 等价于 `targpath`。
  >
  > - 如果成功执行，返回值总是相对于 `basepath` 的，即使 `basepath` 和 `targpath` 没有共享的路径元素。
  >
  > - 如果两个参数一个是相对路径而另一个是绝对路径，或者 `targpath` 无法表示为相对于 `basepath` 的路径，将返回错误。

### 路径的切分和拼接

- **Split**

  ```go
  func Split(path string) (dir, file string)
  ```

  > `Split` 函数根据最后一个路径分隔符将路径 `path` 分隔为目录和文件名两部分（`dir` 和 `file`）。
  >
  > 如果路径中没有路径分隔符，函数返回值 `dir` 为空字符串，`file` 等于 `path`；
  >
  > 反之，如果路径中最后一个字符是 `/`，则 `dir` 等于 `path`，`file` 为空字符串。
  >
  > 返回值满足 `path == dir+file`。`dir` 非空时，最后一个字符总是 `/`。

- **Join**

  ```go
  func Join(elem ...string) string
  ```

  > `Join` 函数可以将任意数量的路径元素放入一个单一路径里，会根据需要添加路径分隔符。
  >
  > 结果是经过 `Clean` 的，所有的空字符串元素会被忽略。
  >
  > 对于拼接路径的需求，我们**应该总是使用** `Join` 函数来处理。

### 规整化路径

```go
func Clean(path string) string
```

> `Clean` 函数通过单纯的词法操作返回和 `path` 代表同一地址的最短路径。
>
> **规则**
>
> 1. 将连续的多个路径分隔符替换为单个路径分隔符
> 2. 剔除每一个 `.` 路径名元素（代表当前目录）
> 3. 剔除每一个路径内的 `..` 路径名元素（代表父目录）和它前面的非 `..` 路径名元素
> 4. 剔除开始于根路径的 `..` 路径名元素，即将路径开始处的 `/..` 替换为 `/`（假设路径分隔符是 `/`）

### 符号链接指向的路径名

```go
func EvalSymlinks(path string) (string, error)
```

> 将所有路径的符号链接都解析出来。除此之外，它返回的路径，是直接可访问的。
>
> `Readlink` 只解一次。

### 文件路径匹配

- **Match**

  ```go
  func Match(pattern, name string) (matched bool, err error)
  ```

  > - pattern
  >
  >   ```go
  >   pattern:
  >       { term }
  >   term:
  >       '*'         匹配 0 或多个非路径分隔符的字符
  >       '?'         匹配 1 个非路径分隔符的字符
  >       '[' [ '^' ] { character-range } ']'  
  >                     字符组（必须非空）
  >       c           匹配字符 c（c != '*', '?', '\\', '['）
  >       '\\' c      匹配字符 c
  >   character-range:
  >       c           匹配字符 c（c != '\\', '-', ']'）
  >       '\\' c      匹配字符 c
  >       lo '-' hi   匹配区间[lo, hi]内的字符
  >   ```
  >
  > `Match` 指示 `name` 是否和 shell 的文件模式匹配。要求 `pattern` 必须和 `name` 全匹配上，不只是子串。
  >
  > **很少使用**

- **Glob**

  ```go
  func Glob(pattern string) (matches []string, err error)
  ```

  > 返回所有匹配了 模式字符串 `pattern` 的文件列表或者 nil（如果没有匹配的文件）。
  >
  > **错误**：忽略任何文件系统相关的错误，如读目录引发的 I/O 错误。唯一的错误和 `Match` 一样，在 `pattern` 不合法时，返回 `filepath.ErrBadPattern`。返回的结果是根据文件名字典顺序进行了排序的。

### 遍历目录

```go
func Walk(root string, walkFn WalkFunc) error
```

> - `Walk` 函数会遍历 `root` 指定的目录下的文件树，对每一个该文件树中的目录和文件都会调用 `walkFn`，包括 `root` 自身。
>
> - 所有访问文件 / 目录时遇到的错误都会传递给 `walkFn` 过滤。
>
> - 文件是按字典顺序遍历的，这让输出更漂亮，但也导致处理非常大的目录时效率会降低。
>
> - `Walk` 函数不会遍历文件树中的符号链接（快捷方式）文件包含的路径。

```go
type WalkFunc func(path string, info os.FileInfo, err error) error
```

> - `Walk` 函数对每一个文件 / 目录都会调用 `WalkFunc` 函数类型值。
> - 调用时 `path` 参数会包含 `Walk` 的 `root` 参数作为前缀；就是说，如果 `Walk` 函数的 `root` 为 "dir"，该目录下有文件 "a"，将会使用 "dir/a" 作为调用 `walkFn` 的参数。
> - `walkFn` 参数被调用时的 `info` 参数是 `path` 指定的地址（文件 / 目录）的文件信息，类型为 `os.FileInfo`。
> - 在回调函数中，如果返回 fs.SkipDir，则会**停止该目录的遍历**。

## io/fs — 抽象文件系统

> 该包定义了一个文件系统需要的相关基础接口，因此我们称之为抽象文件系统。
>
> 该文件系统是层级文件系统或叫树形文件系统，Unix 文件系统就是这种类型。

### FS

```go
type FS interface {
    Open(name string) (File, error)
}
```

> 该接口只有一个方法，即打开一个命名文件。实现要求如下：
>
> - 如果 Open 方法出错，应该返回 `*PathError` 类型的错误。
>
>   ```go
>   type PathError struct {
>     Op   string
>     Path string
>     Err  error
>   }
>   ```
>
>   > - Op 字段设置为 "open"。
>   >
>   > - Path 字段设置为文件名。
>   >
>   > - 而 Err 字段描述错误原因。
>
> - 对于指定的文件名，需要满足 `ValidPath(name)` 函数，如果不满足，则返回 `*PathError` 的 Err 为 fs.ErrInvalid 或 fs.ErrNotExist 的错误。
>
>   ```go
>   func ValidPath(name string) bool
>   ```

### File

```go
type File interface {
  Stat() (FileInfo, error)
  Read([]byte) (int, error)
  Close() error
}
```

### FileInfo

```go
type FileInfo interface {
  Name() string       // base name of the file
  Size() int64        // length in bytes for regular files; system-dependent for others
  Mode() FileMode     // file mode bits
  ModTime() time.Time // modification time
  IsDir() bool        // abbreviation for Mode().IsDir()
  Sys() interface{}   // underlying data source (can return nil)
}
```

> 该接口描述一个文件的元数据信息，它由 Stat 返回。

### DirEntry

```go
type DirEntry interface {
    Name() string
    IsDir() bool
    Type() FileMode
    Info() (FileInfo, error)
}
```

> ```go
> func ReadDir(fsys FS, name string) ([]DirEntry, error)
> ```

> 为了方便遍历文件系统（目录），io/fs 包提供了 ReadDir 函数。
>
> **函数的实现**：
>
> - 如果第一个参数实现了 fs.ReadDirFS 接口，直接调用该接口的 ReadDir 方法。
> - 否则看是否实现了 fs.ReadDirFile 接口，没实现则报错。

#### ReadDirFS

```go
type ReadDirFS interface {
  FS
  ReadDir(name string) ([]DirEntry, error)
}
```

#### ReadDirFile

```go
type ReadDirFile interface {
  File
  ReadDir(n int) ([]DirEntry, error)
}
```

### ReadFileFS

```go
type ReadFileFS interface {
    FS
    ReadFile(name string) ([]byte, error)
}
```

### StatFS

```go
type StatFS interface {
    FS
    Stat(name string) (FileInfo, error)
}
```

### GlobFS

```go
type GlobFS interface {
    FS
    Glob(pattern string) ([]string, error)
}
```

### SubFS

```go
type SubFS interface {
  FS
  Sub(dir string) (FS, error)
}
```

> 这个接口的作用主要是让一个文件系统支持定义子文件系统。

# 进程、线程与goroutine

## 进程

### 创建

> `os` 包及其子包 `os/exec` 提供了创建进程的方法。

#### 系统调用

- **fork**：允许一进程（父进程）创建一新进程（子进程）。具体做法是：

  > 新的子进程几近于对父进程的翻版

  - 子进程获得父进程的栈、数据段、堆和执行文本段的拷贝。
  - 可将此视为把父进程一分为二。

- **exit(status)**：终止一进程，将进程占用的所有资源（内存、文件描述符等）归还内核，交其进行再次分配。

  -  `status` ：一整型变量，表示进程的退出状态。
  - 父进程可使用系统调用 `wait()` 来获取该状态。

- **wait(&status)** ：目的有二：

  - 其一，如果子进程尚未调用 `exit()` 终止，那么 `wait` 会挂起父进程直至子进程终止；
  - 其二，子进程的终止状态通过 `wait` 的 `status` 参数返回。

- **execve(pathname, argv, envp)** ：加载一个新程序到当前进程的内存。

  - 路径名为 pathname，参数列表为 argv，环境变量列表为 envp。
  - 这将丢弃现存的程序文本段，并为新程序重新创建栈、数据段以及堆。通常将这一动作称为执行一个新程序。

#### Process 及其相关方法

> `Process` 提供了四个方法：`Kill`、`Signal`、`Wait` 和 `Release`。

```go
type Process struct {
  Pid    int
  handle uintptr // handle is accessed atomically on Windows
  isdone uint32  // process has been successfully waited on, non zero if true
}

func StartProcess(name string, argv []string, attr *ProcAttr) (*Process, error)
```

- **ProcAttr**

  ```go
  type ProcAttr struct {
      // 如果 Dir 非空，子进程会在创建 Process 实例前先进入该目录。（即设为子进程的当前工作目录）
      Dir string
      // 如果 Env 非空，它会作为新进程的环境变量。必须采用 Environ 返回值的格式。
      // 如果 Env 为 nil，将使用 Environ 函数的返回值。
      Env []string
      // Files 指定被新进程继承的打开文件对象。
      // 前三个绑定为标准输入、标准输出、标准错误输出。
      // 依赖底层操作系统的实现可能会支持额外的文件对象。
      // nil 相当于在进程开始时关闭的文件对象。
      Files []*File
      // 操作系统特定的创建属性。
      // 注意设置本字段意味着你的程序可能会执行异常甚至在某些操作系统中无法通过编译。这时候可以通过为特定系统设置。
      // 看 syscall.SysProcAttr 的定义，可以知道用于控制进程的相关属性。
      Sys *syscall.SysProcAttr
  }
  ```

- **FindProcess**

  ```go
  func FindProcess(pid int) (*Process, error)
  ```

  > 通过 `pid` 查找一个运行中的进程。该函数返回的 `Process` 对象可以用于获取关于底层操作系统进程的信息。

- **Wait**

  ```go
  func (p *Process) Wait() (*ProcessState, error)
  ```

  > `Wait` 方法阻塞直到进程退出，然后返回一个 `ProcessState` 描述进程的状态和可能的错误。`Wait` 方法会释放绑定到 `Process` 的所有资源。

- **ProcessState**

  ```go
  type ProcessState struct {
    pid    int                // The process's id.
    status syscall.WaitStatus // System-dependent status info.
    rusage *syscall.Rusage
  }
  ```

  > **status**：记录了状态原因。
  >
  > **rusage**：用于统计进程的各类资源信息。

- **syscal.WaitStatus**

  - Exited()：是否正常退出，如调用 `os.Exit`。
  - Signaled()：是否收到未处理信号而终止。
  - CoreDump()：是否收到未处理信号而终止，同时生成 coredump 文件，如 SIGABRT。
  - Stopped()：是否因信号而停止（SIGSTOP）。
  - Continued()：是否因收到信号 SIGCONT 而恢复。

### 运行外部命令

#### 查找可执行程序

```go
func LookPath(file string) (string, error)
```

> 如果在 `PATH` 中没有找到可执行文件，则返回 `exec.ErrNotFound`。

#### Cmd 及其相关方法

```go
type Cmd struct {
    // Path 是将要执行的命令路径。
    // 该字段不能为空（也是唯一一个不能为空的字段），如为相对路径会相对于 Dir 字段。
    // 通过 Command 初始化时，会在需要时调用 LookPath 获得完整的路径。
    Path string

    // Args 存放着命令的参数，第一个值是要执行的命令（Args[0])；如果为空切片或者 nil，使用 {Path} 运行。
    // 一般情况下，Path 和 Args 都应被 Command 函数设定。
    Args []string

    // Env 指定进程的环境变量，如为 nil，则使用当前进程的环境变量，即 os.Environ()，一般就是当前系统的环境变量。
    Env []string

    // Dir 指定命令的工作目录。如为空字符串，会在调用者的进程当前工作目录下执行。
    Dir string

    // Stdin 指定进程的标准输入，如为 nil，进程会从空设备读取（os.DevNull）
    // 如果 Stdin 是 *os.File 的实例，进程的标准输入会直接指向这个文件
    // 否则，会在一个单独的 goroutine 中从 Stdin 中读数据，然后将数据通过管道传递到该命令中（也就是从 Stdin 读到数据后，写入管道，该命令可以从管道读到这个数据）。在 goroutine 停止数据拷贝之前（停止的原因如遇到 EOF 或其他错误，或管道的 write 端错误），Wait 方法会一直堵塞。
    Stdin io.Reader

    // Stdout 和 Stderr 指定进程的标准输出和标准错误输出。
    // 如果任一个为 nil，Run 方法会将对应的文件描述符关联到空设备（os.DevNull）
    // 如果两个字段相同，同一时间最多有一个线程可以写入。
    Stdout io.Writer
    Stderr io.Writer

    // ExtraFiles 指定额外被新进程继承的已打开文件，不包括标准输入、标准输出、标准错误输出。
    // 如果本字段非 nil，其中的元素 i 会变成文件描述符 3+i。
    //
    // BUG: 在 OS X 10.6 系统中，子进程可能会继承不期望的文件描述符。
    // http://golang.org/issue/2603
    ExtraFiles []*os.File

    // SysProcAttr 提供可选的、各操作系统特定的 sys 属性。
    // Run 方法会将它作为 os.ProcAttr 的 Sys 字段传递给 os.StartProcess 函数。
    SysProcAttr *syscall.SysProcAttr

    // Process 是底层的，只执行一次的进程。
    Process *os.Process

    // ProcessState 包含一个已经存在的进程的信息，只有在调用 Wait 或 Run 后才可用。
    ProcessState *os.ProcessState
}
```

- **Command**

  ```go
  func Command(name string, arg ...string) *Cmd
  ```

  > 函数返回一个 `*Cmd`，用于使用给出的参数执行 `name` 指定的程序。返回的 `*Cmd` 只设定了 `Path` 和 `Args` 两个字段。
  >
  > 一般的，应该通过 `exec.Command` 函数产生 `Cmd` 实例。

- **Start**

  ```go
  func (c *Cmd) Start() error
  ```

  > 开始执行 `c` 包含的命令，但并不会等待该命令完成即返回。

- **Wait**

  ```go
  func (c *Cmd) Wait() error
  ```

  > `Wait` 会阻塞直到该命令执行完成，该命令必须是先通过 `Start` 执行。
  >
  > `Wait` 方法会在命令返回后释放相关的资源。

- **Output**

  ```go
  func (c *Cmd) Output() ([]byte, error)
  ```

  > 除了 `Run()` 是 `Start`+`Wait` 的简便写法，`Output()` 更是 `Run()` 的简便写法，外加获取外部命令的输出。
  >
  > 要求 `c.Stdout` 必须是 `nil`，内部会将 `bytes.Buffer` 赋值给 `c.Stdout`。

- **CombinedOutput**

  ```go
  func (c *Cmd) CombinedOutput() ([]byte, error)
  ```

  > `Output()` 只返回 `Stdout` 的结果，而 `CombinedOutput` 组合 `Stdout` 和 `Stderr` 的输出，即 `Stdout` 和 `Stderr` 都赋值为同一个 `bytes.Buffer`。

- **StdoutPipe**

  ```go
  func (c *Cmd) StdoutPipe() (io.ReadCloser, error)
  ```

  > 返回一个在命令 `Start` 执行后与命令标准输出关联的管道。
  >
  > `Wait` 方法会在命令结束后会关闭这个管道，所以一般不需要手动关闭该管道。
  >
  > 但是在从管道读取完全部数据之前调用 `Wait` 出错了，则必须手动关闭。

- **StderrPipe**

  ```go
  func (c *Cmd) StderrPipe() (io.ReadCloser, error)
  ```

  > 返回一个在命令 `Start` 执行后与命令标准错误输出关联的管道。
  >
  > `Wait` 方法会在命令结束后会关闭这个管道，一般不需要手动关闭该管道。
  >
  > 但是在从管道读取完全部数据之前调用 `Wait` 出错了，则必须手动关闭。

- **StdinPipe**

  ```go
  func (c *Cmd) StdinPipe() (io.WriteCloser, error)
  ```

  > `StdinPipe` 方法返回一个在命令 `Start` 执行后与命令标准输入关联的管道。
  >
  > `Wait` 方法会在命令结束后会关闭这个管道。必要时调用者可以调用 `Close` 方法来强行关闭管道。
  >
  > 例如，标准输入已经关闭了，命令执行才完成，这时调用者需要显示关闭管道。

### 进程终止

```go
func Exit(code int)
```

> `Exit` 让当前进程以给出的状态码 `code` 退出。一般来说，状态码 0 表示成功，非 0 表示出错。
>
> 进程会立刻终止，defer 的函数不会被执行。

### 进程属性和控制

#### 进程 ID

> - 每个进程都会有一个进程 ID，可以通过 `os.Getpid` 获得。
>
> - 同时，每个进程都有创建自己的父进程，通过 `os.Getppid` 获得。

#### 进程凭证

> Unix 中进程都有一套数字表示的用户 ID(UID) 和组 ID(GID)，有时也将这些 ID 称之为进程凭证。

- **实际用户 ID 和实际组 ID**

  > 实际用户 ID（real user ID）和实际组 ID（real group ID）确定了进程所属的用户和组。

  - *实际用户 ID*

    ```go
    os.Getuid()
    ```

    > 获取当前进程的实际用户 ID

  - *实际组 ID*

    ```go
    os.Getgid()
    ```

    > 获取当前进程的实际组 ID

- **有效用户 ID 和有效组 ID**

  > 通常，有效用户 ID 及组 ID 与其相应的实际 ID 相等，但有两种方法能够致使二者不同。
  >
  > - 一是使用相关系统调用；
  > - 二是执行 set-user-ID 和 set-group-ID 程序。

  - *有效用户 ID*

    ```go
    os.Geteuid()
    ```

    > 获取当前进程的有效用户 ID。

  - *有效组 ID*

    ```go
    os.Getegid()
    ```

    > 获取当前进程的有效组 ID。

#### 操作系统用户 - os/user

```go
type User struct {
  Uid      string // user id
  Gid      string // primary group id
  Username string
  Name     string
  HomeDir  string
}
```

> `Current` 函数可以获取当前用户账号。而 `Lookup` 和 `LookupId` 则分别根据用户名和用户 ID 查询用户。

#### 当前工作目录

- **Getwd**

  ```go
  func Getwd() (dir string, err error)
  ```

  >  返回一个对应当前工作目录的根路径。如果当前目录可以经过多条路径抵达（比如符号链接），`Getwd` 会返回其中一个。

- **Chdir**

  ```go
  func Chdir(dir string) error
  ```

  > 将当前工作目录修改为 `dir` 指定的目录。如果出错，会返回 `*PathError` 类型的错误。

#### 改变进程的根目录

> syscall.Chroot。能改变一个进程的根目录。

```go
func Chroot(path string) (err error)
```

#### 进程环境

- **Environ**

  ```go
  func Environ() []string
  ```

  > 获取环境列表。

- **Getenv**

  ```go
  func Getenv(key string) string
  ```

  > 检索并返回名为 `key` 的环境变量的值。如果不存在该环境变量会返回空字符串。
  >
  > 有时候，可能环境变量存在，只是值刚好是空。

- **LookupEnv**

  ```go
  func LookupEnv(key string) (string, bool)
  ```

  > 如果变量名存在，第二个参数返回 `true`，否则返回 `false`。

- **Setenv**

  ```go
  func Setenv(key, value string) error
  ```

  > 设置名为 `key` 的环境变量，值为 `value`。如果出错会返回该错误。
  >
  > 如果值之前存在，会覆盖。

- **Unsetenv**

  ```go
  func Unsetenv(key string) error
  ```

  > 删除名为 `key` 的环境变量。

- **Clearenv**

  ```go
  func Clearenv()
  ```

  > 删除所有环境变量。

- **Expand**

  ```go
  func Expand(s string, mapping func(string) string) string
  ```

  > 能够将 ${var} 或 $var 形式的变量，经过 mapping 处理，得到结果。

- **ExpandEnv**

  ```go
  func ExpandEnv(s string) string {Expand(s, Getenv)}
  ```

## 线程

> 1. 线程是允许应用程序并发执行多个任务的一种机制。一个进程可以包含多个线程。同一个程序中的所有线程均会独立执行相同程序，且共享同一份全局内存区域。
> 2. 同一进程中的多个线程可以并发执行。在多处理器环境下，多个线程可以同时并行。

## 进程间通信

# 测试

## testing

### 单元测试

1. 单元测试中，传递给测试函数的参数是 `*testing.T` 类型。它用于管理测试状态并支持格式化测试日志。
2. 当测试函数返回时，或者当测试函数调用 `FailNow`、 `Fatal`、`Fatalf`、`SkipNow`、`Skip`、`Skipf` 中的任意一个时，则宣告该测试函数结束。跟 `Parallel` 方法一样，以上提到的这些方法只能在运行测试函数的 goroutine 中调用。
3. 比如 `Log` 以及 `Error` 的变种， 可以在多个 goroutine 中同时进行调用。

### 子测试与子基准测试

```go
func TestFoo(t *testing.T) {
  // <setup code>
  t.Run("A=1", func(t *testing.T) { ... })
  t.Run("A=2", func(t *testing.T) { ... })
  t.Run("B=1", func(t *testing.T) { ... })
  // <tear-down code>
}
```

- 子测试可用于程序并行控制。只有子测试全部执行完毕后，父测试才会完成。

### TestMain

> 有时需要在测试之前或之后进行额外的设置（setup）或拆卸（teardown）。
>
> - 可以在调用 `m.Run` 前后做任何设置和拆卸。注意，在 `TestMain` 函数的最后，应该使用 `m.Run` 的返回值作为参数去调用 `os.Exit`。
> - 则应该显式地调用 `flag.Parse`。

```go
func TestMain(m *testing.M) {
  db.Dns = os.Getenv("DATABASE_DNS")
  if db.Dns == "" {
    db.Dns = "root:123456@tcp(localhost:3306)/?charset=utf8&parseTime=True&loc=Local"
  }

  flag.Parse()
  exitCode := m.Run()

  db.Dns = ""

  // 退出
  os.Exit(exitCode)
}
```

## httptest

> HTTP 测试辅助工具

```go
func TestHandleGet(t *testing.T) {
  mux := http.NewServeMux()
  mux.HandleFunc("/topic/", handleRequest)

  r, _ := http.NewRequest(http.MethodGet, "/topic/1", nil)

  w := httptest.NewRecorder()

  mux.ServeHTTP(w, r)

  resp := w.Result()
  if resp.StatusCode != http.StatusOK {
    t.Errorf("Response code is %v", resp.StatusCode)
  }

  topic := new(Topic)
  json.Unmarshal(w.Body.Bytes(), topic)
  if topic.Id != 1 {
    t.Errorf("Cannot get topic")
  }
}
```

# 应用构件与debug

## flag - 命令行参数解析

> 实现了命令行参数的解析。

### 定义

- **flag.Xxx()**

  ```go
  var ip = flag.Int("flagname", 1234, "help message for flagname")
  ```

  > 返回一个相应类型的指针。

- **flag.XxxVar()**

  ```go
  var flagvar int
  flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
  ```

- **自定义 Value**

  ```go
  flag.Var(&flagVal, "name", "help message for flagname")
  ```

  > flag 中对 Duration 这种非基本类型的支持，使用的就是类似这样的方式。

### 解析

```go
func Parse() {
    // Ignore errors; CommandLine is set for ExitOnError.
    CommandLine.Parse(os.Args[1:])
}
```

> 语法有三种形式：
>
> - -flag // 只支持 bool 类型。
> - -flag=x。
> - -flag x // 只支持非 bool 类型。

- int 类型可以是十进制、十六进制、八进制甚至是负数；
- bool 类型可以是 1, 0, t, f, true, false, TRUE, FALSE, True, False。
- Duration 可以接受任何 time.ParseDuration 能解析的类型。

### 类型

#### ErrorHandling

```go
type ErrorHandling int

const (
  ContinueOnError ErrorHandling = iota
  ExitOnError
  PanicOnError
)
```

> 该类型定义了在参数解析出错时错误处理方式。

#### Flag

```go
// A Flag represents the state of a flag.
type Flag struct {
  Name     string // name as it appears on command line
  Usage    string // help message
  Value    Value  // value as set
  DefValue string // default value (as text); for usage message
}
```

> Flag 类型代表一个 flag 的状态。

#### FlagSet

```go
type FlagSet struct {
  // Usage is the function called when an error occurs while parsing flags.
  // The field is a function (not a method) that may be changed to point to
  // a custom error handler.
  Usage func()

  name string // FlagSet 的名字。CommandLine 给的是 os.Args[0]
  parsed bool // 是否执行过 Parse()
  actual map[string]*Flag // 存放实际传递了的参数（即命令行参数）
  formal map[string]*Flag // 存放所有已定义命令行参数
  args []string // arguments after flags // 开始存放所有参数，最后保留 非 flag（non-flag）参数
  exitOnError bool // does the program exit if there's an error?
  errorHandling ErrorHandling // 当解析出错时，处理错误的方式
  output io.Writer // nil means stderr; use out() accessor
}

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet
```

> 默认：var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

- **Parse**

  ```go
  func (f *FlagSet) Parse(arguments []string) error
  ```

  > 从参数列表中解析定义的 flag。方法参数 arguments 不包括命令名，即应该是 os.Args[1:]。

- **Arg、 Args、NArg、NFlag**

  ```go
  func (f *FlagSet) Arg(i int) string
  func (f *FlagSet) Args() []string
  func (f *FlagSet) NArg() int
  func (f *FlagSet) NFlag() int
  ```

  > - Arg(i int) 和 Args() 这两个方法就是获取 `non-flag` 参数的。
  >
  > - NArg() 获得 `non-flag` 的个数。
  >
  > - NFlag() 获得 FlagSet 中 actual 长度（即被设置了的参数个数）。

- **Visit、VisitAll**

  ```go
  func (f *FlagSet) Visit(fn func(*Flag))
  func (f *FlagSet) VisitAll(fn func(*Flag))
  ```

  > 分别用于访问 FlatSet 的 actual 和 formal 中的 Flag。

- **PrintDefaults**

  ```go
  func (f *FlagSet) PrintDefaults()
  ```

  > 打印所有已定义参数的默认值（调用 VisitAll 实现），默认输出到标准错误，除非指定了 FlagSet 的 output（通过 SetOutput() 设置）。

- **Set**

  ```go
  func (f *FlagSet) Set(name, value string) error
  ```

  > 设置某个 flag 的值（通过 name 查找到对应的 Flag）。

#### Value 接口

```go
// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface {
  String() string
  Set(string) error
}
```

> flag 包中，为 int、float、bool 等实现了该接口。借助该接口，我们可以自定义 flag。

## expvar - 公共变量的标准化接口

> 包 expvar 为公共变量提供了一个标准化的接口，如服务器中的操作计数器。
>
> - 以 JSON 格式通过 `/debug/vars` 接口以 HTTP 的方式公开这些公共变量。
> - 可以非常容易的展示应用程序指标。

```go
import _ "expvar"
```

## log - 日志记录

## runtime/debug - 运行时的调试工具

# 数据持久存储

## database/sql

### DB

> 是一个数据库句柄，代表一个具有零到多个底层连接的连接池，它可以安全的被多个 goroutine 同时使用。
>
> sql 包会自动创建和释放连接；它也会维护一个闲置连接的连接池。
>
> 如果数据库具有单连接状态的概念，该状态只有在事务中被观察时才可信。
>
> 一旦调用了 BD.Begin，返回的 Tx 会绑定到单个连接。
>
> 当调用事务 Tx 的 Commit 或 Rollback 后，该事务使用的连接会归还到 DB 的闲置连接池中。连接池的大小可以用 SetMaxIdleConns 方法控制。

- 实际的 Go 程序，应该在一个 go 文件中的 init 函数中调用 `sql.Open` 初始化全局的 sql.DB 对象，供程序中所有需要进行数据库操作的地方使用。

### 连接池

- **db.Ping()** 会将连接立马返回给连接池。
- **db.Exec()** 会将连接立马返回给连接池，但是它返回的 Result 对象会引用该连接，所以，之后可能会再次被使用。
- **db.Query()** 会传递连接给 sql.Rows 对象，直到完全遍历了所有的行或 Rows 的 Close 方法被调用了，连接才会返回给连接池。
- **db.QueryRow()** 会传递连接给 sql.Row 对象，当该对象的 Scan 方法被调用时，连接会返回给连接池。
- **db.Begin()** 会传递连接给 sql.Tx 对象，当该对象的 Commit 或 Rollback 方法被调用时，该链接会返回给连接池。

### 控制连接池

- **db.SetMaxOpenConns(n int)** 设置连接池中最多保存打开多少个数据库连接。注意，它包括在使用的和空闲的。如果某个方法调用需要一个连接，但连接池中没有空闲的可用，且打开的连接数达到了该方法设置的最大值，该方法调用将堵塞。默认限制是 0，表示最大打开数没有限制。
- **db.SetMaxIdleConns(n int)** 设置连接池中能够保持的最大空闲连接的数量。**默认值是 2**。
