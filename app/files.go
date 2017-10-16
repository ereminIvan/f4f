package app

import (
	"bufio"
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ereminIvan/fffb/model"
)

//readDumps - read list of old fb messages id
func (a *application) readDumps() (map[string]struct{}, map[int64]struct{}, error) {
	fbDump, err := a.readFBDump()
	tgDump, err := a.readTGDump()
	return fbDump, tgDump, err
}

func (a *application) readFBDump() (map[string]struct{}, error) {
	log.Print("Start reading Facebook messages dump ...")
	fbResult := make(map[string]struct{})
	file, err := os.Open(fbDumpFilePath)
	if err != nil {
		return fbResult, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			break
		}
		fbResult[scanner.Text()] = struct{}{}
	}
	log.Printf("End read Facebook messages dump (%d messages)", len(fbResult))
	return fbResult, err
}

func (a *application) readTGDump() (map[int64]struct{}, error) {
	log.Print("Start reading Telegram messages dump ...")
	result := make(map[int64]struct{})
	file, err := os.Open(tgDumpFilePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = scanner.Err(); err != nil {
			break
		}
		id, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		result[id] = struct{}{}
	}
	log.Printf("End read Telegram chats dump (%d subscribers)", len(result))
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

//writeTGDump - write tg chats dump to disk
func (a *application) writeTGDump(dump map[int64]struct{}) error {
	log.Print("Write Telegram chats dump ...")
	if len(dump) == 0 {
		return nil
	}
	dumpLines := make([]string, 0, len(dump))
	for k := range dump {
		dumpLines = append(dumpLines, fmt.Sprintf("%d", k))
	}
	err := ioutil.WriteFile(tgDumpFilePath, []byte(strings.Join(dumpLines, "\n")), 0644)
	log.Printf("Writed %d Telegram chats to dump", len(dump))
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
