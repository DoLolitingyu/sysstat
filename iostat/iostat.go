package iostat

import (
	"errors"
	"math"
	"regexp"
	"time"
	
	"github.com/DoLolitingyu/sysstat/diskStat"
	"github.com/DoLolitingyu/sysstat/systemCall"

)

var LastRawStat = make(map[string]diskStat.DiskStat)
var partition = regexp.MustCompile(`\w.*[^-]\d`)

const (
	oneSecondInMilli = 1000
	BToKB            = 1024
)

//DiskStatDiff struct for all extended io stats storing just the difference between
//the current and the last.
type DiskStatDiff struct {
	Id              int64
	PartId          int64
	Device          string
	ReadsCompleted  float64
	ReadsMerged     float64
	SectorsRead     float64
	MillisReading   float64
	WritesCompleted float64
	WritesMerged    float64
	SectorsWrite    float64
	MillisWriting   float64
	// IoInProgress float64 //not used and calculated diff then all others
	MillisDoingIo         float64
	WeightedMillisDoingIo float64
	RecordTime            float64
	IoTotal               float64
	SectorsTotalRaw       float64
}


func IsPartition(device string) (r bool) {
	return partition.MatchString(device)
}

// GetData iostat implement in Golang.
func GetData() (diskStat.ExtendedIoStats, error) {
	prevStat := diskStat.GetDiskStat()
	time.Sleep(time.Second)
	stat := diskStat.GetDiskStat()

	eIoStat := diskStat.ExtendedIoStats{}
	for _, prePartition := range prevStat {
		// ignore partitions
		if IsPartition(prePartition.Device) {
			continue
		}
		partition, ok := stat[prePartition.Device]
		if !ok {
			continue
		}

		diffStat, err := getDiffDiskStat(&prePartition, &partition)
		if nil != err {
			return eIoStat, err
		}

		timeDiffMilli := getTimeDiffMilli(diffStat.RecordTime)
		eIoStat.Device = diffStat.Device
		eIoStat.ReadsMerged = getOneSecondAvg(diffStat.ReadsMerged, timeDiffMilli)
		eIoStat.WritesMerged = getOneSecondAvg(diffStat.WritesMerged, timeDiffMilli)
		eIoStat.Writes = getOneSecondAvg(diffStat.WritesCompleted, timeDiffMilli)
		eIoStat.Reads = getOneSecondAvg(diffStat.ReadsCompleted, timeDiffMilli)
		// should format round
		eIoStat.SectorsRead = getFloatRound(getOneSecondAvg(diffStat.SectorsRead, timeDiffMilli) / BToKB)
		eIoStat.SectorsWrite = getFloatRound(getOneSecondAvg(diffStat.SectorsWrite, timeDiffMilli) / BToKB)

		eIoStat.Arqsz = getAvgRequestSize(diffStat.SectorsTotalRaw, diffStat.IoTotal)
		eIoStat.AvgQueueSize = getAvgQueueSize(diffStat.WeightedMillisDoingIo, timeDiffMilli)
		eIoStat.Await = getAwait(diffStat.MillisWriting, diffStat.MillisReading, diffStat.IoTotal)
		eIoStat.RAwait = getSingleAwait(diffStat.ReadsCompleted, diffStat.MillisReading)
		eIoStat.WAwait = getSingleAwait(diffStat.WritesCompleted, diffStat.MillisWriting)

		eIoStat.Util = getUtilization(diffStat.MillisDoingIo, timeDiffMilli)
		eIoStat.Svctm = getAvgServiceTime(diffStat.IoTotal, timeDiffMilli, eIoStat.Util)
	}
	return eIoStat, nil
}


func getTimeDiffMilli(diff float64) float64 {
	r := diff / 1000000.0
	return getFloatRound(r)
}

func getOneSecondAvg(diff float64, time float64) float64 {
	r := float64(diff/time) * oneSecondInMilli
	return getFloatRound(r)
}

