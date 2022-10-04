package export

import (
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"sort"
	"strconv"
	"time"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/storage"
)

const minOffset = 4
const defaultPathToPattern = "rsl_pattern.xlsx"

type Exporter interface {
	Export(ctx context.Context, userID int64) (string, ActivityStat, error)
}

type ActivityStat struct {
	DaysFromFisrtStart int
	TotalDays          int
	Cb6Killed          int
	Cb5Killed          int
	Cb4Killed          int
	CbTotalKilled      int
	LegTome            int
	Sacred             int
	EpicTome           int
	Ancient            int
	Void               int
}

type ExcelExporter struct {
	cbStatStorage *storage.CbStatStorage
	pathToPattern string
}

type sheetInfo struct {
	sheet  string
	filter func(stat entities.UserCbStat) bool
}

type itemInfo struct {
	col     string
	extract func(stat entities.UserCbStat) int
}

type monthData struct {
	month time.Month
	year  int
	stats []entities.UserCbStat
}

func (d *monthData) append(stat entities.UserCbStat) {
	d.stats = append(d.stats, stat)
}

func NewExcelExporter(cbStatStorage *storage.CbStatStorage) *ExcelExporter {
	pathToPattern := os.Getenv("EXCEL_PATTERN")
	if pathToPattern == "" {
		pathToPattern = defaultPathToPattern
	}
	return &ExcelExporter{cbStatStorage: cbStatStorage, pathToPattern: pathToPattern}
}

func (e *ExcelExporter) dateRange(stats []entities.UserCbStat) (time.Time, time.Time) {
	min := stats[0].RelatedTo
	max := stats[0].RelatedTo

	for _, stat := range stats {
		dt := stat.RelatedTo
		if dt.After(max) {
			max = dt
		}
		if dt.Before(min) {
			min = dt
		}
	}
	return min.Truncate(24 * time.Hour), max.Truncate(24 * time.Hour)
}

func (e *ExcelExporter) Export(ctx context.Context, userID int64) (string, ActivityStat, error) {
	activity := ActivityStat{}

	stats, err := e.cbStatStorage.LoadAll(ctx, userID)
	if err != nil {
		return "", activity, err
	}

	if len(stats) == 0 {
		return "", activity, nil
	}

	xlsx, err := e.generateXlsx(userID, stats)
	if err != nil {
		return "", ActivityStat{}, err
	}

	return xlsx, e.genActivityStat(stats), nil
}
func (e *ExcelExporter) generateXlsx(userID int64, stats []entities.UserCbStat) (string, error) {
	xls, err := excelize.OpenFile(e.pathToPattern)
	if err != nil {
		return "", nil
	}

	for _, sheetInfo := range []sheetInfo{
		{
			sheet: "4 КБ",
			filter: func(stat entities.UserCbStat) bool {
				return stat.Level == 4
			},
		},
		{
			sheet: "5 КБ",
			filter: func(stat entities.UserCbStat) bool {
				return stat.Level == 5
			},
		},
		{
			sheet: "6 КБ",
			filter: func(stat entities.UserCbStat) bool {
				return stat.Level == 6
			},
		},
		{
			sheet: "Всего",
			filter: func(stat entities.UserCbStat) bool {
				return true
			},
		},
	} {
		activeStats := make([]entities.UserCbStat, 0)
		for _, stat := range stats {
			if sheetInfo.filter(stat) {
				activeStats = append(activeStats, stat)
			}
		}
		if len(activeStats) == 0 {
			xls.DeleteSheet(sheetInfo.sheet)
			continue
		}

		minDt, maxDt := e.dateRange(activeStats)

		err := e.fillDates(xls, sheetInfo.sheet, minDt, maxDt)
		if err != nil {
			return "", nil
		}

		for _, ii := range []itemInfo{
			{
				col: "B",
				extract: func(stat entities.UserCbStat) int {
					return stat.SacredShard
				},
			},
			{
				col: "C",
				extract: func(stat entities.UserCbStat) int {
					return stat.VoidShard
				},
			},
			{
				col: "D",
				extract: func(stat entities.UserCbStat) int {
					return stat.AncientShard
				},
			},
			{
				col: "E",
				extract: func(stat entities.UserCbStat) int {
					return stat.LegTome
				},
			},
			{
				col: "F",
				extract: func(stat entities.UserCbStat) int {
					return stat.EpicTome
				},
			},
		} {
			err := e.fillData(xls, sheetInfo.sheet, activeStats, minDt, maxDt, ii.col, ii.extract)
			if err != nil {
				return "", nil
			}
		}
	}

	xls.SetActiveSheet(xls.GetSheetIndex("Всего"))
	fn := fmt.Sprintf("tmp/stat_%d.xlsx", userID)
	err = xls.SaveAs(fn)
	if err != nil {
		return "", nil
	}
	return fn, nil
}

