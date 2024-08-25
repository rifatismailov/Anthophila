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
	InfoJson            string
}

func (fc *FileChecker) Start() {
	hour := fc.TimeStart[0]
	minute := fc.TimeStart[1]
	fmt.Println("Hour:", hour, "Minute:", minute)

	checker := Checker{fc.Address, fc.Key, fc.Directories, fc.SupportedExtensions, fc.InfoJson}
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
