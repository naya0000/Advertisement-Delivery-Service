package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/naya0000/Advertisement_Manage.git/pkg/api"
	"github.com/naya0000/Advertisement_Manage.git/pkg/db"
)

// Advertisement represents an advertisement object
// type Advertisement struct {
// 	Title      string     `json:"title"`
// 	StartAt    time.Time  `json:"startAt"`
// 	EndAt      time.Time  `json:"endAt"`
// 	Conditions Conditions `json:"conditions"`
// }

// type Gender string

// const (
// 	Male   Gender = "M"
// 	Female Gender = "F"
// )

// // Conditions represents conditions for displaying an advertisement
// type Conditions struct {
// 	AgeStart *AgeStartCondition `json:"ageStart,omitempty"`
// 	AgeEnd   *AgeEndCondition   `json:"ageEnd,omitempty"`
// 	Gender   *GenderCondition   `json:"gender,omitempty"`
// 	Country  *CountryCondition  `json:"country,omitempty"`
// 	Platform *PlatformCondition `json:"platform,omitempty"`
// }

// // AgeCondition represents age condition
// type AgeStartCondition uint8
// type AgeEndCondition uint8

// // GenderCondition represents gender condition
// type GenderCondition string

// // CountryCondition represents country condition
// type CountryCondition []string

// // PlatformCondition represents platform condition
// type PlatformCondition []string

// // AdvertisementStore stores advertisements
// var AdvertisementStore []*Advertisement

// // DB 是數據庫連接的全局變量
// var DB *sql.DB

// // ListAdvertisementHandler handles listing advertisements based on conditions
// func ListAdvertisementHandler(w http.ResponseWriter, r *http.Request) {

// 	age := r.URL.Query().Get("age")
// 	gender := r.URL.Query().Get("gender")
// 	country := r.URL.Query().Get("country")
// 	platform := r.URL.Query().Get("platform")

// 	fmt.Println("Age:", age)
// 	fmt.Println("Gender:", gender)
// 	fmt.Println("Country:", country)
// 	fmt.Println("Platform:", platform)

// 	var matchedAds []*Advertisement

// 	for _, ad := range AdvertisementStore {
// 		if time.Now().After(ad.StartAt) && time.Now().Before(ad.EndAt) {

// 			fmt.Println("In time")

// 			if ad.Conditions.AgeStart != nil && !isAgeStartMatched(age, ad.Conditions.AgeStart) {
// 				continue
// 			}
// 			if ad.Conditions.AgeEnd != nil && !isAgeEndMatched(age, ad.Conditions.AgeEnd) {
// 				continue
// 			}

// 			if ad.Conditions.Gender != nil {
// 				fmt.Println("string(*ad.Conditions.Gender):", string(*ad.Conditions.Gender))
// 			}
// 			fmt.Println("gender:", gender)

// 			if ad.Conditions.Gender != nil && string(*ad.Conditions.Gender) != gender {
// 				continue
// 			}

// 			if ad.Conditions.Country != nil && !isCountryMatched(country, ad.Conditions.Country) {
// 				continue
// 			}
// 			if ad.Conditions.Platform != nil && !isPlatformMatched(platform, ad.Conditions.Platform) {
// 				continue
// 			}
// 			fmt.Println("ad:", ad)

// 			matchedAds = append(matchedAds, ad)

// 		}

// 	}

// 	sort.Slice(matchedAds, func(i, j int) bool {
// 		return matchedAds[i].EndAt.Before(matchedAds[j].EndAt)
// 	})

// 	limit := getQueryInt(r.URL.Query(), "limit", 5)
// 	offset := getQueryInt(r.URL.Query(), "offset", 1)

// 	fmt.Println("limit:", limit)
// 	fmt.Println("offset:", offset)

// 	startIndex := (offset - 1)
// 	endIndex := startIndex + limit
// 	if startIndex > len(matchedAds) {
// 		startIndex = len(matchedAds) - 1
// 	}
// 	if endIndex > len(matchedAds) {
// 		endIndex = len(matchedAds)
// 	}
// 	fmt.Println("startIndex:", startIndex)
// 	fmt.Println("endIndex:", endIndex)

// 	response := struct {
// 		Items []*Advertisement `json:"items"`
// 	}{
// 		Items: matchedAds[startIndex:endIndex],
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
// func isAgeStartMatched(ageStart string, condition *AgeStartCondition) bool {
// 	fmt.Println("ageStart:", ageStart)
// 	if ageStart == "" {
// 		return true
// 	}
// 	age, err := strconv.Atoi(ageStart)
// 	fmt.Println("age:", age, ", err:", err)

// 	if err != nil || age < 1 || age > 100 {
// 		return false
// 	}

// 	return age >= int(*condition)
// }

// func isAgeEndMatched(ageEnd string, condition *AgeEndCondition) bool {
// 	fmt.Println("ageEnd:", ageEnd)
// 	if ageEnd == "" {
// 		return true
// 	}
// 	age, err := strconv.Atoi(ageEnd)
// 	fmt.Println("age:", age, ", err:", err)

// 	if err != nil || age < 1 || age > 100 {
// 		return false
// 	}

// 	return age <= int(*condition)
// }

// func isCountryMatched(countryStr string, condition *CountryCondition) bool {
// 	if countryStr == "" {
// 		return true
// 	}
// 	for _, c := range *condition {
// 		fmt.Println("c:", c)
// 		fmt.Println("countryStr:", countryStr)
// 		if c == countryStr {
// 			return true
// 		}
// 	}
// 	return false
// }

// func isPlatformMatched(platformStr string, condition *PlatformCondition) bool {
// 	if platformStr == "" {
// 		return true
// 	}
// 	for _, p := range *condition {
// 		fmt.Println("p:", p)
// 		fmt.Println("platformStr:", platformStr)
// 		if p == platformStr {
// 			return true
// 		}
// 	}
// 	return false
// }

// func getQueryInt(query url.Values, key string, defaultValue int) int {
// 	valueStr := query.Get(key)
// 	if valueStr == "" {
// 		return defaultValue
// 	}
// 	value, err := strconv.Atoi(valueStr)
// 	if err != nil || value < 1 || value > 100 {
// 		return defaultValue
// 	}
// 	return value
// }

func main() {
	log.Print("server has started")
	//start the db
	bundb := db.StartDB()
	// if err != nil {
	// 	log.Printf("error starting the database %v", err)
	// 	panic("error starting the database")
	// }
	//get the router of the API by passing the db
	router := api.StartAPI(bundb)
	//get the port from the environment variable
	port := os.Getenv("PORT")
	//pass the router and start listening with the server
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router %v\n", err)
		return
	}
	// http.HandleFunc("/api/v1/ad", CreateAdvertisementHandler)
	// http.HandleFunc("/api/v1/getAd", ListAdvertisementHandler)
	// http.ListenAndServe(":8080", nil)
}
