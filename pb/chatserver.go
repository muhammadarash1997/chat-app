package pb

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e6)
	errch := make(chan error)

	// receive message

	// send message

	return <-errch
}

// receive message
func receiveFromStream(csi_ servicesChatServiceServer, clientUniqueCode_ int, errch_ chan error) {

	// implement a loop

	for {
		mssg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errch_ <- err
		} else {
			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        mssg.Name,
				MessageBody:       mssg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUniqueCode_,
			})

			messageHandleObject.mu.Unlock()

			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])
		}
	}
}

func sendToStream(csi_ servicesChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	// implement a loop

	for {

		// loop through message in MQue
		for {
			time.Sleep(500 * time.Millisecond)
			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName4Client := messageHandleObject.MQue[0].ClientName
			message4Client := messageHandleObject.MQue[0].MessageBody

			messageHandleObject.mu.Unlock()

			// send message to designated client (do not send to the same client)
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(&FromServer{Name: senderName4Client, Body: message4Client})
				if err != nil {
					errch_ <- err
				}
				
				messageHandleObject.mu.Lock()
			}
		}
	}
}
