package domain

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ShortDate struct {
	time.Time
}

func NewShortDate(date time.Time) ShortDate {
	return Date(date.Year(), date.Month(), date.Day(), date.Location())
}

func Date(year int, month time.Month, day int, location *time.Location) ShortDate {
	return ShortDate{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (date *ShortDate) MarshalJSON() ([]byte, error) {
	if date.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(date.Format("2006-01-02"))
}

func (date *ShortDate) UnmarshalJSON(data []byte) error {
	unquotedData, _ := strconv.Unquote(string(data))
	time, err := time.Parse("2006-01-02", unquotedData)
	date.Time = time
	return err
}

func (date *ShortDate) String() string {
	return date.Format("2006-01-02")
}

func (s *ShortDate) MarshalText() ([]byte, error) {
	return []byte(s.Format("2006-01-02")), nil
}

func (s *ShortDate) UnmarshalText(text []byte) error {
	time, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return err
	}
	*s = ShortDate{time}
	return nil
}

type Timeframe struct {
	StartDate ShortDate
	EndDate   ShortDate
}

// TimeframeFromDate returns a Timeframe with the StartDate set to date and the EndDate set to today.
// The EndDate will use the same timezone location as provided in StartDate
func TimeframeFromDate(date ShortDate) Timeframe {
	endDate := NewShortDate(time.Now().In(date.Location()))
	return Timeframe{date, endDate}
}

func TimeframeFromQuery(params url.Values) (Timeframe, error) {
	from := params.Get("from")
	to := params.Get("to")
	if from == "" || to == "" {
		return Timeframe{}, fmt.Errorf("'from' and/or 'to' must be set")
	}
	startTime, err1 := time.Parse("20060102", from)
	startDate := ShortDate{startTime}
	endTime, err2 := time.Parse("20060102", to)
	endDate := ShortDate{endTime}
	if err1 != nil || err2 != nil {
		return Timeframe{}, fmt.Errorf("Malformed query params")
	}
	return Timeframe{StartDate: startDate, EndDate: endDate}, nil
}

func (tf *Timeframe) ToQuery() url.Values {
	params := make(url.Values)
	params.Set("from", tf.StartDate.Format("20060102"))
	params.Set("to", tf.EndDate.Format("20060102"))
	return params
}

func (tf *Timeframe) MarshalJSON() ([]byte, error) {
	if tf.StartDate.IsZero() || tf.EndDate.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(fmt.Sprintf("%s,%s", tf.StartDate.Format("2006-01-02"), tf.EndDate.Format("2006-01-02")))
}

func (tf *Timeframe) UnmarshalJSON(data []byte) error {
	unquotedData, _ := strconv.Unquote(string(data))
	dates := strings.Split(unquotedData, ",")
	if len(dates) != 2 {
		*tf = Timeframe{}
		return nil
	}
	startTime, err1 := time.Parse("2006-01-02", dates[0])
	startDate := ShortDate{startTime}
	endTime, err2 := time.Parse("2006-01-02", dates[1])
	endDate := ShortDate{endTime}
	if err1 != nil || err2 != nil {
		*tf = Timeframe{}
		return nil
	}
	*tf = Timeframe{StartDate: startDate, EndDate: endDate}
	return nil
}

func (tf *Timeframe) IsZero() bool {
	return tf.StartDate.IsZero() && tf.EndDate.IsZero()
}

func (tf *Timeframe) String() string {
	return fmt.Sprintf("{%s-%s}", tf.StartDate, tf.EndDate)
}
