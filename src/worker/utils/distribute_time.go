package utils

import (
	"fmt"
	"strconv"
	"time"
)

const secondInDay = 86400

func DistributeTimeByteArgument(args []string) (startIn int64, endIn int64, err error) {
	if len(args) == 0 {
		return 0, 0, nil
	}
	if len(args) != 2 {
		return 0, 0, fmt.Errorf("Distribute info is not valid")
	}
	numDistribute, err0 := strconv.ParseInt(args[0], 10, 64)
	order, err1 := strconv.ParseInt(args[1], 10, 64)
	if err0 != nil || err1 != nil {
		return 0, 0, fmt.Errorf("Distribute info is not valid")
	}

	// get timestamp in day before now
	nowBeforeOneDay, err := time.Parse("2006-01-02", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	if err != nil {
		return 0, 0, err
	}
	nowBeforeOneDayTimeStamp := nowBeforeOneDay.Unix()

	// get distribute time in one day by distribute number
	secondByDistrubute := int64(secondInDay / numDistribute)

	startIn = nowBeforeOneDayTimeStamp + (secondByDistrubute * order) - secondByDistrubute
	endIn = nowBeforeOneDayTimeStamp + (secondByDistrubute * order)

	// check if first or last
	if order == 1 {
		startIn = 0
	}
	if order == numDistribute {
		endIn = 0
	}

	return startIn, endIn, nil
}
