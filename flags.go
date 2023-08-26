package messagix

type ConnectFlags struct {
	Username      bool
	Password      bool
	Retain    bool
	QoS       uint8 // 0, 1, 2, or 3
	CleanSession  bool
}

func CreateConnectFlagByte(flags ConnectFlags) uint8 {
	var connectFlagsByte uint8

	if flags.Username {
		connectFlagsByte |= 0x80
	}
	if flags.Password {
		connectFlagsByte |= 0x40
	}
	if flags.Retain {
		connectFlagsByte |= 0x20
	}

	connectFlagsByte |= (flags.QoS << 3) & 0x18
	if flags.CleanSession {
		connectFlagsByte |= 0x02
	}

	return connectFlagsByte
}