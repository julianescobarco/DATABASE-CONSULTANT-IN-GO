package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Data struct {
	MessageID               string
	Date                    string
	From                    string
	To                      string
	Subject                 string
	MimeVersion             string
	ContentType             string
	ContentTransferEncoding string
	XFrom                   string
	XTo                     string
	Xcc                     string
	Xbcc                    string
	XFolder                 string
	XOrigin                 string
	XFileName               string
	Body                    string
}

func main() {
	filepath.WalkDir("./mailTest", walk)
	fmt.Println("Imprimiendo main prueba")
}

// Leer archivos
func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		leerFile(s)
	}
	return nil
}

//leerFile(folder + "/" + file.Name())

// Permtie leer el archivo
func leerFile(file_name string) {

	//fmt.Println(file_name)
	content, err := ioutil.ReadFile(file_name)

	if err != nil {
		log.Fatal(err)
	}
	data := &Data{}
	extraerData(content, data)
	//fmt.Println(*&data.Subject)
	fmt.Println(data)
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(out))
	send(string(out))
}

// https://onlineutf8tools.com/convert-bytes-to-utf8
func extraerData(content []byte, data *Data) {
	/** convertir arreglo de bytes en io.reader */
	reader := bytes.NewReader(content)
	fileScanner := bufio.NewScanner(reader)
	fileScanner.Split(bufio.ScanLines)

	aux := ""
	auxAnterior := ""
	/** banderas auxiliares para concatenar mensajes */
	MessageIDFlag := true
	DateFlag := true
	XFromFlag := true
	FromFlag := true
	XToFlag := true
	ToFlag := true
	SubjectFlag := true
	MimeVersionFlag := true
	ContentType := true
	ContentTransferEncodingFlag := true
	XFolderFlag := true
	XOriginFlag := true
	XccFlag := true
	XbccFlag := true
	XFileNameFlag := true
	bodyFlag := true

	for fileScanner.Scan() {

		aux = validarRenglonSinKey(fileScanner.Text(), aux)

		if strings.Contains(aux, "---- NoInformation -----") {
			aux = auxAnterior
		}

		if (strings.Contains(fileScanner.Text(), "Message-ID:") || strings.Contains(aux, "Message-ID:")) && bodyFlag == true {
			fmt.Println(auxAnterior)
			if MessageIDFlag == false {
				data.MessageID = data.MessageID + " " + fileScanner.Text()[0:]
				auxAnterior = "Message-ID:"
			} else {
				data.MessageID = fileScanner.Text()[11:]
				MessageIDFlag = false
				auxAnterior = "Message-ID:"
			}
		} else if (strings.Contains(fileScanner.Text(), "Date:") || strings.Contains(aux, "Date:")) && bodyFlag == true {
			if DateFlag == false {
				data.Date = data.Date + " " + fileScanner.Text()[0:]
				auxAnterior = "Date:"
			} else {
				data.Date = fileScanner.Text()[5:]
				DateFlag = false
				auxAnterior = "Date:"
			}
		} else if (strings.Contains(fileScanner.Text(), "X-From:") || strings.Contains(aux, "X-From:")) && bodyFlag == true {
			if XFromFlag == false {
				data.XFrom = data.XFrom + " " + fileScanner.Text()[0:]
				auxAnterior = "X-From:"
			} else {
				data.XFrom = fileScanner.Text()[7:]
				XFromFlag = false
				auxAnterior = "X-From:"
			}
		} else if (strings.Contains(fileScanner.Text(), "From:") || strings.Contains(aux, "From:")) && bodyFlag == true {
			if FromFlag == false {
				data.From = data.From + " " + fileScanner.Text()[0:]
				auxAnterior = "From:"
			} else {
				data.From = fileScanner.Text()[5:]
				FromFlag = false
				auxAnterior = "From:"
			}
		} else if (strings.Contains(fileScanner.Text(), "X-To:") || strings.Contains(aux, "XTo:")) && bodyFlag == true {
			if XToFlag == false {
				data.XTo = data.XTo + " " + fileScanner.Text()[0:]
				auxAnterior = "XTo:"
			} else {
				data.XTo = fileScanner.Text()[5:]
				XToFlag = true
				auxAnterior = "XTo:"
			}
		} else if (strings.Contains(fileScanner.Text(), "To:") || strings.Contains(aux, "To:")) && bodyFlag == true {

			if ToFlag == false {
				data.To = data.To + " " + fileScanner.Text()[0:]
				auxAnterior = "To:"
			} else {
				data.To = fileScanner.Text()[3:]
				ToFlag = false
				auxAnterior = "To:"
			}
		} else if (strings.Contains(fileScanner.Text(), "Subject:") || strings.Contains(aux, "Subject:")) && bodyFlag == true {
			if SubjectFlag == false {
				data.Subject = data.Subject + " " + fileScanner.Text()[0:]
				auxAnterior = "Subject:"
			} else {
				data.Subject = fileScanner.Text()[8:]
				SubjectFlag = false
				auxAnterior = "Subject:"
			}
		} else if (strings.Contains(fileScanner.Text(), "Mime-Version:") || strings.Contains(aux, "Mime-Version:")) && bodyFlag == true {
			if MimeVersionFlag == false {
				data.MimeVersion = data.MimeVersion + " " + fileScanner.Text()[0:]
				auxAnterior = "Mime-Version:"
			} else {
				data.MimeVersion = fileScanner.Text()[13:]
				MimeVersionFlag = false
				auxAnterior = "Mime-Version:"
			}

		} else if (strings.Contains(fileScanner.Text(), "Content-Type:") || strings.Contains(aux, "Content-Type:")) && bodyFlag == true {
			if ContentType == false {
				data.ContentType = data.ContentType + " " + fileScanner.Text()[0:]
				auxAnterior = "Content-Type:"
			} else {
				data.ContentType = fileScanner.Text()[13:]
				ContentType = false
				auxAnterior = "Content-Type:"
			}

		} else if (strings.Contains(fileScanner.Text(), "Content-Transfer-Encoding:") || strings.Contains(aux, "Content-Transfer-Encoding:")) && bodyFlag == true {
			if ContentTransferEncodingFlag == false {
				data.ContentTransferEncoding = data.ContentTransferEncoding + " " + fileScanner.Text()[0:]
				auxAnterior = "Content-Transfer-Encoding:"
			} else {
				data.ContentTransferEncoding = fileScanner.Text()[25:]
				ContentTransferEncodingFlag = false
				auxAnterior = "Content-Transfer-Encoding:"
			}

		} else if (strings.Contains(fileScanner.Text(), "X-Folder:") || strings.Contains(aux, "X-Folder:")) && bodyFlag == true {
			if XFolderFlag == false {
				data.XFolder = data.XFolder + " " + fileScanner.Text()[0:]
				auxAnterior = "X-Folder:"
			} else {
				data.XFolder = fileScanner.Text()[9:]
				XFolderFlag = false
				auxAnterior = "X-Folder:"
			}

		} else if (strings.Contains(fileScanner.Text(), "X-Origin:") || strings.Contains(aux, "X-Origin:")) && bodyFlag == true {
			if XOriginFlag == false {
				data.XOrigin = data.XOrigin + " " + fileScanner.Text()[0:]
				auxAnterior = "X-Origin:"
			} else {
				data.XOrigin = fileScanner.Text()[9:]
				XOriginFlag = false
				auxAnterior = "X-Origin:"
			}

		} else if (strings.Contains(fileScanner.Text(), "X-cc:") || strings.Contains(aux, "X-cc:")) && bodyFlag == true {
			if XccFlag == false {
				data.Xcc = data.Xcc + " " + fileScanner.Text()[0:]
				auxAnterior = "X-cc:"
			} else {
				data.Xcc = fileScanner.Text()[5:]
				XccFlag = false
				auxAnterior = "X-cc:"
			}

		} else if (strings.Contains(fileScanner.Text(), "X-bcc:") || strings.Contains(aux, "X-bcc:")) && bodyFlag == true {
			if XbccFlag == false {
				data.Xbcc = data.Xbcc + " " + fileScanner.Text()[0:]
				auxAnterior = "X-bcc:"
			} else {
				data.Xbcc = fileScanner.Text()[6:]
				XbccFlag = true
				auxAnterior = "X-bcc:"
			}

		} else if (strings.Contains(fileScanner.Text(), "X-FileName:") || strings.Contains(aux, "X-FileName:")) && bodyFlag == true {
			data.XFileName = fileScanner.Text()[11:]
			XFileNameFlag = false
			auxAnterior = "X-FileName:"
			bodyFlag = false

		} else if bodyFlag == false && XFileNameFlag == false {
			data.Body = data.Body + " " + fileScanner.Text()[0:]
		}

	}

}

