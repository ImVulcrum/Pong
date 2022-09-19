package main

import (
	"buttons"
	"fmt"
	gfx "gfx2"
	"math"
	"math/rand"
	"sliders"
	"strconv"
	"sync"
	"time"
)

var left_paddle_x uint16 = 0
var left_paddle_y uint16 = 0
var right_paddle_x uint16 = 0
var right_paddle_y uint16 = 0
var lkn uint16 = 0
var lprssd uint8 = 0
var rkn uint16 = 0
var rprssd uint8 = 0
var ms_x uint16 = 0
var active int = -1
var m sync.Mutex
var start bool = false

func main() {
	
	var w_x uint16 = 950                	//window lenght --> the bigger the window the slower the ball speed
	var w_y uint16 = 680                	//window hight --> the bigger the window the slower the ball speed
	
	gfx.Fenster(w_x, w_y+70)
	gfx.SetzeFont("pong.ttf", 20)

	gfx.Stiftfarbe(0, 0, 0)
	gfx.Vollrechteck(0, 0, w_x, w_y+70)

	var s1, s2, s3, s4, s5, s6, s7, s8, s9, s10 sliders.Slider = *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New(), *sliders.New()
	var b1 buttons.Button = *buttons.New()
	b1.Draw(50, 600, 70, 30, "Start")

	list := [10]sliders.Slider{s1, s2, s3, s4, s5, s6, s7, s8, s9, s10}
	list[0].Draw(50, 70, 300, 20, 10, 2000, 600, "Speed", true)					//indicates how many pixels per secod the ball should move
	list[1].Draw(50, 110, 300, 20, 10, 255, 240, "Tail Length", true)			//indicates how long the tail of the ball will be
	list[2].Draw(50, 150, 300, 20, 10, 200, 10, "Speed Increase", true)			//indicates which number is added to the global speed every time a player scores
	list[3].Draw(50, 190, 300, 20, 10, 100, 10, "Win Count", true)				//indicates up to how many points are played
	list[4].Draw(50, 230, 300, 20, 10, 400, 150, "Paddle Length", true)			//the higher the value the longer the paddles (easier)
	list[5].Draw(50, 270, 300, 20, 10, 16, 2, "Paddle Speed", true)				//the higher the value the faster the movement of the paddles (easier)
	list[6].Draw(50, 310, 300, 20, 10, 10, 2, "Paddle Wait Time", true)			//the higher the value the slower the paddles
	list[7].Draw(50, 350, 300, 20, 10, 100, 20, "Y Randomness", true)			//sets the maxed randomness of the ball in y direction based on the current delta_y value in percent
	list[8].Draw(50, 390, 300, 20, 10, 100, 20, "X Randomness", true)			//sets the maxed randomness of the ball in x direction based on the current delta_x value in percent
	list[9].Draw(50, 430, 300, 20, 10, 100, 35, "Max Delta Y", true)			//sets the max value of delta_y in relation to speed in percent
	go Mouse(list, b1)

	for !start {
		if active != -1 {
			m.Lock()
			list[active].Redraw(ms_x)
			m.Unlock()
		}
	}

	b1.Click()
	time.Sleep(time.Duration(350) * time.Millisecond)

	//starting variables
	var speed int = int(math.Round(float64(list[0].Value)))
	var tail_len uint8 = uint8(math.Round(float64(list[1].Value)))  
	var speed_inc int = int(math.Round(float64(list[2].Value))) 
	var win_count int = int(math.Round(float64(list[3].Value)))
	var paddle_len uint16 = uint16(math.Round(float64(list[4].Value)))        	
	var paddle_speed uint16 = uint16(math.Round(float64(list[5].Value)))          	
	var paddle_wait_time int = int(math.Round(float64(list[6].Value)))  	
	var y_randomness int = int(math.Round(float64(list[7].Value)))     	
	var x_randomness int = int(math.Round(float64(list[8].Value)))      	
	var max_delta_y int = int(math.Round(float64(list[9].Value)))   	
	// var max_randomess float32 = list[10].Value      	
	                 	
	//starting constants
	var c_x float32
	var c_y float32
	var delta_x int
	var delta_y int
	var x_temp float32
	var y_temp float32
	var lcount int
	var rcount int
	var first bool = true

	gfx.SetzeFont("pong.ttf", 50)

	//initializing block
	fmt.Println("initialize")
	c_x, c_y, delta_x, delta_y = initialize(w_x, w_y, speed, tail_len)

	go read_keyboard()
	go left_paddle(w_x, w_y, paddle_len, paddle_speed, paddle_wait_time)
	go right_paddle(w_x, w_y, paddle_len, paddle_speed, paddle_wait_time)

	var c_time int64 = time.Now().UnixMilli()
	var e_time int64 = time.Now().UnixMilli() - c_time

	//main loop
	for {
		e_time = time.Now().UnixMilli() - c_time
		c_time = time.Now().UnixMilli()

		x_temp = c_x
		y_temp = c_y
		gfx.UpdateAus()

		// fmt.Println(c_x, c_y)

		// clear left paddle
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Transparenz(0)
		gfx.Vollrechteck(left_paddle_x, left_paddle_y-8*paddle_speed, 10, paddle_len+16*paddle_speed)

		// clear right paddle
		gfx.Vollrechteck(right_paddle_x, right_paddle_y-8*paddle_speed, 10, paddle_len+16*paddle_speed)

		// clear screen
		gfx.Transparenz(tail_len)
		gfx.Vollrechteck(0, 0, w_x, w_y+70)

		// draw left paddle
		gfx.Stiftfarbe(255, 255, 255)
		gfx.Transparenz(0)
		gfx.Vollrechteck(left_paddle_x, left_paddle_y, 10, paddle_len)

		//draw right paddle
		gfx.Vollrechteck(right_paddle_x, right_paddle_y, 10, paddle_len)

		// draw border
		gfx.Vollrechteck(0, w_y, w_x, 4)

		// middle line
		gfx.Vollrechteck(w_x/2-1, 0, 2, w_y+70)

		// draw lcount
		gfx.SchreibeFont(w_x/2-150-uint16(give_digits_of_value(lcount)*32), w_y+10, strconv.Itoa(lcount))

		// draw rcount
		gfx.SchreibeFont(w_x/2+150, w_y+10, strconv.Itoa(rcount))

		//draw ball
		if first {
			c_x, c_y = exe_lin_func(e_time, c_x, c_y, delta_x/2, delta_y/2)
		}	else {
			c_x, c_y = exe_lin_func(e_time, c_x, c_y, delta_x, delta_y)
			}
		gfx.Vollkreis(uint16(math.Round(float64(c_x))), w_y-uint16(math.Round(float64(c_y))), 10)

		gfx.UpdateAn()
        
          //check win
		if lcount >= win_count {
			win("Player 1", w_x, w_y)
		} else if rcount >= win_count {
			win("Player 2", w_x, w_y)
		}

          //check if ball is left outside the field
		if c_x <= 0 && c_x <= x_temp {
			c_x, c_y, delta_x, delta_y = initialize(w_x, w_y, speed, tail_len)
			rcount = rcount + 1
			first = true
			speed = speed + speed_inc
			c_time = time.Now().UnixMilli()

          //check if ball is right outside the field
		} else if c_x >= float32(w_x) && c_x >= x_temp {
			c_x, c_y, delta_x, delta_y = initialize(w_x, w_y, speed, tail_len)
			lcount = lcount + 1
			first = true
			speed = speed + speed_inc
			c_time = time.Now().UnixMilli()

          //bounce left paddle
		} else if c_x -5 <= 20 && c_x <= x_temp && int(w_y)-int(math.Round(float64(c_y))) >= int(left_paddle_y)-8 && w_y-uint16(math.Round(float64(c_y))) <= left_paddle_y+paddle_len+8 {
			delta_x, delta_y = y_bounce(delta_x, delta_y, x_randomness, speed, max_delta_y)
			first = false

          //bounce right paddle
		} else if c_x +5 >= float32(w_x)-20 && c_x >= x_temp && int(w_y)-int(math.Round(float64(c_y))) >= int(right_paddle_y)-8 && w_y-uint16(math.Round(float64(c_y))) <= right_paddle_y+paddle_len+8 {
			delta_x, delta_y = y_bounce(delta_x, delta_y, x_randomness, speed, max_delta_y)
			first = false

          //bounce up/down
		} else if c_y <= 10 && c_y <= y_temp || c_y >= float32(w_y-10) && c_y >= y_temp {
			delta_x, delta_y = x_bounce(delta_y, delta_x, y_randomness, speed, max_delta_y)
		}
	}
}

