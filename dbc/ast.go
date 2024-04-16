// Package dbc provides a [Parser] and a [Writer] for DBC files.
// It is based on the version 1.0.5 of the DBC file format.
package dbc

// FileExtension is the extension of a DBC file.
const FileExtension = ".dbc"

// DummyNode is the default name of a dummy node used in the DBC file.
const DummyNode = "Vector__XXX"

// GetNewSymbols returns a list of all the new symbols.
func GetNewSymbols() []string {
	return newSymbolsValues
}

// Location is the position in the DBC file of an AST entity.
type Location struct {
	Filename string
	Line     int
	Col      int
}

// File definition:
//
// The DBC file describes the communication of a single CAN network.
// This information is sufficient to monitor and analyze the network and to simulate nodes not
// physically available (remaining bus simulation).
//
// The DBC file can also be used to develop the communication software of an electronic control unit
// which shall be part of the CAN network. The functional behavior
// of the ECU is not addressed by the DBC file.
type File struct {
	Location *Location

	Version             string
	NewSymbols          *NewSymbols
	BitTiming           *BitTiming
	Nodes               *Nodes
	ValueTables         []*ValueTable
	Messages            []*Message
	MessageTransmitters []*MessageTransmitter
	EnvVars             []*EnvVar
	EnvVarDatas         []*EnvVarData
	SignalTypes         []*SignalType
	Comments            []*Comment
	Attributes          []*Attribute
	AttributeDefaults   []*AttributeDefault
	AttributeValues     []*AttributeValue
	ValueEncodings      []*ValueEncoding
	SignalTypeRefs      []*SignalTypeRef
	SignalGroups        []*SignalGroup
	SignalExtValueTypes []*SignalExtValueType
	ExtendedMuxes       []*ExtendedMux
}

// NewSymbols definition:
//
// new_symbols = [ '_NS' ':' ['CM_'] ['BA_DEF_'] ['BA_'] ['VAL_']
// ['CAT_DEF_'] ['CAT_'] ['FILTER'] ['BA_DEF_DEF_'] ['EV_DATA_']
// ['ENVVAR_DATA_'] ['SGTYPE_'] ['SGTYPE_VAL_'] ['BA_DEF_SGTYPE_']
// ['BA_SGTYPE_'] ['SIG_TYPE_REF_'] ['VAL_TABLE_'] ['SIG_GROUP_']
// ['SIG_VALTYPE_'] ['SIGTYPE_VALTYPE_'] ['BO_TX_BU_']
// ['BA_DEF_REL_'] ['BA_REL_'] ['BA_DEF_DEF_REL_'] ['BU_SG_REL_']
// ['BU_EV_REL_'] ['BU_BO_REL_'] [SG_MUL_VAL_'] ];
type NewSymbols struct {
	Location *Location

	Symbols []string
}

// BitTiming definition:
//
// The bit timing section defines the baudrate and the settings of the BTR registers of
// the network. This section is obsolete and not used any more. Nevertheless the
// keyword 'BS_' must appear in the DBC file.
//
// bit_timing = 'BS_:' [baudrate ':' BTR1 ',' BTR2 ] ;
//
// baudrate = unsigned_integer ;
//
// BTR1 = unsigned_integer ;
//
// BTR2 = unsigned_integer ;
type BitTiming struct {
	Location *Location

	Baudrate      uint32
	BitTimingReg1 uint32
	BitTimingReg2 uint32
}

// Nodes definition:
//
// The node section defines the names of all participating nodes. The names defined
// in this section have to be unique within this section.
//
// nodes = 'BU_:' {node_name} ;
//
// node_name = DBC_identifier ;
type Nodes struct {
	Location *Location

	Names []string
}

