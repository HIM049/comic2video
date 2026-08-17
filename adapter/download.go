package adapter

import (
	"comic2video/config"
	"sync"

	"github.com/charmbracelet/log"
)

func DownloadImages(urls []string, download func(i int, url string) (string, error)) ([]string, error) {
	threads := max(config.Cfg.Service.DownloadThreads, 1)

	sem := make(chan struct{}, threads)
	errCh := make(chan error, len(urls))
	paths := make([]string, len(urls))
	var wg sync.WaitGroup

	for i, u := range urls {
		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
				log.Debugf("Finish download image (%d)", i+1)
			}()

			log.Debugf("Start downloading image (%d)", i+1)
			p, err := download(i, u)
			if err != nil {
				errCh <- err
				return
			}
			paths[i] = p
		}()
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}
	return paths, nil
}
