package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/axgle/mahonia"
	"os"
	"os/user"
	_ "os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	defaultSheet = "Sheet1"
)

func main() {
	//filepath
	inputFile := getInputFilePath()

	//check source file exist
	exist, _ := pathExists(inputFile)

	//read excel
	if exist {
		readExcelToFile(inputFile)
	}
}

func getInputFilePath() string {

	list := os.Args

	if len(list) < 2 {
		fmt.Printf("invalid filepath, please input your filepath as arguement")
		return ""
	}

	return list[1]
}


func getCurrentDirectory() string {
	str, _ := os.Getwd()
	return  str
}


func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readExcelToFile(path string) {

	// open file
	xcix, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all the rows in the Sheet1.
	rows := xcix.GetRows(defaultSheet)

	//from left to right
	for column := 1; column < len(rows[0]); column++ {

		langKey := rows[0][column]
		targetFileName := findLocalizedFilePathByLanguage(langKey)
		fmt.Print("开始写入到文件" + targetFileName + "\n")


		outputFile, e := os.OpenFile(targetFileName, os.O_WRONLY|os.O_APPEND, 0666)
		if e != nil {
			fmt.Print(e)
		}
		_, _ = outputFile.WriteString("\n")

		for rowIndex, row := range rows {
			key := row[0]
			value := row[column]
			if len(key) > 0 && len(value) > 0 {
				processorString := processor(key, value)
				if rowIndex == 0 {
					processorString = extraInfo(processorString)

				}
				_, _ = outputFile.WriteString(processorString + "\n")
				fmt.Print(processorString + "\n")
			}
		}
		outputFile.Close()
		fmt.Print("\n\n")
	}
}


func findLocalizedFilePathByLanguage(language string) string {

	languageReflection := map[string]string{
		"Chinese"	:"zh-Hans",      // 简体中文
		"繁体"		:"zh-Hant",    // ios简体中文
		"阿语"		:"ar",         // 阿拉伯语
		"English"	:"en",         // 英语
		"西语"		:"es",         // 西班牙语
		"葡语"		:"pt",         // 葡萄牙语
		"法语"		:"fr",         // 法语
		"印尼语"		:"id",         // 印尼语
		"泰语"		:"th",         // 泰语
		"土耳其语"	:"tr",	      // 土耳其
	}

	standardLanguage := languageReflection[language]

	regString := ".*" + standardLanguage + ".lproj/MicoLocalizable.strings" +  ".*"

	fileName := getCurrentDirectory()

	fmt.Print("搜索文件路径... : " + fileName + "   >>>> 语言：" + language+ "\n")

	return searchInFilePath(getCurrentDirectory(),regString)
}


func searchInFilePath(fileName string, pattern string) string {

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	fi, err := file.Stat()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if !fi.IsDir() {
		fmt.Println(fileName, " is not a dir")
	}

	reg, err := regexp.Compile(pattern)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	targetPath := ""
	// 遍历目录
	_ = filepath.Walk(fileName,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if f.IsDir() {
				return nil
			}

			// 匹配目录
			matched := reg.MatchString(path)

			if matched {
				targetPath = path
			}

			return nil
		})
	return targetPath
}

func extraInfo(lanString string) string {
	timeStamp := "TimeStamp:" + time.Now().Format("2006-01-02 15:04:05") + ";"
	author,_ := user.Current()
	authorString := "Added by:" + author.Username
	lanString = "//" + lanString + timeStamp + authorString
	return lanString
}

func processor(key, lanString string) string {
	lanString = strings.Replace(lanString, "\"", "\\\"", -1) //" -> \"
	lanString = strings.Replace(lanString, "%s", "%@", -1)   //%s -> %@"
	lanString = strings.Replace(lanString, "%S", "%@", -1)   //%s -> %@"#
	lanString = strings.Replace(lanString, "%m", "%@", -1)   //%s -> %@"
	lanString = strings.Replace(lanString, "%g", "%@", -1)   //%g -> %@"
	lanString = strings.Replace(lanString, "%l", "%@", -1)   //%L -> %@"
	lanString = strings.Replace(lanString, "%L", "%@", -1)   //%L -> %@"
	lanString = strings.Replace(lanString, "XXX", "%@", -1)  //xxx -> %@"
	lanString = strings.Replace(lanString, "xxx", "%@", -1)  //xxx -> %@"
	lanString = strings.Replace(lanString, "XX", "%@", -1)   //xx -> %@"
	lanString = strings.Replace(lanString, "yyy", "%@", -1)  //xx -> %@"
	lanString = strings.Replace(lanString, "xx", "%@", -1)   //xx -> %@"
	lanString = strings.Replace(lanString, "NN", "%B", -1)   //xx -> %B"
	lanString = strings.Replace(lanString, "N", "%@", -1)    //xxx -> %@"
	lanString = strings.Replace(lanString, "***", "%@", -1)  //xxx -> %@"
	lanString = strings.Replace(lanString, "\n", "\\n", -1)  //xxx -> %@"
	lanString = strings.Replace(lanString, "\n", "\\n", -1)  //xxx -> %@"

	iosLagText := "\"" + key + "\"" + " = " + "\"" + lanString + "\"" + ";"
	iosLagText = converToUtf8(iosLagText)

	return iosLagText
}

func converToUtf8(lang string) string {
	enc := mahonia.NewEncoder("utf-8")
	return enc.ConvertString(lang)
}

