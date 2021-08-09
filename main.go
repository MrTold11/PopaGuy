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
)

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
// still need to get help on py script
func (s * VoiceServer) handleVoice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request with method", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
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
	for i <= number {
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
			encoded := base64.StdEncoding.EncodeToString(data)
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
	http.ListenAndServe("192.168.234.16:8080", mux)

}

