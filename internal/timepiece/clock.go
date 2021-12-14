package timepiece

import "github.com/rahvin74/spherical-timepiece/internal/timesphere"

// Tick represents a tick through the clock. It holds the next ball and
// deposits it in the right place as defined by the mechanism's logic
type Tick struct {
	NextBall        *timesphere.MinuteBall
	MinuteCount     int
	FiveMinuteCount int
	HourCount       int
}

// Mechanism represents the physical clock. It holds the balls and the tracks.
// It handles the movement of the balls along the tracks with threads that keep
// track of each track.
type Mechanism struct {
	ballBucket   []timesphere.MinuteBall
	minuteTrack  chan Tick
	fiveTrack    chan Tick
	hourTrack    chan Tick
	controlTrack chan Tick
	done         chan bool

	ballCount int
	timeCount int
}

// NewMechanism returns a newly initialized Mechanism objecct
func NewMechanism() Mechanism {
	clock := Mechanism{
		minuteTrack:  make(chan Tick),
		fiveTrack:    make(chan Tick),
		hourTrack:    make(chan Tick),
		controlTrack: make(chan Tick),
		done:         make(chan bool),
		timeCount:    0,
	}

	return clock
}

// Run takes as input a list of MinuteBalls, the runs the clock Mechanism
func (m *Mechanism) Run(balls []timesphere.MinuteBall) int {
	m.ballBucket = balls
	m.ballCount = len(balls)

	// start the go routines to manage the tracks
	go m.minute()
	go m.fiveMinute()
	go m.hour()
	go m.control()

	tick := Tick{}

	// send the first Tick to get things moving
	m.controlTrack <- tick

	// wait until the clock signals it's done.
	// this will only happen once the balls return
	// to their original order
	<-m.done

	return m.timeCount / 1440
}

// manages the minute track. listens on the minute channel for ticks
func (m *Mechanism) minute() {
	minuteBucket := make([]timesphere.MinuteBall, 0, 0)
	for {
		inTick := <-m.minuteTrack
		// if we're at the max for the track, dump the bucket for this track and
		// send the next ball on the the next track with the tick
		if len(minuteBucket) == 4 {
			reversed := reverse(minuteBucket)
			m.ballBucket = append(m.ballBucket, reversed...)
			minuteBucket = make([]timesphere.MinuteBall, 0, 0)
			inTick.MinuteCount = 0
			m.fiveTrack <- inTick
		} else {
			minuteBucket = append(minuteBucket, *inTick.NextBall)
			inTick.NextBall = nil
			inTick.MinuteCount++

			// we don't need to move on to the next track, we can
			// return the tick to the control track now
			m.controlTrack <- inTick
		}
	}
}

// manages the five minute track. listens on the five minute channel for ticks
func (m *Mechanism) fiveMinute() {
	fiveBucket := make([]timesphere.MinuteBall, 0, 0)
	for {
		inTick := <-m.fiveTrack

		// if we're at the max for the track, dump the bucket for this track and
		// send the next ball on the the next track with the tick
		if len(fiveBucket) == 11 {
			reversed := reverse(fiveBucket)
			m.ballBucket = append(m.ballBucket, reversed...)
			fiveBucket = make([]timesphere.MinuteBall, 0, 0)
			inTick.FiveMinuteCount = 0
			m.hourTrack <- inTick
		} else {
			fiveBucket = append(fiveBucket, *inTick.NextBall)
			inTick.NextBall = nil
			inTick.FiveMinuteCount++

			// we don't need to move on to the next track, we can
			// return the tick to the control track now
			m.controlTrack <- inTick
		}
	}
}

// manages the hour track. listens on the hour channel for ticks
func (m *Mechanism) hour() {
	hourBucket := make([]timesphere.MinuteBall, 0, 0)
	for {
		inTick := <-m.hourTrack

		// if we're at the max for the track, dump the bucket for this track and
		// send the next ball on the the next track with the tick
		// then dump the incoming ball back onto the ballBucket
		if len(hourBucket) == 11 {
			reversed := reverse(hourBucket)
			m.ballBucket = append(m.ballBucket, reversed...)
			m.ballBucket = append(m.ballBucket, *inTick.NextBall)
			hourBucket = make([]timesphere.MinuteBall, 0, 0)
			inTick.NextBall = nil
			inTick.HourCount = 0
		} else {
			hourBucket = append(hourBucket, *inTick.NextBall)
			inTick.NextBall = nil
			inTick.HourCount++
		}

		m.controlTrack <- inTick
	}

}

// control listens for ticks, checks the ball bucket and sends the next ball.
// should only be called when the mechanism has finished processing the most recent ball
func (m *Mechanism) control() {
	firstTime := true // the balls will be in order on the first run... don't want to quit early
	for {
		inTick := <-m.controlTrack
		if len(m.ballBucket) == m.ballCount && !firstTime {
			for i := range m.ballBucket {
				// if i and original position are different, we're not in order
				if i != m.ballBucket[i].GetOriginalPosition() {
					break
				}

				// if i is the same as the max length of the ball bucket (minus 1), we are in order
				if i == len(m.ballBucket)-1 {
					// signal the Run function to quit!
					m.done <- true
				}
			}
		}

		firstTime = false

		inTick.NextBall = &m.ballBucket[0] // next ball to run is at the front of the queue
		m.ballBucket = m.ballBucket[1:]    // gotta shrink the size of the queue because we just popped a ball off
		m.timeCount++
		m.minuteTrack <- inTick // send the next tick
	}
}

func reverse(spheres []timesphere.MinuteBall) []timesphere.MinuteBall {
	reversed := make([]timesphere.MinuteBall, 0, 0)
	for i := len(spheres) - 1; i >= 0; i-- {
		reversed = append(reversed, spheres[i])
	}

	return reversed
}
