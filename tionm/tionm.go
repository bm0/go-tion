package tionm

import (
	"errors"
	"time"

	"log"

	"github.com/m-pavel/go-tion/tion"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

type mTion struct {
	addr  string
	cnct  chan error
	debug bool

	d *device.Device1
}

func New(addr string, debug ...bool) tion.Tion {
	nt := mTion{addr: addr}
	nt.cnct = make(chan error)
	if len(debug) > 0 {
		nt.debug = debug[0]
	}
	return &nt
}

func (n *mTion) ReadState(timeout time.Duration) (*tion.Status, error) {
	wc, err := n.d.GetCharByUUID(tion.WRITE_CHARACT)
	if err != nil {
		return nil, err
	}
	log.Println("Got write char")
	if err := wc.WriteValue(tion.StatusRequest, nil); err != nil {
		return nil, err
	}
	log.Println("Written")
	time.Sleep(1 * time.Second)
	rc, err := n.d.GetCharByUUID(tion.READ_CHARACT)
	if err != nil {
		return nil, err
	}
	log.Println("Got read char")
	if data, err := rc.ReadValue(nil); err != nil {
		log.Println(err)
		return nil, err
	} else {
		log.Printf("Read %v\n", data)
		return tion.FromBytes(data)
	}
}

func (n *mTion) Update(s *tion.Status, timeout time.Duration) error {
	return errors.New("not implemented")
}

func (n *mTion) Connect(timeout time.Duration) error {
	ad, err := api.GetDefaultAdapter()
	if err != nil {
		return err
	}
	log.Println("Got adapter")
	n.d, err = ad.GetDeviceByAddress(n.addr)
	if err != nil {
		return err
	}
	log.Println("Got device")
	if err = n.d.Connect(); err != nil {
		return err
	}
	return nil
}

func (n *mTion) Disconnect() error {
	if n.d != nil {
		defer func() {
			n.d = nil
		}()
		return n.d.Disconnect()
	}
	return nil
}