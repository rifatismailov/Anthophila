/*
+----------------------------+
|       FileChecker          |
|                            |
| 1. Create scanning thread  |
|    +-------------------+   |
|    |                   |   |
|    | Start scanning    |   |
|    | thread (Scheduler)|   |
|    +--------+----------+   |
|             |              |
|             v              |
|     +-----------------+    |
|     | Directory Scan  |    |
|     | and Hash Compute |   |
|     +--------+--------+    |
|              |             |
|              v             |
| 2. Record data to JSON file|
|     +-------------------+  |
|     |                   |  |
|     | JSON File Writer  |  |
|     +--------+----------+  |
|              |             |
|              v             |
| 3. Check hash change       |
|     +-------------------+  |
|     |                   |  |
|     | Hash Comparison   |  |
|     | and Send to Server|  |
|     +-------------------+  |
|                            |
+----------------------------+
*/

package checkfile

import (
	"Anthophila/logging"
	"fmt"
	"time"
)

type FileChecker struct {
	Address             string
	Key                 []byte
	Directories         []string
	SupportedExtensions []string
	TimeStart           []int8
}

func (w *FileChecker) Start() {
	hour := w.TimeStart[0]
	minute := w.TimeStart[1]
	fmt.Println("Hour:", hour, "Minute:", minute)
	checker := Checker{w.Address, w.Key, w.Directories, w.SupportedExtensions}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			errorChecker := checker.Checkfile()
			if errorChecker != nil {
				logging.Now().PrintLog(
					"[FileChecker] Помилка:",
					errorChecker.Error())
			}
		}
	}()
}
