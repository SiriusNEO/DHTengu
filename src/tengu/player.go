package tengu

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

type MusicController struct {
	ctrl		*beep.Ctrl
	resampler 	*beep.Resampler
	volume  	*effects.Volume
}

func NewMusicController(sampleRate beep.SampleRate, streamer beep.StreamSeeker) *MusicController {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &MusicController{ctrl, resampler, volume}
}

func Play(path string, songName string, albumName string) {
	fileIO, err := os.Open(path)

	if err != nil {
		red.Println("Music Not Founded")
		return
	}

	streamer, format, err := mp3.Decode(fileIO)

	if err != nil {
		red.Println("MP3 Decode Error")
		return
	}

	controller := NewMusicController(format.SampleRate, streamer)

	err = termbox.Init()

	defer streamer.Close()
	defer termbox.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	magenta.Println("Tengu Player")
	magenta.Println("Now Playing: ", songName, " in ", albumName)
	magenta.Printf("\n")
	magenta.Println("[Left/Right]: Volume Up/Down   [Enter]: Pause/Continue   [Space]: Quit")
	magenta.Printf("\n")

	speaker.Play(controller.volume)

	EndSignal := false

	go func() {
		for !EndSignal {
			position := format.SampleRate.D(streamer.Position())
			length := format.SampleRate.D(streamer.Len())
			magenta.Printf("\r....----....----....	%v / %v", position.Round(time.Second), length.Round(time.Second))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	for {
		switch event := termbox.PollEvent(); event.Type {
			case termbox.EventKey: {
				switch event.Key {
				case termbox.KeySpace:
					{
						EndSignal = true
						time.Sleep(100 * time.Millisecond)
						return
					}
				case termbox.KeyEnter:
					{
						speaker.Lock()
						controller.ctrl.Paused = !controller.ctrl.Paused
						speaker.Unlock()
					}
				case termbox.KeyArrowLeft:
					{
						speaker.Lock()
						controller.volume.Volume -= 0.1
						speaker.Unlock()
					}
				case termbox.KeyArrowRight:
					{
						speaker.Lock()
						controller.volume.Volume += 0.1
						speaker.Unlock()
					}
				}
			}
		}
	}
}