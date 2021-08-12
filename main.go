package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)
var globalCounter = 0
type Voice struct {
	helper int64
}

type Anek struct {
	name string `json:"name"`
	text string `json:"string"`
}

type InmemoryStorage struct {
	VoiceMap []Voice
}

type fileStorage struct {
	fileNames []string
	storage InmemoryStorage
}

type InmemoryStorageAnek struct {
	AnekMap []Anek
}

type fileStorageAnek struct {
	storage InmemoryStorageAnek
}
type VoiceServer struct {
	storage fileStorage
}

type anekServer struct {
	storage fileStorageAnek
}
// unfinished, requires to get anek text to fucking work with the fucking popugay voice
func (s *anekServer) getRandomAnek(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method != "GET" {
		http.Error (
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	an, err :=s.storage.getRandomAnek()
	if err != nil {
		fmt.Println(err)
		http.Error (
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Println(an.text)
	w.WriteHeader(200)
	// here should be code for sending text to ML algorithm to popugay it
	w.Write([]byte (an.text) )
}

func (s *anekServer) getReklamaID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error (
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	str := string (data)
	massiv := strings.Split(str, "\n")
	stry := massiv[3]
	stry = stry[:len(stry) - 1]
	fileName := "C:/Users/prokh/GolandProjects/TelegramPopagayBot/" + stry
	fileUY := fileName + ".ogg"
	fmt.Println(fileUY)
	data, err = os.ReadFile(fileUY)
	if err != nil {
		http.Error (
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
			)
	}
	fileText := base64.StdEncoding.EncodeToString(data)
	w.Write([]byte (fileText))


}
//works as intended
func (s *anekServer) getReklama(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method)
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method != "POST" {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}
	m := make(map[int64]string)
	data, err := ioutil.ReadFile("reklMap.txt")
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println(err)
		fmt.Println("WHY")
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	count := 0

	data, err = ioutil.ReadFile("reklTimeMap.txt")
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	ma := make(map[int64]int64)
	err = json.Unmarshal(data, &ma)
	toWrite := ""
	for counter, name := range m {
		count++
		number := 24 * 3600 * ma[counter]
		copyName := name
		copyName = copyName[2:]
		id, err :=  strconv.ParseInt(copyName, 10, 64)
		if err != nil {
			http.Error (
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Println(id)
		if time.Now().UnixMilli() - id < number * 1000 {
			fmt.Println("HERE")

			fmt.Println(name)
			fmt.Println(toWrite)
			toWrite = toWrite + name
			toWrite += "$$"
		}
		fmt.Println(name)
	}
	if toWrite != "" {
		toWrite = toWrite[0:len(toWrite) - 2]
		fmt.Println(toWrite)
	} else {
		toWrite = "NON"
	}
	w.Write([]byte (toWrite))
	return


}
// still need to get help on py script
func (s * VoiceServer) handleVoice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method)
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method != "POST" {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	rawbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}





	er := os.WriteFile("sound.b64", rawbody, 0)
	if er != nil {
		fmt.Println(er)
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	} else {
		fmt.Println("SUCCESS B64 FILE UPLOADED")


		/* path, _ := os.Executable()
		fmt.Println(path)
		path, _ = filepath.EvalSymlinks(path)
		path = filepath.Join(filepath.Dir(path), "rt.py")
		fmt.Println(path)
		*/

		c1 := exec.Command("python",`./rt.py`)
		 if errr := c1.Run(); errr != nil {
			fmt.Println(errr)
		}
		fmt.Println("SUCCESS PYTHON HAS EXECUTED")
	}

	str, e := os.ReadFile("number.txt")
	if e != nil {
		fmt.Println(e)
		http.Error (
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	st := string (str)
	number, er := strconv.Atoi(st)
	if er != nil {
		http.Error (
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	i := 0
	if number == 0 {
		w.Write([]byte ("NULL"))
	}
	for i < number {
		fileName := fmt.Sprintf("sentence%d", i)
		data, err := ioutil.ReadFile(fileName + ".wav")
		if err != nil {
			fmt.Println(err)
			http.Error (
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		} else {
			fmt.Println(fileName)
			fileNameToStore := fmt.Sprintf("sentence%d", globalCounter)
			globalCounter = globalCounter + 1
			fmt.Println(fileNameToStore)
			err = os.WriteFile("toStoreAllPhrases/" + fileNameToStore + ".wav", data, 0644)
			if err != nil {
				fmt.Println("GO FUCK YOURSELF")
				http.Error (
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
				return
			}
			encoded := base64.StdEncoding.EncodeToString(data)
			err = os.WriteFile("testingForBot.b64",[]byte (encoded) , 0644)
			w.Write([]byte (encoded) )
		}
		i = i + 1
	}

	return
}
// looks like it works
func (s *VoiceServer) getVoice(w http.ResponseWriter, r *http.Request) {
	str, e := os.ReadFile("number.txt")
	if e != nil {
		fmt.Println(e)
		http.Error (
			w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
		)
		return
	}
	st := string (str)
	number, er := strconv.Atoi(st)
	if er != nil {
		http.Error (
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	i := 0
	for i <= number {
		fileName := fmt.Sprintf("sentence%d", i)
		fmt.Println(fileName)
		data, err := ioutil.ReadFile(fileName + ".wav")
		if err != nil {
			fmt.Println(err)
			http.Error (
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		} else {
			fmt.Println(fileName)
			encoded := base64.StdEncoding.EncodeToString(data)
			w.Write([]byte (encoded) )
		}
		i = i + 1
	}
}

func (s *VoiceServer) addVoice (w http.ResponseWriter, r *http.Request)  	{
	fmt.Println("got request with method", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
			)
		return
	}

	rawbody, err :=  ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
			)
		return
	}

	var newVoice Voice
	if err := json.Unmarshal(rawbody, &newVoice); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	if err := s.storage.addVoice(newVoice); err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(200)

}

func main() {
	fmt.Println(os.Executable())
	s := "HELP"
	fmt.Println(s)
	fileStorageNew, _ := NewFileStorage()

	voiceServerNew :=VoiceServer{
		storage : fileStorageNew,
	}

	fileStorageAnekNew, _  := NewFileStorageAnek()

	anekServerNew := anekServer {
		storage : fileStorageAnekNew,
	}
	mux := http.NewServeMux()
	mux.Handle("/add", http.HandlerFunc(voiceServerNew.addVoice))
	mux.Handle("/process", http.HandlerFunc(voiceServerNew.handleVoice))
	mux.Handle("/get", http.HandlerFunc(voiceServerNew.getVoice))
	mux.Handle("/getAnek", http.HandlerFunc(anekServerNew.getRandomAnek))
	mux.Handle("/listad", http.HandlerFunc(anekServerNew.getReklama))
	mux.Handle("/getad", http.HandlerFunc(anekServerNew.getReklamaID))
	http.ListenAndServe("192.168.234.64:8080", mux)
}