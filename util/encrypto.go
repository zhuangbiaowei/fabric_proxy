package util

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
)

func EcdsaSignWithSha256Hex(data []byte, privateKeyPath, signGzipSwitch string) (string, error) {
	keyPEMBlock, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil {
		return "", err
	}

	privateKey, errParsePK := x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	if errParsePK != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}

	h := sha256.New()
	h.Write(data)
	hash := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	var out string
	if signGzipSwitch == "0" {
		out = hex.EncodeToString([]byte(string(rt) + "+" + string(st)))
	} else if signGzipSwitch == "1"{
		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		defer w.Close()
		_, err = w.Write([]byte(string(rt) + "+" + string(st)))
		if err != nil {
			return "", err
		}
		w.Flush()
		out = hex.EncodeToString(b.Bytes())
	}
	return out, nil
}

func getSign(signature string) (rint, sint big.Int, err error) {
	byterun, err := hex.DecodeString(signature)
	if err != nil {
		//err = errors.New("decrypt error, "+ err.Error())
		return
	}
	r, err := gzip.NewReader(bytes.NewBuffer(byterun))
	if err != nil {
		//err = errors.New("decode error,"+err.Error())
		return
	}
	defer r.Close()
	buf := make([]byte, 1024)
	count, err := r.Read(buf)
	if err != nil {
		fmt.Println("decode = ", err)
		//err = errors.New("decode read error," + err.Error())
		return
	}
	rs := strings.Split(string(buf[:count]), "+")
	fmt.Println(">>>>>>>", rs)
	if len(rs) != 2 {
		//err = errors.New("decode fail")
		return
	}
	fmt.Println("rs0>>>>>>>", rs[0])
	fmt.Println("rs1>>>>>>>", rs[1])
	err = rint.UnmarshalText([]byte(rs[0]))

	fmt.Println("rint>>>>>", rint)
	if err != nil {
		//err = errors.New("decrypt rint fail, "+ err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	fmt.Println("sint>>>>>", sint)
	if err != nil {
		//err = errors.New("decrypt sint fail, "+ err.Error())
		return
	}
	return

}

/**
  校验文本内容是否与签名一致
  使用公钥校验签名和文本内容
*/
func verify(text []byte, signature string, key *ecdsa.PublicKey) (bool, error) {
	rint, sint, err := getSign(signature)
	if err != nil {
		return false, err
	}

	h := sha256.New()
	h.Write(text)
	//h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)

	result := ecdsa.Verify(key, hash[:], &rint, &sint)
	return result, nil
}

func parseCert(crt string) {
	//var cert tls.Certificate
	certPEMBlock, err := ioutil.ReadFile(crt)
	if err != nil {
		return
	}
	certDERBlock, _ := pem.Decode(certPEMBlock)
	fmt.Println(string(certDERBlock.Bytes))
	x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	fmt.Println(x509Cert.PublicKeyAlgorithm)
}
