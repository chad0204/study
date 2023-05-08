package main

import "fmt"

type Observer interface {
	Update(temp int)
}

type Subject interface {
	Register(observer Observer)
	notifyAll()
}

type Weather interface {
	setTemp(temp int)
}

type WeatherSubject struct {
	temp      int // 温度
	observers []Observer
}

func (w *WeatherSubject) setTemp(temp int) {
	w.temp = temp
	w.notifyAll()
}

func (w *WeatherSubject) Register(observer Observer) {
	if len(w.observers) == 0 {
		w.observers = make([]Observer, 0)
	}
	w.observers = append(w.observers, observer)
}

func (w *WeatherSubject) notifyAll() {
	for _, ob := range w.observers {
		ob.Update(w.temp)
	}
}

type FarmerObserver struct{}

func (f *FarmerObserver) Update(temp int) {
	if temp > 40 {
		fmt.Println("太热了不干活了")
	} else {
		fmt.Printf("天气不错 接着干")
	}
}

func main() {
	w := &WeatherSubject{temp: 10}

	w.Register(&FarmerObserver{})

	w.setTemp(50)
	w.setTemp(25)

}
