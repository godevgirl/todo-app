package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Todo struct {
	ID     int
	Task   string
	IsDone bool
}

var todos []Todo
var nextID int = 1

const fileName = "todos.txt"

func main() {
	loadTasksFromFile()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nWELCOME TO YOUR DAILY PLANNER")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Delete Task")
		fmt.Println("4. Mark Task as Done")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		switch input {
		case "1":
			addTask(scanner)
		case "2":
			listTasks()
		case "3":
			deleteTask(scanner)
		case "4":
			markTaskDone(scanner)
		case "5":
			saveTasksToFile()
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func addTask(scanner *bufio.Scanner) {
	fmt.Print("Enter a new task: ")
	if !scanner.Scan() {
		return
	}
	task := scanner.Text()

	newTodo := Todo{
		ID:     nextID,
		Task:   task,
		IsDone: false,
	}
	nextID++
	todos = append(todos, newTodo)
	saveTasksToFile()
	fmt.Println("Task added!")
}

func listTasks() {
	fmt.Println("\nCurrent Tasks:")
	for _, todo := range todos {
		status := "Pending"
		if todo.IsDone {
			status = "Done"
		}
		fmt.Printf("ID: %d. Task: %s Status: [%s]\n", todo.ID, todo.Task, status)
	}
}

func deleteTask(scanner *bufio.Scanner) {
	fmt.Print("Enter task ID to delete: ")
	if !scanner.Scan() {
		return
	}
	idStr := scanner.Text()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTasksToFile()
			fmt.Println("Task deleted!")
			return
		}
	}
	fmt.Println("Task not found.")
}

func markTaskDone(scanner *bufio.Scanner) {
	fmt.Print("Enter task ID to mark as done: ")
	if !scanner.Scan() {
		return
	}
	idStr := scanner.Text()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].IsDone = true
			saveTasksToFile()
			fmt.Println("Task marked as done!")
			return
		}
	}
	fmt.Println("Task not found.")
}

func saveTasksToFile() {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, todo := range todos {
		line := fmt.Sprintf("%d|%s|%t\n", todo.ID, todo.Task, todo.IsDone)
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Println("Error writing task:", err)
			return
		}
	}
	writer.Flush()
}

func loadTasksFromFile() {
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return // File doesn't exist, so no tasks to load
		}
		fmt.Println("Error loading tasks:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue // skip invalid lines
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		task := parts[1]
		isDone, err := strconv.ParseBool(parts[2])
		if err != nil {
			continue
		}

		todo := Todo{
			ID:     id,
			Task:   task,
			IsDone: isDone,
		}
		todos = append(todos, todo)

		if id >= nextID {
			nextID = id + 1
		}
	}
}
