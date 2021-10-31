package util

import (
	"bufio"
	"busmap.vn/librarycore/config"
	"busmap.vn/librarycore/config/constants"
	"bytes"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/qor/media/oss"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type DICT map[string]interface{}

func TrimSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func OssToUrl(oss *oss.OSS) string {
	if oss == nil || EmptyOrBlankString(oss.Url) {
		return ""
	}
	return config.Config.CdnUrl + oss.Url
}

func SubUrlToFullUrl(subUrl string) string {
	if EmptyOrBlankString(subUrl) {
		return ""
	}
	return config.Config.Domain + "/media" + subUrl
}

func FormatNumberWithDelimiter(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}

func FormatPhoneNumber(s string) string {
	var str strings.Builder
	for _, v := range s {
		if v == '+' || ('0' <= v && v <= '9') {
			str.WriteRune(v)
		}
	}
	result := str.String()
	if strings.HasPrefix(result, "+") {
		result = strings.Replace(result, "+840", "+84", 1)
	} else if strings.HasPrefix(result, "84") {
		result = "+" + strings.Replace(result, "840", "84", 1)
	} else if strings.HasPrefix(result, "0") {
		result = strings.Replace(result, "0", "+84", 1)
	} else {
		result = "+84" + result
	}
	if len(result) < 10 {
		return ""
	}
	return result
}

var SOURCE_CHARACTERS, LL_LENGTH = stringToRune(`ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝàáâãèéêìíòóôõùúýĂăĐđĨĩŨũƠơƯưẠạẢảẤấẦầẨẩẪẫẬậẮắẰằẲẳẴẵẶặẸẹẺẻẼẽẾếỀềỂểỄễỆệỈỉỊịỌọỎỏỐốỒồỔổỖỗỘộỚớỜờỞởỠỡỢợỤụỦủỨứỪừỬửỮữỰự`)
var DESTINATION_CHARACTERS, _ = stringToRune(`AAAAEEEIIOOOOUUYaaaaeeeiioooouuyAaDdIiUuOoUuAaAaAaAaAaAaAaAaAaAaAaAaEeEeEeEeEeEeEeEeIiIiOoOoOoOoOoOoOoOoOoOoOoOoUuUuUuUuUuUuUu`)

func stringToRune(s string) ([]string, int) {
	ll := utf8.RuneCountInString(s)
	var texts = make([]string, ll+1)
	var index = 0
	for _, runeValue := range s {
		texts[index] = string(runeValue)
		index++
	}
	return texts, ll
}

func binarySearch(sortedArray []string, key string, low int, high int) int {
	var middle int = (low + high) / 2
	if high < low {
		return -1
	}
	if key == sortedArray[middle] {
		return middle
	} else if key < sortedArray[middle] {
		return binarySearch(sortedArray, key, low, middle-1)
	} else {
		return binarySearch(sortedArray, key, middle+1, high)
	}
}

func removeAccentChar(ch string) string {
	var index int = binarySearch(SOURCE_CHARACTERS, ch, 0, LL_LENGTH)
	if index >= 0 {
		ch = DESTINATION_CHARACTERS[index]
	}
	return ch
}

func RemoveAccent(s string) string {
	var buffer bytes.Buffer
	for _, runeValue := range s {
		buffer.WriteString(removeAccentChar(string(runeValue)))
	}
	return buffer.String()
}

func StringAddressFormat(s string) string {
	s = strings.ToLower(s)
	if !strings.HasPrefix(s, "0x") {
		s = "0x" + s
	}
	return s
}

func String2BigInt(s string) *big.Int {
	n := new(big.Int)
	n.SetString(s, 10)
	return n
}

func ShortenHexAddress(s string) string {
	res := strings.Replace(s, "0x", "", -1)
	return res
}

func ExtendHexAddress(s string) string {
	return "0x" + ShortenHexAddress(s)
}

func NullOrBlankString(s *string) bool {
	if s == nil {
		return true
	}
	return len(strings.TrimSpace(*s)) == 0
}

func EmptyOrBlankString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func Int2Byte32(n int) [32]byte {
	var res [32]byte
	for i := 0; i < 32; i++ {
		v := n % (256)
		n /= 256
		res[i] = byte(v)
	}
	return res
}

func Int2Byte32HexString(n int) string {
	arr := Int2Byte32(n)
	return hex.EncodeToString(arr[:])
}

func Int2Byte4(n uint32) [4]byte {
	var res [4]byte
	for i := 0; i < 4; i++ {
		v := n % (256)
		n /= 256
		res[i] = byte(v)
	}
	return res
}

func Int2Byte4HexString(n uint32) string {
	arr := Int2Byte4(n)
	return hex.EncodeToString(arr[:])
}

func Int2Byte8(n uint64) [8]byte {
	var res [8]byte
	for i := 0; i < 8; i++ {
		v := n % (256)
		n /= 256
		res[i] = byte(v)
	}
	return res
}

func Int2Byte8HexString(n uint64) string {
	arr := Int2Byte8(n)
	return hex.EncodeToString(arr[:])
}

func StringPad2Byte32(str string) [32]byte {
	if len(str) > 32 {
		str = str[0:32]
	}
	var res [32]byte
	copy(res[:], str)
	return res
}

func Wei2Ether(a *big.Int) float64 {
	res := 0.0
	s := a.String()
	mul := 1e-18
	for i := len(s) - 1; i >= 0; i-- {
		x := s[i]
		if x != '0' {
			res = res + float64(int(x-'0'))*mul
		}
		mul = mul * 10
	}
	return res
}

func HexString2BigInt(s string) *big.Int {
	s = strings.Replace(s, "0x", "", 1)
	n := new(big.Int)
	n, _ = n.SetString(s, 16)
	if n == nil {
		n = big.NewInt(0)
	}
	return n
}

func HexString2ByteArray(s string) []byte {
	s = strings.Replace(s, "0x", "", 1)
	res, _ := hex.DecodeString(s)
	return res
}

func RoundFloat(f float64, d int) float64 {
	i, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(d)+"f", f), 64)
	return i
}

