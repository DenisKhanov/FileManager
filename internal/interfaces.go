package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileManager interface {
	CreateDirectoryTree(string, int) (string, error)
	FindFileInDirectory(string, string) (string, error)
	OpenFoundFile(string, string) ([]byte, error)
	FilesInfoInDir(string) (string, error)
	CreateNewFile(string) (string, error)
}

type HandlerFile struct{}

// CreateDirectoryTree creates a directory tree with the given number of directories.
func (h HandlerFile) CreateDirectoryTree(rootDirectory string, countPath int) (string, error) {
	var direct strings.Builder
	for i := 1; i < countPath+1; i++ {
		direct.WriteString(fmt.Sprintf("/%d", i))
	}

	var path strings.Builder
	for i := 1; i < countPath+1; i++ {
		path.WriteString(fmt.Sprintf("%s/%d%s", rootDirectory, i, direct.String()))
		err := os.MkdirAll(path.String(), 0600)
		if err != nil {
			return "", err
		}
		path.Reset()
	}
	return fmt.Sprintf("Дерево директорий успешно создано по пути: %s", rootDirectory), nil
}

// FindFileInDirectory finds a file in a directory and returns its absolute path or an error.
func (h HandlerFile) FindFileInDirectory(rootDirectory, targetFile string) (string, error) {
	fmt.Println("Идет поиск файла в директории...")
	var (
		foundPath string // Путь к найденному файлу файлу
	)
	err := filepath.WalkDir(rootDirectory, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			switch {
			case errors.Is(err, os.ErrPermission):
				return fmt.Errorf("ошибка доступа к %s\n", path)
			case errors.Is(err, os.ErrNotExist):
				return fmt.Errorf("директория %s не существует\n", path)
			case err.Error() == " CreateFile :: The filename, directory name, or volume label syntax is incorrect.":
				return fmt.Errorf("неверный формат директории %s\n", path)

			default:

				return errors.New("ошибка в процессе поиска файла в директории! Возможно вам стоит проверить правильность указания пути\n")

			}
		} else if info.IsDir() { //проверяем является ли путь директорией
			return nil
		} else if info.Name() == targetFile {
			foundPath = path
			return nil
		}
		return nil
	})
	if foundPath == "" && err == nil {
		return "", fmt.Errorf("файл \"%s\" не найден в директории %s", targetFile, rootDirectory)
	}
	return foundPath, err
}

// OpenFoundFile Starts the function FindFileInDirectory and reads the file according to the path found by it.
func (h HandlerFile) OpenFoundFile(rootDirectory string, nameFile string) ([]byte, error) {
	foundPath, err := h.FindFileInDirectory(rootDirectory, nameFile)
	if err != nil {
		err.Error()
		return nil, err

	}

	file, err := os.OpenFile(foundPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		fmt.Errorf("проблемы с чтением файла: %s", err).Error()
		return nil, fmt.Errorf("проблемы с чтением файла: %s", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	fileData := make([]byte, 0)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			// Достигнут конец файла, выходим из цикла
			break
		}
		if err != nil {
			fmt.Errorf("Ошибка чтения файла: %s\n", err).Error()
			return nil, err
		}
		fileData = append(fileData, buffer[:bytesRead]...)
	}

	fmt.Println("Файл который вы ищете был найден по данному пути ->->-> ", foundPath)
	fmt.Printf("Содержимое файла:\n\n%+v\n", string(fileData))
	fmt.Println("Вы хотите дописать текст в файл? (y/n): ")
	var input string
	fmt.Scanln(&input)
	fmt.Scanln()
	if input == "y" {
		fmt.Println("Введите ниже текст, который вы хотите добавить в файл:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			_, err = file.WriteString(scanner.Text())
			if err != nil {
				fmt.Errorf("Ошибка записи в файл: %s\n", err).Error()
				return nil, err
			} else {
				fmt.Printf("Текст успешно добавлен в файл: %s\n", foundPath)
			}
		}
	}

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			// Достигнут конец файла, выходим из цикла
			break
		}
		if err != nil {
			fmt.Errorf("Ошибка чтения файла: %s\n", err).Error()
			return nil, err
		}
		fileData = append(fileData, buffer[:bytesRead]...)
	}
	return fileData, err
}

