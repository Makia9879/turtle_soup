package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/jzero-io/jzero-contrib/dynamic_conf"
	"github.com/spf13/cobra"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"

	"turtle-soup/internal/config"
	"turtle-soup/internal/model/t_turtle_soup_stories"
	"turtle-soup/internal/svc"
)

var (
	storiesFile string
)

var storyCmd = &cobra.Command{
	Use:   "story",
	Short: "turtle-soup story",
	Long:  `turtle-soup story`,
}

var addStoryCmd = &cobra.Command{
	Use:   "add",
	Short: "turtle-soup story add",
	Long:  `turtle-soup story add`,
	RunE:  addStory,
}

func addStory(_ *cobra.Command, _ []string) error {
	ss, err := dynamic_conf.NewFsNotify(cfgFile, dynamic_conf.WithUseEnv(true))
	logx.Must(err)
	cc := configurator.MustNewConfigCenter[config.Config](configurator.Config{
		Type: "yaml",
	}, ss)
	c, err := cc.GetConfig()
	logx.Must(err)

	// set up logger
	if err := logx.SetUp(c.Log.LogConf); err != nil {
		logx.Must(err)
	}
	if c.Log.LogConf.Mode != "console" {
		logx.AddWriter(logx.NewWriter(os.Stdout))
	}

	// Read stories from JSON file
	data, err := os.ReadFile(storiesFile)
	if err != nil {
		return err
	}

	var stories []*t_turtle_soup_stories.TTurtleSoupStories
	if err := json.Unmarshal(data, &stories); err != nil {
		return err
	}

	ctx := context.Background()
	svcCtx := svc.NewServiceContext(cc)
	err = svcCtx.Model.TTurtleSoupStories.BulkInsert(ctx, nil, stories)
	logx.Must(err)

	return nil
}

func init() {
	addStoryCmd.Flags().StringVarP(&storiesFile, "file", "f", "./desc/story/stories.json", "path to stories JSON file")
	storyCmd.AddCommand(addStoryCmd)
	rootCmd.AddCommand(storyCmd)
}
