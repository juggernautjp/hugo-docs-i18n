/*
Copyright © 2022 juggernautjp <katsutoshi.harada@gmail.com>

Hugo i18n ListLangPage
*/

package doci18n

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/spf13/cobra"
)

const jsonTempalte = `{
	"path": %s
	"title": %s
	"date": %s
	"publishDate": %s
	"draft": %s
	"permalink": %s
},
`

// from github.com/gohugoio/hugo/commands/commandeer.go
// type commandeer struct

// from github.com/gohugoio/hugo/commands/helpers.go
// type flagsToConfigHandler interface {
// func newSystemError(a ...any) commandError {

// InitializeConfig() cite from github.com/gohugoio/hugo/commands/hugo.go
// func initializeConfig(mustHaveConfigFile, failOnInitErr, running bool,
// func (c *commandeer) createLogger(cfg config.Provider) (loggers.Logger, error) {

// buildSites() cite from github.com/gohugoio/hugo/commands/commands.go
// type commandsBuilder struct {
// type baseCmd struct {
// type baseBuilderCmd struct {


// buildSites() cite from github.com/gohugoio/hugo/commands/list.go
type listCmd struct {
	*baseBuilderCmd
}

func (lc *listCmd) buildSites(config map[string]any) (*hugolib.HugoSites, error) {
	cfgInit := func(c *commandeer) error {
		for key, value := range config {
			c.Set(key, value)
		}
		return nil
	}

	c, err := initializeConfig(true, true, false, &lc.hugoBuilderCommon, lc, cfgInit)
	if err != nil {
		return nil, err
	}

	sites, err := hugolib.NewHugoSites(*c.DepsCfg)
	if err != nil {
		return nil, newSystemError("Error creating sites", err)
	}

	if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		return nil, newSystemError("Error Processing Source Content", err)
	}

	return sites, nil
}

// List pages which page has the same language content
func ListLangPage (lang, outfn string) error {
	cc := &listCmd{}
	sites, err := cc.buildSites(map[string]any{"buildDrafts": true})
	if err != nil {
		return fmt.Errorf("Error building sites: %w", err)
	}

	for _, p := range sites.Pages() {
		if !p.IsPage() {
			continue
		}
		// skip if the language is not target one
		if p.Lang() != lang {
			continue
		}
		// write 
		err := writer.Write([]string{
			strings.TrimPrefix(p.File().Filename(), sites.WorkingDir+string(os.PathSeparator)),
			// p.Slug(),
			p.Title(),
			p.Date().Format(time.RFC3339),
			p.PublishDate().Format(time.RFC3339),
			strconv.FormatBool(p.Draft()),
			p.Permalink(),
		})
		if err != nil {
			return fmt.Errorf("Error writing list to stdout: %w", err)
		}
	}

	return nil
}
