package plvr

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/Walker088/gorealestate/database"
)

const (
	apiUrl    = "https://plvr.land.moi.gov.tw/DownloadSeason?season=%sS%s&type=zip&fileName=lvr_landcsv.zip"
	storePath = "downloaded/plvr"
	storeName = "lvr_landcsv.zip"
)

type PlvrCrawler struct {
	workingDir string
	logger     *zap.SugaredLogger
	pool       *database.PgPool
}

func New(workingDir string, logger *zap.SugaredLogger, pool *database.PgPool) *PlvrCrawler {
	return &PlvrCrawler{
		workingDir: workingDir,
		logger:     logger,
		pool:       pool,
	}
}

func (p *PlvrCrawler) Run() {
	p.logger.Infof("Start crawlering %s", apiUrl)
	start, _ := time.Parse(time.DateOnly, "2013-01-01")
	today := time.Now()
	for yearSeason := start; yearSeason.Before(today); yearSeason = yearSeason.AddDate(0, 3, 0) {
		season := p.monthToSeason(int(yearSeason.Month()))
		query := fmt.Sprintf("%dS%d", p.commonToRocEra(yearSeason.Year()), season)

		storePath := fmt.Sprintf("%s/%s/%s", storePath, query, storeName)
		if exists, _ := p.fileExists(storePath); !exists {
			go p.download(query)
			time.Sleep(3 * time.Second)
		}
	}
}

func (p *PlvrCrawler) Stop() {
	p.logger.Infof("Stop crawlering %s", apiUrl)
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
	p.logger.Info(yearSeason)
}

//func (p *PlvrCrawler) unzip() {}

//func (p *PlvrCrawler) parse() {}

//func (p *PlvrCrawler) export() {}
