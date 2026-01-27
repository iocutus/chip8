package main

type ChipContext struct {
	Memory      [4096]uint8
	V           [16]uint8
	Stack       [16]uint16
	I           uint16
	PC          uint16
	SP          uint16
	DelayReg    uint16
	SoundReg    uint16
	FrameBuffer [ScreenWidth * ScreenHeight]bool
}
