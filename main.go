package main

import "Anthophila/management"

func main() {
	/*infoJson := information.NewInfo().InfoJson()
	directories := []string{"/home/sirius/GolandProjects/Anthophila/"}
	file_cheker := checkfile.FileChecker{
		Address:             "localhost:12345",
		Key:                 []byte("a very very very very secret key"),
		Directories:         directories,
		SupportedExtensions: []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"},
		TimeStart:           []int8{10, 25},
		InfoJson:            infoJson}
	file_cheker.Start()
	for {

		// Цей цикл буде виконуватися вічно
		fmt.Println("Main goroutine continues...")
		time.Sleep(time.Second) // Можна змінити затримку за потреби
	}
	*/
	manager := management.Manager{}
	manager.Start()
}
