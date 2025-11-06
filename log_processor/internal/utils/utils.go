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
    // fmt.Println("parsing-",line)
    if len(matches) == 0 {
        return nil, fmt.Errorf("invalid log format")
    }

    lineNumber, _ := strconv.Atoi(matches[4])
    t, _ := time.Parse("2006/01/02 15:04:05", matches[2])

    return &models.LogEntry{
        Level:      strings.ToUpper(matches[1]),
        Timestamp:  t,
        SourceFile: matches[3],
        LineNumber: lineNumber,
        Message:    matches[5],
    }, nil
}
