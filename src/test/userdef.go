package main

import (
	"chord"
	"github.com/sirupsen/logrus"
	"strconv"
)

/* In this file, you should implement function "NewNode" and
 * a struct which implements the interface "dhtNode".
 */

func NewNode(port int) dhtNode {
	// Todo: create a node and then return it.
	ret := NewPubNode(":" + strconv.Itoa(port))
	return ret
}

// Todo: implement a struct which implements the interface "dhtNode".

/* PubNode which implements the interface "dhtNode"
 * PubNode -> Receiver, Receiver (RCVR) is for RPC
 * Receiver -> Node, Node is the true node in DHT network
 */

type PubNodeType struct {
	receiver *chord.ReceiverType
}

func NewPubNode(_addr string) *PubNodeType{
	return &PubNodeType{
		receiver: chord.NewReceiver(_addr),
	}
}

func (this *PubNodeType) Run() {
	var err error

	err = this.receiver.Server.Register(this.receiver)

	if err != nil {
		chord.Log.WithFields(logrus.Fields{
		}).Error("Register Error. " + err.Error())
	}
	this.receiver.Node.Running = true
	go this.receiver.Node.Updating()

	chord.Log.WithFields(logrus.Fields{
		"addr" : this.receiver.Node.Addr,
	}).Info("Running Success!")

	this.receiver.Server.Accept(this.receiver.Listener)
}

func (this *PubNodeType) Create() {
	this.receiver.Node.Create()
}

func (this PubNodeType) Join(addr string) bool {
	return this.receiver.Node.Join(addr)
}

func (this PubNodeType) Quit() {
	this.receiver.Node.Quit()
	this.receiver.Node.Running = false
	err := this.receiver.Listener.Close()
	if err == nil {
		chord.Log.WithFields(logrus.Fields{
			"ip" : this.receiver.Node.Addr,
		}).Info("Quit Success.")
	}
}

func (this PubNodeType) ForceQuit() {
	return
}

func (this PubNodeType) Ping(addr string) bool {
	err := chord.Ping(addr)
	return err == nil
}

func (this PubNodeType) Put(key string, value string) bool {
	return this.receiver.Node.Put(key, value)
}

func (this PubNodeType) Get(key string) (bool, string) {
	return this.receiver.Node.Get(key)
}

func (this *PubNodeType) Delete(key string) bool {
	return this.receiver.Node.Delete(key)
}