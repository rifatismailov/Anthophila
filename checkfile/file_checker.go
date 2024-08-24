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
	"fmt"
	"time"
)

type FileChecker struct {
	Address     string
	Key         []byte
	Directories []string
}

func (w *FileChecker) Start() {

	checker := Checker{w.Address, w.Key, w.Directories}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			err := checker.Checkfile()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
		fmt.Println("FileChecker: All iterations finished")
	}()
}