func Sha256Data(answer []byte) [32]byte {
	sum := sha256.Sum256(answer)
	return sum
}

func HexSha256String(arr []byte) string {
	h := sha256.New()
	h.Write(arr)
	hashedData := h.Sum(nil)
	return hex.EncodeToString(hashedData)
}

func FinalizePhoneNumber(number, countryPrefix string) string {
	// Finalize Vietnamese phone number
	if countryPrefix == "84" && len(number) == 13 {
		// Mobifone
		number = strings.Replace(number, "+84120", "+8470", 1)
		number = strings.Replace(number, "+84121", "+8479", 1)
		number = strings.Replace(number, "+84122", "+8477", 1)
		number = strings.Replace(number, "+84126", "+8476", 1)
		number = strings.Replace(number, "+84128", "+8478", 1)
		// Vinaphone
		number = strings.Replace(number, "+84123", "+8483", 1)
		number = strings.Replace(number, "+84124", "+8484", 1)
		number = strings.Replace(number, "+84125", "+8485", 1)
		number = strings.Replace(number, "+84127", "+8481", 1)
		number = strings.Replace(number, "+84129", "+8482", 1)
		// Viettel
		number = strings.Replace(number, "+84162", "+8432", 1)
		number = strings.Replace(number, "+84163", "+8433", 1)
		number = strings.Replace(number, "+84164", "+8434", 1)
		number = strings.Replace(number, "+84165", "+8435", 1)
		number = strings.Replace(number, "+84166", "+8436", 1)
		number = strings.Replace(number, "+84167", "+8437", 1)
		number = strings.Replace(number, "+84168", "+8438", 1)
		number = strings.Replace(number, "+84169", "+8439", 1)
		// Vietnamobile
		number = strings.Replace(number, "+84186", "+8456", 1)
		number = strings.Replace(number, "+84188", "+8458", 1)
		// Gmobile
		number = strings.Replace(number, "+84199", "+8459", 1)
	}

	return number
}

