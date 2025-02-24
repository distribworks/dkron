package shell

import (
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	namespace = "dkron_job"
)

var (
	cpuUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cpu_usage",
		Help:      "CPU usage by job",
	},
		[]string{"job_name"})

	memUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mem_usage_bytes",
		Help:      "Current memory consumed by job",
	},
		[]string{"job_name"})

	jobExecutionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "execution_time_seconds",
		Help:      "Job Execution Time",
	},
		[]string{"job_name"})

	jobDoneCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "execution_done_count",
		Help:      "Job Execution Counter",
	},
		[]string{"job_name", "exit_code"})

	jobExitCode = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "exit_code",
		Help:      "Exit code of a job",
	},
		[]string{"job_name"})
)

func CollectProcessMetrics(jobname string, pid int, quit chan int) {
	start := time.Now()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case exitCode, ok := <-quit:
			if !ok {
				log.Println("Exit code received and quit channel closed.")
				return
			}
			exitCodeStr := strconv.Itoa(exitCode)
			cpuUsage.WithLabelValues(jobname).Set(0)
			memUsage.WithLabelValues(jobname).Set(0)
			jobExecutionTime.WithLabelValues(jobname).Set(0)
			jobDoneCount.WithLabelValues(jobname, exitCodeStr).Inc()
			jobExitCode.WithLabelValues(jobname).Set(float64(exitCode))
		case <-ticker.C:
			cpu, mem, err := GetTotalCPUMemUsage(pid)
			if err != nil {
				log.Printf("Error getting pid statistics: %v", err)
				return
			}
			cpuUsage.WithLabelValues(jobname).Set(cpu)
			memUsage.WithLabelValues(jobname).Set(mem)
			jobExecutionTime.WithLabelValues(jobname).Set(time.Since(start).Seconds())
		}
	}
}

func GetTotalCPUMemUsage(pid int) (float64, float64, error) {
	var totalCPU, totalMem float64

	parentProc, err := process.NewProcess(int32(pid))
	if err != nil {
		log.Printf("NewProcess err: %v", err)
		return totalCPU, totalMem, err
	}

	allProc := append(GetChildrenProcesses(parentProc), parentProc)

	for _, p := range allProc {
		cpu, err := p.Times()
		if err != nil {
			// log.Printf("p.Times() err: %v", err)
			continue
		}
		mem, err := p.MemoryInfo()
		if err != nil {
			// log.Printf("p.MemoryInfo() err: %v", err)
			continue
		}

		// log.Printf("Pid: %d, CPU: %f, Mem: %d", p.Pid, cpu.Total(), mem.RSS)

		totalCPU = totalCPU + cpu.Total()
		totalMem = totalMem + float64(mem.RSS)
	}

	return totalCPU, totalMem, nil
}

func GetChildrenProcesses(pp *process.Process) []*process.Process {
	var ret []*process.Process

	c, err := pp.Children()
	if err != nil || len(c) == 0 {
		// log.Printf("pp.Children() err: %v", err)
		return ret
	}

	for _, cc := range c {
		ret = append(ret, cc)

		cp := GetChildrenProcesses(cc)
		if len(cp) > 0 {
			ret = append(ret, cp...)
		}
	}
	return ret
}