// FilesInfoInDir Provides information about the number of files of types (Audio,Video, Image, Documents, Other) in the directory.
func (h HandlerFile) FilesInfoInDir(rootDirectory string) (string, error) {
	fmt.Println("Выполняется анализ...")
	var result strings.Builder
	result.WriteString(fmt.Sprintf("В директории %s найдены файлы следующих типов:\n\n", rootDirectory))
	CountTypes := make(map[string]map[string]int64)
	CountTypes["Image"] = make(map[string]int64)
	FileTypes := map[string]string{
		"Image": ".png.jpg.jpeg.gif.webp.bmp.tiff.ico.svg.psd.ai.eps.raw.indd.tif.tga.pbm.pgm.ppm.pam" +
			".ico.heif.heic.webp.jxr.hdp.wdp.hpcx.jpe.jif.jfif.jfi.arw.sr2.dcr.k25.kdc.erf.cr2.crw.bay.sr" +
			".srf.srw.x3f.raf.3fr.qtk.qpwr.ptx.pef.ope.orf.nr2.nrw.nik.dng.cine.dng.gpr.fff.nef.nrw.orf.raw" +
			".rw2.rwl.rwz.sr2.srf.srw.3fr.k25.kdc.dcr.cr2.crw.bay.hdr.kdc.erf.mos.pxn.rdc.dng.fff.nef.nrw.psd" +
			".x3f.rw2.rwl.rwz.sr2.srf.srw.3fr.k25.mdc.mos.pxn.fff.nef.nrw.psd.rw2.rwl.rwz.sr2.srf.srw.3fr.k25.kdc" +
			".erf.mos.px.ari.fff.nef.nrw.orf.psd.raw.rw2.rwl.rwz.sr2.srf.srw.3fr.k25.kdc.erf.cr2.crw.bay.hdr.kdc.mef" +
			".mos.pxn.sqf.fff.mos.nef.nrw.orf.ari.fff.nef.nrw.orf.tiff.tif",
		"Audio":    ".mp3.m4a.aac.wav.flac.ogg.wma.mid.midi.aiff.aif.aifc.au.snd.rmf.ra.rm.ram.ape.wv.mpc.oga.opus.ac3.aa3.mka.dts.caf.aa.aa3",
		"Video":    ".mp4.webm.avi.flv.mkv.mpeg.mpg.mov.wmv.asf.asx.asf.asf.3gp.3gpp.3g2.m4v.f4v.f4p.ogv.mts.m2ts.ts.vob.divx.xvid.m2v.m1v.mpv",
		"Document": ".txt.rtf.md.pdf.doc.docx.xls.xlsx.ppt.pptx.odt.ods.odp.pages.numbers.key.xml.html.htm.odp.csv",
		"Other": ".exe.zip.rar.7z.7z.exe.zip.rar.json.xml.yaml.sqlite.json.csv.txt.log.log4j.log4js.log4cxx.log4net.log4php.log4rb.log4php." +
			"log4net.xml.gitignore.dockerignore.htaccess.htpasswd.bashrc.profile.dbf.html.htm.xhtml.asp.aspx.php.jsp.js.jsx.c.cpp.h.cs.go.java." +
			"py.rb.php3.php4.php5.php7.sh.bat.ps1.css.scss.less.styl.js.map.ts.tsx.json5.wasm.class.jar.pom.xml.gradle.kts.gradle.properties.bat"}
	for key, _ := range FileTypes {
		CountTypes[key] = make(map[string]int64)
	}
	err := filepath.WalkDir(rootDirectory, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Что-то не так с родительской директорией %s", err)
		} else if info.IsDir() { //проверяем является ли путь директорией
			return nil
		}
		for k, v := range FileTypes {
			if strings.Contains(v, filepath.Ext(info.Name())) {
				CountTypes[k]["Количество"]++
				fileInfo, _ := info.Info()
				CountTypes[k]["Размер"] += fileInfo.Size()
			}
		}
		return err
	})
	if err != nil {
		return "", err
	}
	categoriesNames := []string{"Image", "Audio", "Video", "Document", "Other"}
	for _, v := range categoriesNames {
		result.WriteString(fmt.Sprintf("%s: %s, %dшт   %s, %d mb\n"+
			"---------------------------------------------------------\n",
			v, "Количество", CountTypes[v]["Количество"], "Размер", CountTypes[v]["Размер"]/(1024^2)))
	}

	return result.String(), nil
}

