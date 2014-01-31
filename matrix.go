package main

import ( "fmt"; "os"; "strconv"; "crypto/rand"; "time"; "os/signal"; "syscall")

const speed = 70
const fastSpeed = 10
const lineTime = 200
const eraseTime = 100
const fastEraseTime = 4

func printOneChar(line uint, col uint, char byte) {
  fmt.Printf("\033[%d;%dH\033[32m%c", line, col, char)
}
func printTwoChar(line uint, col uint, char_1 byte, char_2 byte) {
  fmt.Printf("\033[%d;%dH\033[92m%c\033[%d;%dH\033[97m%c", line, col, char_2, line + 1, col, char_1)
}
func printEraseReady(col uint, lines uint, fast bool, ready chan bool) {
  printErase(col, lines, fast)
  ready <- true
}
func printErase(col uint, lines uint, fast bool) {
  var line uint = 1
  for line < lines {
    printOneChar(line, col, 32)
    line += 1
    if fast {
      time.Sleep(fastSpeed * time.Millisecond)
    } else {
      time.Sleep(speed * time.Millisecond)
    }
  }
  printOneChar(line, col, 32)
}
func printALineReady(col uint, lines uint, sto_p chan bool, ready chan bool) {
  printALine(col, lines, sto_p)
  ready <- true
}
func printALine(col uint, lines uint, sto_p chan bool) {
  var line uint = 1
  gen := getNextByte()
  c := gen()
  d := gen()
  for line < lines {
    select {
    case <- sto_p:
      return
    default:
      printTwoChar(line, col, c, d)
      c = d
      d = gen()
      line += 1
      time.Sleep(speed * time.Millisecond)
    }
  }
}
func generateRandomBytes(c int) []byte {
  b := make([]byte, c)
  _, err := rand.Read(b)
  if err != nil { panic(err); }
  return b
}
func byteIsNotPrintable(b byte) bool {
  if b < 33 { return true; }
  if b > 159 { return false; }
  if b > 126 { return true; }
  return false
}
func getNextByte() func() byte {
  const cbytes = 10
  var b []byte = generateRandomBytes(cbytes)
  var bp uint = 0
  var theByte uint8
  theByte = b[bp]
  return func() byte {
    for byteIsNotPrintable(theByte) {
      bp += 1
      if bp >= cbytes {
	bp = 0
	b = generateRandomBytes(cbytes)
      }
      theByte = b[bp]
    }
    bp += 1
    if bp >= cbytes {
      bp = 0
      b = generateRandomBytes(cbytes)
    }
    tb := theByte
    theByte = b[bp]
    return tb
  }
}
func getNextInt(max uint) func() uint {
  const cbytes = 10
  var b []byte
  var bp uint = cbytes
  return func() uint {
    bp += 1
    if bp >= cbytes-1 {
      bp = 0
      b = generateRandomBytes(cbytes)
    }
    return (uint(b[bp]) + uint(b[bp + 1])*256)%(max + 1)
  }
}
func running(geni func() uint, lines uint) {
  c := make(chan os.Signal, 1)
  sto_p := make(chan bool, 200)
  signal.Notify(c, os.Interrupt)
  for {
    select {
    case <-c:
      for i := 200; i > 0; {
	sto_p <- true
	i -= 1
      }
      return
    default:
      col := geni()
      go printALine(col, lines, sto_p)
      time.Sleep(lineTime * time.Millisecond)
      col = geni()
      go printErase(col, lines, false)
      time.Sleep(eraseTime * time.Millisecond)
      col = geni()
      go printErase(col, lines, false)
      time.Sleep(eraseTime * time.Millisecond)
    }
  }
}

func main() {
  lines_t := os.Getenv("lines_t")
  cols_t := os.Getenv("cols_t")
  lines64, err := strconv.ParseUint(lines_t, 10, 0)
  if err != nil { panic(err); }
  cols64, err := strconv.ParseUint(cols_t, 10, 0)
  if err != nil { panic(err); }
  lines := uint(lines64)
  cols := uint(cols64)
  fmt.Print("\0337\033[?47h\033[40m\033[97m\033[2J")
  fmt.Print("\033[?25l")
  ready := make(chan bool)
  geni := getNextInt(cols)
  running(geni, lines)
  for i := cols; i > 0; i--{
    go printEraseReady(i, lines, true, ready)
    time.Sleep(fastEraseTime * time.Millisecond)
  }
  for i := cols; i > 0; i--{
    _ = <-ready
  }
  fmt.Print("\033[?25h")
  fmt.Print("\033107m\03330m\033[?47l\0338")
  fmt.Printf("\033[%d;%dH", lines, 1)
  fmt.Printf("\033[30mNumber of lines %3d\n", lines)
  fmt.Printf("Number of cols  %3d\n", cols)
  ppid := os.Getppid()
  fmt.Printf("get_ppid %d\n", ppid)
  if ppid < 2 {
    ppid_t := os.Getenv("ppid_t")
    ppid, err := strconv.Atoi(ppid_t)
    if err != nil { panic(err); }
    fmt.Printf("env_ppid %d\n", ppid)
    err = syscall.Kill(ppid,2)
    if err != nil { panic(err); }
  } else {
    fmt.Printf("realppid %d\n", ppid)
    err = syscall.Kill(ppid,2)
    if err != nil { panic(err); }
  }
}
