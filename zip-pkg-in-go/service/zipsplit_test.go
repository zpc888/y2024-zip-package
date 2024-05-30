package service

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	tm := time.Now()
	// In Go's time package rounds down to the nearest second by default when formatting.
	t1 := tm.Format("060102150405000")
	fmt.Println()
	fmt.Println("year: ", tm.Year())
	fmt.Println("month: ", tm.Month())
	fmt.Println("day: ", tm.Day())
	fmt.Println("hour: ", tm.Hour())
	fmt.Println("minute: ", tm.Minute())
	fmt.Println("second: ", tm.Second())
	fmt.Println("millisecond: ", tm.Nanosecond()/1000000)
	t2 := fmt.Sprintf("%02d%02d%02d%02d%02d%02d%03d", tm.Year()%100, tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond()/1000000)
	fmt.Println("=============================")
	fmt.Println("t1: ", t1)
	fmt.Println("t2: ", t2)

}
