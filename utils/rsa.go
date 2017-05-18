package utils

import (
	"crypto"
	"bytes"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error
	SignM(src [][]byte, hash crypto.Hash) ([]byte, error)
	VerifyM(src [][]byte, sign []byte, hash crypto.Hash) error
}

func pkcs1Padding(src []byte, keySize int) [][]byte {

	srcSize := len(src)

	blockSize := keySize - 11

	var v [][]byte

	if srcSize <= blockSize {
		v = append(v, src)
	} else {
		groups := len(src) / blockSize
		for i := 0; i < groups; i++ {
			block := src[:blockSize]

			v = append(v, block)
			src = src[blockSize:]

			if len(src) < blockSize {
				v = append(v, src)
			}
		}
	}
	return v
}

func unPadding(src []byte, keySize int) [][]byte {

	srcSize := len(src)

	blockSize := keySize

	var v [][]byte

	if srcSize == blockSize {
		v = append(v, src)
	} else {
		groups := len(src) / blockSize
		for i := 0; i < groups; i++ {
			block := src[:blockSize]

			v = append(v, block)
			src = src[blockSize:]
		}
	}
	return v
}

type pkcsClient struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (this *pkcsClient) Encrypt(plaintext []byte) ([]byte, error) {
	if this.publicKey == nil{
		panic("public key not exist")
	}
	blocks := pkcs1Padding(plaintext, this.publicKey.N.BitLen()/8)

	buffer := bytes.Buffer{}
	for _, block := range blocks {
		ciphertextPart, err := rsa.EncryptPKCS1v15(rand.Reader, this.publicKey, block)
		if err != nil {
			return nil, err
		}
		buffer.Write(ciphertextPart)
	}

	return buffer.Bytes(), nil
}

func (this *pkcsClient) Decrypt(ciphertext []byte) ([]byte, error) {
	if this.privateKey == nil{
		panic("private key not exist")
	}
	ciphertextBlocks := unPadding(ciphertext, this.privateKey.N.BitLen()/8)

	buffer := bytes.Buffer{}
	for _, ciphertextBlock := range ciphertextBlocks {
		plaintextBlock, err := rsa.DecryptPKCS1v15(rand.Reader, this.privateKey, ciphertextBlock)
		if err != nil {
			return nil, err
		}
		buffer.Write(plaintextBlock)
	}

	return buffer.Bytes(), nil
}

func (this *pkcsClient) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	if this.privateKey == nil{
		panic("private key not exist")
	}
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, hash, hashed)
}

func (this *pkcsClient) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	if this.publicKey == nil{
		panic("public key not exist")
	}
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.publicKey, hash, hashed, sign)
}

func (this *pkcsClient) SignM(src [][]byte, hash crypto.Hash) ([]byte, error) {
	if this.privateKey == nil{
		panic("private key not exist")
	}
	if len(src) < 1{
		panic("empty src")
	}
	h := hash.New()
	for _,data := range src{
		h.Write(data)
	}
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, hash, hashed)
}

func (this *pkcsClient) VerifyM(src [][]byte, sign []byte, hash crypto.Hash) error {
	if this.publicKey == nil{
		panic("public key not exist")
	}
	if len(src) < 1{
		panic("empty src")
	}
	h := hash.New()
	for _,data := range src{
		h.Write(data)
	}
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.publicKey, hash, hashed, sign)
}

type Type int64

const (
	PKCS1 Type = iota
	PKCS8
)

//默认客户端，pkcs8私钥格式，pem编码
func NewDefaultCipher(privateKey,publicKey string) (Cipher, error) {
	blockPri, _ := pem.Decode([]byte(privateKey))
	if blockPri == nil {
		return nil, errors.New("private key error")
	}

	blockPub, _ := pem.Decode([]byte(publicKey))
	if blockPub == nil {
		return nil, errors.New("public key error")
	}

	return NewCipher(blockPri.Bytes, blockPub.Bytes, PKCS8)
}

func NewCipher(privateKey, publicKey []byte, privateKeyType Type) (Cipher, error) {

	priKey, err := genPriKey(privateKey, privateKeyType)
	if err != nil {
		return nil, err
	}
	pubKey, err := genPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &pkcsClient{privateKey: priKey, publicKey: pubKey}, nil
}

