package id

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/facades"
)

/**
 * Rat算法
 * 格式：日期（去掉前2位）+ 当日已过的毫秒数（不足左补0） + 2位机器码 + 3位自增数
 * 例如：230403 00509559 00 000

 */

// RatID 结构体定义
type RatID struct {
	mu            sync.Mutex // 互斥锁
	nodeID        int        // 节点ID
	sequence      uint       // 序列号
	lastTimestamp uint       // 上次生成ID的时间戳
}

// NewRatID 创建一个RatID实例
func NewRatID() *RatID {
	return &RatID{
		nodeID:        facades.Config.GetInt("id.node"),
		sequence:      0,
		lastTimestamp: 0,
	}
}

// Generate 生成一个唯一ID
func (rat *RatID) Generate() (uint, error) {
	rat.mu.Lock()
	defer rat.mu.Unlock()

	timestamp := rat.timeGen()

	// 如果当前时间小于上次生成ID的时间戳，表示时间回退了，抛出错误
	if timestamp < rat.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards, rejecting requests until %d", rat.lastTimestamp)
	}

	second := carbon.Now().StartOfDay().DiffAbsInSeconds(carbon.Now())
	milli := carbon.Now().Millisecond()
	// 不足8位时分别在左边补0
	secondStr := fmt.Sprintf("%05d", second)
	milliStr := fmt.Sprintf("%03d", milli)

	// 如果是同一毫秒生成的ID，则自增序列号
	if timestamp == rat.lastTimestamp {
		rat.sequence += 1
		if rat.sequence > 999 {
			// 如果序列号溢出了，则等待下一毫秒再生成ID
			timestamp = rat.tilNextMillis(rat.lastTimestamp)
		}
	} else {
		rat.sequence = 0 // 不是同一毫秒生成的ID，序列号归零
	}

	rat.lastTimestamp = timestamp // 更新上次生成ID的时间戳

	id, err := strconv.ParseUint(carbon.Now().Format("ymd")+secondStr+milliStr+strconv.Itoa(rat.nodeID)+fmt.Sprintf("%03d", rat.sequence), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// timeGen 获取当前时间戳，单位：毫秒
func (rat *RatID) timeGen() uint {

	return uint(carbon.Now().TimestampMilli())
}

// tilNextMillis 等待下一毫秒
func (rat *RatID) tilNextMillis(lastTimestamp uint) uint {
	timestamp := rat.timeGen()
	for timestamp <= lastTimestamp {
		timestamp = rat.timeGen()
	}

	return timestamp
}