/** Permite saber cuando la siguiente linea cuando es la contuniacion de la anterior */
func validarRenglonSinKey(fileScanner string, aux string) string {
	if strings.Contains(fileScanner, "Message-ID:") {
		aux = "Message-ID:"
	} else if strings.Contains(fileScanner, "Date:") {
		aux = "Date:"
	} else if strings.Contains(fileScanner, "X-From:") {
		aux = "X-From:"
	} else if strings.Contains(fileScanner, "From:") {
		aux = "From:"
	} else if strings.Contains(fileScanner, "X-To:") {
		aux = "X-To:"
	} else if strings.Contains(fileScanner, "To:") {
		aux = "To:"
	} else if strings.Contains(fileScanner, "Subject:") {

		aux = "Subject:"
	} else if strings.Contains(fileScanner, "Mime-Version:") {
		aux = "Mime-Version:"

	} else if strings.Contains(fileScanner, "Content-Type:") {
		aux = "Content-Type:"

	} else if strings.Contains(fileScanner, "Content-Transfer-Encoding:") {
		aux = "Content-Transfer-Encoding:"
	} else if strings.Contains(fileScanner, "X-Folder:") {
		aux = "X-Folder:"
	} else if strings.Contains(fileScanner, "X-Origin:") {
		aux = "X-Origin:"
	} else if strings.Contains(fileScanner, "X-cc:") {
		aux = "X-cc:"
	} else if strings.Contains(fileScanner, "X-bcc:") {
		aux = "X-bcc:"
	} else if strings.Contains(fileScanner, "X-FileName:") {
		aux = "X-FileName:"
	} else {

		aux = "---- NoInformation -----"
	}

	return aux
}

/** Funcion encargada de subir información a zincsearch */
func send(data string) {
	req, err := http.NewRequest("POST", "http://localhost:4080/api/mail_dir/_doc", strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
