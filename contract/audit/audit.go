package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/orbs-network/contract-external-libraries-go/v1/keys"
	"github.com/orbs-network/contract-external-libraries-go/v1/structs"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/address"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/env"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
)

var PUBLIC = sdk.Export(recordEvent, getEventsByHash, setEventSourceContractAddress)
var SYSTEM = sdk.Export(_init)

var COUNTER_KEY = "counter"
var OWNER_KEY = []byte("owner")

func _init() {
	state.WriteBytes(OWNER_KEY, address.GetSignerAddress())
}

type Event struct {
	Action string
	From   string
	To     string

	SignerAddress string
	Timestamp     uint64
}

func recordEvent(hash string, action string, from string, to string) {
	_verifyEventSource()
	event := Event{
		Action:        action,
		From:          from,
		To:            to,
		SignerAddress: hex.EncodeToString(address.GetSignerAddress()),
		Timestamp:     env.GetBlockTimestamp(),
	}

	structs.WriteStruct("events."+hash+"."+strconv.FormatUint(_value(hash), 10), event)
	_inc(hash)
}

func getEventsByHash(hash string) string {
	var events []Event

	events_total := _value(hash)
	for i := uint64(0); i < events_total; i++ {
		event := Event{}
		structs.ReadStruct("events."+hash+"."+strconv.FormatUint(i, 10), &event)
		events = append(events, event)
	}

	rawJson, _ := json.Marshal(events)
	return string(rawJson)
}

var EVENT_SOURCE_CONTRACT_ADDRESS = []byte("event_source_contract_address")

func setEventSourceContractAddress(addr string) {
	_ownerOnly()
	state.WriteString(EVENT_SOURCE_CONTRACT_ADDRESS, addr)
}

func getEventSourceContractAddress() string {
	return state.ReadString(EVENT_SOURCE_CONTRACT_ADDRESS)
}

func _verifyEventSource() {
	eventSourceContractAddress := getEventSourceContractAddress()
	if !bytes.Equal(address.GetCallerAddress(), address.GetContractAddress(eventSourceContractAddress)) {
		panic("event source contract address is not set or is wrong")
	}
}

func _inc(hash string) uint64 {
	v := _value(hash) + 1
	state.WriteUint64(keys.Key(COUNTER_KEY, ".", hash), v)
	return v
}

func _value(hash string) uint64 {
	return state.ReadUint64(keys.Key(COUNTER_KEY, ".", hash))
}

func _ownerOnly() {
	if !bytes.Equal(state.ReadBytes(OWNER_KEY), address.GetSignerAddress()) {
		panic("not allowed!")
	}
}
