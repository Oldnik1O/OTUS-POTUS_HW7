//  Queue - это потокобезопасная очередь команд с условными переменными для сигнализации о наличии новых команд
// Run - метод, который выполняется в отдельной goroutine, и он вытягивает команды из очереди и выполняет их
// HardStop и SoftStop - методы для остановки выполнения очереди
// PrintCommand - это простой пример команды, которая просто печатает сообщение
// В функции main демонстрируется пример использования
go
package main

import (
  "fmt"
  "sync"
  "time"
)

type Command interface {
  Execute()
}

type Queue struct {
  commands  []Command
  mutex     sync.Mutex
  notEmpty  *sync.Cond
  stopHard  bool
  stopSoft  bool
  processing bool
}

func NewQueue() *Queue {
  q := &Queue{}
  q.notEmpty = sync.NewCond(&q.mutex)
  return q
}

func (q *Queue) AddCommand(c Command) {
  q.mutex.Lock()
  q.commands = append(q.commands, c)
  q.notEmpty.Signal()
  q.mutex.Unlock()
}

func (q *Queue) GetCommand() (Command, bool) {
  q.mutex.Lock()
  for len(q.commands) == 0 && !q.stopHard && !(q.stopSoft && !q.processing) {
    q.notEmpty.Wait()
  }
  if len(q.commands) == 0 || q.stopHard {
    q.mutex.Unlock()
    return nil, false
  }
  cmd := q.commands[0]
  q.commands = q.commands[1:]
  q.processing = true
  q.mutex.Unlock()
  return cmd, true
}

func (q *Queue) Run() {
  for {
    cmd, ok := q.GetCommand()
    if ok {
      cmd.Execute()
      q.mutex.Lock()
      q.processing = false
      q.mutex.Unlock()
    } else {
      return
    }
  }
}

func (q *Queue) HardStop() {
  q.mutex.Lock()
  q.stopHard = true
  q.notEmpty.Signal()
  q.mutex.Unlock()
}

func (q *Queue) SoftStop() {
  q.mutex.Lock()
  q.stopSoft = true
  q.notEmpty.Signal()
  q.mutex.Unlock()
}

type PrintCommand struct {
  msg string
}

func (p *PrintCommand) Execute() {
  fmt.Println(p.msg)
  time.Sleep(1 * time.Second)
}

func main() {
  q := NewQueue()
  go q.Run()

  q.AddCommand(&PrintCommand{"Hello 1"})
  q.AddCommand(&PrintCommand{"Hello 2"})
  q.AddCommand(&PrintCommand{"Hello 3"})

  time.Sleep(2 * time.Second)
  q.SoftStop()

  q.AddCommand(&PrintCommand{"Hello 4"})
  q.AddCommand(&PrintCommand{"Hello 5"})
  time.Sleep(10 * time.Second)
}