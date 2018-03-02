package grpool

import "github.com/silentred/toolkit/util/timeutil"

// 开始任务
func (j *PoolJob) start() {
	go func() {
		for {
			if f := <-j.job; f != nil {
				// 执行任务
				f()
				// 更新活动时间
				j.update = timeutil.Second()
				// 执行完毕后添加到空闲队列
				if !j.pool.addJob(j) {
					break
				}
			} else {
				break
			}
		}
	}()
}

// 关闭当前任务
func (j *PoolJob) stop() {
	j.setJob(nil)
}

// 设置当前任务的执行函数
func (j *PoolJob) setJob(f func()) {
	j.job <- f
}
