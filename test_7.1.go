// Тестируем базовую функциональность CommandQueue (добавление и извлечение команд из очереди) 
package main

import (
  "sync"
  "testing"
  "time"
)

type MockCommand struct {
  executed bool
}

func (m *MockCommand) Execute() {
  m.executed = true
}

func TestCommandQueueBasic(t *testing.T) {
  queue := NewCommandQueue()

  mockCmd := &MockCommand{}
  queue.Enqueue(mockCmd)

  dequeuedCmd := queue.Dequeue()
  dequeuedCmd.Execute()

  if !mockCmd.executed {
    t.Errorf("Expected command to be executed")
  }
}
// Проверяем, что команды выполняются в отдельной горутине
func TestCommandQueueMultithreaded(t *testing.T) {
  queue := NewCommandQueue()
  startCmd := &StartCommand{queue: queue}
  startCmd.Execute()

  var wg sync.WaitGroup

  wg.Add(1)
  go func() {
    mockCmd := &MockCommand{}
    queue.Enqueue(mockCmd)
    time.Sleep(100 * time.Millisecond) // Give some time for the command to be dequeued and executed
    if !mockCmd.executed {
      t.Errorf("Expected command to be executed")
    }
    wg.Done()
  }()

  wg.Wait()
}

func TestCommandQueueStop(t *testing.T) {
  queue := NewCommandQueue()

  stopCmd := &HardStopCommand{queue: queue}
  stopCmd.Execute()

  if !queue.HasStopSignal() {
    t.Errorf("Expected stop signal to be sent")
  }
}