// ValueTable definition:
//
// The value table section defines the global value tables. The value descriptions in
// value tables define value encodings for signal raw values. In commonly used DBC
// files the global value tables aren't used, but the value descriptions are defined for
// each signal independently.
//
// value_tables = {value_table} ;
//
// value_table = 'VAL_TABLE_' value_table_name {value_description} ';' ;
//
// value_table_name = DBC_identifier ;
type ValueTable struct {
	Location *Location

	Name   string
	Values []*ValueDescription
}

// ValueDescription definition:
//
// A value description defines a textual description for a single value. This value may
// either be a signal raw value transferred on the bus or the value of an environment
// variable in a remaining bus simulation.
//
// value_description = unsigned_integer char_string ;
type ValueDescription struct {
	Location *Location

	ID   uint32
	Name string
}

// ValueEncodingKind defines the kind of a [ValueEncoding].
type ValueEncodingKind uint

const (
	// ValueEncodingSignal defines a signal value encodings.
	ValueEncodingSignal ValueEncodingKind = iota
	// ValueEncodingEnvVar defines an envvar value encodings.
	ValueEncodingEnvVar
)

// ValueEncoding definition:
//
// Signal value descriptions define encodings for specific signal raw values.
//
// value_descriptions = { value_descriptions_for_signal | value_descriptions_for_env_var } ;
//
// value_descriptions_for_signal = 'VAL_' message_id signal_name { value_description } ';' ;
//
// The value descriptions for environment variables provide textual representations of
// specific values of the variable.
//
// value_descriptions_for_env_var = 'VAL_' env_var_name { value_description } ';' ;
type ValueEncoding struct {
	Location *Location

	Kind       ValueEncodingKind
	MessageID  uint32
	SignalName string
	EnvVarName string
	Values     []*ValueDescription
}

// Message definition:
//
// The message section defines the names of all frames in the cluster as well as their
// properties and the signals transferred on the frames.
//
// messages = {message} ;
//
// message = BO_ message_id message_name ':' message_size transmitter {signal} ;
//
// message_id = unsigned_integer ;
//
// The message's CAN-ID. The CAN-ID has to be unique within the DBC file. If the
// most significant bit of the CAN-ID is set, the ID is an extended CAN ID.
// The extended CAN ID can be determined by masking out the most significant bit with the
// mask 0x7FFFFFFF.
//
// message_name = DBC_identifier ;
//
// The names defined in this section have to be unique within the set of messages.
//
// message_size = unsigned_integer ;
//
// The message_size specifies the size of the message in bytes.
//
// transmitter = node_name | 'Vector__XXX' ;
//
// The transmitter name specifies the name of the node transmitting the message.
// The sender name has to be defined in the set of node names in the node section.
// If the massage shall have no sender, the string 'Vector__XXX' has to be given
// here.
type Message struct {
	Location *Location

	ID          uint32
	Name        string
	Size        uint32
	Transmitter string
	Signals     []*Signal
}

// SignalByteOrder defines the order of a [Signal].
type SignalByteOrder uint

const (
	// SignalLittleEndian defines a little endian signal.
	SignalLittleEndian SignalByteOrder = iota
	// SignalBigEndian defines a little endian signal.
	SignalBigEndian
)

// SignalValueType defines the value type of a [Signal].
type SignalValueType uint

const (
	// SignalUnsigned defines an unsigned signal.
	SignalUnsigned SignalValueType = iota
	// SignalSigned defines a signed signal.
	SignalSigned
)

