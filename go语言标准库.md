# 输入输出

## io - 基本的 IO 接口

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

## ioutil - 方便的IO操作函数集

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

## fmt 格式化IO

### Printing

#### 占位符

**普通占位符**

| 占位符 | 说明                                                         | 举例                                              | 输出                                |
| ------ | ------------------------------------------------------------ | ------------------------------------------------- | ----------------------------------- |
| %v     | 相应值的默认格式。<br>**在打印结构体时，“加号”标记（%+v）会添加字段名** | **Printf("%v", site)**<br>**Printf("%+v", site)** | {studygolang}<br>{Name:studygolang} |
| %#v    | 相应值的Go语法表示                                           | Printf("%#v", site)                               | main.Website{Name:"studygolang"}    |
| %T     | 相应值的类型的Go语法表示                                     | Printf("%T", site)                                | main.Website                        |
| %%     | 字面上的百分号，并非值的占位符                               | Printf("%%")                                      | %                                   |

**布尔占位符**

| 占位符 | 说明                 | 举例               | 输出 |
| ------ | -------------------- | ------------------ | ---- |
| %t     | 单词 true 或 false。 | Printf("%t", true) | true |

**整数占位符**

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

**浮点数和复数的组成部分（实部和虚部）**

| 占位符 | 说明                                                         | 举例                   | 输出         |
| ------ | ------------------------------------------------------------ | ---------------------- | ------------ |
| %b     | 无小数部分的，指数为二的幂的科学计数法，与 strconv.FormatFloat 的 'b' 转换格式一致。例如 -123456p-78 |                        |              |
| %e     | 科学计数法，例如 -1234.456e+78                               | Printf("%e", 10.2)     | 1.020000e+01 |
| %E     | 科学计数法，例如 -1234.456E+78                               | Printf("%E", 10.2)     | 1.020000E+01 |
| %f     | 有小数点而无指数，例如 123.456                               | Printf("%f", 10.2)     | 10.200000    |
| %g     | 根据情况选择 %e 或 %f 以产生更紧凑的（无末尾的0）输出        | Printf("%g", 10.20)    | 10.2         |
| %G     | 根据情况选择 %E 或 %f 以产生更紧凑的（无末尾的0）输出        | Printf("%G", 10.20+2i) | (10.2+2i)    |

**字符串与字节切片**

| 占位符 | 说明                                   | 举例                                 | 输出           |
| ------ | -------------------------------------- | ------------------------------------ | -------------- |
| %s     | 输出字符串表示（string类型或[]byte)    | Printf("%s", []byte("Go语言中文网")) | Go语言中文网   |
| %q     | 双引号围绕的字符串，由Go语法安全地转义 | Printf("%q", "Go语言中文网")         | "Go语言中文网" |
| %x     | 十六进制，小写字母，每字节两个字符     | Printf("%x", "golang")               | 676f6c616e67   |
| %X     | 十六进制，大写字母，每字节两个字符     | Printf("%X", "golang")               | 676F6C616E67   |

**指针**

| 占位符 | 说明                  | 举例                | 输出     |
| ------ | --------------------- | ------------------- | -------- |
| %p     | 十六进制表示，前缀 0x | Printf("%p", &site) | 0x4f57f0 |

## bufio

