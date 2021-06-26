package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func getFlagsAndEnvValues() (string, string, string, []string) {
	flag.Parse()
	// need to pass file as command args or env variable
	fileLocation := flag.Arg(0)
	if fileLocation == "" {
		fileLocation = os.Getenv("FILE")
	}
	// by default its queries on 10120 as mentioned in assignment
	queryPostcode := os.Getenv("QUERY_POSTCODE")
	if queryPostcode == "" {
		queryPostcode = "10120"
	}
	// by default its queries between 10AM - 3PM as mentioned
	queryTime := os.Getenv("QUERY_DELIVERY_TIME")
	if queryTime == "" {
		queryTime = "10AM - 3PM"
	}
	// by default match by name will have Potato/Veggie and Mushroom as mentioned
	queryMatchByName := os.Getenv("MATCH_BY_NAME")
	queryMatchByNameArgs := strings.Split(queryMatchByName, " ")
	if queryMatchByName == "" {
		queryMatchByNameArgs = []string{"Potato", "Veggie", "Mushroom"}
	}
	return fileLocation, queryPostcode, queryTime, queryMatchByNameArgs
}

func main() {
	fileLocation, queryPostcode, queryTime, queryMatchByName := getFlagsAndEnvValues()
	// Todo read in chunks and use streams
	// would use goroutines for further improvisation
	recipes, err := getRecipesFrom(fileLocation)
	if err != nil {
		exitGracefully(err)
	}
	// for processing delivery count per postcode per time
	_, timeInAM, timeInPM := extractDayAndTime(queryTime, false)
	deliveryCountPerPostcodeAndTime := CountPerPostcodeAndTime{
		Postcode:      queryPostcode,
		From:          timeInAM + "AM",
		To:            timeInPM + "PM",
		DeliveryCount: 0,
	}
	recipeNameVsFreqCount := map[string]int{}
	postcodeVsDeliveryCount := map[string]int{}
	busiestPostcode := BusiestPostcode{"", 0}
	// to sort the recipe names on alphabetic order
	distinctRecipeNames := []string{}

	for _, r := range recipes {
		//todo trim and lower recipe name
		var newRecipeCount int
		if currRecipeCount, ok := recipeNameVsFreqCount[r.Recipe]; !ok {
			newRecipeCount = 1
			distinctRecipeNames = append(distinctRecipeNames, r.Recipe)
		} else {
			newRecipeCount = currRecipeCount + 1
		}
		recipeNameVsFreqCount[r.Recipe] = newRecipeCount

		var newDeliveryCount int
		if currDeliveryCount, ok := postcodeVsDeliveryCount[r.Postcode]; !ok {
			newDeliveryCount = 1
		} else {
			newDeliveryCount = currDeliveryCount + 1
		}
		postcodeVsDeliveryCount[r.Postcode] = newDeliveryCount

		if newDeliveryCount > busiestPostcode.DeliveryCount {
			busiestPostcode.DeliveryCount = newDeliveryCount
			busiestPostcode.Postcode = r.Postcode
		}

		if queryPostcode == r.Postcode && isDeliveryInRangeOfQueryDelivery(r.Delivery, queryTime) {
			deliveryCountPerPostcodeAndTime.DeliveryCount = deliveryCountPerPostcodeAndTime.DeliveryCount + 1
		}
	}
	prepareJsonOutput(distinctRecipeNames, recipeNameVsFreqCount, busiestPostcode, deliveryCountPerPostcodeAndTime, queryMatchByName)
}

func prepareJsonOutput(distinctRecipeNames []string, recipeNameVsCountMap map[string]int, busiestPostcode BusiestPostcode, queryPostcodeAndTimeDeliveryCount CountPerPostcodeAndTime, matchByNameArgs []string) {
	// alphabetic order sorting
	sort.Strings(distinctRecipeNames)
	recipeMatchByName := []string{}
	recipeNameAndCount := []RecipeVsCount{}
	for _, recipeName := range distinctRecipeNames {
		recipeNameAndCount = append(recipeNameAndCount, RecipeVsCount{recipeName, recipeNameVsCountMap[recipeName]})
		if stringInSlice(recipeName, matchByNameArgs) {
			recipeMatchByName = append(recipeMatchByName, recipeName)
		}
	}

	// output
	recipeOutput := RecipeOutput{}
	recipeOutput.UniqueRecipeCount = len(distinctRecipeNames)
	recipeOutput.RecipeVsCount = recipeNameAndCount
	recipeOutput.BusiestPostcode = busiestPostcode
	recipeOutput.MatchByName = recipeMatchByName
	recipeOutput.CountPerPostcodeAndTime = queryPostcodeAndTimeDeliveryCount
	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(&recipeOutput)
	if err != nil {
		exitGracefully(err)
	}
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func isDeliveryInRangeOfQueryDelivery(delivery string, queryDeliveryTime string) bool {
	_, currDeliveryTimeAM, currDeliveryTimePM := extractDayAndTime(delivery, true)
	_, queryDeliveryTimeAM, queryDeliveryTimePM := extractDayAndTime(queryDeliveryTime, false)
	return isTimeRangeWithIn(queryDeliveryTimeAM, queryDeliveryTimePM, currDeliveryTimeAM, currDeliveryTimePM)
}

func isTimeRangeWithIn(queryTimeInAM, queryTimeInPM, currDeliveryInAM, currDeliveryInPM string) bool {
	t1, t2 := getTimeRange(queryTimeInAM, queryTimeInPM)
	t3, t4 := getTimeRange(currDeliveryInAM, currDeliveryInPM)

	// delivery times should be within query times
	// Hence below two conditons
	// t3 should be greater than equal to t1
	// t4 should be less than or equal to t2
	if t3.Before(t1) {
		return false
	}
	if t4.After(t2) {
		return false
	}
	return true
}

func getTimeRange(timeInAM, timeInPM string) (time.Time, time.Time) {
	layout1 := "3:04:05PM"
	t1, err1 := time.Parse(layout1, timeInAM+":00:00AM")
	if err1 != nil {
		exitGracefully(err1)
	}
	t2, err2 := time.Parse(layout1, timeInPM+":00:00PM")
	if err2 != nil {
		exitGracefully(err2)
	}
	return t1, t2
}

// this assumes the format of delivery as "{weekday} {h}AM - {h}PM"
// for query string it assumes format as {h}AM - {h}PM
func extractDayAndTime(delivery string, considerWeekday bool) (day string, amTime string, pmTime string) {
	delivery = strings.TrimSpace(delivery)
	if considerWeekday {
		firstSpaceIndex := strings.Index(delivery, " ")
		day = delivery[:firstSpaceIndex]
		delivery = delivery[firstSpaceIndex+1:]
	}
	amIndex := strings.Index(delivery, "AM - ")
	amTime = delivery[:amIndex]

	delivery = delivery[amIndex+len("AM - "):]
	pmIndex := strings.Index(delivery, "PM")
	pmTime = delivery[:pmIndex]
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func getRecipesFrom(fileLocation string) ([]RecipeInput, error) {
	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open Json File.")
	}
	defer file.Close()

	// Todo use streams/pipe
	dataInBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred in reading bytes from file.")
	}

	var recipes []RecipeInput

	err = json.Unmarshal(dataInBytes, &recipes)
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred in unmarshing json.")
	}

	return recipes, nil
}
