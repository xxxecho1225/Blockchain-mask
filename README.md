# Blockchain-mask
使用go语言创建一个简单的本地区块链钱包插件，可以在浏览器插件使用。本插件功能包括区块模块、区块链模块、账户模块、交易模块和钱包构建模块。区块模块包括对区块数据结构的构建 ，例如随机数生成、区块时间戳、区块号、区块难度、夫区块哈希、本区块哈希、本区块交易数量。区块链模块包括区块链数据结构、区块链查询、和挖矿。功能例如将新的交易加入交易池、从交易池中删除已经被确认的交易等、根据区块ID输出该区块结构体内容、根据区块hash输出该区块结构体内容和实现难度调整的功能，非前导零的方式、区块奖励发放（1e+20）。账户模块包括账户模块、账户加载、账户查询、账户转账。交易模块包括交易数据结构、交易验证和交易查询。钱包构建模块包括对钱包适应插件页面的构建、钱包与后端交互以及对钱包功能的实现。