// Signal definition:
//
// The message's signal section lists all signals placed on the message, their position
// in the message's data field and their properties.
//
// signal = 'SG_' signal_name multiplexer_indicator ':' start_bit '|'
// signal_size '@' byte_order value_type '(' factor ',' offset ')'
// '[' minimum '|' maximum ']' unit receiver {',' receiver} ;
//
// signal_name = DBC_identifier ;
//
// The names defined here have to be unique for the signals of a single message.
//
// multiplexer_indicator = ' ' | [m multiplexer_switch_value] [M] ;
//
// The multiplexer indicator defines whether the signal is a normal signal,
// a multiplexer switch for multiplexed signals, or a multiplexed signal. A 'M' (uppercase)
// character defines the signal as the multiplexer switch. A 'm' (lowercase) character
// followed by an unsigned integer defines the signal as being multiplexed by the
// multiplexer switch. A multiplexed signal is transferred in the message if the switch
// value of the multiplexer signal is equal to its multiplexer_switch_value.
//
// Note: A signal may be a multiplexed signal and a multiplexor switch signal at the
// same time. And further: more than one signal within a single message can be a
// multiplexer switch. In both cases the extended multiplexing section (see below)
// mustn’t be empty then.
//
// start_bit = unsigned_integer ;
//
// The start_bit value specifies the position of the signal within the data field of the
// frame. For signals with byte order Intel (little endian) the position of the leastsignificant bit is given.
// For signals with byte order Motorola (big endian) the position of the most significant bit is given.
// The bits are counted in a sawtooth manner.
// The startbit has to be in the range of 0 to (8 * message_size - 1).
//
// signal_size = unsigned_integer ;
//
// The signal_size specifies the size of the signal in bits.
//
// byte_order = '0' | '1' ; (* 0=big endian, 1=little endian *)
//
// The byte_format is 0 if the signal's byte order is Motorola (big endian) or 1 if the
// byte order is Intel (little endian).
//
// value_type = '+' | '-' ; (* +=unsigned, -=signed *)
//
// The value_type defines the signal as being of type unsigned (-) or signed (-).
//
// factor = double ;
// offset = double ;
//
// The factor and offset define the linear conversion rule to convert the signals raw
// value into the signal's physical value and vice versa:
//
// physical_value = raw_value * factor + offset
//
// raw_value = (physical_value – offset) / factor
//
// As can be seen in the conversion rule formulas the factor must not be 0.
//
// minimum = double ;
//
// maximum = double ;
//
// The minimum and maximum define the range of valid physical values of the signal.
//
// unit = char_string ;
//
// receiver = node_name | 'Vector__XXX' ;
//
// The receiver name specifies the receiver of the signal. The receiver name has to
// be defined in the set of node names in the node section. If the signal shall have no
// receiver, the string 'Vector__XXX' has to be given here.
type Signal struct {
	Location *Location

	Name           string
	IsMultiplexor  bool
	IsMultiplexed  bool
	MuxSwitchValue uint32
	Size           uint32
	StartBit       uint32
	ByteOrder      SignalByteOrder
	ValueType      SignalValueType
	Factor         float64
	Offset         float64
	Min            float64
	Max            float64
	Unit           string
	Receivers      []string
}

// SignalExtValueTypeType defines the type of a [SignalExtValueType].
type SignalExtValueTypeType uint

const (
	// SignalExtValueTypeInteger defines an integer signal ext value type.
	SignalExtValueTypeInteger SignalExtValueTypeType = iota
	// SignalExtValueTypeFloat defines a float signal ext value type.
	SignalExtValueTypeFloat
	// SignalExtValueTypeDouble defines a double signal ext value type.
	SignalExtValueTypeDouble
)

// SignalExtValueType definition:
//
// Signals with value types 'float' and 'double' have additional entries in the signal_valtype_list section.
//
// signal_extended_value_type_list = 'SIG_VALTYPE_' message_id signal_name signal_extended_value_type ';' ;
//
// signal_extended_value_type = '0' | '1' | '2' | '3' ;
// (* 0=signed or unsigned integer, 1=32-bit IEEE-float, 2=64-bit IEEE-double *)
type SignalExtValueType struct {
	Location *Location

	MessageID    uint32
	SignalName   string
	ExtValueType SignalExtValueTypeType
}

// MessageTransmitter definition:
//
// The message transmitter section enables the definition of multiple transmitter
// nodes of a single message. This is used to describe communication data for
// higher-layer protocols. This is not used to define CAN layer-2 communication.
//
// message_transmitters = {message_transmitter} ;
//
// Message_transmitter = 'BO_TX_BU_' message_id ':' {transmitter} ';' ;
type MessageTransmitter struct {
	Location *Location

	MessageID    uint32
	Transmitters []string
}

