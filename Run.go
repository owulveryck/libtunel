package libtunnel

import (
	// Because of the ciphers... I know it is bad
	"github.com/ScriptRock/crypto/ssh"
	"io/ioutil"
	"os/user"
)

func getKeyFile() (key ssh.Signer, err error) {
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		panic(err)
		return
	}
	return
}
func Run(username string, password string, serverAddrString string, localAddrString string, remoteAddrString string, currentRetriesLocal int) {

	// Get the key file
	key, err := getKeyFile()
	if err != nil {
		panic(err)
	}
	// Setup SSH config (type *ssh.ClientConfig)
	// Create the SSH configuration:
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Config: ssh.Config{
			Ciphers: ssh.AllSupportedCiphers(), // include cbc ciphers (even though they're bad, mmmkay)
		},
	}

	// Create the local end-point:
	localListener := CreateLocalEndPoint(localAddrString, currentRetriesLocal)

	// Accept client connections (will block forever):
	AcceptClients(localListener, config, serverAddrString, remoteAddrString, password)
}
