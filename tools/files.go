package tools

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

/*
func FileVerify(url string) bool {

	response, err := http.Get(url)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		return true
	}
	return false
}

func FileDownload(url string, path string) error {

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Usa io.Copy para copiar los streams. No hay problema de tamano. (se puede usar para video o lo que sea)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
*/

var typesImagesAccepted = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// ImageDownload : Receive the url for get the image and the path of the folder target
func ImageDownload(url, path, name string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	extension := typesImagesAccepted[response.Header.Get("Content-Type")] // verificar la extension
	if extension == "" {
		return "", errors.New("Error: La imagen no tiene un formato aceptable: " + response.Header.Get("Content-Type")) //
	}

	// verificar el directorio
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			return "", err
		}
	}

	//open a file for writing
	file, err := os.Create(path + name + extension)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Usa io.Copy para copiar los streams. No hay problema de tamano. (se puede usar para video o lo que sea)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	return extension, nil
}

// MoveFile : copy the file fron source path to destiny dir, it is possible delete the file from the source dir
func MoveFile(sourcePath, names, destPath, named string, deleteold bool) error {

	inputFile, err := os.Open(sourcePath + names)
	if err != nil {
		return errors.New("Couldn't open source file: " + err.Error())
	}

	// verificar directorio destino
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err = os.MkdirAll(destPath, 0644)
		if err != nil {
			return err
		}
	}

	outputFile, err := os.Create(destPath + named)
	if err != nil {
		inputFile.Close()
		return errors.New("Couldn't open dest file: " + err.Error())
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return errors.New("Writing to output file failed: " + err.Error())
	}
	if deleteold {
		err = os.Remove(sourcePath + names)
		if err != nil {
			return errors.New("Failed removing original file: " + err.Error())
		}
	}
	return nil
}

// FilePutContents : create the file if not exist and put a new line
func FilePutContents(path, filename string, content string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModeDir)
		if err != nil {
			//fmt.Println("error on open dir: ", path)
			return err
		}
	}
	fp, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		//fmt.Println("error on open file: ", path+filename)
		return err
	}
	defer fp.Close()
	_, err = fp.WriteString(content + "\n")
	//fmt.Fprintln(fp, content)
	return err
}

func FileGetContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

// CreateDirectoryObject crea directorio objeto retornar prefix + "/" + dirnivel
func CreateDirectoryObject(clave, nivel int, prefix string) (string, error) {

	switch nivel {
	case 1:
		return prefix + "/" + fmt.Sprint(clave) + "/", nil
	case 2:
		sclave := fmt.Sprintf("%02d", clave)
		subd := sclave[0:2] + "/" + fmt.Sprint(clave)
		return prefix + "/" + subd + "/", nil
	case 3:
		return prefix + "/" + GetMD5Hash("kl123kl"+fmt.Sprint(clave)) + "/", nil
	}
	return "", errors.New("error nivel")
}

func IsFile(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil && !f.IsDir() {
		return true, nil
	}
	return false, err
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil && f.IsDir() {
		return true, nil
	}
	return false, err
}

// VerifyImage
// first parameter returns the list of messages of things done, second parameter returns changed name of file if needed, or empty
func VerifyImage(path, filename string, expectedname string, status bool) ([]string, string) {
	result := []string{}
	newname := ""
	ext := filepath.Ext(filename)
	if expectedname != "" && filename[:len(filename)-len(ext)] != expectedname {
		// file error, try to find file and rename file if possible

	}
	return result, newname
}
