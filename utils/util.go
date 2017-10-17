/**
 * @author liangbo
 * @email  liangbogopher87@gmail.com
 * @date   2017/9/24 21:28 
 */
package utils

import (
    "reflect"
    "bytes"
    "fmt"
    "runtime"
    "io/ioutil"
    "strings"
)

var (
    dunno     = []byte("???")
    centerDot = []byte("·")
    dot       = []byte(".")
    slash     = []byte("/")
)

func IsStructPtr(t reflect.Type) bool {
    return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
    buf := new(bytes.Buffer) // the returned data
    // As we loop, we open files and read them. These variables record the currently
    // loaded file.
    var lines [][]byte
    var lastFile string
    for i := skip; ; i++ { // Skip the expected number of frames
        pc, file, line, ok := runtime.Caller(i)
        if !ok {
            break
        }
        // Print this much at least.  If we can't find the source, it won't show.
        fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
        if file != lastFile {
            data, err := ioutil.ReadFile(file)
            if err != nil {
                continue
            }
            lines = bytes.Split(data, []byte{'\n'})
            lastFile = file
        }
        fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
    }
    return buf.Bytes()
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return dunno
    }
    name := []byte(fn.Name())
    // The name includes the path name to the package, which is unnecessary
    // since the file name is already included.  Plus, it has center dots.
    // That is, we see
    //	runtime/debug.*T·ptrmethod
    // and want
    //	*T.ptrmethod
    // Also the package path might contains dot (e.g. code.google.com/...),
    // so first eliminate the path prefix
    if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
        name = name[lastslash+1:]
    }
    if period := bytes.Index(name, dot); period >= 0 {
        name = name[period+1:]
    }
    name = bytes.Replace(name, centerDot, dot, -1)
    return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
    n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
    if n < 0 || n >= len(lines) {
        return dunno
    }
    return bytes.TrimSpace(lines[n])
}

func GetBetweenStr(str, start, end string) string {
    n := strings.Index(str, start) + len(start)
    if n == -1 {
        n = 0
    }
    if n > len(str) {
        return ""
    }
    str = string([]byte(str)[n:])
    m := strings.Index(str, end)
    if m == -1 {
        m = len(str)
    }
    str = string([]byte(str)[:m])
    return str
}

//Be careful to use, from,to must be pointer
func DumpStruct(to interface{}, from interface{}) {
    fromv := reflect.ValueOf(from)
    tov := reflect.ValueOf(to)
    if fromv.Kind() != reflect.Ptr || tov.Kind() != reflect.Ptr {
        return
    }

    from_val := reflect.Indirect(fromv)
    to_val := reflect.Indirect(tov)

    for i := 0; i < from_val.Type().NumField(); i++ {
        fdi_from_val := from_val.Field(i)
        fd_name := from_val.Type().Field(i).Name
        fdi_to_val := to_val.FieldByName(fd_name)

        if fdi_to_val.IsValid() && fdi_from_val.Type() == fdi_to_val.Type() {
            fdi_to_val.Set(fdi_from_val)
        }
    }
}