// EnvVarType defines the type of an [EnvVar].
type EnvVarType uint

const (
	// EnvVarInt defines an int envvar.
	EnvVarInt EnvVarType = iota
	// EnvVarFloat defines a float envvar.
	EnvVarFloat
	// EnvVarString defines a string envvar.
	EnvVarString
)

// EnvVarAccessType defines the access type of an [EnvVar].
type EnvVarAccessType uint

const (
	// EnvVarDummyNodeVector0 defines a dummy vector node 0 access type.
	EnvVarDummyNodeVector0 EnvVarAccessType = iota
	// EnvVarDummyNodeVector1 defines a dummy vector node 1 access type.
	EnvVarDummyNodeVector1
	// EnvVarDummyNodeVector2 defines a dummy vector node 2 access type.
	EnvVarDummyNodeVector2
	// EnvVarDummyNodeVector3 defines a dummy vector node 3 access type.
	EnvVarDummyNodeVector3
	// EnvVarDummyNodeVector8000 defines a dummy vector node 8000 access type.
	EnvVarDummyNodeVector8000
	// EnvVarDummyNodeVector8001 defines a dummy vector node 8001 access type.
	EnvVarDummyNodeVector8001
	// EnvVarDummyNodeVector8002 defines a dummy vector node 8002 access type.
	EnvVarDummyNodeVector8002
	// EnvVarDummyNodeVector8003 defines a dummy vector node 8003 access type.
	EnvVarDummyNodeVector8003
)

// EnvVar definition:
//
// In the environment variables section the environment variables for the usage in
// system simulation and remaining bus simulation tools are defined.
//
// environment_variables = {environment_variable}
//
// environment_variable = 'EV_' env_var_name ':' env_var_type '[' minimum '|' maximum ']'
// unit initial_value ev_id access_type access_node {',' access_node } ';' ;
//
// env_var_name = DBC_identifier ;
//
// env_var_type = '0' | '1' | '2' ; (* 0=integer, 1=float, 2=string *)
//
// minimum = double ;  maximum = double ;  initial_value = double ;
//
// ev_id = unsigned_integer ; (* obsolete *)
//
// access_type = 'DUMMY_NODE_VECTOR0' | 'DUMMY_NODE_VECTOR1' |
// 'DUMMY_NODE_VECTOR2' | 'DUMMY_NODE_VECTOR3' |
// 'DUMMY_NODE_VECTOR8000' | 'DUMMY_NODE_VECTOR8001' |
// 'DUMMY_NODE_VECTOR8002' | 'DUMMY_NODE_VECTOR8003';
// (* 0=unrestricted, 1=read, 2=write, 3=readWrite, if the value behind 'DUMMY_NODE_VECTOR'
// is OR-ed with 0x8000, the value type is always string. *)
//
// access_node = node_name | 'VECTOR__XXX' ;
type EnvVar struct {
	Location *Location

	Name         string
	Type         EnvVarType
	Min          float64
	Max          float64
	Unit         string
	InitialValue float64
	ID           uint32
	AccessType   EnvVarAccessType
	AccessNodes  []string
}

// EnvVarData definition:
//
// The entries in the environment variables data section define the environments
// listed here as being of the data type "Data". Environment variables of this type can
// store an arbitrary binary data of the given length. The length is given in bytes.
//
// environment_variables_data = environment_variable_data ;
//
// environment_variable_data = 'ENVVAR_DATA_' env_var_name ':' data_size ';' ;
//
// data_size = unsigned_integer ;
type EnvVarData struct {
	Location *Location

	EnvVarName string
	DataSize   uint32
}

