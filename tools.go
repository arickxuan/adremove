package main

import (
	"crypto/rsa"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tjfoc/gmsm/pkcs12"
	"github.com/tjfoc/gmsm/x509"
	matepkcs12 "software.sslmate.com/src/go-pkcs12"
)

func base64P12ToP12File(base64P12 string) {
	// Base64 编码的 P12 数据
	//passphrase := "8D6CB874" // P12 文件密码

	// 解码 Base64
	p12Data, err := base64.StdEncoding.DecodeString(base64P12)
	if err != nil {
		fmt.Println("Base64 解码失败:", err)
		return
	}

	// 保存为 .p12 文件
	filePath := "ca_certificate.p12"
	err = os.WriteFile(filePath, p12Data, 0644)
	if err != nil {
		fmt.Println("保存文件失败:", err)
		return
	}

	fmt.Println("P12 文件已保存:", filePath)
}

func p12FileToPem() { //p12FileToPem
	// 读取 P12 文件
	p12Data, err := os.ReadFile("ca_certificate.p12")
	if err != nil {
		fmt.Println("无法读取 P12 文件:", err)
		return
	}

	// 解析 P12 文件（提供 P12 文件密码）
	password := "8D6CB874"
	privateKey, certificate, caCerts, err := matepkcs12.DecodeChain(p12Data, password)
	if err != nil {
		fmt.Println("解析 P12 文件失败:", err)
		return
	}

	// 创建 PEM 编码的私钥
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey)),
	})

	// 创建 PEM 编码的证书
	certificatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	})

	// 创建 PEM 编码的 CA 证书
	var caPEM []byte
	for _, caCert := range caCerts {
		caPEM = append(caPEM, pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caCert.Raw,
		})...)
	}

	// 保存为 .pem 文件
	os.WriteFile("private_key.pem", privateKeyPEM, 0644)
	os.WriteFile("certificate.pem", certificatePEM, 0644)
	os.WriteFile("ca_cert.pem", caPEM, 0644)

	log.Println("转换成功！已生成 private_key.pem, certificate.pem 和 ca_cert.pem")
}

func pemToP12() {
	// 读取私钥文件
	privateKeyData, err := ioutil.ReadFile("mitmproxy-ca.pem")
	if err != nil {
		log.Fatalf("读取私钥文件失败: %v", err)
	}

	// 解码私钥
	privateKeyBlock, _ := pem.Decode(privateKeyData)
	if privateKeyBlock == nil {
		log.Fatal("解析私钥失败")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	// 读取证书文件
	certData, err := ioutil.ReadFile("mitmproxy-ca-cert.pem")
	if err != nil {
		log.Fatalf("读取证书文件失败: %v", err)
	}

	// 解码证书
	certBlock, _ := pem.Decode(certData)
	if certBlock == nil {
		log.Fatal("解析证书失败")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		log.Fatalf("解析证书失败: %v", err)
	}

	// 转换为 P12
	pfxData, err := pkcs12.Encode(privateKey, cert, nil, "adremove")
	if err != nil {
		log.Fatalf("生成 P12 失败: %v", err)
	}

	// 保存 P12 文件
	err = ioutil.WriteFile("mitmproxy-ca.p12", pfxData, 0644)
	if err != nil {
		log.Fatalf("保存 P12 文件失败: %v", err)
	}

	fmt.Println("成功将证书转换为 P12 格式")
}

// 计算 OpenSSL 风格的 subject_hash_old
func calculateSubjectHashOld(subject pkix.Name) uint32 {
	// OpenSSL 的旧版本哈希算法
	var result uint32
	if len(subject.CommonName) > 0 {
		result = rotateLeft(result, 1) + uint32(subject.CommonName[0])
	}
	for _, org := range subject.Organization {
		if len(org) > 0 {
			result = rotateLeft(result, 1) + uint32(org[0])
		}
	}
	for _, country := range subject.Country {
		if len(country) > 0 {
			result = rotateLeft(result, 1) + uint32(country[0])
		}
	}
	return result
}

func rotateLeft(value uint32, bits uint) uint32 {
	return (value << bits) | (value >> (32 - bits))
}

func hashCa() {
	// 读取证书文件
	certPEMData, err := os.ReadFile("mitmproxy-ca-cert.pem")
	if err != nil {
		fmt.Printf("读取证书文件失败: %v\n", err)
		return
	}

	// 解析 PEM 格式
	block, _ := pem.Decode(certPEMData)
	if block == nil {
		fmt.Println("解析 PEM 格式失败")
		return
	}

	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("解析证书失败: %v\n", err)
		return
	}

	// 计算 subject_hash_old
	hashValue := calculateSubjectHashOld(cert.Subject)
	fmt.Printf("%08x\n", hashValue)
	re := fmt.Sprintf("%08x\n", hashValue)
	log.Println("re:", re, hashValue)

}

// android 可用
func pemToDer() {
	// 读取证书文件
	certPEMData, err := os.ReadFile("mitmproxy-ca-cert.pem")
	if err != nil {
		fmt.Printf("读取证书文件失败: %v\n", err)
		return
	}

	// 解析 PEM 格式
	block, _ := pem.Decode(certPEMData)
	if block == nil {
		fmt.Println("解析 PEM 格式失败")
		return
	}

	// 解析证书以验证其有效性
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("解析证书失败: %v\n", err)
		return
	}

	// 将 DER 格式的证书保存为 .cer 文件
	err = os.WriteFile("mitmproxy-ca-cert.cer", cert.Raw, 0644)
	if err != nil {
		fmt.Printf("保存证书文件失败: %v\n", err)
		return
	}

	log.Println("证书转换成功！已保存为 mitmproxy-ca-cert.cer")
}

func isInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func extractContent(s string) string {
	start := strings.Index(s, "(")   // 找到第一个 (
	end := strings.LastIndex(s, ")") // 找到最后一个 )

	if start == -1 || end == -1 || start >= end {
		return "" // 如果找不到或者位置不正确，则返回空字符串
	}
	return s[start+1 : end] // 提取括号中的内容
}

func trimBlank(str string) string {
	regex := regexp.MustCompile(`\s+`)
	noSpaceStr := regex.ReplaceAllString(str, "")
	return noSpaceStr
}