//function to read the mouse - only relevant for the settings menu
func Mouse(list [10]sliders.Slider, b1 buttons.Button)  () {
	var m_x uint16
	var m_y uint16
	var ms_bttn uint8
	var ms_prssd int8
	for !start{
		ms_bttn, ms_prssd, m_x, m_y = gfx.MausLesen1()

		if ms_bttn == 1 && ms_prssd == 1 || ms_bttn == 1 && ms_prssd == 0 	{
			ms_x = m_x
			if active == -1 {
				if m_x >= b1.X && m_x <= b1.X + b1.Length && m_y >= b1.Y && m_y <= b1.Y + b1.Height {
					start = true
				}	else {
						for i := 0; i < len(list); i++ {
							if m_x >= list[i].X && m_x <= list[i].X + list[i].Length && m_y >= list[i].Y && m_y <= list[i].Y + list[i].Height {
								m.Lock()
								active = i
								m.Unlock()
							}			
						}
					}
			}
		}	else {
				m.Lock()
				active = -1
				m.Unlock()
			}
	}
}

//initializing function is executd every time the ball gets out of bounds
func initialize(w_x uint16, w_y uint16, speed int, tail_len uint8) (float32, float32, int, int) {
	gfx.Transparenz(0)
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Vollrechteck(0, 0, w_x, w_y)
	gfx.Transparenz(tail_len)
	gfx.Stiftfarbe(0, 0, 0)
	for i := 0; i <= 50; i++ {
		gfx.Vollrechteck(0, 0, w_x, w_y)
	}
	var delta_x int
	var delta_y int

	c_x := float32(w_x) / 2
	c_y := float32(w_y) / 2

	rand.Seed(time.Now().UTC().UnixNano())
    delta_x = rand.Intn(speed - speed/2) + speed/2
    rand.Seed(time.Now().UTC().UnixNano())
    if rand.Intn(2) == 0 {
        delta_x = -1*delta_x
    }
	delta_y = (speed) - int(math.Abs(float64(delta_x)))
    rand.Seed(time.Now().UTC().UnixNano()*time.Now().UTC().UnixNano())
    if rand.Intn(2) == 0 {
        delta_y = -1*delta_y
    }

	gfx.Transparenz(0)
	gfx.Stiftfarbe(0, 0, 0)
	gfx.Vollrechteck(0, 0, w_x, w_y+70)
	gfx.Stiftfarbe(255, 255, 255)
	gfx.Vollkreis(uint16(math.Round(float64(c_x))), w_y-uint16(math.Round(float64(c_y))), 10)
	return c_x, c_y, delta_x, delta_y
}

