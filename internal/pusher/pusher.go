package pusher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sync"

	"github.com/ip75/rutubeuploader/internal/log"
	"github.com/ip75/rutubeuploader/internal/rutube"
	"github.com/ip75/rutubeuploader/internal/token"
)

const DefaultTasksCount = 5

type Pusher struct {
	TasksCount        int
	pool              chan rutube.Video
	once              sync.Once
	waitForCompletion bool
	chain             []rutube.Video
	token             token.Token
}

func New(tCnt int, wfc bool, j io.Reader, tkf string) (*Pusher, error) {
	if tCnt == 0 {
		tCnt = DefaultTasksCount
	}

	p := Pusher{
		TasksCount:        tCnt,
		pool:              make(chan rutube.Video, tCnt),
		waitForCompletion: wfc,
	}

	if err := p.token.LoadToken(tkf); err != nil {
		return nil, fmt.Errorf("load token from file: %w", err)
	}

	if err := json.NewDecoder(j).Decode(&p.chain); err != nil {
		return nil, fmt.Errorf("unmarshal json with video: %w", err)
	}

	return &p, nil
}

func (p *Pusher) Run() {
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	defer close(errCh)
	finish := make(chan struct{})
	defer close(finish)
	c := make(chan os.Signal, 1)

	// run thread that fills queue
	go func() {
		for _, v := range p.chain {
			if _, err := url.ParseRequestURI(v.Url); err != nil {
				log.Logger.Error().Err(err).Str("url", v.Url).Msgf("parse url")
				continue
			}
			p.pool <- v
		}
		finish <- struct{}{}
	}()

	for {
		select {
		case v := <-p.pool:
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Logger.Info().Str("url", v.Url).Msgf("start uploading file to rutube")
				if err := v.Upload(p.token.Token, p.waitForCompletion); err != nil {
					errCh <- err
				}
			}()
		case err := <-errCh:
			log.Logger.Err(err).Msg("send video to upload")
		case s := <-c:
			log.Logger.Info().Msgf("got a signal: %s", s)
			p.closePool()
			return
		case <-finish:
			log.Logger.Info().Msgf("all URLs are set to processing. Let's stop the loop")
			wg.Wait()
			p.closePool()
			return
		}
	}
}

func (p *Pusher) closePool() {
	p.once.Do(func() {
		close(p.pool)
	})
}
