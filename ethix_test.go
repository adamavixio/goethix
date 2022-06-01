package goethix

import "testing"

const (
	ADDRESS = "http://127.0.0.1:7545"

	ACC_1_PRIVATE = "84bf34442cd7d3dee01465398716034734aa445ac1fdef3dd4836cea6b5f4e37"
	ACC_1_ADDRESS = "0xbA3A2950a2178153716b8f1917992746AA166A31"

	ACC_2_PRIVATE = "15699e8369f6c1cb3019aa5119a6fd8491a2e259dc32c1fc00b1bbd018876c1b"
	ACC_2_ADDRESS = "0xe40bBa4959E9CB1bb614eC2228B15b7c554bc44A"

	CONTRACT = "0x7DcC70f68EDaEd301969382BD87cE7ddd899417B"
)

func TestBalance(t *testing.T) {
	ethix := NewEthix()
	ethix.Dial(ADDRESS)

	val, err := ethix.Balance(ACC_1_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	t.Logf("ACCOUNT 1 BALANCE: %s", val)
}

func TestTransfer(t *testing.T) {
	ethix := NewEthix()
	ethix.Dial(ADDRESS)

	bal1, err := ethix.Balance(ACC_1_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	bal2, err := ethix.Balance(ACC_2_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	t.Logf("ACCOUNT 1 BALANCE: %v, ACCOUNT 2 BALANCE: %v", bal1, bal2)

	_, err = ethix.Transfer("10", ACC_1_PRIVATE, ACC_2_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	bal1, err = ethix.Balance(ACC_1_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	bal2, err = ethix.Balance(ACC_2_ADDRESS)
	if err != nil {
		t.Error(err)
	}

	t.Logf("ACCOUNT 1 BALANCE: %v, ACCOUNT 2 BALANCE: %v", bal1, bal2)
}
