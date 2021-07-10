package cryutil

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
)

//出处:https://blog.csdn.net/hbshhb/article/details/93404539
//加密
func EncryptDesCBC(plainText /*明文*/, key []byte) ([]byte, error) {
	//第一步：创建des密码接口, 输入秘钥，返回接口
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//第二步：创建cbc分组
	// 返回一个密码分组链接模式的、底层用b解密的BlockMode接口
	// func NewCBCEncrypter(b Block, iv []byte) BlockMode
	blockSize := block.BlockSize()

	//创建一个8字节的初始化向量
	iv := bytes.Repeat([]byte("1"), blockSize)

	mode := cipher.NewCBCEncrypter(block, iv)

	//第三步：填充
	//TODO
	plainText, err = paddingNumber(plainText, blockSize)
	if err != nil {
		return nil, err
	}

	//第四步：加密
	// type BlockMode interface {
	// 	// 返回加密字节块的大小
	// 	BlockSize() int
	// 	// 加密或解密连续的数据块，src的尺寸必须是块大小的整数倍，src和dst可指向同一内存地址
	// 	CryptBlocks(dst, src []byte)
	// }

	//密文与明文共享空间，没有额外分配
	mode.CryptBlocks(plainText /*密文*/, plainText /*明文*/)

	return plainText, nil
}

//输入密文，得到明文
func DecryptDesCBC(encryptData, key []byte) ([]byte, error) {
	//TODO
	//第一步：创建des密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//第二步：创建cbc分组
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	mode := cipher.NewCBCDecrypter(block, iv)

	//第三步：解密
	mode.CryptBlocks(encryptData /*明文*/, encryptData /*密文*/)

	//第四步: 去除填充
	//TODO
	encryptData, err = unPaddingNumber(encryptData)
	if err != nil {
		return nil, err
	}

	// return []byte("Hello world"), nil
	return encryptData, nil
}

//填充数据
func paddingNumber(src []byte, blockSize int) ([]byte, error) {

	if src == nil {
		return nil, errors.New("src长度不能小于0")
	}

	//1. 得到分组之后剩余的长度 5
	leftNumber := len(src) % blockSize //5

	//2. 得到需要填充的个数 8 - 5 = 3
	needNumber := blockSize - leftNumber //3

	//3. 创建一个slice，包含3个3
	b := byte(needNumber)
	newSlice := bytes.Repeat([]byte{b}, needNumber) //newSlice  ==》 []byte{3,3,3}

	fmt.Printf("newSclie : %v\n", newSlice)
	//4. 将新切片追加到src
	src = append(src, newSlice...)

	return src, nil
}

//解密后去除填充数据
func unPaddingNumber(src []byte) ([]byte, error) {
	//1. 获取最后一个字符
	lastChar := src[len(src)-1] //byte(3)

	//2. 将字符转换为数字
	num := int(lastChar) //int(3)

	//3. 截取切片(左闭右开)

	return src[:len(src)-num], nil
}
