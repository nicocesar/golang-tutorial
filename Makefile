all: golang-tutorial

tools/solc-0.8.10:
	mkdir -p tools
	wget -O tools/solc-0.8.10 https://github.com/ethereum/solidity/releases/download/v0.8.10/solc-static-linux
	chmod +x tools/solc-0.8.10

tools/solc-0.7.6:
	mkdir -p tools
	wget -O tools/solc-0.7.6 https://github.com/ethereum/solidity/releases/download/v0.7.6/solc-static-linux
	chmod +x tools/solc-0.7.6

tools/abigen:
	mkdir -p tools
	wget https://gethstore.blob.core.windows.net/builds/geth-alltools-linux-amd64-1.10.13-7a0c19f8.tar.gz
	md5sum --check geth.md5sum
	tar xvzf geth-alltools-linux-amd64-1.10.13-7a0c19f8.tar.gz -C tools --strip-components 1

solidity_contracts/openzeppelin-contracts:
	mkdir -p solidity_contracts
	cd solidity_contracts && if [ ! -d openzeppelin-contracts ] ; then git clone https://github.com/OpenZeppelin/openzeppelin-contracts.git ; else cd openzeppelin-contracts ; git pull ; fi

lib/contracts/erc20/erc20.go: tools/solc-0.8.10 tools/abigen solidity_contracts/openzeppelin-contracts
	mkdir -p lib/contracts/erc20
	tools/abigen --solc ./tools/solc-0.8.10 --sol solidity_contracts/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol --pkg erc20 --out lib/contracts/erc20/erc20.go
	go mod tidy

golang-tutorial: lib/contracts/erc20/erc20.go main.go
	go build -o golang-tutorial main.go 

clean:
	rm -rf geth-* tools	solidity_contracts golang-tutorial