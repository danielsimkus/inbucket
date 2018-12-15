package webui

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jhillyerd/inbucket/pkg/config"
	"github.com/jhillyerd/inbucket/pkg/server/web"
)

// RootGreeting serves the Inbucket greeting.
func RootGreeting(w http.ResponseWriter, req *http.Request, ctx *web.Context) (err error) {
	greeting, err := ioutil.ReadFile(ctx.RootConfig.Web.GreetingFile)
	if err != nil {
		return fmt.Errorf("Failed to load greeting: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(greeting)
	return err
}

// RootStatus renders portions of the server configuration as JSON.
func RootStatus(w http.ResponseWriter, req *http.Request, ctx *web.Context) (err error) {
	root := ctx.RootConfig
	retPeriod := ""
	if root.Storage.RetentionPeriod > 0 {
		retPeriod = root.Storage.RetentionPeriod.String()
	}

	return web.RenderJSON(w,
		&jsonServerConfig{
			Version:      config.Version,
			BuildDate:    config.BuildDate,
			POP3Listener: root.POP3.Addr,
			WebListener:  root.Web.Addr,
			SMTPConfig: jsonSMTPConfig{
				Addr:           root.SMTP.Addr,
				DefaultAccept:  root.SMTP.DefaultAccept,
				AcceptDomains:  root.SMTP.AcceptDomains,
				RejectDomains:  root.SMTP.RejectDomains,
				DefaultStore:   root.SMTP.DefaultStore,
				StoreDomains:   root.SMTP.StoreDomains,
				DiscardDomains: root.SMTP.DiscardDomains,
			},
			StorageConfig: jsonStorageConfig{
				MailboxMsgCap:   root.Storage.MailboxMsgCap,
				StoreType:       root.Storage.Type,
				RetentionPeriod: retPeriod,
			},
		})
}
