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
	"regexp"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	e "github.com/Walker088/gorealestate/error"
	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	HttpStatusError           = "PV00001"
	CreateZipFileError        = "PV00002"
	CopyZipContentToFileError = "PV00003"
	OpenZippedFileError       = "PV00004"
	ReadZippedFileError       = "PV00005"
	CheckRecordExistsError    = "PV00006"
	HttpRequestError          = "PV00007"
	ReadZipFileFromLocalError = "PV00008"
	CreateZipReaderError      = "PV00009"
	UnmarshalCsvError         = "PV00010"

	currentPackage = "github.com/Walker088/gorealestate/crawler/plvr"
	apiUrl         = "https://plvr.land.moi.gov.tw/DownloadSeason?season=%s&type=zip&fileName=lvr_landcsv.zip"
	startDate      = "2013-01-01"
	storePath      = "downloaded/plvr"
	storeName      = "lvr_landcsv.zip"
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
		ResultsCh:  make(chan string),
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
		yearSeason := fmt.Sprintf("%dS%d", commonToRocEra(yearSeason.Year()), season)

		os.MkdirAll(fmt.Sprintf("%s/%s", p.workingDir, yearSeason), 0644)
		zipFilePath := fmt.Sprintf("%s/%s/%s", p.workingDir, yearSeason, storeName)

		p.wg.Add(1)
		go p.crawl(yearSeason, zipFilePath)
		break
	}

	p.wg.Wait()
	p.cancel()
}

func (p *PlvrCrawler) Stop() {
	p.logger.Infof("stop crawlering %s", apiUrl)
}

func (p *PlvrCrawler) crawl(yearSeason string, zipFilePath string) {
	recordExists := func(yearSeason string) (bool, error) {
		var exists bool
		query := `
		SELECT FALSE
		`
		err := p.pool.QueryRow(context.Background(), query /*, yearSeason*/).Scan(&exists)
		if err != nil {
			return false, err
		}
		return exists, nil
	}

	defer p.wg.Done()

	r := rand.Intn(10)
	select {
	case <-p.ctx.Done():
		p.logger.Debugf("download terminated %s", yearSeason)
		return
	case <-time.After(time.Duration(r) * time.Second):
		hasRecord, err := recordExists(yearSeason)
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
			zipReader, errorData := p.readZipFile(yearSeason, zipFilePath)
			if errorData != nil {
				p.ErrorsCh <- errorData
				return
			}
			if err := p.exportZipToDb(zipReader); err != nil {
				p.ErrorsCh <- err
				return
			}
			p.ResultsCh <- yearSeason
		}
	}
}

func (p *PlvrCrawler) readZipFile(yearSeason string, zipFilePath string) (*zip.Reader, *e.ErrorData) {
	fileExists := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	exists, _ := fileExists(zipFilePath)
	if exists {
		p.logger.Debugf("found downloaded zip file %s", zipFilePath)
		zipBytes, err := os.ReadFile(zipFilePath)
		if err != nil {
			return nil, e.NewErrorData(
				ReadZipFileFromLocalError,
				err.Error(),
				fmt.Sprintf("%s.readZipFile", currentPackage),
				nil,
				nil,
			)
		}
		zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
		if err != nil {
			return nil, e.NewErrorData(
				CreateZipReaderError,
				err.Error(),
				fmt.Sprintf("%s.readZipFile", currentPackage),
				nil,
				nil,
			)
		}
		return zipReader, nil
	} else {
		remoteZip := fmt.Sprintf(apiUrl, yearSeason)
		p.logger.Debugf("zip file not found, trying to download it from %s", remoteZip)
		resp, err := http.Get(remoteZip)
		if err != nil {
			return nil, e.NewErrorData(
				HttpRequestError,
				err.Error(),
				fmt.Sprintf("%s.readZipFile", currentPackage),
				nil,
				nil,
			)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			var inner interface{} = fmt.Sprintf("[%d] body %s", resp.StatusCode, string(body))
			return nil, e.NewErrorData(
				HttpStatusError,
				err.Error(),
				fmt.Sprintf("%s.readZipFile", currentPackage),
				nil,
				&inner,
			)
		}
		zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			return nil, e.NewErrorData(
				CopyZipContentToFileError,
				err.Error(),
				fmt.Sprintf("%s.readZipFile", currentPackage),
				nil,
				nil,
			)
		}
		return zipReader, nil
	}
}

func (p *PlvrCrawler) exportZipToDb(zip *zip.Reader) *e.ErrorData {
	r := regexp.MustCompile(`^[a-z]_lvr_land_[a-c]\.csv$`)
	for _, zf := range zip.File {
		fileName := zf.FileHeader.Name
		if !r.MatchString(fileName) {
			p.logger.Debugf("file %s is omitted", fileName)
			continue
		}
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
		head := strings.Split(string(content), "\n")[1:2]
		body := strings.Split(string(content), "\n")[2:]
		csvBytes := head[0] + "\n" + strings.Join(body, "\n")

		p.logger.Debugf("Opened file %s", fileName)
		p.parse([]byte(csvBytes))
	}
	return nil
}

func (p *PlvrCrawler) parse(csvBytes []byte) ([]*RealEstateItem, *e.ErrorData) {
	items := []*RealEstateItem{}

	if err := gocsv.UnmarshalBytes(csvBytes, &items); err != nil {
		return nil, e.NewErrorData(
			UnmarshalCsvError,
			err.Error(),
			fmt.Sprintf("%s.parse", currentPackage),
			nil,
			nil,
		)
	}

	return items, nil
}