// SignalType definition:
//
// Signal types are used to define the common properties of several signals. They
// are normally not used in DBC files.
//
// signal_types = {signal_type} ;
//
// signal_type = 'SGTYPE_' signal_type_name ':' signal_size '@'
// byte_order value_type '(' factor ',' offset ')' '[' minimum '|'
// maximum ']' unit default_value ',' value_table ';' ;
//
// signal_type_name = DBC_identifier ; default_value = double ;
//
// value_table = value_table_name ;
type SignalType struct {
	Location *Location

	TypeName       string
	Size           uint32
	ByteOrder      SignalByteOrder
	ValueType      SignalValueType
	Factor         float64
	Offset         float64
	Min            float64
	Max            float64
	Unit           string
	DefaultValue   float64
	ValueTableName string
}

// SignalTypeRef definition:
//
// signal_type_refs = {signal_type_ref} ;
//
// signal_type_ref = 'SGTYPE_' message_id signal_name ':' signal_type_name ';' ;
type SignalTypeRef struct {
	Location *Location

	TypeName   string
	MessageID  uint32
	SignalName string
}

// SignalGroup definition:
//
// Signal groups are used to define a group of signals within a messages,
// e.g. to define that the signals of a group have to be updated in common.
//
// signal_groups = 'SIG_GROUP_' message_id signal_group_name repetitions ':' { signal_name } ';' ;
//
// signal_group_name = DBC_identifier ; repetitions = unsigned_integer ;
type SignalGroup struct {
	Location *Location

	MessageID   uint32
	GroupName   string
	Repetitions uint32
	SignalNames []string
}

// CommentKind defines the kind of a [Comment].
type CommentKind uint

const (
	// CommentGeneral defines a general comment.
	CommentGeneral CommentKind = iota
	// CommentNode defines a node comment.
	CommentNode
	// CommentMessage defines a message comment.
	CommentMessage
	// CommentSignal defines a signal comment.
	CommentSignal
	// CommentEnvVar defines an envvar comment.
	CommentEnvVar
)

// Comment definition:
//
// The comment section contains the object comments. For each object having a
// comment, an entry with the object's type identification is defined in this section.
//
// comments = {comment} ;
//
// comment = 'CM_' (char_string |
// 'BU_' node_name char_string |
// 'BO_' message_id char_string |
// 'SG_' message_id signal_name char_string |
// 'EV_' env_var_name char_string)
// ';' ;
type Comment struct {
	Location *Location

	Kind       CommentKind
	Text       string
	NodeName   string
	MessageID  uint32
	SignalName string
	EnvVarName string
}

// AttributeKind defines the kind of an [Attribute].
type AttributeKind uint

const (
	// AttributeGeneral defines a general attribute.
	AttributeGeneral AttributeKind = iota
	// AttributeNode defines a node attribute.
	AttributeNode
	// AttributeMessage defines a message attribute.
	AttributeMessage
	// AttributeSignal defines a signal attribute.
	AttributeSignal
	// AttributeEnvVar defines an envvar attribute.
	AttributeEnvVar
)

// AttributeType defines the type of an [Attribute].
type AttributeType uint

const (
	// AttributeInt defines an int attribute.
	AttributeInt AttributeType = iota
	// AttributeFloat defines a float attribute.
	AttributeFloat
	// AttributeString defines a string attribute.
	AttributeString
	// AttributeEnum defines an enum attribute.
	AttributeEnum
	// AttributeHex defines an hex attribute.
	AttributeHex
)

