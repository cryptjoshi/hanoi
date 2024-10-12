package encrypt

import (
	"crypto/md5"
	"crypto/des"
	"crypto/cipher"
	"bytes"
	"encoding/base64"
	//"encoding/hex"
	//"pkd/handler"
	"encoding/json"
	"os"
	"fmt"
	"time"
	"errors"
	//"log"
	"strconv"
	// "time"
	//"github.com/valyala/fasthttp"
)
var jwtKey  = []byte(os.Getenv("PASSWORD_SECRET"))
var CLIENT_ID = "bfaae307-613a-424d-a60b-04b1b8a2bc62" //"6342e1be-fa03-456f-8d2d-8e1c9513c351" //[]byte(os.Getenv("CLIENT_ID"))
var CLIENT_SECRET = "46052b8b"//"6d83ac42" //[]byte(os.Getenv("CLIENT_SECRET"))
var DESKEY = "3d68fd30" //"9c62a148"
var DESIV =	"2a492233"//"8e014099"
var SYSTEMCODE = "DocDemoSystem"//"tsxthb"
var WEBID = "DocDemoWeb" //"tsxthb"
type ECResult struct {
	Enc string `json:"enc"`
	Unx int64 `json:"unx"` // ค่า unx คุณสามารถกำหนดเอง
	Des string `json:"des"` // ค่า dex คุณสามารถกำหนดเอง
}
type GData struct {
	Token string `json:"token"`
}
type GResult struct {
	MsgId int `json:"msgid"`
	Message string `json:"message"`
	Data GData `json:"data"`
	Timestamp int `json:"timestamp"`
}
type Data struct {
	Token  string  `json:"token"`
	Amount float64 `json:"amount"`
	TranID string  `json:"tranId"`
}
type DecryptedData struct {
	MsgId     int         `json:"msgid"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int         `json:"timestamp"`
}

 

func DESEncrypt(data, key, iv []byte) (string, error) {
	// สร้าง DES cipher block
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Error creating DES cipher: %v", err)
	}

	// เติมข้อมูลด้วย PKCS5 ให้พอดีกับบล็อก
	paddedData := pkcs5Padding(data, block.BlockSize())

	// สร้าง CBC encrypter
	blockMode := cipher.NewCBCEncrypter(block, iv)

	// เข้ารหัสข้อมูล
	cipherText := make([]byte, len(paddedData))
	blockMode.CryptBlocks(cipherText, paddedData)

	// แปลง cipherText เป็น Base64
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func EncryptDESAndMD5(data, key, iv, clientID, clientSecret string) (*ECResult, error) {
	// เข้ารหัสข้อมูลด้วย DES
	encryptedData, err := DESEncrypt([]byte(data), []byte(key), []byte(iv))
	if err != nil {
		return nil, fmt.Errorf("Error encrypting data: %v", err)
	}

	// รับ timestamp ปัจจุบันเป็น Unix time
	unx := time.Now().UnixNano() / int64(time.Millisecond)

	// สร้าง MD5 hash
	hash := CreateMD5(fmt.Sprintf("%s%s%d%s", clientID, clientSecret, unx, encryptedData))

	// คืนค่าผลลัพธ์เป็น map
	return &ECResult{
		Enc: hash,
		Unx: unx,
		Des: encryptedData,
	}, nil
}

// TripleDESEncrypt เข้ารหัสด้วย 3DES
func TripleDESEncrypt(data, key, iv []byte) (string, error) {
	// สร้าง 3DES cipher block
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", fmt.Errorf("Error creating 3DES cipher: %v", err)
	}

	// เติมข้อมูลด้วย PKCS5 ให้พอดีกับบล็อก
	paddedData := pkcs5Padding(data, block.BlockSize())

	// สร้าง CBC encrypter
	blockMode := cipher.NewCBCEncrypter(block, iv)

	// เข้ารหัสข้อมูล
	cipherText := make([]byte, len(paddedData))
	blockMode.CryptBlocks(cipherText, paddedData)

	// แปลง cipherText เป็น Base64
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// CreateMD5 สร้าง MD5 จากข้อมูล
func CreateMD5(data string) string {
	hash := md5.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// encryptTripleDESAndMD5 เข้ารหัสข้อมูล JSON ด้วย Triple DES และสร้าง MD5
func encryptTripleDESAndMD5(data, key, iv, clientID, clientSecret string) (map[string]string, error) {
	// เข้ารหัสข้อมูลด้วย Triple DES
	encryptedData, err := TripleDESEncrypt([]byte(data), []byte(key), []byte(iv))
	if err != nil {
		return nil, fmt.Errorf("Error encrypting data: %v", err)
	}

	// รับ timestamp ปัจจุบันเป็น Unix time
	unx := time.Now().UnixNano() / int64(time.Millisecond)

	// สร้าง MD5 hash
	hash := CreateMD5(fmt.Sprintf("%s%s%d%s", clientID, clientSecret, unx, encryptedData))

	// คืนค่าผลลัพธ์เป็น map
	return map[string]string{
		"enc": hash,
		"des": encryptedData,
		"unx": fmt.Sprintf("%d", unx),
	}, nil
}

func DecryptDES5(encryptedData string, key []byte, iv []byte) (GResult, error) {
	// สร้าง DES cipher block
	block, err := des.NewCipher(key)
	if err != nil {
		return GResult{}, fmt.Errorf("Error creating cipher: %w", err)
	}

	// แปลงข้อมูลที่เข้ารหัสจาก Base64 เป็นไบต์
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return GResult{}, fmt.Errorf("Error decoding base64: %w", err)
	}

	// สร้าง CBC mode decrypter
	blockMode := cipher.NewCBCDecrypter(block, iv)

	// เตรียมตัวแปรสำหรับเก็บข้อมูลที่ถอดรหัส
	decrypted := make([]byte, len(encryptedBytes))

	// ถอดรหัสข้อมูล
	blockMode.CryptBlocks(decrypted, encryptedBytes)

	// ลบ padding
	unpaddedData := pkcs5Unpadding(decrypted)

	// สร้าง struct DecryptedData จาก unpaddedData
	var data DecryptedData
	if err := json.Unmarshal(unpaddedData, &data); err != nil {
		return GResult{}, fmt.Errorf("Error unmarshalling data: %w", err)
	}

	// สร้าง GResult โดยใช้ข้อมูลที่ได้จาก unpaddedData
	return GResult{
		MsgId:     data.MsgId,
		Message:   data.Message,
		Data:      GData{
			Token: ""},
		Timestamp: data.Timestamp,
	}, nil
}
// ฟังก์ชัน unpadding (PKCS5)
func pkcs5Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

//// padding 7

// ฟังก์ชันสำหรับเติม Padding แบบ PKCS7
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// ฟังก์ชันสำหรับลบ Padding แบบ PKCS7
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data length is zero")
	}
	unpadding := int(data[length-1])
	return data[:(length - unpadding)], nil
}

// ฟังก์ชันเข้ารหัส DES แบบ CBC พร้อม Padding แบบ PKCS7
func EncryptDES(jsonString string, desKey, desIV string) (string, error) {
	key := []byte(desKey)
	iv := []byte(desIV)

	if len(key) != des.BlockSize {
		return "", errors.New("invalid key size")
	}


	// สร้างบล็อกเข้ารหัส DES
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	// เติม Padding แบบ PKCS7 ให้กับข้อมูล
	data := pkcs7Padding([]byte(jsonString), block.BlockSize())

	// เข้ารหัสแบบ CBC
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)

	// แปลงผลลัพธ์เป็น Base64
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// ฟังก์ชันถอดรหัส DES แบบ CBC
func DecryptDES(encrypted string, desKey, desIV string) (string, error) {
	key := []byte(desKey)
	iv := []byte(desIV)

	// แปลงผลลัพธ์จาก Base64 กลับมาเป็น byte
	encryptedData, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	// สร้างบล็อกถอดรหัส DES
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	// ถอดรหัสแบบ CBC
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(decrypted, encryptedData)

	// ลบ Padding แบบ PKCS7
	decrypted, err = pkcs7Unpadding(decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func CreateSignature(clientID, clientSecret, encryptData string) (bool, ECResult) {
	ecresult := ECResult{}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CreateSignature:", r)
		}
	}()
	now := time.Now()
	unx := now.Unix() //UnixNano() / int64(time.Millisecond)
	// รวมข้อมูลทั้งหมด
	data := fmt.Sprintf("%s%s%d%s", clientID, clientSecret, unx, encryptData)
	//clientID + clientSecret + unx + encryptData
	//fmt.Println("data:",data)
	// เข้ารหัส MD5
	hash := md5.New()
	_, err := hash.Write([]byte(data))
	if err != nil {
		return false, ecresult
	}

	md5Hash := hash.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(md5Hash)
	ecresult.Des= encryptData
	ecresult.Unx= unx
	ecresult.Enc= signature
	return true, ecresult

	// แปลงผลลัพธ์เป็น string แบบ hex
	// signature := hex.EncodeToString(hash.Sum(nil))

	// return true, signature
}


func encryptTripleDES(data, desKey, desIV string) (string, error) {
	key := []byte(desKey)
	iv := []byte(desIV)
	fmt.Println(key)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}

	paddedData := pkcs5Padding([]byte(data), block.BlockSize())
	encrypted := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, paddedData)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// PKCS5 padding
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// ฟังก์ชันคำนวณ MD5 hash และแปลงเป็น base64
func CalculateMD5Base64(clientID, clientSecret, encryptedData string) string {
	unx := strconv.FormatInt(time.Now().UnixMilli(), 10)
	hash := md5.Sum([]byte(clientID + clientSecret + unx + encryptedData))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// ฟังก์ชันหลักสำหรับเข้ารหัสและคำนวณ MD5
// func Xencryption(data, desKey, desIV string) (map[string]interface{}, error) {
// 	unx := strconv.FormatInt(time.Now().UnixMilli(), 10)
	
// 	des, err := encryptTripleDES(data, desKey, desIV)
// 	if err != nil {
// 		return nil, err
// 	}

// 	md5Base64 := CalculateMD5Base64(CLIENT_ID, CLIENT_SECRET, des)

// 	return map[string]interface{}{
// 		"enc": md5Base64,
// 		"des": des,
// 		"unx": unx,
// 	}, nil
// }