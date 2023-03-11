package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// CLI args
	f, closeFile, err := openProcessingFile(os.Args...)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile()

	// Load and parse processes
	processes, err := loadProcesses(f)
	if err != nil {
		log.Fatal(err)
	}

	// First-come, first-serve scheduling
	FCFSSchedule(os.Stdout, "First-come, first-serve", processes)

	//SJFSchedule(os.Stdout, "Shortest-job-first", processes)
	//
	//SJFPrioritySchedule(os.Stdout, "Priority", processes)
	//
	//RRSchedule(os.Stdout, "Round-robin", processes)
}

func openProcessingFile(args ...string) (*os.File, func(), error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf("%w: must give a scheduling file to process", ErrInvalidArgs)
	}
	// Read in CSV process CSV file
	f, err := os.Open(args[1])
	if err != nil {
		return nil, nil, fmt.Errorf("%v: error opening scheduling file", err)
	}
	closeFn := func() {
		if err := f.Close(); err != nil {
			log.Fatalf("%v: error closing scheduling file", err)
		}
	}

	return f, closeFn, nil
}

type (
	Process struct {
		ProcessID     int64
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
	}
	TimeSlice struct {
		PID   int64
		Start int64
		Stop  int64
	}
)

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) { 
   
	var processes []Process
	for _, row := range rows {
		process := Process{}
		process.ProcessID, _ = strconv.ParseInt(row[0], 10, 64)
		process.ArrivalTime, _ = strconv.ParseInt(row[1], 10, 64)
		process.BurstDuration, _ = strconv.ParseInt(row[2], 10, 64)
		process.Priority, _ = strconv.ParseInt(row[3], 10, 64)
		processes = append(processes, process)
	}

    // Sort the processes by arrival time
    sort.Slice(processes, func(i, j int) bool {
        return processes[i].ArrivalTime < processes[j].ArrivalTime
    })

    // Initialize the ready queue, current time, and current process
    var readyQueue []Process
    currentTime := 0
    var currentProcess *Process

    // Loop until all processes have completed
    for len(processes) > 0 || len(readyQueue) > 0 || currentProcess != nil {
        // Move any arriving processes to the ready queue
        for len(processes) > 0 && processes[0].ArrivalTime == currentTime {
            readyQueue = append(readyQueue, processes[0])
            processes = processes[1:]
        }

        // If there is no current process, select the shortest process from the ready queue
        if currentProcess == nil && len(readyQueue) > 0 {
            shortestIndex := 0
            for i := 1; i < len(readyQueue); i++ {
                if readyQueue[i].RemainingTime < readyQueue[shortestIndex].RemainingTime {
                    shortestIndex = i
                }
            }
            currentProcess = &readyQueue[shortestIndex]
            readyQueue = append(readyQueue[:shortestIndex], readyQueue[shortestIndex+1:]...)
        }

        // If there is a current process, check if any shorter process has arrived
        if currentProcess != nil {
            for i := range processes {
                if processes[i].ArrivalTime > currentTime {
                    break
                }
                if processes[i].RemainingTime < currentProcess.RemainingTime {
                    readyQueue = append(readyQueue, *currentProcess)
                    *currentProcess = processes[i]
                    processes = append(processes[:i], processes[i+1:]...)
                    break
                }
            }
        }

        // Execute the current process for one time unit
        if currentProcess != nil {
            currentProcess.RemainingTime--
            if currentProcess.RemainingTime == 0 {
                currentProcess = nil
            }
        }

        // Increment the waiting time of processes in the ready queue
        for i := range readyQueue {
            readyQueue[i].Waiting
            Time++
        }
    }

    // Calculate the total waiting time and print the results
    var totalWaitingTime int
    for i := range processes {
        totalWaitingTime += processes[i].WaitingTime
    }
    
    averageWaitingTime := float64(totalWaitingTime) / float64(len(processes))
    fmt.Printf("Average waiting time: %.2f\n", averageWaitingTime)
}
//
//func SJFSchedule(w io.Writer, title string, processes []Process) { }
//
func RRSchedule(w io.Writer, title string, processes []Process) { 
	ar processes []Process
	for _, row := range rows {
		process := Process{}
		process.ProcessID, _ = strconv.ParseInt(row[0], 10, 64)
		process.ArrivalTime, _ = strconv.ParseInt(row[1], 10, 64)
		process.BurstDuration, _ = strconv.ParseInt(row[2], 10, 64)
		process.Priority, _ = strconv.ParseInt(row[3], 10, 64)
		processes = append(processes, process)
	}

	// Round Robin algorithm
	timeQuantum := int64(2) // Set the time quantum to 2
	currentTime := int64(0) // Start at time 0
	totalWaitTime := int64(0)

	for {
		done := true
		for _, process := range processes {
			if process.BurstDuration > 0 {
				done = false
				if process.ArrivalTime <= currentTime {
					fmt.Printf("Processing process %v at time %v with burst duration %v\n", process.ProcessID, currentTime, process.BurstDuration)
					if process.BurstDuration <= timeQuantum {
						// The process finishes within the time quantum
						currentTime += process.BurstDuration
						totalWaitTime += currentTime - process.ArrivalTime - process.BurstDuration
						process.BurstDuration = 0
					} else {
						// The process does not finish within the time quantum
						currentTime += timeQuantum
						process.BurstDuration -= timeQuantum
					}
				}
			}
		}
		if done {
			break
		}
	}

	// Print the average waiting time
	averageWaitTime := float64(totalWaitTime) / float64(len(processes))
	fmt.Printf("Average waiting time: %v\n", averageWaitTime)
}
}



//endregion

//region Output helpers

func outputTitle(w io.Writer, title string) {
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
	_, _ = fmt.Fprintln(w, strings.Repeat(" ", len(title)/2), title)
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
	_, _ = fmt.Fprintln(w, "Gantt schedule")
	_, _ = fmt.Fprint(w, "|")
	for i := range gantt {
		pid := fmt.Sprint(gantt[i].PID)
		padding := strings.Repeat(" ", (8-len(pid))/2)
		_, _ = fmt.Fprint(w, padding, pid, padding, "|")
	}
	_, _ = fmt.Fprintln(w)
	for i := range gantt {
		_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Start), "\t")
		if len(gantt)-1 == i {
			_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Stop))
		}
	}
	_, _ = fmt.Fprintf(w, "\n\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
	_, _ = fmt.Fprintln(w, "Schedule table")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
	table.AppendBulk(rows)
	table.SetFooter([]string{"", "", "", "",
		fmt.Sprintf("Average\n%.2f", wait),
		fmt.Sprintf("Average\n%.2f", turnaround),
		fmt.Sprintf("Throughput\n%.2f/t", throughput)})
	table.Render()
}

//endregion

//region Loading processes.

var ErrInvalidArgs = errors.New("invalid args")

func loadProcesses(r io.Reader) ([]Process, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: reading CSV", err)
	}

	processes := make([]Process, len(rows))
	for i := range rows {
		processes[i].ProcessID = mustStrToInt(rows[i][0])
		processes[i].BurstDuration = mustStrToInt(rows[i][1])
		processes[i].ArrivalTime = mustStrToInt(rows[i][2])
		if len(rows[i]) == 4 {
			processes[i].Priority = mustStrToInt(rows[i][3])
		}
	}

	return processes, nil
}

func mustStrToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return i
}

//endregion