// CreateNewFile Creates a new file in the specified directory.
func (h HandlerFile) CreateNewFile(name string) (string, error) {
	//проверяем существует ли директория, если нет то создаём ее
	const dir = "c:/NewFiles"
	_, e := os.Stat(dir)
	if e != nil {
		fmt.Printf("Директория %s создана!\n", dir)
		er := os.Mkdir(dir, 0600)
		if er != nil {
			return "", er
		}
	}
	//создаем новый файл избавляясь от суффикса введенного пользователем и заменяем расширение на .txt
	newFile, err := os.Create("c:/NewFiles/" + strings.TrimSuffix(name, filepath.Ext(name)) + ".txt")
	if err != nil {
		return "", err
	}
	fmt.Println("Тут можно написать текст в файл или оставить его пустым нажав Enter")
	reader := bufio.NewReader(os.Stdin)
	text, err2 := reader.ReadBytes('\n')
	if err2 != nil {
		return "", err2
	}
	//записываем текст в файл
	_, err = newFile.Write(text)
	defer newFile.Close()
	msg := fmt.Sprintf("Файл %s успешно создан!\n", newFile.Name())
	return msg, nil
}

// CreateDirectoryTree creates a directory tree with the given number of directories.
func CreateDirectoryTree(fm FileManager) {
	var rootDirectory string
	var countPath int
	fmt.Println("Введите директорию в которой хотите создать новые директории, в формате С:/\n" +
		"(Обращаем внимание, что если директория начинается с символа /, то она будет создана в корневой директории диска,\n" +
		"если она не начинается с символа /, то будет создана в директории откуда вы произвели запуск программы)")
	_, err := fmt.Scanln(&rootDirectory)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Введите количество директорий и поддиректорий которое хотите " +
		"создать (число должно быть больше 0)!")
	for {
		_, err1 := fmt.Scanln(&countPath)
		if err1 != nil {
			fmt.Println("Неверный ввод. Пожалуйста, введите целое положительное число.")
			fmt.Scanln()
			continue // Перейти к следующей итерации цикла
		}
		if countPath <= 0 {
			fmt.Println("Введите положительное число больше нуля.")
			continue // Перейти к следующей итерации цикла
		}
		// Если данные верны, выходим из цикла
		break
	}

	result, err := fm.CreateDirectoryTree(rootDirectory, countPath)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

// FindFileInDirectory finds a file in a directory and returns its absolute path or an error.
func FindFileInDirectory(hf FileManager) {
	var rootDirectory string
	var targetFile string
	fmt.Println("Введите директорию в которой хотите искать в формате С:/")
	_, err := fmt.Scanln(&rootDirectory)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Введите название файла которые вы хотите найти!")
	_, err = fmt.Scanln(&targetFile)
	if err != nil {
		fmt.Println(err)
	}

	timeStart := time.Now()
	result, err := hf.FindFileInDirectory(rootDirectory, targetFile)
	if err != nil {
		fmt.Println(err)
		timeEnd := time.Now()
		fmt.Printf("Поиск выполнен за: %s\n", timeEnd.Sub(timeStart))
	} else {
		timeEnd := time.Now()
		fmt.Printf("Поиск выполнен за: %s\n", timeEnd.Sub(timeStart))
		fmt.Printf("Файл который вы искали находится в этой директории: %s\n", result)
	}
}

// OpenFindedFile Starts the function FindFileInDirectory and reads the file according to the path found by it.
func OpenFindedFile(hf FileManager) {
	var rootDirectory string
	var nameFile string
	fmt.Println("Введите директорию в которой хотите искать в формате С:/")
	_, err := fmt.Scanln(&rootDirectory)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Введите название файла которые вы хотите найти и открыть !")
	_, err = fmt.Scanln(&nameFile)
	if err != nil {
		fmt.Println(err)
	}

	hf.OpenFoundFile(rootDirectory, nameFile)
}

// FilesInfoInDir Provides information about the number of files of types (Audio,Video, Image, Documents, Other) in the directory.
func FilesInfoInDir(hf FileManager) {
	var rootDirectory string
	fmt.Println("Введите директорию в которой хотите искать в формате С:/")
	_, err := fmt.Scanln(&rootDirectory)
	if err != nil {
		fmt.Println(err)
	}

	result, err := hf.FilesInfoInDir(rootDirectory)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

// CreateNewFile Creates a new file in the specified directory.
func CreateNewFile(fm FileManager) {
	var nameFile string
	//var path string

	fmt.Println("Введите название будущего файла с желаемым расширением (файл будет создан в директории c:/NewFiles):")
	for {
		_, err1 := fmt.Scanln(&nameFile)
		if err1 != nil {
			fmt.Println("Ошибка ввода данных, название не может быть пустым!")
			fmt.Scanln()
			continue // Перейти к следующей итерации цикла
		}
		// Если данные верны, выходим из цикла
		break
	}
	// Создаем новый файл сообщение с результатом или ошибку
	msg, err := fm.CreateNewFile(nameFile)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(msg)
	}

	var input string
	fmt.Println("Вы хотите повторить? (y/n):")
	fmt.Scanln(&input)
	if input == "y" {
		CreateNewFile(fm)
	}
}