var base31Alphabet = []byte("12345abcdefghijklmnopqrstuvwxyz")

func IdToDSBAccountName(id uint64) string {
	a := []byte{'u', 's', 'e', 'r', '.', '.', '.', '.', '.', '.', '.', '.'}
	for i := 11; i >= 4; i-- {
		v := id % 31
		a[i] = base31Alphabet[v]
		id /= 31
	}
	return string(a)
}

func IdToDSBChallengeName(id uint64) string {
	a := []byte{'p', 'r', 'o', 'b', '.', '.', '.', '.', '.', '.', '.', '.'}
	for i := 11; i >= 4; i-- {
		v := id % 31
		a[i] = base31Alphabet[v]
		id /= 31
	}
	return string(a)
}

func IdToProblemName(id uint64) string {
	a := []byte{'c', 'o', 'd', 'e', '.', '.', '.', '.', '.', '.', '.', '.'}
	for i := 11; i >= 4; i-- {
		v := id % 31
		a[i] = base31Alphabet[v]
		id /= 31
	}
	return string(a)
}

func IdToAttemptName(id uint64) string {
	a := []byte{'s', 'u', 'b', '.', '.', '.', '.', '.', '.', '.', '.', '.'}
	for i := 11; i >= 3; i-- {
		v := id % 31
		a[i] = base31Alphabet[v]
		id /= 31
	}
	return string(a)
}

func FormatOutput(output string) string {
	output = strings.ToLower(output)
	tab := regexp.MustCompile(`\t+`)
	output = tab.ReplaceAllString(output, " ")

	scanner := bufio.NewScanner(strings.NewReader(output))
	var str strings.Builder
	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
		} else {
			str.WriteString("\n")
		}
		str.WriteString(strings.TrimSpace(scanner.Text()))
	}
	return strings.TrimSpace(str.String())
}

func ToStructFieldSlice(slice interface{}, fieldName string) ([]interface{}, error) {
	res := make([]interface{}, 0)

	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)

		for i := 0; i < s.Len(); i++ {
			val := s.Index(i).Interface()
			value, found := GetFieldValue(val, fieldName)
			if !found {
				message := fmt.Sprintf("%s field is not found at index [%d]", fieldName, i)
				return res, errors.New(message)
			}
			res = append(res, value)
		}
	default:
		return res, errors.New("input must be slice")
	}
	return res, nil
}

func GetFieldValue(value interface{}, fieldName string) (interface{}, bool) {
	if _, found := reflect.TypeOf(value).FieldByName(fieldName); !found {
		return nil, false
	}
	res := reflect.Indirect(reflect.ValueOf(value)).FieldByName(fieldName)
	return res.Interface(), true
}

type isMn struct{}

func (isMn) Contains(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func GenerateUrl(value string) string {
	isMn := isMn{}
	t := transform.Chain(norm.NFD, runes.Remove(isMn), norm.NFC)
	value, _, _ = transform.String(t, value)
	re := regexp.MustCompile(`\W+`)
	return re.ReplaceAllLiteralString(strings.ToLower(value), "-")
}

func ParseBool(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.String:
		return v.Len() > 0
	case reflect.Int:
		return v.Int() > 0
	case reflect.Float32, reflect.Float64:
		return v.Float() > 0
	default:
		return false
	}
}

// Haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
// SOURCE: "https://gist.github.com/cdipaolo/d3f8db3848278b49db68"
func CalculateDistanceBetweenPoint(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func CalculateDeg(lat1, lon1, lat2, lon2 float64) float64 {
	if math.Abs(lat1-lat2) < 1e-8 && math.Abs(lon1-lon2) < 1e-8 {
		return -100
	}
	res := (-math.Atan2(lat2-lat1, lon2-lon1) / math.Pi * 180) + 90
	if res < 0 {
		res = 360 + res
	}
	return res
}

func CalculateSpeedAndHeading(lat1, lon1, lat2, lon2 float64, t1, t2 int64) (speed, heading float64) {
	if t1 >= t2 {
		return 0, -100
	}
	distance := CalculateDistanceBetweenPoint(lat1, lon1, lat2, lon2)
	delta := t2 - t1
	speed = 3.6 * distance / float64(delta)
	heading = CalculateDeg(lat1, lon1, lat2, lon2)
	return speed, heading
}

func ContainsString(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsInt(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsUInt(a []uint, x uint) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func GetLocation() *time.Location {
	locale, err := time.LoadLocation(constants.LocationName)
	if err != nil {
		panic(err)
	}
	return locale
}

func GetHourAndMinuteInGMT() (hour, minute int) {
	now := time.Now().In(GetLocation())
	return now.Hour(), now.Minute()
}

func Add0ToFirstPosition(number int) string {
	// Add0ToFirstPosition "8 => '08'"
	_str := strconv.Itoa(number)
	if len(_str) == 1 {
		return "0" + _str
	}
	return _str
}

func GetHourInTimeString() string {
	hour, _ := GetHourAndMinuteInGMT()
	return Add0ToFirstPosition(hour)
}

func GetMinuteInTimeString() string {
	_, minutes := GetHourAndMinuteInGMT()
	return Add0ToFirstPosition(minutes)
}

func PrefixHourMinuteString() string {
	return fmt.Sprintf("[%s:%s]", GetHourInTimeString(), GetMinuteInTimeString())
}

func FormatString(rawString string) string {
	return strings.ToLower(strings.TrimSpace(rawString))
}

func ConvertHexToDec(hex string) int64 {
	// Replace 0x or 0X with empty String
	numberStr := strings.Replace(strings.ToLower(hex), "0x", "", -1)
	output, _ := strconv.ParseInt(numberStr, 16, 64)
	return output
}

func Reverse(s string) string {
	// Reverse returns its argument string reversed rune-wise left to right.
	if len(s)%2 == 1 {
		s = "0" + s
	}
	r := []rune(s)
	for i, j := 0, len(r)-2; i < len(r)/2; i, j = i+2, j-2 {
		r[i], r[j] = r[j], r[i]
		r[i+1], r[j+1] = r[j+1], r[i+1]
	}
	return string(r)
}

func ReverseDec(s string) string {
	v, _ := strconv.ParseInt(s, 10, 64)
	hexValue := fmt.Sprintf("%x", v)
	revHex := Reverse(hexValue)
	dec := ConvertHexToDec(revHex)
	return fmt.Sprintf("%v", dec)
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

func GetBeginDayInUnix(year, month, day int) int64 {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Unix()
}

func GetEndDayInUnix(year, month, day int) int64 {
	return time.Date(year, time.Month(month), day, 23, 59, 59, 0, time.UTC).Unix()
}

func GetBeginDate(now time.Time) time.Time {
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 1, 0, time.UTC)
}

func GetEndDate(now time.Time) time.Time {
	year, month, day := now.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
}

func GetMorningOrAfternoonWithGMT() int {
	hour := time.Now().In(GetLocation()).Hour()
	if hour < constants.NoonTimeBy24Hour {
		return constants.InMorning
	} else {
		return constants.InAfternoon
	}
}

func ReadXLXS(fullPath string) *excelize.File {
	f, err := excelize.OpenFile(fullPath)
	if err != nil {
		panic(err)
	}
	return f
}

func RemoveDuplicateUintValues(intSlice []uint) []uint {
	keys := make(map[uint]bool)
	list := make([]uint, 0)

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func PointerTimeToUnix(time *time.Time) int64 {
	unixTime := int64(0)
	if time != nil {
		unixTime = time.Unix()
	}
	return unixTime
}

func UIntArrayToStringWithCommas(array []uint) string {
	if len(array) == 0 {
		return ""
	}

	// Example: [1 2 3 4 5] => 1,2,3,4,5
	arrayInString := fmt.Sprintf("%v", array)
	arrayInString = strings.Replace(arrayInString, "[", "", 1)
	arrayInString = strings.Replace(arrayInString, "]", "", 1)
	arrayInString = strings.ReplaceAll(arrayInString, " ", ",")

	return arrayInString
}

func JoinWithUintArray(uintArr []uint, separate string) string {
	strArr := make([]string, 0)
	for _, uintEl := range uintArr {
		strArr = append(strArr, fmt.Sprintf("%d", uintEl))
	}
	return strings.Join(strArr, separate)
}

func GetNumberDayOfMonth(year, month int) int {
	locale, err := time.LoadLocation(constants.LocationName)
	if err != nil {
		panic(err)
	}
	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, locale)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	return lastDayOfMonth.Day()
}

func HashingPassword(password string) string {
	bPassword := []byte(password)
	hashPassword, err := bcrypt.GenerateFromPassword(bPassword, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashPassword)
}

func RefineUIDDec(key string) string {
	for len(key) < 10 {
		key = "0" + key
	}
	return key
}

func RefineHex(key string) string {
	for len(key) < 8 {
		key = "0" + key
	}
	return key
}

func GetWhereConditionWithUIDHexTypes(uidHex string) string {
	// mã hex
	// mã hex đảo
	// mã dec
	// mã dec đảo
	_code := FormatString(uidHex)
	// ---
	_codeRefined := RefineHex(_code)
	// ---
	fCode := strconv.Itoa(int(ConvertHexToDec(_codeRefined)))
	fCodeRefined := RefineUIDDec(strconv.Itoa(int(ConvertHexToDec(_codeRefined))))
	// ---
	reversedFCode := strconv.Itoa(int(ConvertHexToDec(Reverse(_codeRefined))))
	reversedFCodeRefined := RefineUIDDec(strconv.Itoa(int(ConvertHexToDec(Reverse(_codeRefined)))))
	// ---
	whereCondition := fmt.Sprintf("uid ='%s' or uid ='%s' or uid ='%s' or uid ='%s' or uid ='%s' or uid = '%s'", _code, _codeRefined, fCode, fCodeRefined, reversedFCode, reversedFCodeRefined)
	return whereCondition
}

func GetWhereConditionWithUIDDecTypes(uidDec string) string {
	// chuyển mã dec về int
	// mã dec
	// mã dec đảo
	_codeInt, err := strconv.ParseInt(uidDec, 10, 64)
	codeIntCondition := ""
	if err == nil {
		codeIntCondition = fmt.Sprintf("or uid = '%d'", _codeInt)
	}
	// ---
	_code := RefineUIDDec(FormatString(uidDec))
	reversedCode := RefineUIDDec(ReverseDec(_code))
	whereCondition := fmt.Sprintf("uid ='%s' or uid ='%s' %s", reversedCode, _code, codeIntCondition)
	return whereCondition
}

func GetBeginAndEndUnix(selectedDateInUnix int64) (int64, int64) {
	beginUnix := int64(0)
	endUnix := int64(0)

	if selectedDateInUnix > 0 {
		var location *time.Location
		var err error
		if location, err = time.LoadLocation(constants.LocationName); err != nil {
			return beginUnix, endUnix
		}
		// ---
		dt := time.Unix(selectedDateInUnix, 0).In(location)
		day, month, year := dt.Day(), dt.Month(), dt.Year()
		// ---
		beginUnix = time.Date(year, month, day, 0, 0, 0, 0, location).Unix()
		endUnix = time.Date(year, month, day, 23, 59, 59, 0, location).Unix()
	}

	return beginUnix, endUnix
}

func IsEmailValid(e string) bool {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

type MultiThreadHandler func(startIndex, endIndex int, wg *sync.WaitGroup)

func RunFunctionWithMultiThreads(threadQuantity, objectArraySize int, handler MultiThreadHandler) {
	if objectArraySize <= 0 {
		return
	}

	if threadQuantity <= 0 {
		threadQuantity = 2
	}

	if objectArraySize <= threadQuantity {
		startIndex := 0
		endIndex := objectArraySize

		var waitGroup sync.WaitGroup
		waitGroup.Add(objectArraySize)

		handler(startIndex, endIndex, &waitGroup)

		waitGroup.Wait()
	} else {
		loopQuantity := objectArraySize / threadQuantity
		if objectArraySize%threadQuantity != 0 {
			loopQuantity += 1
		}

		for k := 0; k < loopQuantity; k++ {
			startIndex := k * threadQuantity
			endIndex := (k + 1) * threadQuantity
			if endIndex > objectArraySize {
				endIndex = objectArraySize
			}

			var waitGroup sync.WaitGroup
			waitGroup.Add(endIndex - startIndex)

			handler(startIndex, endIndex, &waitGroup)

			waitGroup.Wait()
		}
	}
}

func DateInRange(date, rangeStartDate, rangeEndDate time.Time) bool {
	if rangeStartDate.After(rangeEndDate) {
		panic("Invalid RangeStartDate and RangeEndDate!")
	}

	if (date.Equal(rangeStartDate) || date.After(rangeStartDate)) &&
		(date.Equal(rangeEndDate) || date.Before(rangeEndDate)) {
		return true
	}
	return false
}

func DateOutRange(date, rangeStartDate, rangeEndDate time.Time) bool {
	if rangeStartDate.After(rangeEndDate) {
		panic("Invalid RangeStartDate and RangeEndDate!")
	}

	if date.Before(rangeStartDate) || date.After(rangeEndDate) {
		return true
	}
	return false
}

func StandardizedString(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Join(strings.Split(str, " "), " ")

	return str
}

func SliceRemoveAtIndex(slice []interface{}, index int) []interface{} {
	return append(slice[:index], slice[index+1:]...)
}

func Min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func GetUnix(t *time.Time) int64 {
	unix := t.Unix()
	if unix <= 0 {
		return 1
	}
	return unix
}

func IsZeroTime(t time.Time) bool {
	return (t.Year() == 1 &&
		t.Month() == 1 &&
		t.Day() == 1 &&
		t.Hour() == 0 &&
		t.Minute() == 0 &&
		t.Second() == 0) || t.Unix() == constants.TimeZeroValueUnix
}

func IsT1EqualT2(t1 *time.Time, t2 *time.Time) bool {
	if t1 == nil || t2 == nil {
		return false
	}

	return t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day()
}

func IsT1AfterT2(t1 time.Time, t2 time.Time) bool {
	t1DateUnix := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local).Unix()
	t2DateUnix := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local).Unix()

	return t1DateUnix > t2DateUnix && t1.Unix() > t2.Unix()
}

func AppendValueInsert(dbFields []string, values []string, fieldName string, value interface{}) ([]string, []string) {
	dbFields = append(dbFields, fieldName)

	switch value.(type) {
	case time.Time:
		values = append(values, fmt.Sprintf("'%s'", value.(time.Time).Format("2006-01-02 15:04:05")))
	default:
		values = append(values, fmt.Sprintf("'%v'", value))
	}

	return dbFields, values
}

func GetNowDate() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 17, 0, 0, 0, time.UTC)
}

func GetNowDateString() string {
	return fmt.Sprintf("%d-%d-%d 11:00:00", time.Now().Year(), time.Now().Month(), time.Now().Day()-1)
}

func FormatDate(t time.Time) time.Time {
	y, m, d := t.Date()

	if d <= 12 {
		return time.Date(y, time.Month(d), int(m), 0, 0, 0, 0, t.Location())
	}

	return t
}
