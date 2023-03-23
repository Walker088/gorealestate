package plvr

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	e "github.com/Walker088/gorealestate/error"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	currentPackage            = "github.com/Walker088/gorealestate/crawler/plvr"
	HttpStatusError           = "PV00001"
	CreateZipFileError        = "PV00002"
	CopyZipContentToFileError = "PV00003"
	OpenZippedFileError       = "PV00004"
	ReadZippedFileError       = "PV00005"
	CheckRecordExistsError    = "PV00006"

	apiUrl    = "https://plvr.land.moi.gov.tw/DownloadSeason?season=%s&type=zip&fileName=lvr_landcsv.zip"
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
	pool       *pgxpool.Pool
	ResultsCh  chan string
	ErrorsCh   chan *e.ErrorData
}

func New(ctx context.Context, cancel context.CancelFunc, rootDir string, logger *zap.SugaredLogger, pool *pgxpool.Pool) *PlvrCrawler {
	return &PlvrCrawler{
		ctx:        ctx,
		cancel:     cancel,
		workingDir: fmt.Sprintf("%s/%s", rootDir, storePath),
		logger:     logger,
		pool:       pool,
		ResultsCh:  make(chan string, 10),
		ErrorsCh:   make(chan *e.ErrorData),
	}
}

func (p *PlvrCrawler) Start() {
	p.logger.Infof("start crawlering %s", apiUrl)
	monthToSeason := func(month int) int {
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
	commonToRocEra := func(era int) int {
		return era - 1911
	}

	start, _ := time.Parse(time.DateOnly, startDate)
	today := time.Now()

	for yearSeason := start; yearSeason.Before(today); yearSeason = yearSeason.AddDate(0, 3, 0) {
		season := monthToSeason(int(yearSeason.Month()))
		query := fmt.Sprintf("%dS%d", commonToRocEra(yearSeason.Year()), season)

		os.MkdirAll(fmt.Sprintf("%s/%s", p.workingDir, query), 0644)
		zipPath := fmt.Sprintf("%s/%s/%s", p.workingDir, query, storeName)

		p.wg.Add(1)
		go p.crawl(query, zipPath)
		break
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

func (p *PlvrCrawler) recordExists(yearSeason string) (bool, error) {
	var exists bool
	query := `
	
	`
	err := p.pool.QueryRow(context.Background(), query, yearSeason).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *PlvrCrawler) crawl(yearSeason string, zipFile string) {
	defer p.wg.Done()

	r := rand.Intn(10)
	select {
	case <-p.ctx.Done():
		p.logger.Debugf("download terminated %s", yearSeason)
		return
	case <-time.After(time.Duration(r) * time.Second):
		exists, err := p.fileExists(zipFile)
		if err != nil {
			p.ErrorsCh <- e.NewErrorData(
				HttpStatusError,
				err.Error(),
				fmt.Sprintf("%s.download", currentPackage),
				nil,
				nil,
			)
			return
		}
		resp, err := http.Get(fmt.Sprintf(apiUrl, yearSeason))
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			var inner interface{}
			inner = fmt.Sprintf("[%d] body %s", resp.StatusCode, string(body))
			p.ErrorsCh <- e.NewErrorData(
				HttpStatusError,
				err.Error(),
				fmt.Sprintf("%s.download", currentPackage),
				nil,
				&inner,
			)
			return
		}
		hasRecord, err := p.recordExists(yearSeason)
		if err != nil {
			p.ErrorsCh <- e.NewErrorData(
				CheckRecordExistsError,
				err.Error(),
				fmt.Sprintf("%s.download", currentPackage),
				nil,
				nil,
			)
			return
		}
		if !hasRecord {
			zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				p.ErrorsCh <- e.NewErrorData(
					CopyZipContentToFileError,
					err.Error(),
					fmt.Sprintf("%s.download", currentPackage),
					nil,
					nil,
				)
				return
			}
			if err := p.exportZipToDb(zipReader); err != nil {
				p.ErrorsCh <- err
				return
			}
		}
		p.ResultsCh <- yearSeason
	}
}

func (p *PlvrCrawler) exportZipToDb(zip *zip.Reader) *e.ErrorData {
	for _, zf := range zip.File {
		fileName := zf.FileHeader.Name
		f, err := zf.Open()
		if err != nil {
			return e.NewErrorData(
				OpenZippedFileError,
				err.Error(),
				fmt.Sprintf("%s.exportZipToDb", currentPackage),
				nil,
				nil,
			)
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return e.NewErrorData(
				ReadZippedFileError,
				err.Error(),
				fmt.Sprintf("%s.exportZipToDb", currentPackage),
				nil,
				nil,
			)
		}
	}
	return nil
}

//func (p *PlvrCrawler) parse() {}
