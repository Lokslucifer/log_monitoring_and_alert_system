package logstreamer
import(
	"sync"
)

type LogConsumer interface {
	StartConsuming(ch chan<- string, wg *sync.WaitGroup) 
	Stop()
}