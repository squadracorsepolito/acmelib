package acmelib

// func Test_Network_AddTemplate(t *testing.T) {
// 	assert := assert.New(t)

// 	net := NewNetwork("net")

// 	bus0 := NewBus("bus_0")
// 	assert.NoError(net.AddBus(bus0))

// 	node0 := NewNode("node_0", 0)
// 	assert.NoError(bus0.AddNode(node0))

// 	msg0 := NewMessage("msg_0", 1, 1)
// 	assert.NoError(node0.AddMessage(msg0))

// 	sigTypeSize4, err := NewIntegerSignalType("sig_type_size_4", 4, false)
// 	assert.NoError(err)

// 	stdSig0, err := NewStandardSignal("std_sig_0", sigTypeSize4)
// 	assert.NoError(err)
// 	assert.NoError(msg0.AppendSignal(stdSig0))
// 	stdSig1, err := NewStandardSignal("std_sig_1", sigTypeSize4)
// 	assert.NoError(err)
// 	assert.NoError(msg0.AppendSignal(stdSig1))

// 	sigUnit := NewSignalUnit("sig_unit", SignalUnitKindCustom, "unit")
// 	stdSig0.SetUnit(sigUnit)
// 	stdSig1.SetUnit(sigUnit)

// 	msg1 := NewMessage("msg_1", 2, 1)
// 	assert.NoError(node0.AddMessage(msg1))

// 	sigEnum := NewSignalEnum("sig_enum")
// 	assert.NoError(sigEnum.AddValue(NewSignalEnumValue("val_15", 15)))

// 	enumSig0, err := NewEnumSignal("enum_sig_0", sigEnum)
// 	assert.NoError(err)
// 	assert.NoError(msg1.AppendSignal(enumSig0))
// 	enumSig1, err := NewEnumSignal("enum_sig_1", sigEnum)
// 	assert.NoError(err)
// 	assert.NoError(msg1.AppendSignal(enumSig1))

// 	assert.NoError(net.AddTemplate(sigTypeSize4))
// 	assert.NoError(net.AddTemplate(sigUnit))
// 	assert.NoError(net.AddTemplate(sigEnum))

// 	// Testing signal type templates
// 	sigTypeTpls := net.SignaTypeTemplates()
// 	assert.Len(sigTypeTpls, 1)
// 	stdSigRefs := sigTypeTpls[0].References()
// 	assert.Len(stdSigRefs, 2)
// 	assert.Contains(stdSigRefs, stdSig0)
// 	assert.Contains(stdSigRefs, stdSig1)

// 	// Testing signal unit templates
// 	sigUnitTpls := net.SignaUnitTemplates()
// 	assert.Len(sigUnitTpls, 1)
// 	stdSigRefs = sigUnitTpls[0].References()
// 	assert.Len(stdSigRefs, 2)
// 	assert.Contains(stdSigRefs, stdSig0)
// 	assert.Contains(stdSigRefs, stdSig1)

// 	// Testing signal enum templates
// 	sigEnumTpls := net.SignaEnumTemplates()
// 	assert.Len(sigEnumTpls, 1)
// 	enumSigRefs := sigEnumTpls[0].References()
// 	assert.Len(enumSigRefs, 2)
// 	assert.Contains(enumSigRefs, enumSig0)
// 	assert.Contains(enumSigRefs, enumSig1)
// }
