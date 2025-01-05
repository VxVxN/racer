package statisticer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	Name   string
	Points int
}

func NewRecord(name string, points int) Record {
	return Record{
		Name:   name,
		Points: points,
	}
}

type Statisticer struct {
}

func NewStatisticer() *Statisticer {
	return &Statisticer{}
}

func (s *Statisticer) Load() ([]Record, error) {
	file, err := os.Open("statistics.txt")
	if err != nil {
		return []Record{}, nil
	}
	defer file.Close()
	var records []Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, ",")
		if len(splitLine) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		points, err := strconv.Atoi(splitLine[1])
		if err != nil {
			return nil, err
		}
		records = append(records, Record{
			Name:   splitLine[0],
			Points: points,
		})
	}
	return records, nil
}

func (s *Statisticer) Save(records []Record) error {
	var data string
	for _, recotd := range records {
		data += fmt.Sprintf("%s,%d\n", recotd.Name, recotd.Points)
	}
	if err := os.WriteFile("statistics.txt", []byte(data), 0644); err != nil {
		return err
	}
	return nil
}
