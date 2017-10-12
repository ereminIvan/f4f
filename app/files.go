package app

import (
	"bufio"
	"errors"
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
	"strings"

	"github.com/ereminIvan/fffb/model"
)

//readDumps - read list of old fb messages id
func (a *application) readDumps() (map[string]struct{}, error) {
	log.Print("Start reading Facebook messages dump ...")
	result := make(map[string]struct{})
	file, err := os.Open(fbDumpFilePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			break
		}
		result[scanner.Text()] = struct{}{}
	}
	log.Printf("End read Facebook messages dump (%d messages)", len(result))
	return result, err
}


//writeFBDump - write fb messages dump to disk
func (a *application) writeFBDump(dump []string) error {
	log.Print("Write Facebook messages dump ...")
	if len(dump) == 0 {
		return nil
	}
	err := ioutil.WriteFile(fbDumpFilePath, []byte(strings.Join(dump, "\n")), 0644)
	log.Printf("Writed %d Facebook message to dump", len(dump))
	return err
}

func (a *application) readConfig() (model.Config, error) {
	log.Print("Read configs ...")
	file, err := ioutil.ReadFile(configFilePath)

	config := model.Config{}

	if err != nil {
		return config, errors.New("Error during read config.json: " + err.Error())
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, errors.New("Error during unmarshal config.json: " + err.Error() + " | " + string(file))
	}

	log.Printf("Starting Application with config: %+v", config)

	return config, nil
}
