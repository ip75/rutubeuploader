package pusher

import (
	"net/url"
	"os"
	"sync"

	"github.com/ip75/rutubeuploader/internal/log"
	"github.com/ip75/rutubeuploader/internal/rutube"
)

const defaultTasksCount = 5

type Pusher struct {
	TasksCount        int `json:"tasks_count"`
	pool              chan rutube.Video
	once              sync.Once
	waitForCompletion bool
}

func New(tCnt int, wfc bool) *Pusher {
	if tCnt == 0 {
		tCnt = defaultTasksCount
	}

	return &Pusher{
		TasksCount:        tCnt,
		pool:              make(chan rutube.Video, tCnt),
		waitForCompletion: wfc,
	}
}

func (p *Pusher) Run(list []string) {
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	c := make(chan os.Signal, 1)
	finish := make(chan struct{})

	// run thread that fills queue
	go func() {
		for _, v := range list {
			if u, err := url.Parse(v); err == nil {
				p.pool <- rutube.Video{
					Url: u,
				}
			} else {
				log.Logger.Error().Err(err).Str("url", v).Msgf("parse url")
			}
		}
		finish <- struct{}{}
	}()

	for {
		select {
		case v := <-p.pool:
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Logger.Info().Str("url", v.Url.String()).Msgf("start uploading file to rutube")
				if err := v.Upload(p.waitForCompletion); err != nil {
					errCh <- err
				}
			}()
		case err := <-errCh:
			log.Logger.Err(err).Msg("send video to upload")
		case s := <-c:
			log.Logger.Info().Msgf("got a signal: %s", s.String())
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
