/*
+----------------------------+
|       FileChecker          |
|                            |
| 1.      Сканування
|    +-------------------+   |
|    |                   |   |
|    | Почати сканування |   |
|    | thread (Scheduler)|   |
|    +--------+----------+   |
|             |              |
|             v              |
|     +-----------------+    |
|    |Сканування каталогу|   |
|    | і хеш-обчислення  |   |
|     +--------+--------+    |
|              |             |
|              v             |
| 2. Запис даних у файл JSON|
|     +-------------------+  |
|     |                   |  |
|     | JSON File Writer  |  |
|     +--------+----------+  |
|              |             |
|              v             |
| 3. Перевірте зміну хешу    |
|     +-------------------+  |
|     |                   |  |
|     | Порівняння хешів  |  |
|     |Надіслати на сервер|  |
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
	//вставити функцію яка буде робити затримку до окремого часу
	checker := Checker{fc.Address, fc.Key, fc.Directories, fc.SupportedExtensions, fc.InfoJson}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			checker.CheckFile()
		}
	}()
}