func getDiffDiskStat(old *diskStat.DiskStat, cur *diskStat.DiskStat) (r DiskStatDiff, err error) {
	r.Id = cur.Id
	r.PartId = cur.PartId
	r.Device = cur.Device
	r.ReadsCompleted, err = getDiff(old.ReadsCompleted, cur.ReadsCompleted)
	if nil != err {
		return
	}
	r.ReadsMerged, err = getDiff(old.ReadsCompleted, cur.ReadsCompleted)
	if nil != err {
		return
	}
	// SectorsRead uint64
	r.SectorsRead, err = getDiffUint64(old.SectorsRead, cur.SectorsRead)
	if nil != err {
		return
	}
	// MillisReading int64
	r.MillisReading, err = getDiff(old.MillisReading, cur.MillisReading)
	if nil != err {
		return
	}
	// WritesCompleted int64
	r.WritesCompleted, err = getDiff(old.WritesCompleted, cur.WritesCompleted)
	if nil != err {
		return
	}
	// WritesMerged int64
	r.WritesMerged, err = getDiff(old.WritesMerged, cur.WritesMerged)
	if nil != err {
		return
	}
	// SectorsWrite uint64
	r.SectorsWrite, err = getDiffUint64(old.SectorsWrite, cur.SectorsWrite)
	if nil != err {
		return
	}
	// MillisWriting int64
	r.MillisWriting, err = getDiff(old.MillisWriting, cur.MillisWriting)
	if nil != err {
		return
	}
	// IoInProgress int64 //this stat seems to have old large then cur most of the time???
	// r.IoInProgress, err = getDiff(old.IoInProgress, cur.IoInProgress);
	if nil != err {
		return
	}
	// MillisDoingIo int64
	r.MillisDoingIo, err = getDiff(old.MillisDoingIo, cur.MillisDoingIo)
	if nil != err {
		return
	}
	// WeightedMillisDoingIo 64
	r.WeightedMillisDoingIo, err = getDiff(old.WeightedMillisDoingIo, cur.WeightedMillisDoingIo)
	if nil != err {
		return
	}
	// RecordTime int64
	r.RecordTime, err = getDiff(old.RecordTime, cur.RecordTime)
	if nil != err {
		return
	}
	//    IoTotal int64
	r.IoTotal, err = getDiff(old.IoTotal, cur.IoTotal)
	if nil != err {
		return
	}
	// SectorsTotalRaw uint64
	r.SectorsTotalRaw, err = getDiffUint64(old.SectorsTotalRaw, cur.SectorsTotalRaw)
	if nil != err {
		return
	}
	return
}

//getDiff gets old and current stat.
func getDiff(old int64, cur int64) (r float64, err error) {
	if old > cur {
		err = errors.New("Old is newer then current... impressive!")
		return
	}
	r = float64(cur - old)
	return getFloatRound(r), nil
}

//getDiffUint64 gets old and current stat.
func getDiffUint64(old uint64, cur uint64) (r float64, err error) {
	if old > cur {
		err = errors.New("Old is newer then current... impressive!")
		return
	}
	r = float64(cur - old)
	return getFloatRound(r), nil
}

//getAvgRequestSize get the avg request size for a disk.
func getAvgRequestSize(diffSectorsTotalRaw float64, diffIoTotal float64) (r float64) {
	if diffIoTotal <= 0 {
		r = 0.00
		return
	}
	r = float64(diffSectorsTotalRaw) / float64(diffIoTotal)
	return getFloatRound(r)
}

//getAvgQueueSize get the avg queue size for a disk.
func getAvgQueueSize(diffWeightedMillisDoingIo float64, time float64) (r float64) {
	r = diffWeightedMillisDoingIo / time
	return getFloatRound(r)
}

//getAwait get average wait time for a disk.
func getAwait(diffMillisWriting float64, diffMillisReading float64, diffIoTotal float64) (r float64) {
	if diffIoTotal <= 0 {
		r = 0.00
		return
	}
	totalRW := diffMillisWriting + diffMillisReading
	r = totalRW / diffIoTotal
	return getFloatRound(r)

}

//getSingleAwait get single await time used for parsing read or write await times
func getSingleAwait(diffIo float64, diffMillis float64) (r float64) {
	if diffIo <= 0 {
		r = 0.00
		return
	}
	r = diffMillis / diffIo
	return getFloatRound(r)
}

//getAvgServiceTime get the avg service time for a disk.
func getAvgServiceTime(diffIoTotal float64, time float64, util float64) (r float64) {
	hz := systemCall.GetClockTicksPerSecond()
	tput := diffIoTotal * float64(hz) / time

	if tput <= 0 {
		r = 0.0
		return
	}
	r = util / tput
	return getFloatRound(r)
}

//getUtilization get a single disks utilization percentage
func getUtilization(diffMillisDoingIo float64, time float64) (r float64) {
	r = (float64(diffMillisDoingIo) / (time * 100) * 10.0) * oneSecondInMilli
	if r > 100.00 {
		r = 100.00
	}
	return getFloatRound(r)
}

func getFloatRound(r float64) float64 {
	return math.Round(r*100) / 100
}