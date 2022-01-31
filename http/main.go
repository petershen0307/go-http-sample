package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/chunked", chunkedHandler)
	// go router.RunTLS(":4343", "./certificate.pem", "./key.pem")
	go https(router)
	router.Run(":8080")
}

func https(g *gin.Engine) {
	srv := http.Server{
		Addr:    ":4343",
		Handler: g,
		// disable automatic enable http 2
		// http2: TLSConfig.CipherSuites is missing an HTTP/2-required AES_128_GCM_SHA256 cipher (need at least one of TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 or TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256).
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			MaxVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				// tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				// tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				// tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				// tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}
	err := srv.ListenAndServeTLS("./certificate.pem", "./key.pem")
	if err != nil {
		log.Println(err)
	}
}

func chunkedHandler(ctx *gin.Context) {
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ctx.Writer.Header().Set("Content-Type", "text")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(`
<html>
	<body>
`))
	ctx.Writer.Flush()
	for i := 0; i < 30; i++ {
		ctx.Writer.Write([]byte(fmt.Sprintf("<h2>%d</h2>", i)))
		ctx.Writer.Flush()
		time.Sleep(500 * time.Millisecond)
	}
	ctx.Writer.Write([]byte(`
	</body>
</html>
`))
	ctx.Writer.Flush()
}