// Attribute definition:
//
// User defined attributes are a means to extend the object properties of the DBC
// file. These additional attributes have to be defined using an attribute definition with
// an attribute default value. For each object having a value defined for the attribute
// an attribute value entry has to be defined. If no attribute value entry is defined for
// an object the value of the object's attribute is the attribute's default.
//
// attribute_definitions = { attribute_definition } ;
//
// attribute_definition = 'BA_DEF_' object_type attribute_name attribute_value_type ';' ;
//
// object_type = ” | 'BU_' | 'BO_' | 'SG_' | 'EV_' ;
// attribute_name = '"' DBC_identifier '"' ;
//
// attribute_value_type = 'INT' signed_integer signed_integer |
// 'HEX' signed_integer signed_integer |
// 'FLOAT' double double |
// 'STRING' |
// 'ENUM' [char_string {',' char_string}]
type Attribute struct {
	Location *Location

	Kind       AttributeKind
	Type       AttributeType
	Name       string
	MinInt     int
	MaxInt     int
	MinHex     int
	MaxHex     int
	MinFloat   float64
	MaxFloat   float64
	EnumValues []string
}

// AttributeDefaultType defines the type of an [AttributeDefault].
type AttributeDefaultType uint

const (
	// AttributeDefaultInt defines an int or enum attribute default.
	AttributeDefaultInt AttributeDefaultType = iota
	// AttributeDefaultString defines a string attribute default.
	AttributeDefaultString
	// AttributeDefaultFloat defines a float attribute default.
	AttributeDefaultFloat
	// AttributeDefaultHex defines an hex attribute default.
	AttributeDefaultHex
)

// AttributeDefault definition:
//
// attribute_defaults = { attribute_default } ;
//
// attribute_default = 'BA_DEF_DEF_' attribute_name attribute_value ';' ;
//
// attribute_value = unsigned_integer | signed_integer | double | char_string ;
type AttributeDefault struct {
	Location *Location

	Type          AttributeDefaultType
	AttributeName string
	ValueString   string
	ValueInt      int
	ValueHex      int
	ValueFloat    float64
}

// AttributeValueType defines the type of an [AttributeValue].
type AttributeValueType uint

const (
	// AttributeValueInt defines an int or enum attribute value.
	AttributeValueInt AttributeValueType = iota
	// AttributeValueString defines a string attribute value.
	AttributeValueString
	// AttributeValueFloat defines a float attribute value.
	AttributeValueFloat
	// AttributeValueHex defines an hex attribute value.
	AttributeValueHex
)

// AttributeValue definition:
//
// attribute_values = { attribute_value_for_object } ;
//
// attribute_value_for_object = 'BA_' attribute_name (attribute_value |
// 'BU_' node_name attribute_value |
// 'BO_' message_id attribute_value |
// 'SG_' message_id signal_name attribute_value |
// 'EV_' env_var_name attribute_value)
// ';' ;
type AttributeValue struct {
	Location *Location

	AttributeKind AttributeKind
	Type          AttributeValueType
	AttributeName string
	NodeName      string
	MessageID     uint32
	SignalName    string
	EnvVarName    string
	ValueString   string
	ValueInt      int
	ValueHex      int
	ValueFloat    float64
}

// ExtendedMux definition:
//
// Extended multiplexing allows defining more than one multiplexer switch within one
// message. Further it allows the usage of more than one multiplexer switch value for
// each multiplexed signal.
// The extended multiplexing section contains multiplexed signals for which following
// conditions were fulfilled:
//   - The multiplexed signal is multiplexed by more than one multiplexer switch value
//   - The multiplexed signal belongs to a message which contains more than one multiplexor switch
//
// extended multiplexing = {multiplexed signal} ;
//
// multiplexed signal = SG_MUL_VAL_ message_id multiplexed_signal_name
// multiplexor_switch_name multiplexor_value_ranges ';' ;
//
// message_id = unsigned_integer ; multiplexed_signal_name = DBC_identifier ; multiplexor_switch_name = DBC_identifier ;
type ExtendedMux struct {
	Location *Location

	MessageID       uint32
	MultiplexorName string
	MultiplexedName string
	Ranges          []*ExtendedMuxRange
}

// ExtendedMuxRange definition:
//
// multiplexor_value_ranges = {multiplexor_value_range} ;
//
// multiplexor_value_range = unsigned_integer '-' unsigned_integer ;
type ExtendedMuxRange struct {
	Location *Location

	From uint32
	To   uint32
}
