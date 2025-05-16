//go:build windows
// +build windows
// create by iwlb@outlook.com at 20250516
package mondir

import (
	"encoding/binary"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
	"golang.org/x/sys/windows"
)

func MonDir(dir s, dof func(string)) {
	dir = ToAbsolutePath(dir)

	dirhand, dirhander := windows.CreateFile(windows.StringToUTF16Ptr(dir), windows.FILE_LIST_DIRECTORY, windows.FILE_SHARE_WRITE|windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if dirhander != nil {
		return
	}
	dirchangeevent, er := windows.FindFirstChangeNotification(dir, true, windows.FILE_NOTIFY_CHANGE_CREATION|windows.FILE_NOTIFY_CHANGE_FILE_NAME|windows.FILE_NOTIFY_CHANGE_SIZE|windows.FILE_NOTIFY_CHANGE_LAST_WRITE|windows.FILE_NOTIFY_CHANGE_ATTRIBUTES|windows.FILE_NOTIFY_CHANGE_DIR_NAME)
	if er != nil {
		return
	}
	rdbuf := make([]byte, 1024*1024)

	pas := make(map[s]int64, 0)
	var pasmu sync.Mutex
	go func() {
		for {
			tms := nms()
			pasmu.Lock()
			for pa, ms := range pas {
				if tms-ms > 40 {
					go dof(pa)
					delete(pas, pa)
				}
			}
			pasmu.Unlock()
			time.Sleep(20 * time.Millisecond)
		}
	}()

	for true {
		var rdbuflen uint32
		for i := 0; i < 1024*1024; i += 1 {
			rdbuf[i] = 0
		}
		rder := windows.ReadDirectoryChanges(
			dirhand,
			&rdbuf[0],
			1024*1024,
			true,
			windows.FILE_NOTIFY_CHANGE_CREATION|windows.FILE_NOTIFY_CHANGE_FILE_NAME|windows.FILE_NOTIFY_CHANGE_SIZE|windows.FILE_NOTIFY_CHANGE_LAST_WRITE|windows.FILE_NOTIFY_CHANGE_ATTRIBUTES|windows.FILE_NOTIFY_CHANGE_DIR_NAME,
			&rdbuflen,
			nil,
			0)
		if rder == nil {
			prdbuf := rdbuf
			var preno, preno1 byte
			for true {
				if preno != 0 {
					prdbuf[0] = preno
				}
				if preno1 != 0 {
					prdbuf[0+1] = preno1
				}
				var nextoff uint32
				nextoff = binary.LittleEndian.Uint32(prdbuf[:4])
				//cpathlen := binary.LittleEndian.Uint32(prdbuf[4:8])
				preno = prdbuf[nextoff]
				preno1 = prdbuf[nextoff+1]
				prdbuf[nextoff] = 0
				prdbuf[nextoff+1] = 0
				subpath := windows.UTF16PtrToString((*uint16)(unsafe.Pointer(&prdbuf[12])))
				if len(subpath) > 4 && subpath[len(subpath)-4:] == ".TMP" && (strings.Index(subpath, ".cpp~") != -1 && strings.Index(subpath, ".cpp~") > strings.Index(subpath, "\\") || strings.Index(subpath, ".cc~") != -1 && strings.Index(subpath, ".cc~") > strings.Index(subpath, "\\") || strings.Index(subpath, ".h~") != -1 && strings.Index(subpath, ".h~") > strings.Index(subpath, "\\") || strings.Index(subpath, ".hpp~") != -1 && strings.Index(subpath, ".hpp~") > strings.Index(subpath, "\\") || strings.Index(subpath, ".c~") != -1 && strings.Index(subpath, ".c~") > strings.Index(subpath, "\\")) {
					subpath = subpath[0:strings.Index(subpath, "~")]
				} else if len(subpath) > 12 && subpath[len(subpath)-11] == '.' && subpath[len(subpath)-10] == 'c' && subpath[len(subpath)-9] == 'p' && subpath[len(subpath)-8] == 'p' && subpath[len(subpath)-7] == '.' && subpath[len(subpath)-4] == '3' && subpath[len(subpath)-3] == '4' && subpath[len(subpath)-2] == '4' && subpath[len(subpath)-1] == '4' {
					subpath = subpath[0:strings.Index(subpath, ".")]
				} else if len(subpath) > 10 && subpath[len(subpath)-9] == '.' && subpath[len(subpath)-8] == 'h' && subpath[len(subpath)-7] == '.' && subpath[len(subpath)-4] == '3' && subpath[len(subpath)-3] == '4' && subpath[len(subpath)-2] == '4' && subpath[len(subpath)-1] == '4' {
					subpath = subpath[0:strings.Index(subpath, ".")]
				} else if len(subpath) > 12 && subpath[len(subpath)-11] == '.' && subpath[len(subpath)-10] == 'p' && subpath[len(subpath)-9] == 'r' && subpath[len(subpath)-8] == 'o' && subpath[len(subpath)-7] == '.' && subpath[len(subpath)-4] == '3' && subpath[len(subpath)-3] == '4' && subpath[len(subpath)-2] == '4' && subpath[len(subpath)-1] == '4' {
					subpath = subpath[0:strings.Index(subpath, ".")]
				} else if len(subpath) > 12 && subpath[len(subpath)-11] == '.' && subpath[len(subpath)-10] == 'q' && subpath[len(subpath)-9] == 'r' && subpath[len(subpath)-8] == 'c' && subpath[len(subpath)-7] == '.' && subpath[len(subpath)-4] == '3' && subpath[len(subpath)-3] == '4' && subpath[len(subpath)-2] == '4' && subpath[len(subpath)-1] == '4' {
					subpath = subpath[0:strings.Index(subpath, ".")]
				} else if len(subpath) > 10 && subpath[len(subpath)-9] == '.' && subpath[len(subpath)-9] == 'u' && subpath[len(subpath)-8] == 'i' && subpath[len(subpath)-7] == '.' && subpath[len(subpath)-4] == '3' && subpath[len(subpath)-3] == '4' && subpath[len(subpath)-2] == '4' && subpath[len(subpath)-1] == '4' {
					subpath = subpath[0:strings.Index(subpath, ".")]
				}
				akp := StdPath(dir + subpath)
				st, ste := os.Stat(akp)
				if ste == nil && st.IsDir() == false && strings.HasSuffix(subpath, ".exe~") == false {
					//fmt.Println(subpath)
					//dof(subpath)
					pasmu.Lock()
					pas[subpath] = nms()
					pasmu.Unlock()

				}

				if nextoff != 0 {
					prdbuf = prdbuf[nextoff:]
				} else {
					break
				}
			}
		}
		if windows.FindNextChangeNotification(dirchangeevent) != nil {
			break
		}
	}
	windows.FindCloseChangeNotification(dirchangeevent)
	windows.CloseHandle(dirhand)
}
