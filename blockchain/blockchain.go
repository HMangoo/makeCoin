/*
	Singleton pattern
	-> 우리의 application 내에서 언제든지 blockchain의 단 하나의 instance만을 공유하는 방법
	-> 이 변수의 instance를 직접 공유하지 않고, 이 변수의 instance를 대신해서 드러내주는 function을 생성하는 것
	-> 다른 패키지에서 우리의 blockchain이 어떻게 드러날 지 제어할 수 있음
*/
/*
	Sync package
	 -> 동기적으로 처리해야하는 부분을 도와주는 패키지

	 we use "Once"
	 	- 병렬로 실행하고 있는 프로그램이 몇 개이던 간에 (thread, goroutine이 몇 개든)
		 코드를 단 한 번만 실행시키고 싶을 때
		- Do(function): 단 한 번만 호출되도록 해주는 함수
*/

package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string	`json:"data"`
	Hash     string	`json:"hash"`
	PrevHash string `json:"PreHash,omitempty"`
	Height	 int	`json:"height"`
}

type bloackchain struct {
	blocks []*Block // slice of block
}

var b *bloackchain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data+b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlcoks := len(GetBlockchain().blocks)
	if totalBlcoks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlcoks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks)+1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *bloackchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *bloackchain {
	if b == nil {
		once.Do(func() {
			b = &bloackchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func (b *bloackchain)AllBlocks() []*Block {
	return b.blocks
}

var ErrNotFound = errors.New("block not found")

func (b *bloackchain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}