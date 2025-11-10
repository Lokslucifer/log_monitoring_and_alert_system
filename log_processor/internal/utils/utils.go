package utils
import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "time"
	"log_processor/internal/models"

)

var logPattern = regexp.MustCompile(`^(?P<level>\w+): (?P<datetime>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (?P<file>\S+):(?P<line>\d+): (?P<message>.+)\n$`)

func ParseLogLine(line string) (*models.LogEntry, error) {
    matches := logPattern.FindStringSubmatch(line)

    if len(matches) == 0 {
        return nil, fmt.Errorf("invalid log format")
    }

    lineNumber, _ := strconv.Atoi(matches[4])
    t, err := time.Parse("2006/01/02 15:04:05", matches[2])
    if(err!=nil){
        return nil,fmt.Errorf("invalid time format")
    }

    return &models.LogEntry{
        Level:      strings.ToUpper(matches[1]),
        Timestamp:  t,
        SourceFile: matches[3],
        LineNumber: lineNumber,
        Message:    matches[5],
    }, nil
}

func ParseTimeString(timeStr string)(time.Time,error){
    if(timeStr==""){
        return time.Time{},nil
    }
    timeVal, err := time.Parse(time.RFC3339, timeStr)
    if err != nil {
        return time.Time{},err
    }
    return timeVal,nil
}

func ParseLevels(levelsParam string)([]string){
    var levels []string
    if(levelsParam!=""){
    for _, lvl := range strings.Split(levelsParam, ",") {
        levels = append(levels, strings.ToUpper(strings.TrimSpace(lvl)))
    }
}
    return levels

}

func ParseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}