//function is relevant for the calculation of the exact position of the point counter
func give_digits_of_value(value int) int {
	count := 0
	for value > 0 {
		value = value / 10
		count++
	}
	if count == 0 {
		count++
	}
	return count
}

//function will be executed when the win_count == lcount or rcount
func win(winner string, w_x uint16, w_y uint16) {
	var r uint8
	var g uint8
	var b uint8
	gfx.Transparenz(0)
	gfx.SetzeFont("pong.ttf", 100)
	for {
		gfx.UpdateAus()
		rand.Seed(time.Now().UTC().UnixNano())
		r = uint8(50 + rand.Intn(255-50))
		time.Sleep(time.Duration(1) * time.Millisecond)
		rand.Seed(time.Now().UTC().UnixNano())
		g = uint8(50 + rand.Intn(255-50))
		time.Sleep(time.Duration(1) * time.Millisecond)
		rand.Seed(time.Now().UTC().UnixNano())
		b = uint8(50 + rand.Intn(255-50))
		gfx.Stiftfarbe(r, g, b)
		gfx.Vollrechteck(0, 0, w_x, w_y+70)
		gfx.Stiftfarbe(255, 255, 255)
		gfx.SchreibeFont(225, 300, winner)
		gfx.UpdateAn()
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
}

//function is responsible for reading the keyboard which means it is also responsible for controlling the paddles
func read_keyboard() {
	var kn uint16
	var prssd uint8
	for {
		kn, prssd, _ = gfx.TastaturLesen1()
		if kn == 115 || kn == 119 {
			lkn = kn
			lprssd = prssd
		} else if kn == 273 || kn == 274 {
			rkn = kn
			rprssd = prssd
		}
	}
}

//controlling the left paddle
func left_paddle(w_x uint16, w_y uint16, paddle_len uint16, paddle_speed uint16, paddle_wait_time int) {
	left_paddle_x = 10

	for {
		if lkn == 119 && lprssd == 1 {
			for {
				if left_paddle_y > 0 {
					left_paddle_y = left_paddle_y - paddle_speed
				}
				time.Sleep(time.Duration(paddle_wait_time) * time.Millisecond)
				if lkn == 119 && lprssd == 0 || lkn == 115 {
					break
				}
			}
		} else if lkn == 115 && lprssd == 1 {
			for {
				if left_paddle_y < w_y-paddle_len {
					left_paddle_y = left_paddle_y + paddle_speed
				}
				time.Sleep(time.Duration(paddle_wait_time) * time.Millisecond)
				if lkn == 115 && lprssd == 0 || lkn == 119 {
					break
				}
			}
		}
	}
}

//controlling the right paddle
func right_paddle(w_x uint16, w_y uint16, paddle_len uint16, paddle_speed uint16, paddle_wait_time int) {
	right_paddle_x = w_x - 20

	for {
		if rkn == 273 && rprssd == 1 {
			for {
				if right_paddle_y > 0 {
					right_paddle_y = right_paddle_y - paddle_speed
				}
				time.Sleep(time.Duration(paddle_wait_time) * time.Millisecond)
				if rkn == 273 && rprssd == 0 || rkn == 274 {
					break
				}
			}
		} else if rkn == 274 && rprssd == 1 {
			for {
				if right_paddle_y < w_y-paddle_len {
					right_paddle_y = right_paddle_y + paddle_speed
				}
				time.Sleep(time.Duration(paddle_wait_time) * time.Millisecond)
				if rkn == 274 && rprssd == 0 || rkn == 273 {
					break
				}
			}
		}
	}
}

//function to calculate the next cords
func exe_lin_func(e_time int64, c_x float32, c_y float32, delta_x int, delta_y int) (float32, float32) {
	c_x = c_x + float32(delta_x) * float32(e_time) / 1000
	c_y = c_y + float32(delta_y) * float32(e_time) / 1000
	return c_x, c_y
}

func x_bounce(delta_y int, delta_x int, y_randomness int, speed int, max_delta_y int) (int, int) {
	fmt.Println("x axe bounce -------- delta_y =", delta_y, "--- delta_x =", delta_x)

	delta_y = - delta_y
	y_max_random := y_randomness * int(math.Abs(float64(delta_y))) / 100
	rand.Seed(time.Now().UTC().UnixNano())
	r_add := rand.Intn(y_max_random - -y_max_random) -y_max_random

	if int(math.Abs(float64(delta_y + r_add))) > max_delta_y * speed / 100 {
		if int(math.Abs(float64(delta_y + r_add))) > int(math.Abs(float64(delta_y - r_add)))	{
			r_add = -1*r_add
		}
	}

	if delta_y + r_add >= speed || delta_y + r_add <= -speed{
		r_add = -1*r_add
	}

	delta_y = delta_y + r_add

	var s_x int = delta_x / int(math.Abs(float64(delta_x)))
	delta_x = speed - int(math.Abs(float64(delta_y)))
	delta_x = delta_x * s_x

	fmt.Println("after calculation --- delta_y =", delta_y, "--- delta_x =", delta_x)
	return delta_x, delta_y
}

func y_bounce(delta_x int, delta_y int, x_randomness int, speed int, max_delta_y int) (int, int) {
	fmt.Println("y axe bounce -------- delta_x =", delta_x, "--- delta_y =", delta_y)

	delta_x = - delta_x
	x_max_random := x_randomness * int(math.Abs(float64(delta_x))) / 100
	rand.Seed(time.Now().UTC().UnixNano())
	r_add := rand.Intn(x_max_random - -x_max_random) -x_max_random

	if int(math.Abs(float64(delta_x + r_add))) < (100-max_delta_y) * speed / 100 {
		if int(math.Abs(float64(delta_x + r_add))) < int(math.Abs(float64(delta_x - r_add)))	{
			r_add = -1*r_add
		}
	}

	if delta_x + r_add >= speed || delta_x + r_add <= -speed{
		r_add = -1*r_add
	}

	delta_x = delta_x + r_add

	var s_y int = delta_y / int(math.Abs(float64(delta_y)))
	delta_y = speed - int(math.Abs(float64(delta_x)))
	delta_y = delta_y * s_y
	
	fmt.Println("after calculation --- delta_x =", delta_x, "--- delta_y =", delta_y)
	return delta_x, delta_y
}

//feature: -Ball kann angeschnitten werden      Slider Anschneidungsstärke in Prozent