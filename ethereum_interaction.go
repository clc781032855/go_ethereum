package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//连接到Sepolia 测试网络。
	fmt.Println("🚀 开始连接Sepolia 测试网络...")
	client ,err := ethclient.Dial("https://sepolia.infura.io/v3/931e1d7d357b47f2a1015ceed2746901")
	if err != nil {
		log.Fatal("❌ 连接失败:", err)
	}
	defer client.Close()
	fmt.Println("✅ 成功连接到Sepolia 测试网络！")

	//查询指定区块号的区块信息
	blockNumber := big.NewInt(666)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	block , err := client.BlockByNumber(ctx, blockNumber)

	if err != nil {
		log.Fatal("❌ 获取区块信息失败:", err)
	}

	//访问区块的各种属性
	fmt.Println("🔍 区块信息:")
	fmt.Printf("🔗 区块号: %s\n", block.Number().String())
	fmt.Printf("🔗 区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("🔗 区块时间戳: %s\n", block.Time())
	fmt.Printf("🔗 交易数量: %d\n", len(block.Transactions()))

	//发送交易
	fmt.Println("🚀 开始发送交易...")

	// 注意：实际使用时请替换为您自己的私钥和接收地址
	privateKeyHex := "11fd99fad093d5c8eab5f5ab3af1e263d7245df7ebb10a19954ef154fd3d6ac7"
	recipientAddress := "0x4558fa23D70a875b78C295f885AD718D8B6f7110"

	// 解析私钥，将私钥从十六进制字符串转换为ECDSA私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("❌ 私钥解析失败:", err)
	}

	// 获取发送方地址，从公钥中提取地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("❌ 无法获取公钥")
	}
	senderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("📤 发送方地址: %s\n", senderAddress.Hex())

	// 解析接收地址
	receiverAddr := common.HexToAddress(recipientAddress)
	fmt.Printf("📥 接收方地址: %s\n", receiverAddr.Hex())

	// 设置转账金额 (0.001 ETH = 1000000000000000 Wei)
	amount := big.NewInt(1000000000000000)
	fmt.Printf("💰 转账金额: %s Wei\n", amount.String())

	// 创建上下文，设置5秒超时
	txCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取链ID
	chainID, err := client.ChainID(txCtx)
	if err != nil {
		log.Fatal("❌ 获取链ID失败:", err)
	}

	// 获取Nonce
	nonce, err := client.PendingNonceAt(txCtx, senderAddress)
	if err != nil {
		log.Fatal("❌ 获取Nonce失败:", err)
	}

	// 获取Gas价格
	gasPrice, err := client.SuggestGasPrice(txCtx)
	if err != nil {
		log.Fatal("❌ 获取Gas价格失败:", err)
	}

	// 设置Gas限制
	gasLimit := uint64(21000) // 标准转账交易的Gas限制

	// 创建交易对象
	tx := types.NewTransaction(
		nonce,
		receiverAddr,
		amount,
		gasLimit,
		gasPrice,
		[]byte{}, // 空数据，因为是简单转账
	)

	// 签名交易，需要的参数，交易对象，签名算法，私钥
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("❌ 签名交易失败:", err)
	}

	// 发送交易到网络
	err = client.SendTransaction(txCtx, signedTx)
	if err != nil {
		log.Fatal("❌ 发送交易失败:", err)
	}

	// 输出交易哈希
	txHash := signedTx.Hash().Hex()
	fmt.Printf("✅ 交易发送成功！\n")
	fmt.Printf("📊 交易哈希: %s\n", txHash)
	fmt.Printf("🔗 交易链接: https://sepolia.etherscan.io/tx/%s\n", txHash)

}

