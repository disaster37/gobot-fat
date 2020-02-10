package models

// FAT describe the current state of pond filter
type FAT struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	IsWashed    bool   `json:"is_running"`
	IsStarted    bool   `json:"is_started"`
	IsStopped    bool   `json:"is_stopped"`
	IsSecurity   bool   `json:"is_security"`
	IsEmergencyStopped bool   `json:"is_emmergency"`
}


func NewFATState() *FAT {
	return &FAT{}
}