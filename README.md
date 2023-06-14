# blockchain-and-blockwallet
一个用go语言开发的建议区块链及区块链钱包
#### 一、启动区块链节点，命令为：go run ./BlockServer

![1686717667788](C:\Users\89517\AppData\Roaming\Typora\typora-user-images\1686717667788.png)

我对区块链及区块做了持久化，如果blockchain.db没有保存的数据，则创建一条新链，如果有则加载区块链及区块，继续挖矿

#### 二、启动钱包服务，命令为：go run ./WalletServer

![1686717889566](C:\Users\89517\AppData\Roaming\Typora\typora-user-images\1686717889566.png)

#### 三、导入钱包插件到浏览器

![1686717955469](C:\Users\89517\AppData\Roaming\Typora\typora-user-images\1686717955469.png)

#### 核心功能：

1.区块链持久化

2.发起交易，交易的签名验证（使用以太坊源码的验证方式）

3.pow工作量证明（使用计算难度值的方式，实现了难度调整）

5.钱包的账户生成与加载（使用以太坊源码生成账户的方式）

6.查询账户余额

7.实现了getBlockByHash，getBlockByNumber，getTransactionByHash

8.区块链节点之间的数据同步与通信（待完成）
