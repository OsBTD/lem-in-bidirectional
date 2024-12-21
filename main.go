package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type farm struct {
	ants_number int
	rooms       map[string][]int
	start       map[string][]int
	end         map[string][]int
	links       map[string][]string
}

func main() {
	var myFarm farm
	myFarm.Read("test.txt")
	BiBFS(&myFarm) // Passing pointer to farm to avoid copying it
	fmt.Println("number of ants is : ", myFarm.ants_number)
	fmt.Println("rooms are : ", myFarm.rooms)
	fmt.Println("start is : ", myFarm.start)
	fmt.Println("end is : ", myFarm.end)
	fmt.Println("links are : ", myFarm.links)
	fmt.Println("adjacent is : ", Graph(myFarm))
}

func (myFarm *farm) Read(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Println("error reading", err)
	}
	content := strings.Split(string(bytes), "\n")

	myFarm.rooms = make(map[string][]int)
	myFarm.start = make(map[string][]int)
	myFarm.end = make(map[string][]int)
	myFarm.links = make(map[string][]string)

	var st, en int
	number, err := strconv.Atoi(content[0])
	if err != nil {
		log.Println("couldn't convert", err)
	}
	myFarm.ants_number = number

	for index := range content {
		if strings.TrimSpace(content[index]) == "##start" {
			st++
			if index+1 <= len(content)-1 {
				split := strings.Split(strings.TrimSpace(content[index+1]), " ")
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil && err2 == nil {
					myFarm.start[split[0]] = []int{x, y}
				}
			}
		} else if strings.TrimSpace(content[index]) == "##end" {
			en++
			if index+1 <= len(content)-1 {
				split := strings.Split(strings.TrimSpace(content[index+1]), " ")
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil && err2 == nil {
					myFarm.end[split[0]] = []int{x, y}
				}
			}
		} else if strings.Contains(content[index], "-") {
			split := strings.Split(strings.TrimSpace(content[index]), "-")
			if len(split) == 2 {
				myFarm.links[split[0]] = append(myFarm.links[split[0]], split[1])
			}
		} else if strings.Count(content[index], " ") == 2 {
			split := strings.Split(strings.TrimSpace(content[index]), " ")
			if len(split) == 3 {
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil || err2 == nil {
					myFarm.rooms[split[0]] = []int{x, y}
				}
			}
		} else if (strings.HasPrefix(strings.TrimSpace(content[index]), "#") || strings.HasPrefix(strings.TrimSpace(content[index]), "L")) && (strings.TrimSpace(content[index]) != "##start" && strings.TrimSpace(content[index]) != "##end") {
			continue
		}
	}
	if en != 1 || st != 1 {
		log.Println("rooms setup is incorrect", err)
	}
}

func Graph(farm farm) map[string][]string {
	adjacent := make(map[string][]string)
	for room := range farm.rooms {
		adjacent[room] = []string{}
	}
	for room, links := range farm.links {
		for _, link := range links {
			adjacent[room] = append(adjacent[room], link)
			adjacent[link] = append(adjacent[link], room)
		}
	}

	return adjacent
}

func BiBFS(myFarm *farm) {
	adjacent := Graph(*myFarm)
	var QueueStart, QueueEnd []string
	var start, end string
	VisitedStart := make(map[string]bool)
	VisitedEnd := make(map[string]bool)
	ParentsStart := make(map[string]string)
	ParentsEnd := make(map[string]string)

	// Initialize Start and End rooms
	for key := range myFarm.start {
		start = key
		QueueStart = append(QueueStart, start)
		VisitedStart[start] = true
	}
	for key := range myFarm.end {
		end = key
		QueueEnd = append(QueueEnd, end)
		VisitedEnd[end] = true
	}

	fmt.Println("\n=== Bi-Directional BFS Initialization ===")
	fmt.Println("Start room:", start)
	fmt.Println("End room:", end)

	stepCount := 1
	// Perform Bi-Directional BFS
	for len(QueueStart) > 0 && len(QueueEnd) > 0 {
		// BFS from start side
		if meetingRoom := bfsStep(adjacent, &QueueStart, VisitedStart, VisitedEnd, ParentsStart); meetingRoom != "" {
			// If we find a meeting point, reconstruct and print the path
			fmt.Println("\nFound a path from start to end!")
			printPath(meetingRoom, ParentsStart, ParentsEnd)
			return
		}

		// BFS from end side
		if meetingRoom := bfsStep(adjacent, &QueueEnd, VisitedEnd, VisitedStart, ParentsEnd); meetingRoom != "" {
			// If we find a meeting point, reconstruct and print the path
			fmt.Println("\nFound a path from end to start!")
			printPath(meetingRoom, ParentsEnd, ParentsStart)
			return
		}

		stepCount++
	}

	// If we finish and haven't found a path
	fmt.Println("\n=== No path found ===")
}

func bfsStep(adjacent map[string][]string, Queue *[]string, Visited, OppositeVisited map[string]bool, Parents map[string]string) string {
	// Process the front element in the queue
	current := (*Queue)[0]
	*Queue = (*Queue)[1:] // Remove the front element of the queue

	// Explore connected rooms
	for _, link := range adjacent[current] {
		if !Visited[link] {
			Visited[link] = true
			Parents[link] = current
			*Queue = append(*Queue, link)

			// If the opposite search has already visited this room, we found a meeting point
			if OppositeVisited[link] {
				fmt.Printf("!!! Found meeting point at room '%s' !!!\n", link)
				return link // Return the meeting room where both searches meet
			}
		}
	}

	return "" // Return an empty string if no meeting point is found
}

func printPath(meetingRoom string, ParentsStart, ParentsEnd map[string]string) {
	// Reconstruct path from start to meeting room
	pathStart := []string{meetingRoom}
	current := meetingRoom
	for ParentsStart[current] != "" {
		current = ParentsStart[current]
		pathStart = append([]string{current}, pathStart...)
	}

	// Reconstruct path from end to meeting room
	pathEnd := []string{}
	current = meetingRoom
	for ParentsEnd[current] != "" {
		current = ParentsEnd[current]
		pathEnd = append([]string{current}, pathEnd...)
	}

	pathEnd = append(pathEnd, current)

	// Combine both paths
	fullPath := append(pathStart, pathEnd[1:]...)
	fmt.Printf("\nFull path from start to end: %v\n", fullPath)
	fmt.Println(pathEnd, pathStart, meetingRoom)
}
