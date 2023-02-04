package models

import (
	"errors"
	"fmt"
	"history-rewinded-regan/pb"
	"log"
	"math"
)

const maximumStackSize = math.MaxInt64

type IncidentsStack struct {
	incidents []pb.Incident	
	stackPointer int64
}

func NewIncidentsStack() IncidentsStack {
	return IncidentsStack{
		incidents: []pb.Incident{},
		stackPointer: -1,
	}
}

func (s *IncidentsStack) Pop() (pb.Incident, error) {
	log.Println(s.stackPointer)
	if s.stackPointer == - 1 {
		return pb.Incident{}, errors.New("stack is empty")
	}
	poppedElement := s.incidents[s.stackPointer]
	s.incidents[s.stackPointer] = pb.Incident{} 
	s.stackPointer = s.stackPointer - 1
	return poppedElement, nil
}

func (s *IncidentsStack) Push(incident pb.Incident) error {
	s.incidents = append(s.incidents, incident)
	s.stackPointer++
	return nil	
}

func (s *IncidentsStack) Peek(position int64) (pb.Incident, error) {
	if s.stackPointer < position - 1 {
		return pb.Incident{}, fmt.Errorf("no element present at %dth position", position)
	}
	return s.incidents[position - 1], nil
}

func (s *IncidentsStack) Top() (pb.Incident) {
	return s.incidents[s.stackPointer - 1]
}

func (s IncidentsStack) IsEmpty() bool {
	return s.stackPointer == -1
}

func (s *IncidentsStack) IsFull() bool {
	return s.stackPointer == maximumStackSize
}

func (s *IncidentsStack) Size() int64 {
	return s.stackPointer
}

func (s *IncidentsStack) Display()  {
	for i :=0; i < len(s.incidents); i++ {
		log.Println(s.incidents[i])
	}
}