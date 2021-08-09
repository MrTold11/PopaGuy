package main
import (
	"encoding/json"
	// "github.com/google/uuid" //
	"math/rand"
	"os"
)

func NewFileStorage() (fileStorage, error) {
	FileStorage := fileStorage{

	}
	return FileStorage, nil
}

func (s *InmemoryStorage) getRandomVoice() Voice{
	index := rand.Intn(len(s.VoiceMap))
	return s.VoiceMap[index]
}

func (s *InmemoryStorage) addVoice(v Voice) error {
	s.VoiceMap= append(s.VoiceMap, v)
	return nil
}

func (s *fileStorage) addVoice(v Voice) error {
	err := s.storage.addVoice(v)
	if err != nil {
		return err
	}
	return s.updateStorage()
}

func(s *fileStorage) getRandomVoice() Voice {
	return s.storage.getRandomVoice()
}

func (s *fileStorage) updateStorage() error {
	str := "HELP"
	// str := uuid.NewString() //
	bytes, err := json.Marshal(s.storage.VoiceMap[len(s.storage.VoiceMap) - 1])
	if err != nil {
		return err
	}
	er := os.WriteFile(str + ".json", bytes, 0)
	if er != nil  {
		return er
	}
	return nil
}




