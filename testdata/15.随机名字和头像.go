package main

import (
	"fmt"
	"github.com/DanPlayer/randomname"
	"github.com/disintegration/letteravatar"
	"github.com/golang/freetype"
	"image/png"
	"os"
	"path"
	"unicode/utf8"
)

func main() {
	// 随机名称生成
	//name := randomname.GenerateName()

	// 提前随机头像生成
	GenerateNameAvatar()
}

// GenerateNameAvatar 根据两个名字表，提前生成对应名字的第一个汉字对应的头像
func GenerateNameAvatar() {
	dir := "uploads/chat_avatar"
	for _, s := range randomname.AdjectiveSlice {
		DrawImage(string([]rune(s)[0]), dir)
	}
	for _, s := range randomname.PersonSlice {
		DrawImage(string([]rune(s)[0]), dir)
	}
}

// DrawImage 随机头像生成，指定中文字体参数
func DrawImage(name string, dir string) {
	fontFile, err := os.ReadFile("uploads/system/方正综艺简体.TTF")
	font, err := freetype.ParseFont(fontFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	options := &letteravatar.Options{
		Font: font,
	}
	// 绘制文字
	firstLetter, _ := utf8.DecodeRuneInString(name)
	img, err := letteravatar.Draw(140, firstLetter, options)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 保存
	filePath := path.Join(dir, name+".png")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = png.Encode(file, img)
	if err != nil {
		fmt.Println(err)
		return
	}
}
