package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"./counter"
)

func main() {
	// 连接到Sepolia测试网络
	fmt.Println("🚀 开始连接Sepolia测试网络...")
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		log.Fatal("❌ 连接Sepolia测试网络失败:", err)
	}
	fmt.Println("✅ 成功连接到Sepolia测试网络！")

	// 私钥和合约地址（需要替换为实际的私钥和部署后的合约地址）
	privateKeyHex := "您的私钥（不含0x前缀）"
	contractAddress := "0x合约地址"

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("❌ 解析私钥失败:", err)
	}

	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("❌ 获取公钥失败")
	}

	// 获取发送方地址
	senderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("📤 发送方地址: %s\n", senderAddress.Hex())

	// 创建交易选项
	txCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	// 设置交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("❌ 创建交易选项失败:", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // 不发送以太币
	auth.GasLimit = uint64(300000)  // 设置Gas限制
	auth.GasPrice = gasPrice

	// 部署合约
	fmt.Println("🔧 开始部署合约...")
	contractAddr, tx, instance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatal("❌ 部署合约失败:", err)
	}
	fmt.Printf("✅ 合约部署成功！\n")
	fmt.Printf("📊 交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("🏠 合约地址: %s\n", contractAddr.Hex())

	// 等待交易确认
	fmt.Println("⏳ 等待交易确认...")
	txReceipt, err := bind.WaitMined(txCtx, client, tx)
	if err != nil {
		log.Fatal("❌ 等待交易确认失败:", err)
	}
	fmt.Printf("✅ 交易已确认，区块高度: %d\n", txReceipt.BlockNumber)

	// 与合约交互
	fmt.Println("🔄 开始与合约交互...")

	// 获取当前计数
	count, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal("❌ 获取计数失败:", err)
	}
	fmt.Printf("📊 当前计数: %d\n", count)

	// 调用increment方法
	fmt.Println("➕ 调用increment方法...")
	auth.Nonce = big.NewInt(int64(nonce + 1))  // 更新nonce
	incrementTx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal("❌ 调用increment方法失败:", err)
	}
	fmt.Printf("✅ increment调用成功，交易哈希: %s\n", incrementTx.Hash().Hex())

	// 等待increment交易确认
	_, err = bind.WaitMined(txCtx, client, incrementTx)
	if err != nil {
		log.Fatal("❌ 等待increment交易确认失败:", err)
	}

	// 再次获取计数
	updatedCount, err := instance.GetCount(nil)
	if err != nil {
		log.Fatal("❌ 获取更新后的计数失败:", err)
	}
	fmt.Printf("📊 更新后的计数: %d\n", updatedCount)

	// 获取合约所有者
	owner, err := instance.Owner(nil)
	if err != nil {
		log.Fatal("❌ 获取合约所有者失败:", err)
	}
	fmt.Printf("👑 合约所有者: %s\n", owner.Hex())

	fmt.Println("🎉 合约交互完成！")
}
