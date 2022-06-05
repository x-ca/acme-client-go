package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

// AcmeUser You'll need a user or account type that implements acme.User
type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}

func (u AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func write(commonName string, key, cert []byte) error {
	// mkdir
	var dir = strings.Replace(commonName, "*.", "", -1)
	err := os.MkdirAll(fmt.Sprintf("certs/%s", dir), 0700)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// write key
	keyPath := fmt.Sprintf("certs/%s/%s.key", dir, commonName)
	keyFile, err := os.OpenFile(keyPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	_, err = keyFile.Write(key)
	if err != nil {
		return err
	}

	// write cert
	certPath := fmt.Sprintf("certs/%s/%s.crt", dir, commonName)
	certFile, err := os.OpenFile(certPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer certFile.Close()

	_, err = certFile.Write(cert)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// set env
	//os.Setenv("LEGO_CA_SERVER_NAME", "dev.xiexianbin.cn")

	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	AcmeUser := AcmeUser{
		Email: "me@xiexianbin.cn",
		key:   privateKey,
	}

	config := lego.NewConfig(&AcmeUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = "https://dev.xiexianbin.cn:14000/dir"
	config.Certificate.KeyType = certcrypto.RSA2048

	// disable http ssl/tls key

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an HTTP port of 5002 and an TLS port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	//err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "5001"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	AcmeUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{"xiexianbin.cn", "www.xiexianbin.cn", "dev.xiexianbin.cn", "*.xiexianbin.cn"},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	write(certificates.Domain, certificates.PrivateKey, certificates.Certificate)
	fmt.Println("Domain:", certificates.Domain)

	// ... all done.
}