func (e *ExcelExporter) splitData(data []entities.UserCbStat) []*monthData {
	splitted := make(map[int]*monthData)
	md := make([]*monthData, 0)
	for _, stat := range data {
		month := stat.RelatedTo.Month()
		year := stat.RelatedTo.Year()
		key := year*100 + int(month)
		if _, ok := splitted[year]; !ok {
			splitted[key] = newMonthData(year, month)
			md = append(md, splitted[key])
		}
		splitted[key].append(stat)
	}

	sort.SliceStable(md, func(i, j int) bool {
		if md[i].year == md[j].year {
			return md[i].month < md[j].month
		}
		return md[i].year < md[j].year
	})
	return md
}

func (e *ExcelExporter) fillDates(xls *excelize.File, sheet string, min time.Time, max time.Time) error {
	cur := min
	offset := minOffset
	for !cur.After(max) {
		err := xls.SetCellValue(sheet, e.axis("A", offset), cur)
		if err != nil {
			return err
		}
		offset++
		cur = cur.Add(24 * time.Hour)
	}
	return nil
}

func (e *ExcelExporter) fillData(xls *excelize.File, sheet string, stats []entities.UserCbStat, minDt time.Time, maxDt time.Time, col string, extractData func(stat entities.UserCbStat) int) error {
	minDt = minDt.Truncate(24 * time.Hour)
	maxDt = maxDt.Truncate(24 * time.Hour)

	for _, stat := range stats {
		statTime := stat.RelatedTo.Truncate(24 * time.Hour)
		diff := statTime.Sub(minDt)
		offset := minOffset + int(diff.Hours()/24)
		strVal, err := xls.GetCellValue(sheet, e.axis(col, offset))
		if err != nil {
			return err
		}

		if strVal == "" {
			strVal = "0"
		}

		intVal, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return err
		}

		intVal += int64(extractData(stat))

		err = xls.SetCellInt(sheet, e.axis(col, offset), int(intVal))

		if err != nil {
			return err
		}
	}

	l := int(maxDt.Sub(minDt).Hours() / 24)
	//formulaType := excelize.STCellFormulaTypeDataTable
	err := xls.SetCellFormula(
		sheet,
		e.axis(col, minOffset-1),
		//"=SUM(D6:D11)",
		fmt.Sprintf("=SUM(%s:%s)", e.axis(col, minOffset), e.axis(col, minOffset+l)),
		//excelize.FormulaOpts{Type: &formulaType},
	)
	print(fmt.Sprintf("=SUM(%s:%s);", e.axis(col, minOffset), e.axis(col, minOffset+l)))
	if err != nil {
		return err
	}
	return xls.UpdateLinkedValue()
}

func (e *ExcelExporter) axis(col string, row int) string {
	return fmt.Sprintf("%s%d", col, row)
}

func newMonthData(year int, month time.Month) *monthData {
	return &monthData{
		month: month,
		year:  year,
		stats: make([]entities.UserCbStat, 0),
	}
}

func (e *ExcelExporter) genActivityStat(stats []entities.UserCbStat) ActivityStat {
	activity := ActivityStat{}

	days := make(map[time.Time]interface{})
	for _, stat := range stats {
		days[stat.RelatedTo.Truncate(24*time.Hour)] = true
		if stat.Level == 4 {
			activity.Cb4Killed++
		}
		if stat.Level == 5 {
			activity.Cb5Killed++
		}
		if stat.Level == 6 {
			activity.Cb6Killed++
		}
		activity.LegTome += stat.LegTome
		activity.EpicTome += stat.EpicTome
		activity.Sacred += stat.SacredShard
		activity.Void += stat.VoidShard
		activity.Ancient += stat.AncientShard
	}

	activity.TotalDays = len(days)
	activity.CbTotalKilled = activity.Cb4Killed + activity.Cb5Killed + activity.Cb6Killed
	minDt, _ := e.dateRange(stats)
	activity.DaysFromFisrtStart = int(time.Now().Sub(minDt).Hours() / 24)

	return activity
}

func (a ActivityStat) IsActive(rate float64) bool {
	return (float64(a.TotalDays) / (float64(a.DaysFromFisrtStart) + 0.1)) > rate
}