func NewDefaultCipherEx(privateKey,publicKey string) (Cipher, error) {
	var bytePri []byte = nil
	if len(privateKey) > 0{
		blockPri, _ := pem.Decode([]byte(privateKey))
		if blockPri == nil {
			return nil, errors.New("private key error")
		}
		bytePri = blockPri.Bytes
	}

	var bytePub []byte = nil
	if len(publicKey) > 0{
		blockPub, _ := pem.Decode([]byte(publicKey))
		if blockPub == nil {
			return nil, errors.New("public key error")
		}
		bytePub = blockPub.Bytes
	}

	return NewCipherEx(bytePri, bytePub, PKCS8)
}

func NewCipherEx(privateKey, publicKey []byte, privateKeyType Type) (Cipher, error) {
	var priKey *rsa.PrivateKey = nil
	var err error
	if privateKey != nil{
		priKey, err = genPriKey(privateKey, privateKeyType)
		if err != nil {
			return nil, err
		}
	}

	var pubKey *rsa.PublicKey = nil
	if publicKey != nil{
		pubKey, err = genPubKey(publicKey)
		if err != nil {
			return nil, err
		}
	}
	return &pkcsClient{privateKey: priKey, publicKey: pubKey}, nil
}

func genPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func genPriKey(privateKey []byte, privateKeyType Type) (*rsa.PrivateKey, error) {
	var priKey *rsa.PrivateKey
	var err error
	switch privateKeyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, errors.New("unsupport private key type")
		}
	}
	return priKey, nil
}


//
//func test(){
//	privateKey := `-----BEGIN PRIVATE KEY-----
//MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKQFwnWl65dugYGW
//D0LrrQgeKEK3QrNhRw1I1nLvYHWcAL3cDGxmmHAOJcJuBww4m51eXR7zGH0uR2Y2
//GDpXzCJB35HfUCw7NTYb+Q7ve4uXQaaHDgn2A76qNdky05++PnmSHT9CzkZMWqhW
//wp8HAsFvu57OWU2T0RdsFxTbXR73AgMBAAECgYBd26XhKKbdqrCU9Mea5b3IDWnA
//c5nJh/rekTWV44DxC+oousipJzRHuvDEh62kwqfZr2veEAGNcHQO+xl2GVOHxDDV
//3UzlVVnU57Dd/gqHSxOd4kJdZfHGCNeqPP5zCMIHpLyq26ciXSD13qnsnPwKkZJQ
//cjlrQK7sYJRKUPLA4QJBANVWQLkUeWTdPMbAIrvI/A20BPMVAsfdlyuVUQ95Qtdz
//HVJ8DkFTVUkIZzwFnpgx+K7ZaaOHYgHqyx9mmES4vO8CQQDE0tz6eK4xlvcf2xqn
//0rCXzOQKLOILBstc1sQmhARNMf910JGYhOFkgVJZqiAEz/DXE6TGDREdJ6xSXUSX
//NU55AkEAlSJkwH1Vl3MpZ28tWMTZnuK3iw6nEP0RDoClWAHW/jIUz3K1rGkK97EO
//KeFrys00IVcPCCg+FUUDlgHsdC4ItQJBAJMhyf0W/6ikYMIYiRmRX19q08Fjgeqa
//PqV9Co58O7b1PDF3I4+vLcpy/ft3OI5AX5p33cILfJKd2KyNejvKpokCQGuyv6vy
//0fvf/ATWjFZQ4zHi+LVVtN9/XBpyIDMlM+zNxVUSxTwv52bKn6TQlXTv6XDazU00
//xL2iMhub6Conkls=
//-----END PRIVATE KEY-----`
//
//	publicKey := `-----BEGIN PUBLIC KEY-----
//MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCkBcJ1peuXboGBlg9C660IHihC
//t0KzYUcNSNZy72B1nAC93AxsZphwDiXCbgcMOJudXl0e8xh9LkdmNhg6V8wiQd+R
//31AsOzU2G/kO73uLl0Gmhw4J9gO+qjXZMtOfvj55kh0/Qs5GTFqoVsKfBwLBb7ue
//zllNk9EXbBcU210e9wIDAQAB
//-----END PUBLIC KEY-----`
//
//	a := "http://www.baidu.com/?a=b:c=d:e=ffasdf"
//	cipher,err := NewDefault(privateKey,publicKey)
//	if err != nil{
//		panic(err)
//	}
//	var data []byte
//	data,err = cipher.Sign([]byte(a),crypto.SHA256)
//	if err != nil{
//		panic(err)
//	}
//	sign := hex.EncodeToString(data)
//	fmt.Println(sign)
//
//	data, err = hex.DecodeString(sign)
//	if err != nil{
//		panic(err)
//	}
//	errV := cipher.Verify([]byte(a), data, crypto.SHA256)
//	if errV != nil {
//		panic(err)
//	}
//}
