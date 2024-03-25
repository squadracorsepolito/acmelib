package acmelib

// func Test_SignalEnum_AddValue(t *testing.T) {
// 	assert := assert.New(t)

// 	msg := NewMessage("message", "", 1)

// 	sigEnum := NewSignalEnum("signal_enum", "")

// 	sig, err := NewEnumSignal("signal", "", sigEnum)
// 	assert.NoError(err)

// 	assert.NoError(msg.AppendSignal(sig))

// 	sigEnumVal0 := NewSignalEnumValue("val_0", "", 8)
// 	sigEnumVal1 := NewSignalEnumValue("val_1", "", 0)
// 	sigEnumVal2 := NewSignalEnumValue("val_2", "", 512)
// 	sigEnumVal3 := NewSignalEnumValue("val_3", "", 127)

// 	assert.NoError(sigEnum.AddValue(sigEnumVal0))
// 	assert.NoError(sigEnum.AddValue(sigEnumVal1))
// 	assert.Error(sigEnum.AddValue(sigEnumVal2))
// 	assert.NoError(sigEnum.AddValue(sigEnumVal3))

// 	sigEnumVal4 := NewSignalEnumValue("val_4", "", 8)
// 	assert.Error(sigEnum.AddValue(sigEnumVal4))

// 	correctOrder := []string{"val_1", "val_0", "val_3"}
// 	for idx, val := range sigEnum.GetValuesByIndex() {
// 		assert.Equal(correctOrder[idx], val.GetName())
// 	}

// 	t.Log(sig.String())
// }

// func Test_SignalEnumValue_UpdateIndex(t *testing.T) {
// 	assert := assert.New(t)

// 	msg0 := NewMessage("msg_0", "", 1)

// 	sigEnum := NewSignalEnum("signal_enum", "")

// 	sig0, err := NewEnumSignal("sig_0", "", sigEnum)
// 	assert.NoError(err)

// 	assert.NoError(msg0.AppendSignal(sig0))

// 	sigEnumVal := NewSignalEnumValue("val", "", 0)

// 	assert.NoError(sigEnum.AddValue(sigEnumVal))

// 	assert.NoError(sigEnumVal.UpdateIndex(127))
// 	assert.Error(sigEnumVal.UpdateIndex(512))

// 	msg1 := NewMessage("msg_1", "", 1)

// 	sig1, err := NewEnumSignal("sig_1", "", sigEnum)
// 	assert.NoError(err)

// 	assert.NoError(msg1.AppendSignal(sig1))

// 	assert.NoError(sigEnumVal.UpdateIndex(8))
// 	assert.Error(sigEnumVal.UpdateIndex(512))

// 	t.Log(msg0)
// 	t.Log(msg1)
// }
