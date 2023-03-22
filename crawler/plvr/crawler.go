package plvr

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Walker088/gorealestate/database"
)

const (
	apiUrl    = "https://plvr.land.moi.gov.tw/DownloadSeason?season=%sS%s&type=zip&fileName=lvr_landcsv.zip"
	startDate = "2013-01-01"
	storePath = "downloaded/plvr"
	storeName = "lvr_landcsv.zip"
)

type PlvrCrawler struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	workingDir string
	logger     *zap.SugaredLogger
	pool       *database.PgPool
	ResultsCh  chan string
}

func New(ctx context.Context, cancel context.CancelFunc, rootDir string, logger *zap.SugaredLogger, pool *database.PgPool) *PlvrCrawler {
	return &PlvrCrawler{
		ctx:        ctx,
		cancel:     cancel,
		workingDir: fmt.Sprintf("%s/%s", rootDir, storePath),
		logger:     logger,
		pool:       pool,
		ResultsCh:  make(chan string, 10),
	}
}

func (p *PlvrCrawler) Start() {
	p.logger.Infof("start crawlering %s", apiUrl)

	start, _ := time.Parse(time.DateOnly, startDate)
	today := time.Now()

	for yearSeason := start; yearSeason.Before(today); yearSeason = yearSeason.AddDate(0, 3, 0) {
		season := p.monthToSeason(int(yearSeason.Month()))
		query := fmt.Sprintf("%dS%d", p.commonToRocEra(yearSeason.Year()), season)

		zipPath := fmt.Sprintf("%s/%s/%s", p.workingDir, query, storeName)
		if exists, _ := p.fileExists(zipPath); !exists {
			p.wg.Add(1)
			go p.download(query)
		}
	}

	p.wg.Wait()
	close(p.ResultsCh)
	p.cancel()
}

func (p *PlvrCrawler) Stop() {
	p.logger.Infof("stop crawlering %s", apiUrl)
}

func (p *PlvrCrawler) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (p *PlvrCrawler) commonToRocEra(era int) int {
	return era - 1911
}

func (p *PlvrCrawler) monthToSeason(month int) int {
	if month < 4 {
		return 1
	}
	if month < 7 {
		return 2
	}
	if month < 10 {
		return 3
	}
	return 4
}

func (p *PlvrCrawler) download(yearSeason string) {
	defer p.wg.Done()

	r := rand.Intn(10)
	select {
	case <-p.ctx.Done():
		p.logger.Debugf("download terminated %s", yearSeason)
		return
	case <-time.After(time.Duration(r) * time.Second):

		p.ResultsCh <- yearSeason
	}
}

//func (p *PlvrCrawler) unzip() {}

//func (p *PlvrCrawler) parse() {}

//func (p *PlvrCrawler) export() {}
