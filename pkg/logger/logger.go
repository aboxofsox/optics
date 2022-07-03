package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogStore struct {
	LogItems []LogItem
}

type LogItem struct {
	Timestamp         string
	Endpoint          string
	StatusCode        int
	StatusCodeMessage string
	Elapsed           time.Duration
}

func New() *LogStore {
	return &LogStore{
		LogItems: []LogItem{},
	}
}

func (ls *LogStore) Stash(li *LogItem) { ls.LogItems = append(ls.LogItems, *li) }
func (ls *LogStore) First() LogItem    { return ls.LogItems[0] }
func (ls *LogStore) Last() LogItem     { return ls.LogItems[len(ls.LogItems)-1] }

func (ls *LogStore) FilterByStatusCode(statusCode int) []LogItem {
	lis := []LogItem{}

	for _, li := range ls.LogItems {
		if li.StatusCode == statusCode {
			lis = append(lis, li)
		}
	}

	return lis
}

func (ls *LogStore) Write(p string) {
	file, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0740)
	if err != nil {
		fmt.Printf("unable to open file: %s\n", err.Error())
		return
	}
	defer file.Close()

	for _, li := range ls.LogItems {
		sb := &strings.Builder{}
		s := fmt.Sprintf(
			"%s:\t%s %d %s %.2dms\n",
			li.Timestamp,
			li.Endpoint,
			li.StatusCode,
			strings.ToUpper(li.StatusCodeMessage),
			li.Elapsed.Milliseconds(),
		)
		sb.WriteString(s)
		file.WriteString(sb.String())
	}
}

func (ls *LogStore) Print(p string) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		fmt.Printf("log file does not exist: %s\n", p)
		return
	}

	file, err := os.Open(p)
	if err != nil {
		fmt.Printf("unable to open log file: %s\n", err.Error())
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("unable to read log file: %s\n", err.Error())
		return
	}

	fmt.Println(string(data))
}
