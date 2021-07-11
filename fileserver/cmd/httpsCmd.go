package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/topcoder520/gfsr/fileserver/config"
	"github.com/topcoder520/gfsr/fileserver/handler"
	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/middleware"
)

var httpsCommand = &cobra.Command{
	Use:     "https",
	Short:   "HTTPs protocol file server",
	Example: "gofileserver https [-p 80] [-d path] --ca ca.pem --ser server.pem --key server.key",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := startHttpsServer(args); err != nil {
			logs.Error(err)
			return err
		}
		return nil
	},
}

func init() {
	RootCommand.AddCommand(httpsCommand)
	httpsCommand.Flags().IntVarP(&config.Port, "port", "p", 80, "The port of server listening on")
	httpsCommand.Flags().StringVarP(&config.Dir, "dir", "d", ".", "server directory")
	httpsCommand.Flags().StringVarP(&config.CACertPath, "ca", "c", "", "ca cert")
	httpsCommand.Flags().StringVarP(&config.ServerCertPath, "ser", "s", "", "server cert")
	httpsCommand.Flags().StringVarP(&config.ServerKeyPath, "key", "k", "", "key cert")
}

func startHttpsServer(args []string) error {
	absPath, err := filepath.Abs(filepath.Clean(config.Dir))
	if err != nil {
		return err
	}
	f, err := os.Open(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(absPath, 0600)
		if err != nil {
			return err
		}
	} else {
		fileInfo, err := f.Stat()
		f.Close()
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			return errors.New("dir not directory")
		}
	}
	config.AbsDir = absPath

	//cert
	caCrt, err := ioutil.ReadFile(config.CACertPath)
	if err != nil {
		log.Println("caCertPath ReadFile err :", err)
		return err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: middleware.FileServerMiddleWare(handler.GetServeMux()),
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert, //强制校验client端证书
			ClientCAs:  pool,                           //ca池 验证客户端证书的ca
		},
	}

	logs.Printf("listening port: %d \n", config.Port)
	logs.Printf("server directory: %s\n", config.AbsDir)
	logs.Println("https fileserver started successfully.")

	return server.ListenAndServeTLS(config.ServerCertPath, config.ServerKeyPath)
}
