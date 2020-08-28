package log

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strconv"
  "sync"
  "time"
  "strings"
)

var fileIndex int
type RotateWriter struct {
  lock          sync.Mutex
  processName   string
  prefix        string
  dirName       string
  fileName      string
  fp            *os.File
  maxFileSize   int64
  dateStr       string
}

// Make a RootateWriter. Return nil if error occurs
func NewRotateWriter(dirName string, processName string, prefix string, maxFileSize int64 ) (*RotateWriter, error) {
  currFileIndex, err := getCurrFileIndex(dirName, processName)
  if err != nil {
    fmt.Errorf("Error while writing %s", err.Error())
    return nil, err
  }
  currDate := time.Now().Format("2006-01-02")
  var fileName string
  if currFileIndex == -1 {
      fileName = filepath.Join(dirName, processName + prefix + currDate + "." + strconv.Itoa(currFileIndex + 1) + ".log")
  } else {
      fileName = filepath.Join(dirName, processName + prefix + currDate + "." + strconv.Itoa(currFileIndex ) + ".log")
  }

  // check if the currentfile is already full
  currFileSize , err := getFileSize(fileName)
  if err != nil {
    return nil, err
  }

  if currFileSize >= maxFileSize {
      fileName = filepath.Join(dirName, processName + prefix + currDate + "." + strconv.Itoa(currFileIndex + 1) + ".log")
  }

  w := &RotateWriter{dirName: dirName, processName: processName, prefix: prefix, fileName: fileName, maxFileSize: maxFileSize, dateStr: currDate}
  logFileHandle , err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
  if err != nil {
    return nil, err
  } else {
    w.fp = logFileHandle
  }
  return w, nil
}



// Get file size in bytes
func getFileSize(fileName string)(int64, error) {
  var err error
  file , err := os.Stat(fileName)
  if err != nil {
    return -1, err
  }
  size := file.Size()
  return size, nil
}

func (w * RotateWriter) Rotate() (bool, error) {
  var err error
  file, err := os.Stat(w.fileName)
  if err != nil {
    return false, err
  }
  size := file.Size()
  currDate := time.Now().Format("2006-01-02")
  if size >= w.maxFileSize || currDate != w.dateStr {
    w.lock.Lock()
    defer w.lock.Unlock()

    // close existing file if open
    if w.fp != nil {
        err := w.fp.Close()
        w.fp = nil
        if err != nil {
          return false, err
        }
    }
    currFileIndex, err := getCurrFileIndex (w.dirName, w.processName)
    if err != nil {
      return false, err
    }
    currFileIndex = currFileIndex + 1
    fileName := filepath.Join(w.dirName, w.processName + w.prefix + currDate + "." + strconv.Itoa(currFileIndex)+ ".log")
    w.fileName = fileName
    w.dateStr  = currDate
    w.fp, err =  os.OpenFile(w.fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    if err != nil {
      fmt.Errorf("Error while rotating. New File: %s , Err: %s", w.fileName, err.ERROR())
      return false, err
    } else {
        return true, nil
    }
  } else {
      return true, nil
    }
}


func getCurrFileIndex(dirName string processName string) (int, error) {
  var err error
  files, err := ioutil.ReadDir(dirName)
  if err := nil {
    fmt.Errorf("Error while reading the directory: %s. Error: %s",dirName,err)
    return -1, err
  }
  maxVal := -1
  currDate := time.Now().Format("2006-01-02")
  for _, file := range files {
    fileName := file.Name()
    if strings.Contains(fileName, processName + "Log-" + currDate) {
      indexOfDot := strings.Index(fileName, ".")
      if indexOfDot != -1 {
        indexOfNextDot  := strings.Index(fileName[(indexOfDot +1):], ".")

        if indexOfNextDot != -1 {
          lastIndex := fileName[(indexOfDot + 1 ):(indexOfDot + indexOfNextDot + 1)]
          i1, err := strconv.Atoi(lastIndex)
          if err == nil {
            if i1 >maxVal {
              maxVal = i1
            }
          }
        }
      }
    }
  }
  return maxVal, nil
}


// write implements the io.Writer interface
func (w * RotateWriter) Write(output [] byte) (int, error){
  _, err := w.Rotate()
  if err != nil {
    fmt.Errorf("Error while writing %s", err)
  }
  w.lock.Lock()
  defer w.lock.Unlock()
  return w.fp.WriteString(time.Now().Format("2006-01-02 15:04:05.999999") + " " + string(output))
}
