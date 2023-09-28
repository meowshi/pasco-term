package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Readln() string {
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	return s
}

func ReadlnTrimmed() string {
	s := Readln()
	return strings.Trim(s, " \n")
}

var RuToEn = map[rune]rune{
	'Ё': '`',
	'Й': 'Q',
	'Ц': 'W',
	'У': 'E',
	'К': 'R',
	'Е': 'T',
	'Н': 'Y',
	'Г': 'U',
	'Ш': 'I',
	'Щ': 'O',
	'З': 'P',
	'Х': '[',
	'Ъ': ']',
	'Ф': 'A',
	'Ы': 'S',
	'В': 'D',
	'А': 'F',
	'П': 'G',
	'Р': 'H',
	'О': 'J',
	'Л': 'K',
	'Д': 'L',
	'Ж': ';',
	'Э': '\'',
	'Я': 'Z',
	'Ч': 'X',
	'С': 'C',
	'М': 'V',
	'И': 'B',
	'Т': 'N',
	'Ь': 'M',
	'Б': ',',
	'Ю': '.',
}

func RfidToKey(id string) string {
	var key string

	for _, ch := range id {
		if letter, ok := RuToEn[ch]; ok {
			key += string(letter)
		} else {
			key += string(ch)
		}
	}

	return key
}

func AppFolderPath() string {
	folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("ошибка при получении пути программы: %s", err))
	}

	return folder
}
