package config

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v2"
)

var (
	jwtPublicKey  *rsa.PublicKey
	jwtPrivateKey *rsa.PrivateKey
	C             Config
)

type Config struct {
	Mysql  Mysql  `yaml:"mysql"`
	RsaKey RsaKey `yaml:"rsakey"`
}
type RsaKey struct {
	Public  string `yaml:"public"`
	Private string `yaml:"private"`
}

type Mysql struct {
	DBname   string `yaml:"dbname"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Init 用于初始化全局配置。
func Init(configFilePath string) error {
	if configFilePath == "" {
		return fmt.Errorf("empty configuring filepath")
	}
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	rawContent := make([]byte, len(yamlFile))
	copy(rawContent, yamlFile)
	if err := yaml.Unmarshal(yamlFile, &C); err != nil {
		return err
	}
	initJwtPublicKey(C.RsaKey.Public)
	initJwtPrivateKey(C.RsaKey.Private)
	return nil
}

func initJwtPrivateKey(privatePath string) {
	signBytes, err := ioutil.ReadFile(privatePath)
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("ParseRSAPrivateKeyFromPEM: %v", err)
	}
	return
}

func initJwtPublicKey(keyPath string) {
	verifyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}
	jwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatalf("ParseRSAPublicKeyFromPEM: %v", err)
	}
	return
}

func GetJwtPublicKey() *rsa.PublicKey {
	return jwtPublicKey
}

func GetJwtPrivateKey() *rsa.PrivateKey {
	return jwtPrivateKey
}
