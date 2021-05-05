package start

import (
	"fmt"
	"minter-hub-sentinel/config"
	"minter-hub-sentinel/services/minterhub"
	"minter-hub-sentinel/services/telegram"
	"sync"
	"time"

	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Command struct {
	log    *logrus.Logger
	config *config.Config

	wg sync.WaitGroup

	lastBlockHeight int

	minterHub *minterhub.Service
	telegram  *tgbotapi.BotAPI
}

func New(log *logrus.Logger, config *config.Config) *Command {
	return &Command{
		log:    log,
		config: config,
	}
}

func (cmd *Command) Command() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start watching for missed blocks",
		Action: func(ctx *cli.Context) error {
			if len(cmd.config.MinterHub.Api) == 0 {
				return errors.New("define at least one api_url in configuration file")
			}

			if mh, err := minterhub.New(cmd.config.MinterHub.Api, cmd.log); err != nil {
				return err
			} else {
				cmd.minterHub = mh
			}

			if len(cmd.config.Telegram.Token) > 0 {
				b, err := telegram.New(cmd.config.Telegram.Token)

				if err != nil {
					return err
				}

				cmd.telegram = b
			} else {
				cmd.log.Warn("Telegram token not set. Notifications will not be sent")
			}

			lastBlock, err := cmd.minterHub.GetLatestBlock()

			if err != nil {
				return err
			}

			cmd.lastBlockHeight = lastBlock.Block.Header.Height

			return cmd.run()
		},
	}
}

func (cmd *Command) run() error {
	cmd.newLogEntry(cmd.lastBlockHeight).
		WithField("sleep", cmd.config.MinterHub.Sleep).
		Println("Watcher started")

	ticker := time.NewTicker(time.Duration(cmd.config.MinterHub.Sleep) * time.Second)

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				nextHeight := cmd.lastBlockHeight + 1

				signed, err := cmd.isSigned(nextHeight)

				if err != nil {
					if _, ok := err.(*minterhub.BlockNotFound); ok {
						cmd.newLogEntry(nextHeight).Debugln("Block not created yet.")
						continue
					}

					go cmd.sendBotMessage(fmt.Sprintf("⚠️ Failed to detect if block is signed: %s", err))

					quit <- true

					return
				}

				cmd.lastBlockHeight = nextHeight

				if signed {
					cmd.newLogEntry(nextHeight).Println("Block signed")

					continue
				}

				cmd.newLogEntry(nextHeight).Warnln("Block missed")

				go cmd.sendBotMessage(fmt.Sprintf("⚠️ Block %d missed", nextHeight))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit

	cmd.wg.Wait()

	return nil
}

func (cmd *Command) isSigned(height int) (bool, error) {
	block, err := cmd.minterHub.GetBlock(height)

	if err != nil {
		return false, err
	}

	for _, s := range block.Block.LastCommit.Signatures {
		if s.ValidatorAddress == cmd.config.MinterHub.ValidatorAddress && len(s.Signature) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (cmd *Command) sendBotMessage(message string) {
	if cmd.telegram == nil {
		return
	}

	cmd.wg.Add(1)

	for _, id := range cmd.config.Telegram.Admins {
		msg := tgbotapi.NewMessage(int64(id), message)

		_, _ = cmd.telegram.Send(msg)
	}

	cmd.wg.Done()
}

func (cmd *Command) newLogEntry(height int) *logrus.Entry {
	return cmd.log.WithField("height", height)
}
