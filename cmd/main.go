package main

import (
	"GoProgects/FileMan/cmd/fileSorter/interfaces"
	"fmt"
)

func main() {
	fileHandler := interfaces.HandlerFile{}
	var choice int
	// Выбор метода
	for {
		fmt.Println("Введите цифру соответствующую нужному действию:")
		fmt.Println(
			"1. Создать директории\n" +
				"2. Найти файл в директории\n" +
				"3. Найти файл в директории и открыть его\n" +
				"4. Подсчитать файлы разных типов в директории\n" +
				"5. Создать новый файл")

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch choice {
		case 1:
			interfaces.CreateDirectoryTree(fileHandler)
		case 2:
			interfaces.FindFileInDirectory(fileHandler)
		case 3:
			interfaces.OpenFindedFile(fileHandler)
		case 4:
			interfaces.FilesInfoInDir(fileHandler)
		case 5:
			interfaces.CreateNewFile(fileHandler)
		default:
			fmt.Println("Неверный выбор")
		}
		fmt.Println("Нажмите Enter для повтора...")
		_, err1 := fmt.Scanln()
		if err1 != nil {
			fmt.Scanln()
		}
	}
}
