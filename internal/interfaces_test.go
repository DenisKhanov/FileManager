package internal

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerFile_CreateDirectoryTree(t *testing.T) {
	type args struct {
		rootDirectory string
		countPath     int
	}
	tests := []struct {
		name    string
		f       HandlerFile
		args    args
		want    string
		wantErr bool
	}{
		{name: "Test Positive 1",
			args: args{
				rootDirectory: "/tmp",
				countPath:     3,
			},
			want:    fmt.Sprintf("Дерево директорий успешно создано по пути: %s", "/tmp"),
			wantErr: false,
		},
		{name: "Test Negative-проблема с директорией",
			args: args{
				rootDirectory: ":",
				countPath:     3,
			},
			wantErr: true,
		},
		{name: "Test Negative-не верно указано количество путей",
			args: args{
				rootDirectory: ":",
				countPath:     3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		defer removeTempDirectory(tt.args.rootDirectory)
		t.Run(tt.name, func(t *testing.T) {
			f := HandlerFile{}
			got, err := f.CreateDirectoryTree(tt.args.rootDirectory, tt.args.countPath)
			assert.NoError(t, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerFile.CreateDirectoryTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HandlerFile.CreateDirectoryTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerFile_FindFileInDirectory(t *testing.T) {
	type args struct {
		rootDirectory string
		targetFile    string
	}
	tests := []struct {
		name    string
		f       HandlerFile
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := HandlerFile{}
			got, err := f.FindFileInDirectory(tt.args.rootDirectory, tt.args.targetFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerFile.FindFileInDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HandlerFile.FindFileInDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestHandlerFile_OpenFoundFile(t *testing.T) {
//	type args struct {
//		rootDirectory func()
//		nameFile      string
//	}
//	tests := []struct {
//		name string
//		f    HandlerFile
//		args args
//		want string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			f := HandlerFile{}
//			if got := f.OpenFoundFile(tt.args.rootDirectory, tt.args.nameFile); got != tt.want {
//				t.Errorf("HandlerFile.OpenFoundFile() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestHandlerFile_CountTypeFiles(t *testing.T) {
//	type args struct {
//		rootDirectory string
//	}
//	tests := []struct {
//		name    string
//		f       HandlerFile
//		args    args
//		want    string
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			f := HandlerFile{}
//			got, err := f.FilesInfoInDir(tt.args.rootDirectory)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("HandlerFile.FilesInfoInDir() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("HandlerFile.FilesInfoInDir() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestCreateDirectoryTree(t *testing.T) {
//	type args struct {
//		hf FileManager
//	}
//	tests := []struct {
//		name string
//		args args
//	}{{
//		name: "Test Positive 1",
//		args: args{
//			hf: HandlerFile{},
//		},
//	},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			CreateDirectoryTree(tt.args.hf)
//		})
//	}
//}
//
//func TestFindFileInDirectory(t *testing.T) {
//	type args struct {
//		hf FileManager
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			FindFileInDirectory(tt.args.hf)
//		})
//	}
//}
//
//func TestOpenFindedFile(t *testing.T) {
//	type args struct {
//		hf FileManager
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			OpenFindedFile(tt.args.hf)
//		})
//	}
//}
//
//func TestCountTypeFiles(t *testing.T) {
//	type args struct {
//		hf FileManager
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			FilesInfoInDir(tt.args.hf)
//		})
//	}
//}

// Вспомогательная функция для создания временной директории
func createTempDirectory(path string) string {
	os.Mkdir("/tmp/test", 0600)
	return "/tmp/test"
}

// Вспомогательная функция для удаления временной директории
func removeTempDirectory(dirPath string) {
	os.RemoveAll(dirPath)
}

// Вспомогательная функция для создания тестового файла
func createTestFile(dirPath, fileName string) string {
	data := "Просто какой-то текст"
	os.WriteFile(dirPath+"/"+fileName, []byte(data), 0600)
	return fileName
}
